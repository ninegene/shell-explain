package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/detect"
)

func run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	c, err := catalog.Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	status := Run(context.Background(), c, args, &out, &stderr, Config{Version: "test", Width: 40, Detector: detect.Detector{OS: "linux", Path: t.TempDir()}})
	return status, out.String(), stderr.String()
}
func TestCommandInterface(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"--version"}, {"list"}, {"search", "logs"}, {"grep"}, {"ggrep", "--implementation", "gnu"}, {"--implementation=apple", "sed", "--option=-i"}, {"cd", "--platform=bash"}} {
		status, out, err := run(t, args...)
		if status != 0 || out == "" || err != "" {
			t.Errorf("%v: %d %s %s", args, status, out, err)
		}
	}
	for _, args := range [][]string{{"unknown"}, {"grep", "--implementation=freebsd"}, {"grep", "--option=--missing"}, {"grep", "--width=-1"}, {"grep", "--format=xml"}, {"grep", "--implementation=gnu", "--platform=linux"}, {"grep", "foo"}, {"search"}, {"list", "foo"}} {
		status, _, err := run(t, args...)
		if status == 0 || err == "" {
			t.Errorf("invalid request succeeded: %v", args)
		}
	}
}
func TestJSONParityAndErrors(t *testing.T) {
	c, err := catalog.Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range c.Entries {
		for _, impl := range append([]string{"common"}, e.Implementations...) {
			status, out, stderr := run(t, e.ID, "--implementation="+impl, "--format=json", "--color=always")
			if status != 0 || stderr != "" || strings.Contains(out, "\x1b") {
				t.Fatal(out, stderr)
			}
			var result response
			if err := json.Unmarshal([]byte(out), &result); err != nil {
				t.Fatal(err)
			}
			expected, err := e.Select(impl, "")
			if err != nil {
				t.Fatal(err)
			}
			a, _ := json.Marshal(expected)
			b, _ := json.Marshal(result.Sections)
			if !bytes.Equal(a, b) || result.ContentRevision != c.Revision {
				t.Fatal("CLI model diverged from exported content")
			}
		}
	}
	status, out, stderr := run(t, "missing", "--format=json")
	if status != 2 || stderr != "" || !strings.Contains(out, `"error"`) {
		t.Fatal("JSON error contract", out, stderr)
	}
}
func TestOptionLookupAndColor(t *testing.T) {
	_, out, _ := run(t, "grep", "--implementation=gnu", "--option=-n")
	if !strings.Contains(out, "Print each matching line") || strings.Contains(out, "GNU examples") {
		t.Fatal("wrong option view")
	}
	_, out, _ = run(t, "grep", "--color=always")
	if !strings.Contains(out, "\x1b[") {
		t.Fatal("color override missing")
	}
}
