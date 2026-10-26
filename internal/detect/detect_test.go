package detect

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// A native test binary stands in for a real executable, without a shell wrapper.
	if mode := os.Getenv("SHELL_EXPLAIN_PROBE_FIXTURE"); mode != "" && len(os.Args) == 2 && os.Args[1] == "--version" {
		switch mode {
		case "gnu":
			fmt.Println("grep (GNU grep) 3.12")
		case "apple":
			fmt.Println("grep (BSD grep, GNU compatible) 2.6.0-FreeBSD")
		case "busybox":
			fmt.Println("grep: unrecognized option --version\nBusyBox v1.37.0 multi-call binary")
		case "timeout":
			time.Sleep(3 * time.Second)
		case "flood":
			fmt.Println("grep (GNU grep) 3.12")
			fmt.Print(strings.Repeat("x", 100000))
		default:
			fmt.Println("unknown implementation")
			os.Exit(2)
		}
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestNativeProbeResultsAndTimeout(t *testing.T) {
	t.Setenv("GORACE", "atexit_sleep_ms=0")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err = os.Symlink(executable, filepath.Join(dir, "grep")); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ mode, impl, confidence string }{
		{"gnu", "gnu", "detected"}, {"apple", "apple", "detected"}, {"busybox", "busybox", "detected"}, {"unknown", "common", "unknown"}, {"timeout", "common", "unknown"}, {"flood", "gnu", "detected"},
	} {
		t.Run(tt.mode, func(t *testing.T) {
			t.Setenv("SHELL_EXPLAIN_PROBE_FIXTURE", tt.mode)
			start := time.Now()
			s, err := (Detector{OS: "darwin", Path: dir}).Select(context.Background(), Request{Name: "grep", Command: "grep"})
			if err != nil || s.Implementation != tt.impl || s.Confidence != tt.confidence {
				t.Fatalf("%+v %v", s, err)
			}
			if tt.mode == "timeout" && (time.Since(start) > 2*time.Second || !strings.Contains(strings.Join(s.Warnings, " "), "timed out")) {
				t.Fatal("probe timeout did not bound execution")
			}
		})
	}
}

func TestSelectionBoundaries(t *testing.T) {
	d := Detector{OS: "linux", Path: t.TempDir()}
	for _, tt := range []struct {
		r                Request
		impl, confidence string
	}{
		{Request{Name: "grep", Command: "grep"}, "common", "unknown"},
		{Request{Name: "grep", Command: "grep", Platform: "linux"}, "common", "target"},
		{Request{Name: "grep", Command: "grep", Platform: "macos"}, "apple", "target"},
		{Request{Name: "grep", Command: "grep", Implementation: "gnu"}, "gnu", "explicit"},
		{Request{Name: "cd", Command: "cd"}, "common", "unknown"},
		{Request{Name: "cd", Command: "cd", Implementation: "bash"}, "bash", "explicit"},
	} {
		s, err := d.Select(context.Background(), tt.r)
		if err != nil || s.Implementation != tt.impl || s.Confidence != tt.confidence {
			t.Fatalf("%+v: %+v %v", tt.r, s, err)
		}
	}
	if _, err := d.Select(context.Background(), Request{Command: "grep", Platform: "linux", Implementation: "gnu"}); err == nil {
		t.Fatal("conflicting flags accepted")
	}
	if _, err := d.Select(context.Background(), Request{Command: "grep", Platform: "other"}); err == nil {
		t.Fatal("unknown platform")
	}
}

func TestScriptNotExecuted(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "executed")
	script := filepath.Join(dir, "grep")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nprintf ran > '"+marker+"'\necho 'GNU grep 99.0'\n"), 0755); err != nil {
		t.Fatal(err)
	}
	s, err := (Detector{OS: "darwin", Path: dir}).Select(context.Background(), Request{Name: "grep", Command: "grep"})
	if err != nil || s.Implementation != "common" {
		t.Fatal(s, err)
	}
	if _, err = os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("untrusted script was executed")
	}
	if !strings.Contains(strings.Join(s.Warnings, " "), "skipped") {
		t.Fatal("probe policy not reported")
	}
}

func TestBusyBoxSymlinkAndNoProbe(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "busybox")
	if err := os.WriteFile(binary, []byte("not runnable"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(binary, filepath.Join(dir, "grep")); err != nil {
		t.Fatal(err)
	}
	s, err := (Detector{OS: "linux", Path: dir}).Select(context.Background(), Request{Name: "grep", Command: "grep", NoProbe: true})
	if err != nil || s.Implementation != "busybox" || s.Confidence != "inferred" {
		t.Fatal(s, err)
	}
}

func TestRealInstalledTools(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS mixed PATH integration")
	}
	d := Detector{OS: runtime.GOOS, Path: "/usr/bin"}
	for _, command := range []string{"grep", "sed", "find"} {
		s, err := d.Select(context.Background(), Request{Name: command, Command: command})
		if err != nil || s.Implementation != "apple" {
			t.Errorf("%s: %+v %v", command, s, err)
		}
	}
	for _, command := range []string{"grep", "sed", "find"} {
		executable := "/opt/homebrew/bin/g" + command
		if _, err := os.Stat(executable); err != nil {
			t.Skip("Homebrew GNU tools unavailable")
		}
		s, err := d.Select(context.Background(), Request{Name: executable, Command: command})
		if err != nil || s.Implementation != "gnu" || s.Version == "" {
			t.Errorf("%s: %+v %v", executable, s, err)
		}
		dir := t.TempDir()
		if err := os.Symlink(executable, filepath.Join(dir, command)); err != nil {
			t.Fatal(err)
		}
		mixed := Detector{OS: "darwin", Path: dir + string(os.PathListSeparator) + "/usr/bin"}
		s, err = mixed.Select(context.Background(), Request{Name: command, Command: command})
		if err != nil || s.Implementation != "gnu" {
			t.Errorf("GNU-first PATH %s: %+v %v", command, s, err)
		}
		s, err = mixed.Select(context.Background(), Request{Name: command, Command: command, NoProbe: true})
		if err != nil || s.Confidence == "detected" {
			t.Fatal("--no-probe ignored")
		}
	}
}

func TestProbeBuffer(t *testing.T) {
	b := &limitedBuffer{}
	n, err := b.Write([]byte(strings.Repeat("x", 20000)))
	if err != nil || n != 20000 || len(b.String()) != 8192 {
		t.Fatal("buffer not bounded")
	}
	if got := firstVersion("grep (GNU grep) 3.12\nCopyright 2026"); got != "3.12" {
		t.Fatal(got)
	}
	if got := firstVersion("grep: unrecognized option --version\nBusyBox v1.37.0 multi-call binary"); got != "1.37.0" {
		t.Fatal(got)
	}
}
