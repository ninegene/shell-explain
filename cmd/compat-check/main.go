// compat-check runs reviewed examples and the built CLI in a disposable Linux VM.
// It is a test helper, not a shipped command.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"time"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/examplecheck"
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
		for _, name := range applicable(c, impl) {
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
		for _, name := range applicable(c, impl) {
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
					if examplecheck.NotRun(code) {
						continue
					}
					dir, err := os.MkdirTemp("", "example-")
					if err != nil {
						return err
					}
					defer os.RemoveAll(dir)
					if err = examplecheck.Setup(dir, name); err != nil {
						return err
					}
					output, err := run("/bin/sh", []string{"-c", code}, dir, tools)
					if err != nil {
						return fmt.Errorf("%s/%s example %s: %w %s", impl, name, code, err, output)
					}
					if err = examplecheck.Verify(code, dir, string(output)); err != nil {
						return fmt.Errorf("%s/%s example: %w", impl, name, err)
					}
				}
			}
		}
	}
	return nil
}

// applicable lists the entries that declare impl and are not shell builtins.
func applicable(c *catalog.Catalog, impl string) []string {
	var names []string
	for _, e := range c.Entries {
		if !e.ShellOnly() && slices.Contains(e.Implementations, impl) {
			names = append(names, e.ID)
		}
	}
	return names
}
func run(name string, args []string, dir, tools string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + tools + ":/usr/bin:/bin", "LC_ALL=C", "HOME=/tmp"}
	return cmd.CombinedOutput()
}
