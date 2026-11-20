// Package startup decides what `explain` does when it is run with no
// arguments on a terminal: show the usage once and ask whether to keep showing
// it, then start the interactive browser.
package startup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const key = "show_help_on_start"

// Path returns the preferences file: explain/config in the user's config
// directory ($XDG_CONFIG_HOME or ~/.config on Linux, ~/Library/Application
// Support on macOS).
func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "explain", "config"), nil
}

// Load reads the saved answer. known is false until the user has answered.
func Load(path string) (show, known bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, false
	}
	for _, line := range strings.Split(string(data), "\n") {
		k, v, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok || strings.TrimSpace(k) != key {
			continue
		}
		switch strings.TrimSpace(v) {
		case "true":
			return true, true
		case "false":
			return false, true
		}
	}
	return false, false
}

// Save stores the answer, creating the directory if needed.
func Save(path string, show bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	text := fmt.Sprintf("# explain preferences. Delete this file to be asked again.\n%s = %t\n", key, show)
	return os.WriteFile(path, []byte(text), 0o644)
}

// Flow holds everything Run needs, so tests can replace the terminal.
type Flow struct {
	// ConfigPath is the preferences file; empty means preferences cannot be saved.
	ConfigPath string
	In         io.Reader
	Err        io.Writer
	// ShowHelp displays the usage (through the pager when there is one).
	ShowHelp func() error
	// Browse starts the interactive browser.
	Browse func() error
}

// Run shows the usage when it is wanted, asks the first time, and starts the
// browser. A preference that cannot be read or saved is reported and ignored.
func (f Flow) Run() error {
	show, known := false, false
	if f.ConfigPath != "" {
		show, known = Load(f.ConfigPath)
	}
	if !known || show {
		if err := f.ShowHelp(); err != nil {
			return err
		}
	}
	if !known {
		answer := f.ask()
		if f.ConfigPath != "" {
			if err := Save(f.ConfigPath, answer); err != nil {
				fmt.Fprintln(f.Err, "explain: could not save your choice:", err)
			}
		} else {
			fmt.Fprintln(f.Err, "explain: no config directory; your choice is not saved")
		}
	}
	return f.Browse()
}

// ask returns true for yes (the default) and false for no.
func (f Flow) ask() bool {
	fmt.Fprint(f.Err, "Show this usage every time you run explain without arguments? [Y/n] ")
	line, _ := bufio.NewReader(f.In).ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "n", "no":
		fmt.Fprintln(f.Err, "OK. explain now opens the browser directly. Delete "+f.ConfigPath+" to be asked again.")
		return false
	}
	return true
}
