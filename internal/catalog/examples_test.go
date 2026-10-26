package catalog

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ninegene/shell-explain/content"
)

// Executes only this repository's reviewed examples in temporary fixtures.
// It is not an execution feature exposed by explain.
func TestAuthoredExamples(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	implementations := []string{"bash"}
	if runtime.GOOS == "darwin" {
		implementations = append(implementations, "apple")
	}
	if _, err := os.Stat("/opt/homebrew/bin/ggrep"); err == nil {
		implementations = append(implementations, "gnu")
	} else if runtime.GOOS == "linux" {
		implementations = append(implementations, "gnu")
	}
	for _, impl := range implementations {
		for _, entry := range c.Entries {
			selected, err := entry.Select(impl, "")
			if err != nil {
				continue
			}
			for _, section := range selected {
				for i, code := range section.Code {
					t.Run(entry.ID+"/"+impl+"/"+section.ID+string(rune('a'+i)), func(t *testing.T) {
						dir := t.TempDir()
						for _, sub := range []string{"projects", "logs", "cache"} {
							if err := os.Mkdir(filepath.Join(dir, sub), 0755); err != nil {
								t.Fatal(err)
							}
						}
						for name, value := range map[string]string{"app.log": "INFO ready\nERROR broken\nerror small\nexample.org\n", "logs/server.log": "ERROR broken\n", "notes.txt": "old old\nsecond\nthird\nfourth\nfifth\n"} {
							if err := os.WriteFile(filepath.Join(dir, name), []byte(value), 0644); err != nil {
								t.Fatal(err)
							}
						}
						toolPath := "/usr/bin:/bin"
						if impl == "gnu" && runtime.GOOS == "darwin" {
							tools := t.TempDir()
							for _, name := range []string{"grep", "sed", "find"} {
								if err := os.Symlink("/opt/homebrew/bin/g"+name, filepath.Join(tools, name)); err != nil {
									t.Fatal(err)
								}
							}
							toolPath = tools + ":" + toolPath
						}
						cmd := exec.Command("/bin/bash", "--noprofile", "--norc", "-c", "cd .\n"+code)
						cmd.Dir = dir
						cmd.Env = []string{"PATH=" + toolPath, "HOME=" + dir, "OLDPWD=" + dir, "LC_ALL=C"}
						output, err := cmd.CombinedOutput()
						if err != nil {
							t.Fatalf("%s: %v: %s", code, err, output)
						}
						if entry.ID == "sed" && strings.Contains(code, "-i") {
							got, err := os.ReadFile(filepath.Join(dir, "notes.txt"))
							if err != nil || !strings.HasPrefix(string(got), "new new") {
								t.Fatal("in-place result", err, string(got))
							}
							backup, err := os.ReadFile(filepath.Join(dir, "notes.txt.bak"))
							if err != nil || !strings.HasPrefix(string(backup), "old old") {
								t.Fatal("backup not preserved")
							}
						}
						if entry.ID == "grep" && strings.Contains(code, "-n -i") && string(output) != "2:ERROR broken\n3:error small\n" {
							t.Fatalf("unexpected match result: %s", output)
						}
						if entry.ID == "sed" && strings.Contains(code, "2,4p") && string(output) != "second\nthird\nfourth\n" {
							t.Fatal("range output", string(output))
						}
					})
				}
			}
		}
	}
}
