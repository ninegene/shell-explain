// compat-check runs reviewed examples and the built CLI in a disposable Linux VM.
// It is a test helper, not a shipped command.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
)

func main() {
	if err := check(); err != nil {
		fmt.Println("LINUX_SMOKE_FAIL:", err)
		os.Exit(1)
	}
	fmt.Println("LINUX_SMOKE_PASS")
}
func check() error {
	c, err := catalog.Load(content.Files)
	if err != nil {
		return err
	}
	for _, impl := range []string{"busybox", "gnu"} {
		tools, err := os.MkdirTemp("", "tools-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tools)
		for _, name := range []string{"grep", "sed", "find"} {
			target := "/bin/busybox"
			if impl == "gnu" {
				target, err = exec.LookPath(name)
				if err != nil {
					return err
				}
			}
			if err = os.Symlink(target, filepath.Join(tools, name)); err != nil {
				return err
			}
		}
		for _, name := range []string{"grep", "sed", "find"} {
			data, err := run("/explain", []string{name, "--format=json"}, "/", tools)
			if err != nil {
				return fmt.Errorf("%s/%s CLI: %w %s", impl, name, err, data)
			}
			var result struct {
				Selection struct{ Implementation, Confidence, Version string }
				Sections  []catalog.Section
			}
			if err = json.Unmarshal(data, &result); err != nil {
				return err
			}
			if result.Selection.Implementation != impl {
				return fmt.Errorf("%s detected as %+v", impl, result.Selection)
			}
			fmt.Printf("Linux detection: %s -> %s (%s, version %s)\n", name, impl, result.Selection.Confidence, result.Selection.Version)
			entry, _ := c.Find(name)
			sections, err := entry.Select(impl, "")
			if err != nil {
				return err
			}
			for _, section := range sections {
				for _, code := range section.Code {
					dir, err := os.MkdirTemp("", "example-")
					if err != nil {
						return err
					}
					defer os.RemoveAll(dir)
					for _, sub := range []string{"logs", "cache"} {
						if err = os.Mkdir(filepath.Join(dir, sub), 0755); err != nil {
							return err
						}
					}
					for file, value := range map[string]string{"app.log": "INFO ready\nERROR broken\nerror small\nexample.org\n", "logs/server.log": "ERROR broken\n", "notes.txt": "old old\nsecond\nthird\nfourth\nfifth\n"} {
						if err = os.WriteFile(filepath.Join(dir, file), []byte(value), 0644); err != nil {
							return err
						}
					}
					output, err := run("/bin/sh", []string{"-c", code}, dir, tools)
					if err != nil {
						return fmt.Errorf("%s/%s example %s: %w %s", impl, name, code, err, output)
					}
					if name == "sed" && bytes.Contains([]byte(code), []byte("-i")) {
						changed, err := os.ReadFile(filepath.Join(dir, "notes.txt"))
						if err != nil || !bytes.HasPrefix(changed, []byte("new new")) {
							return fmt.Errorf("%s sed did not edit expected contents", impl)
						}
						backup, err := os.ReadFile(filepath.Join(dir, "notes.txt.bak"))
						if err != nil || !bytes.HasPrefix(backup, []byte("old old")) {
							return fmt.Errorf("%s sed backup incorrect", impl)
						}
					}
					if name == "grep" && bytes.Contains([]byte(code), []byte("-n -i")) && string(output) != "2:ERROR broken\n3:error small\n" {
						return fmt.Errorf("%s grep matching output incorrect: %s", impl, output)
					}
				}
			}
		}
	}
	return nil
}
func run(name string, args []string, dir, tools string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + tools + ":/usr/bin:/bin", "LC_ALL=C", "HOME=/tmp"}
	return cmd.CombinedOutput()
}
