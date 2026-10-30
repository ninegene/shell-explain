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

func TestExamplesDefaultAndDetails(t *testing.T) {
	c, err := catalog.Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range c.Entries {
		for _, impl := range append([]string{"common"}, e.Implementations...) {
			t.Run(e.ID+"/"+impl, func(t *testing.T) {
				args := []string{e.ID, "--implementation=" + impl, "--width=0"}
				status, examples, stderr := run(t, args...)
				if status != 0 || stderr != "" {
					t.Fatalf("examples: status=%d stderr=%s", status, stderr)
				}
				status, details, stderr := run(t, append(args, "--details")...)
				if status != 0 || stderr != "" {
					t.Fatalf("details: status=%d stderr=%s", status, stderr)
				}
				for _, omitted := range []string{"\nOverview\n", "\nCommon mistakes\n", "\nPortability\n", "\nOption: ", "\nReferences\n", "\nSelection:"} {
					if strings.Contains(examples, omitted) {
						t.Errorf("default output contains detail %q", omitted)
					}
				}
				for _, required := range []string{"Overview\n", "Common mistakes\n", "Portability\n", "References\n", "Selection:"} {
					if !strings.Contains(details, required) {
						t.Errorf("full output missing %q", required)
					}
				}
				selected, _ := e.Select(impl, "")
				for _, section := range selected {
					for _, code := range section.Code {
						for _, line := range strings.Split(strings.TrimSuffix(code, "\n"), "\n") {
							if !strings.Contains(examples, "\n"+line+"\n") {
								t.Errorf("missing or changed example: %q", line)
							}
						}
					}
				}
				for _, ref := range e.References {
					if strings.Contains(examples, ref) || !strings.Contains(details, ref) {
						t.Errorf("reference visibility: %s", ref)
					}
				}
			})
		}
	}
}

func TestDetailsFlagPositionsAndJSONContract(t *testing.T) {
	_, after, _ := run(t, "ggrep", "--implementation=gnu", "--details")
	status, before, stderr := run(t, "--details", "ggrep", "--implementation=gnu")
	if status != 0 || stderr != "" || before != after {
		t.Fatal("details flag must work before and after the command")
	}
	_, normalJSON, _ := run(t, "sed", "--implementation=apple", "--format=json")
	status, detailedJSON, stderr := run(t, "sed", "--implementation=apple", "--format=json", "--details")
	if status != 0 || stderr != "" || detailedJSON != normalJSON {
		t.Fatal("details changed the full JSON contract")
	}
	for _, args := range [][]string{{"list", "--details"}, {"search", "logs", "--details"}} {
		if status, _, _ := run(t, args...); status != 2 {
			t.Fatalf("details must be limited to command lookup: %v", args)
		}
	}
}

func TestGroupedOptions(t *testing.T) {
	status, out, stderr := run(t, "grep", "--implementation=common", "--details", "--width=72")
	if status != 0 || stderr != "" || strings.Count(out, "\nOptions\n") != 1 || strings.Contains(out, "Option: -") {
		t.Fatalf("options were not grouped: %d %s %s", status, out, stderr)
	}
	if !strings.Contains(out, "\n-n\nPrint each matching line") || !strings.Contains(out, "\n\n-i\nIgnore letter case") {
		t.Fatalf("options are missing their spelling or description: %s", out)
	}
	status, out, stderr = run(t, "grep", "--implementation=common", "--option=-i")
	if status != 0 || stderr != "" || !strings.Contains(out, "Options\n-i\n") || strings.Contains(out, "\n-n\n") {
		t.Fatalf("single option lookup included unrelated options: %d %s %s", status, out, stderr)
	}
}
