package cli

import (
	"regexp"
	"strings"
	"testing"

	"github.com/ninegene/shell-explain/internal/catalog"
)

func TestRenderExamplesAlignsShortLinesAndStacksTheRest(t *testing.T) {
	examples := []catalog.Example{
		{Description: "Log in", Code: "ssh user@host"},
		{Description: "Use another port", Code: "ssh -p 2222 user@host"},
		{Description: "A description that is much too long to sit beside its command on one line", Code: "ssh -v user@host"},
		{Description: "Two lines", Code: "a \\\n  b"},
	}
	want := "ssh user@host          # Log in\n" +
		"ssh -p 2222 user@host  # Use another port\n" +
		"\n" +
		"# A description that is much too long to sit beside its\n" +
		"# command on one line\n" +
		"ssh -v user@host\n" +
		"\n" +
		"# Two lines\n" +
		"a \\\n  b\n"
	if got := RenderExamples(examples, 60, false, true); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderExamplesKeepsCommandsUnchangedWithColor(t *testing.T) {
	for _, code := range []string{"ssh -G example.org 2> /dev/null | grep '^hostname '", "ssh user@host 'cd app && git pull'", "ls -la"} {
		got := Highlight(code, true)
		if strings.Contains(got, "\x1b") == false {
			t.Errorf("no color in %q", code)
		}
		stripped := ansiPattern.ReplaceAllString(got, "")
		if stripped != code {
			t.Errorf("highlight changed the command: %q != %q", stripped, code)
		}
		if Highlight(code, false) != code {
			t.Errorf("color off must return the command as it is")
		}
	}
}

func TestOperatorsAndFlagsGetDifferentColors(t *testing.T) {
	got := Highlight("ls -l | wc", true)
	for _, want := range []string{ansiCommand + "ls", ansiFlag + "-l", ansiOperator + "|", ansiCommand + "wc"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Every style must stay readable on both dark and light terminal themes: only
// the 8 standard colors, no cyan or yellow, no bright colors, no 256-color or
// truecolor values, and bold never combined with a color.
func TestPaletteIsThemeSafe(t *testing.T) {
	styles := map[string]string{"heading": ansiHeading, "command": ansiCommand, "flag": ansiFlag, "string": ansiString,
		"operator": ansiOperator, "comment": ansiComment, "label": ansiLabel, "link": ansiLink}
	for name, style := range styles {
		codes := strings.Split(strings.TrimSuffix(strings.TrimPrefix(style, "\x1b["), "m"), ";")
		hasColor, hasBold := false, false
		for _, c := range codes {
			switch c {
			case "1":
				hasBold = true
			case "2", "4":
			case "31", "32", "34", "35":
				hasColor = true
			default:
				t.Errorf("%s uses code %q, which is not theme-safe", name, c)
			}
		}
		if hasColor && hasBold {
			t.Errorf("%s combines bold with a color", name)
		}
	}
}

func TestShortLabelOnlyInDefaultLayout(t *testing.T) {
	examples := []catalog.Example{{Description: "Log in to a host; the first time, check the fingerprint it shows", Short: "Log in", Code: "ssh user@host"}}
	if got := RenderExamples(examples, 100, false, true); got != "ssh user@host  # Log in\n" {
		t.Errorf("short layout: %q", got)
	}
	if got := RenderExamples(examples, 100, false, false); !strings.Contains(got, "# Log in to a host; the first time") {
		t.Errorf("full layout must keep the full description: %q", got)
	}
}
