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

// unidentifiable lists GNU commands that the probe cannot identify, with the reason.
var unidentifiable = map[string]string{
	"test":       "GNU test prints no version banner, because a lone argument is a string test",
	"zcat":       "GNU zcat is a shell script, and scripts are never executed by the probe",
	"zgrep":      "GNU zgrep is a shell script, and scripts are never executed by the probe",
	"zless":      "GNU zless is a shell script, and scripts are never executed by the probe",
	"zmore":      "GNU zmore is a shell script, and scripts are never executed by the probe",
	"zcmp":       "GNU zcmp is a shell script, and scripts are never executed by the probe",
	"zdiff":      "GNU zdiff is a shell script, and scripts are never executed by the probe",
	"zforce":     "GNU zforce is a shell script, and scripts are never executed by the probe",
	"znew":       "GNU znew is a shell script, and scripts are never executed by the probe",
	"uncompress": "GNU uncompress is a shell script, and scripts are never executed by the probe",
	"gzexe":      "GNU gzexe is a shell script, and scripts are never executed by the probe",
}

func check() error {
	c, err := catalog.Load(content.Files)
	if err != nil {
		return err
	}
	for _, impl := range []string{"busybox", "gnu", "mawk", "procps"} {
		tools, err := os.MkdirTemp("", "tools-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tools)
		missing := map[string]bool{}
		for _, name := range applicable(c, impl) {
			target := "/bin/busybox"
			if impl != "busybox" {
				target, err = exec.LookPath(provider(impl, name))
				if err == nil && isBusyBox(target) {
					// The name resolves to a BusyBox applet, not to a GNU program.
					err = fmt.Errorf("%s is a BusyBox applet", target)
				}
				if err != nil {
					// Alpine's coreutils package leaves out a few commands that other
					// packages provide, such as groups and uptime. Say so and go on.
					fmt.Printf("Linux %s: %s is not installed in the guest, skipped\n", impl, name)
					missing[name] = true
					continue
				}
			}
			if err = os.Symlink(target, filepath.Join(tools, name)); err != nil {
				return err
			}
			for _, companion := range examplecheck.Companions[name] {
				// A program the command runs by name must come from the same implementation.
				companionTarget := "/bin/busybox"
				if impl != "busybox" {
					if companionTarget, err = exec.LookPath(provider(impl, companion)); err != nil || isBusyBox(companionTarget) {
						return fmt.Errorf("%s needs %s from the %s guest packages", name, companion, impl)
					}
				}
				if err = os.Symlink(companionTarget, filepath.Join(tools, companion)); err != nil && !os.IsExist(err) {
					return err
				}
			}
		}
		for _, name := range applicable(c, impl) {
			if missing[name] {
				continue
			}
			data, err := run("/explain", []string{name, "--format=json"}, "/", tools)
			if err != nil && impl == "gnu" && strings.Contains(string(data), "has no procps documentation") {
				// Alpine's procps-ng package provides this command, not GNU coreutils.
				fmt.Printf("Linux gnu: %s in the guest is procps-ng, skipped\n", name)
				continue
			}
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
			if reason, ok := unidentifiable[name]; ok && impl == "gnu" && result.Selection.Implementation == "common" {
				fmt.Printf("Linux detection: %s -> not identified (%s)\n", name, reason)
			} else if result.Selection.Implementation != impl {
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
					script := code
					if slices.Contains(entry.Implementations, "bash") {
						// Run the executable under test, not the shell's builtin of the same name.
						script = name + "() { \"" + filepath.Join(tools, name) + "\" \"$@\"; }\n" + code
					}
					output, err := run("/bin/sh", []string{"-c", script}, dir, tools)
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
	if err := checkMissingApplets(c); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + tools + ":/usr/bin:/bin", "LC_ALL=C", "HOME=/tmp"}
	return cmd.CombinedOutput()
}

// checkMissingApplets verifies the claim a page makes by leaving out BusyBox: for
// every entry that is documented for GNU or Apple but not BusyBox, the guest's
// BusyBox must not provide an applet of that name.
func checkMissingApplets(c *catalog.Catalog) error {
	out, err := run("/bin/busybox", []string{"--list"}, "/", "/nonexistent")
	if err != nil {
		return fmt.Errorf("busybox --list: %w %s", err, out)
	}
	applets := strings.Fields(string(out))
	fmt.Printf("Linux BusyBox applets (%d): %s\n", len(applets), strings.Join(applets, " "))
	count := 0
	for _, e := range c.Entries {
		if slices.Contains(e.Implementations, "busybox") || !slices.Contains(e.Implementations, "gnu") {
			continue
		}
		if slices.Contains(applets, e.ID) {
			return fmt.Errorf("%s omits busybox, but the guest BusyBox has a %s applet", e.ID, e.ID)
		}
		count++
	}
	fmt.Printf("Linux BusyBox: %d pages without a BusyBox view have no applet of that name\n", count)
	return nil
}

// isBusyBox reports whether path resolves to the BusyBox binary.
func isBusyBox(path string) bool {
	real, err := filepath.EvalSymlinks(path)
	return err == nil && strings.Contains(filepath.Base(real), "busybox")
}
