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
	"strings"
	"time"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/detect"
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
	for _, impl := range []string{"busybox", "gnu", "mawk"} {
		tools, err := os.MkdirTemp("", "tools-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tools)
		for _, name := range applicable(c, impl) {
			target := "/bin/busybox"
			if impl != "busybox" {
				target, err = exec.LookPath(provider(impl, name))
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
	if err := checkGuides(c); err != nil {
		return err
	}
	return checkAwk()
}

// checkGuides runs every command a guide asks a reader to try, in a temporary
// folder, on the guest's own tools (BusyBox ps, GNU coreutils and so on). Each
// must be safe, succeed, and print what the guide shows.
func checkGuides(c *catalog.Catalog) error {
	count := 0
	for _, g := range c.Guides {
		for _, s := range g.Sections {
			if s.Title != "Try it" {
				continue
			}
			for _, b := range examplecheck.TryItBlocks(s.Markdown) {
				if why := examplecheck.UnsafeTryIt(b.Code); why != "" {
					return fmt.Errorf("guide %s: %q is not allowed in a Try it block", g.ID, why)
				}
				dir, err := os.MkdirTemp("", "guide-")
				if err != nil {
					return err
				}
				defer os.RemoveAll(dir)
				if err = examplecheck.Setup(dir, g.ID); err != nil {
					return err
				}
				out, err := run("/bin/sh", []string{"-c", b.Code}, dir, "/nonexistent")
				if err != nil {
					return fmt.Errorf("guide %s: %q failed: %w %s", g.ID, b.Code, err, out)
				}
				if err = examplecheck.CheckTryIt(b, string(out)); err != nil {
					return fmt.Errorf("guide %s: %q: %w", g.ID, b.Code, err)
				}
				count++
			}
		}
	}
	fmt.Printf("Linux guides: %d Try it commands ran and matched\n", count)
	return nil
}

// checkAwk runs implementation detection against the guest's real awk
// programs: BusyBox, gawk and mawk. No awk page exists yet, so this calls the
// detector directly instead of the CLI.
func checkAwk() error {
	for _, tt := range []struct{ command, provider, impl string }{
		{"awk", "busybox", "busybox"},
		{"awk", "gawk", "gnu"},
		{"gawk", "gawk", "gnu"},
		{"awk", "mawk", "mawk"},
		{"mawk", "mawk", "mawk"},
	} {
		tools, err := os.MkdirTemp("", "awk-tools-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tools)
		target := "/bin/busybox"
		if tt.provider != "busybox" {
			if target, err = exec.LookPath(tt.provider); err != nil {
				return err
			}
		}
		link := filepath.Join(tools, tt.command)
		if err = os.Symlink(target, link); err != nil {
			return err
		}
		banner, _ := run(link, []string{"--version"}, "/", tools)
		s, err := (detect.Detector{OS: "linux", Path: tools}).Select(context.Background(), detect.Request{Name: tt.command, Command: "awk"})
		if err != nil {
			return err
		}
		if s.Implementation != tt.impl || s.Confidence != "detected" {
			return fmt.Errorf("awk %s (%s) detected as %+v", tt.command, tt.provider, s)
		}
		first, _, _ := strings.Cut(string(banner), "\n")
		fmt.Printf("Linux awk detection: %s via %s -> %s (version %s); banner: %s\n", tt.command, tt.provider, tt.impl, s.Version, first)
	}
	return nil
}

// provider names the guest executable that implements name for impl. The tool
// is linked under name, so awk can be gawk or mawk, whatever plain awk is.
func provider(impl, name string) string {
	switch {
	case impl == "mawk":
		return "mawk"
	case impl == "gnu" && name == "awk":
		return "gawk"
	}
	return name
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
