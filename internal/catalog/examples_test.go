package catalog

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/detect"
	"github.com/ninegene/shell-explain/internal/examplecheck"
)

// Executes only this repository's reviewed examples in temporary fixtures.
// It is not an execution feature exposed by explain.
func TestAuthoredExamples(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	implementations := []string{"bash", "apple", "gnu", "busybox", "mawk"}
	for _, impl := range implementations {
		for _, entry := range c.Entries {
			selected, err := entry.Select(impl, "")
			if err != nil {
				continue
			}
			t.Run(entry.ID+"/"+impl, func(t *testing.T) {
				shell := "/bin/sh"
				executable := ""
				if impl == "bash" {
					var err error
					shell, err = exec.LookPath("bash")
					if err != nil {
						t.Skip("Bash is unavailable")
					}
				} else {
					executable = exampleTool(detect.Detector{OS: runtime.GOOS, Path: os.Getenv("PATH")}, entry.ID, impl)
					if executable == "" {
						t.Skipf("%s %s is unavailable or unidentified", impl, entry.ID)
					}
				}
				for _, section := range selected {
					for i, code := range section.Code {
						t.Run(section.ID+string(rune('a'+i)), func(t *testing.T) {
							if examplecheck.NotRun(code) {
								t.Skip("reviewed only: the example never ends")
							}
							dir := t.TempDir()
							if err := examplecheck.Setup(dir, entry.ID); err != nil {
								t.Fatal(err)
							}
							toolPath := "/usr/bin:/bin"
							toolDir := ""
							if executable != "" {
								tools := t.TempDir()
								toolDir = tools
								if err := os.Symlink(executable, filepath.Join(tools, entry.ID)); err != nil {
									t.Fatal(err)
								}
								for _, companion := range examplecheck.Companions[entry.ID] {
									path := exampleTool(detect.Detector{OS: runtime.GOOS, Path: os.Getenv("PATH")}, companion, impl)
									if path == "" {
										t.Skipf("%s %s is unavailable or unidentified", impl, companion)
									}
									if err := os.Symlink(path, filepath.Join(tools, companion)); err != nil {
										t.Fatal(err)
									}
								}
								toolPath = tools + ":" + toolPath
							}
							script := "cd .\n" + code
							if executable != "" && slices.Contains(entry.Implementations, "bash") {
								// The shell would run its own builtin; a function makes the
								// example run the executable that this view documents.
								script = entry.ID + "() { \"" + filepath.Join(toolDir, entry.ID) + "\" \"$@\"; }\n" + script
							}
							cmd := exec.Command(shell, "-c", script)
							cmd.Dir = dir
							cmd.Env = []string{"PATH=" + toolPath, "HOME=" + dir, "OLDPWD=" + dir, "LC_ALL=C"}
							output, err := cmd.CombinedOutput()
							if err != nil {
								t.Fatalf("%s: %v: %s", code, err, output)
							}
							if err := examplecheck.Verify(code, dir, string(output)); err != nil {
								t.Fatalf("%s: %v", code, err)
							}
						})
					}
				}
			})
		}
	}
}

// Identify each command independently; a prefix or the host OS is not proof of GNU.
func exampleTool(d detect.Detector, command, implementation string) string {
	candidates := []string{command}
	if implementation == "gnu" {
		candidates = []string{"g" + command, command}
	}
	if implementation == "mawk" {
		// mawk is installed as mawk; a plain awk may be any implementation.
		candidates = []string{"mawk", command}
	}
	if implementation == "apple" {
		if d.OS != "darwin" {
			return ""
		}
		candidates = []string{"/usr/bin/" + command, "/bin/" + command, "/sbin/" + command, "/usr/sbin/" + command}
	}
	for _, name := range candidates {
		selected, err := d.Select(context.Background(), detect.Request{Name: name, Command: command})
		if err == nil && selected.Implementation == implementation && selected.Executable != "" {
			return selected.Executable
		}
	}
	return ""
}

func TestExampleToolPartialInstallation(t *testing.T) {
	dir := t.TempDir()
	// A script is never probed, but its BusyBox-named target is labeled inference.
	// This models Linux without GNU tools, plus a misleading GNU-prefixed name.
	binary := filepath.Join(dir, "busybox")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 1\n"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"grep", "ggrep"} {
		if err := os.Symlink(binary, filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	d := detect.Detector{OS: "linux", Path: dir}
	if got := exampleTool(d, "grep", "busybox"); got != filepath.Join(dir, "grep") {
		t.Fatalf("BusyBox grep: %q", got)
	}
	for _, command := range []string{"grep", "sed", "find"} {
		if got := exampleTool(d, command, "gnu"); got != "" {
			t.Fatalf("%s falsely selected as GNU: %s", command, got)
		}
	}
	if got := exampleTool(d, "sed", "busybox"); got != "" {
		t.Fatalf("missing sed inherited grep coverage: %s", got)
	}
}

// Every authored example needs a recorded expected result, even when the host
// lacks the tool to run it, so new examples cannot ship unchecked.
func TestEveryExampleHasExpectation(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range c.Entries {
		for _, s := range e.Sections {
			for _, code := range s.Code {
				if !examplecheck.HasExpectation(code) {
					t.Errorf("%s/%s: no expected result for %q", e.ID, s.ID, code)
				}
			}
		}
	}
}

// A page that leaves out the Apple view claims macOS ships no such tool. Check
// that on a Mac, so the claim cannot go stale unnoticed.
func TestMissingAppleViewsMatchMacOS(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("only macOS can show what it ships")
	}
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range c.Entries {
		if !slices.Contains(e.Implementations, "gnu") || slices.Contains(e.Implementations, "apple") {
			continue
		}
		for _, dir := range []string{"/usr/bin/", "/bin/", "/usr/sbin/", "/sbin/"} {
			if _, err := os.Stat(dir + e.ID); err == nil {
				t.Errorf("%s has no Apple view, but macOS ships %s%s", e.ID, dir, e.ID)
			}
		}
	}
}
