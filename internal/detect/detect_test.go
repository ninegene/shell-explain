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
		case "apple-version":
			fmt.Println("2.3-Apple (199)")
		case "apple-gzip":
			fmt.Println("Apple gzip 479.100.4")
		case "gnu-gzip":
			fmt.Println("gzip 1.13\nCopyright (C) 2023 Free Software Foundation, Inc.")
		case "gnu-coreutils":
			fmt.Println("sort (GNU coreutils) 9.11")
		case "busybox":
			fmt.Println("grep: unrecognized option --version\nBusyBox v1.37.0 multi-call binary")
		case "gawk":
			fmt.Println("GNU Awk 5.3.2, API 4.0, PMA Avon 8-g1, (GNU MPFR 4.2.2, GNU MP 6.3.0)\nCopyright (C) 1989, 1991-2026 Free Software Foundation.")
		case "mawk":
			fmt.Println("mawk 1.3.4 20260302\nCopyright 2008-2024,2026, Thomas E. Dickey")
		case "mawk-old":
			fmt.Println("mawk: not an option: --version")
			os.Exit(2)
		case "bwk-awk":
			fmt.Println("awk version 20200816")
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

func TestCoreutilsProbesAndAppleSystemPaths(t *testing.T) {
	t.Setenv("GORACE", "atexit_sleep_ms=0")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for _, name := range []string{"sort", "gzip", "touch", "ln", "du", "basename", "dirname", "tee", "rmdir", "mktemp", "readlink", "realpath", "seq", "paste", "comm", "join", "nl", "fold", "unlisted-tool"} {
		if err = os.Symlink(executable, filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	d := Detector{OS: "darwin", Path: dir}
	for _, tt := range []struct{ mode, impl, version string }{
		{"apple-version", "apple", "2.3"}, {"gnu-coreutils", "gnu", "9.11"},
	} {
		t.Setenv("SHELL_EXPLAIN_PROBE_FIXTURE", tt.mode)
		s, err := d.Select(context.Background(), Request{Name: "sort", Command: "sort"})
		if err != nil || s.Implementation != tt.impl || s.Confidence != "detected" || s.Version != tt.version {
			t.Fatalf("%s: %+v %v", tt.mode, s, err)
		}
	}
	for _, tt := range []struct{ mode, impl, version string }{
		{"apple-gzip", "apple", "479.100.4"}, {"gnu-gzip", "gnu", "1.13"},
	} {
		t.Setenv("SHELL_EXPLAIN_PROBE_FIXTURE", tt.mode)
		s, err := d.Select(context.Background(), Request{Name: "gzip", Command: "gzip"})
		if err != nil || s.Implementation != tt.impl || s.Confidence != "detected" || s.Version != tt.version {
			t.Fatalf("%s: %+v %v", tt.mode, s, err)
		}
	}
	// Each command added for the file-path batch is probed like the other coreutils.
	for _, command := range []string{"touch", "ln", "du", "basename", "dirname", "tee", "rmdir", "mktemp", "readlink", "realpath", "seq", "paste", "comm", "join", "nl", "fold"} {
		for _, tt := range []struct{ mode, impl string }{{"apple-version", "apple"}, {"gnu-coreutils", "gnu"}} {
			t.Setenv("SHELL_EXPLAIN_PROBE_FIXTURE", tt.mode)
			s, err := d.Select(context.Background(), Request{Name: command, Command: command})
			if err != nil || s.Implementation != tt.impl || s.Confidence != "detected" {
				t.Fatalf("%s %s: %+v %v", command, tt.mode, s, err)
			}
		}
	}
	// Commands outside the probe list are never executed, even with a GNU-looking banner.
	t.Setenv("SHELL_EXPLAIN_PROBE_FIXTURE", "gnu-coreutils")
	if s, err := d.Select(context.Background(), Request{Name: "unlisted-tool", Command: "unlisted-tool"}); err != nil || s.Implementation != "common" || s.Confidence != "unknown" {
		t.Fatalf("unlisted command probed: %+v %v", s, err)
	}
	for _, path := range []string{"/usr/bin/grep", "/bin/cat"} {
		if !isAppleSystemPath(path, filepath.Base(path)) {
			t.Errorf("%s not recognized as an Apple system path", path)
		}
	}
	if isAppleSystemPath("/opt/homebrew/bin/cat", "cat") || isAppleSystemPath("/usr/local/bin/cat", "cat") {
		t.Error("third-party path treated as Apple system")
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

func TestShellBuiltinAndSingleImplementation(t *testing.T) {
	d := Detector{OS: "darwin", Path: t.TempDir()}
	s, err := d.Select(context.Background(), Request{Name: "type", Command: "type", ShellBuiltin: true})
	if err != nil || s.Implementation != "common" || s.Confidence != "unknown" || len(s.Warnings) == 0 {
		t.Fatal(s, err)
	}
	if s, err = d.Select(context.Background(), Request{Name: "type", Command: "type", ShellBuiltin: true, Platform: "macos"}); err != nil || s.Implementation != "common" {
		t.Fatal("macOS platform must not select Apple for shell builtins", s, err)
	}
	if s, err = d.Select(context.Background(), Request{Name: "type", Command: "type", ShellBuiltin: true, Platform: "bash"}); err != nil || s.Implementation != "bash" {
		t.Fatal(s, err)
	}
	zsh := Request{Name: "where", Command: "where", ShellBuiltin: true, Shells: []string{"zsh"}}
	if s, err = d.Select(context.Background(), zsh); err != nil || !strings.Contains(strings.Join(s.Warnings, " "), "--implementation zsh") || strings.Contains(strings.Join(s.Warnings, " "), "bash") {
		t.Fatal("zsh-only entry must not mention bash", s, err)
	}
	zsh.Platform = "bash"
	if _, err = d.Select(context.Background(), zsh); err == nil {
		t.Fatal("bash platform accepted for a zsh-only entry")
	}
	for _, platform := range []string{"", "macos", "linux"} {
		if s, err = d.Select(context.Background(), Request{Name: "jq", Command: "jq", Single: "upstream", Platform: platform}); err != nil || s.Implementation != "upstream" || len(s.Warnings) != 0 {
			t.Fatalf("platform %q: %+v %v", platform, s, err)
		}
	}
	if _, err = d.Select(context.Background(), Request{Name: "jq", Command: "jq", Single: "upstream", Platform: "bash"}); err == nil {
		t.Fatal("bash platform accepted for a non-builtin")
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

func TestAwkImplementationsAreProbedUnderAnyTypedName(t *testing.T) {
	t.Setenv("GORACE", "atexit_sleep_ms=0")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for _, name := range []string{"awk", "gawk", "mawk"} {
		if err = os.Symlink(executable, filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	// Name is what the user typed; Command is the page ID, which is what the probe list matches.
	for _, tt := range []struct{ command, mode, impl, confidence string }{
		{"awk", "gawk", "gnu", "detected"},
		{"gawk", "gawk", "gnu", "detected"},
		{"awk", "mawk", "mawk", "detected"},
		{"mawk", "mawk", "mawk", "detected"},
		{"awk", "busybox", "busybox", "detected"},
		// BWK awk is not one of the registered implementations; on Linux it stays unknown.
		{"awk", "bwk-awk", "common", "unknown"},
		// Old mawk has no --version, so the banner proves nothing.
		{"mawk", "mawk-old", "common", "unknown"},
	} {
		t.Run(tt.command+"/"+tt.mode, func(t *testing.T) {
			t.Setenv("SHELL_EXPLAIN_PROBE_FIXTURE", tt.mode)
			s, err := (Detector{OS: "linux", Path: dir}).Select(context.Background(), Request{Name: tt.command, Command: "awk"})
			if err != nil || s.Implementation != tt.impl || s.Confidence != tt.confidence {
				t.Fatalf("%+v %v", s, err)
			}
			if tt.impl == "mawk" && s.Version != "1.3.4" {
				t.Fatalf("version %q", s.Version)
			}
		})
	}
}
