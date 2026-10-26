package catalog

import (
	"strings"
	"testing"
)

func TestInlineCodePreservesWhitespace(t *testing.T) {
	command := "printf '%s\\n' 'a  b'"
	for _, width := range []int{0, 12, 80} {
		for _, markdown := range []string{
			"Keep `" + command + "`, unchanged.",
			"- Run `" + command + "` here.",
			"Keep **`" + command + "`**, unchanged.",
		} {
			got := Plain(markdown, width)
			if !strings.Contains(got, command) {
				t.Fatalf("width %d corrupted inline code: %q", width, got)
			}
			if strings.Contains(markdown, "`,") && !strings.Contains(got, command+",") {
				t.Fatalf("punctuation detached: %q", got)
			}
		}
	}
	if got := Plain("Use `a\nb` and `  x  `.", 80); got != "Use a b and  x .\n" {
		t.Fatalf("CommonMark code span normalization: %q", got)
	}
}
