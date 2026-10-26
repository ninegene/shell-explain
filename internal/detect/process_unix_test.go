//go:build darwin || linux

package detect

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func init() {
	mode := os.Getenv("SHELL_EXPLAIN_GROUP_FIXTURE")
	if mode == "" || len(os.Args) != 2 {
		return
	}
	marker := os.Getenv("SHELL_EXPLAIN_GROUP_MARKER")
	if os.Args[1] == "--probe-child" {
		_ = os.WriteFile(marker, []byte(strconv.Itoa(os.Getpid())), 0600)
		for i := 0; i < 500; i++ {
			_ = os.WriteFile(marker+".heartbeat", []byte(strconv.Itoa(i)), 0600)
			time.Sleep(20 * time.Millisecond)
		}
		os.Exit(0)
	}
	if os.Args[1] != "--version" {
		return
	}
	executable, _ := os.Executable()
	child := exec.Command(executable, "--probe-child")
	child.Stdout, child.Stderr = os.Stdout, os.Stderr
	if strings.HasPrefix(mode, "escaped") {
		child.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	}
	if err := child.Start(); err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		if _, err := os.Stat(marker + ".heartbeat"); err == nil {
			break
		}
		time.Sleep(time.Millisecond)
	}
	fmt.Println("grep (GNU grep) 3.12")
	if strings.HasSuffix(mode, "timeout") {
		time.Sleep(5 * time.Second)
	}
	os.Exit(0)
}

func TestProbeDescendantCleanupAndWaitDelay(t *testing.T) {
	t.Setenv("GORACE", "atexit_sleep_ms=0")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"group-timeout", "group-exit", "escaped-timeout", "escaped-exit"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			marker := filepath.Join(dir, "child")
			t.Setenv("SHELL_EXPLAIN_GROUP_FIXTURE", mode)
			t.Setenv("SHELL_EXPLAIN_GROUP_MARKER", marker)
			if err := os.Symlink(executable, filepath.Join(dir, "grep")); err != nil {
				t.Fatal(err)
			}
			// Always clean up an escaped fixture, even if an assertion fails.
			defer func() {
				data, _ := os.ReadFile(marker)
				pid, _ := strconv.Atoi(string(data))
				if pid > 0 {
					_ = syscall.Kill(pid, syscall.SIGKILL)
				}
			}()
			start := time.Now()
			s, err := (Detector{OS: "linux", Path: dir}).Select(context.Background(), Request{Name: "grep", Command: "grep"})
			if err != nil {
				t.Fatal(err)
			}
			if elapsed := time.Since(start); elapsed > 1500*time.Millisecond {
				t.Fatalf("probe blocked for %s", elapsed)
			}
			if s.Confidence == "detected" || !strings.Contains(strings.Join(s.Warnings, " "), "timed out") {
				t.Fatalf("incomplete probe trusted: %+v", s)
			}
			before, err := os.ReadFile(marker + ".heartbeat")
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(100 * time.Millisecond)
			after, _ := os.ReadFile(marker + ".heartbeat")
			if strings.HasPrefix(mode, "group") && string(before) != string(after) {
				t.Fatal("ordinary descendant survived group cleanup")
			}
			if strings.HasPrefix(mode, "escaped") && string(before) == string(after) {
				t.Fatal("escaped fixture did not keep running; WaitDelay case not exercised")
			}
		})
	}
}
