package catalog

import (
	"html"
	"strings"
	"testing"
)

func TestProseDecodingMatchesHTML(t *testing.T) {
	for _, tt := range []struct{ markdown, want string }{
		{`Use \* and A &amp; B.`, "Use * and A & B."},
		{`A &#35; sign and &#x23; sign.`, "A # sign and # sign."},
		{`&copy; &NotEqualTilde; &unknown;`, "© ≂̸ &unknown;"},
		{`\&amp; &amp;amp; &#38;amp;`, "&amp; &amp; &amp;"},
		{`\\* and &#92;*`, `\* and \*`},
	} {
		t.Run(tt.markdown, func(t *testing.T) {
			rendered, _, err := RenderHTML(tt.markdown)
			if err != nil {
				t.Fatal(err)
			}
			web := html.UnescapeString(strings.TrimSuffix(strings.TrimPrefix(rendered, "<p>"), "</p>\n"))
			if web != tt.want {
				t.Fatalf("HTML text = %q, want %q", web, tt.want)
			}
			for _, width := range []int{0, 80} {
				if got := Plain(tt.markdown, width); got != web+"\n" {
					t.Fatalf("width %d: terminal = %q, HTML text = %q", width, got, web)
				}
			}
		})
	}
	if got := Plain("A &amp; B", 5); got != "A & B\n" {
		t.Fatalf("wrapping must use decoded text: %q", got)
	}
}

func TestDecodingPreservesCodeAndBlocksTerminalControls(t *testing.T) {
	literal := `\* &amp; &#35;`
	markdown := "`" + literal + "`\n\n```sh\n" + literal + "\n```\n"
	if got := Plain(markdown, 80); got != literal+"\n\n    "+literal+"\n" {
		t.Fatalf("code changed: %q", got)
	}
	if got := Plain("&#27;[31mred&#x9b;31m", 80); got != "�[31mred�31m\n" {
		t.Fatalf("entity introduced terminal controls: %q", got)
	}
}

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

func TestListContinuationIndentation(t *testing.T) {
	for _, tt := range []struct {
		name, markdown, want string
		width                int
	}{
		{"unordered", "- alpha beta gamma delta", "- alpha beta\n  gamma delta\n", 14},
		{"ordered", "10. alpha beta gamma delta", "10. alpha beta\n    gamma delta\n", 16},
		{"nested", "- outer\n  - alpha beta gamma delta", "- outer\n  - alpha beta\n    gamma delta\n", 16},
		{"loose", "- alpha beta gamma delta\n\n- second item", "- alpha beta\n  gamma delta\n\n- second item\n", 14},
		{"unwrapped", "- outer\n  - alpha beta gamma delta", "- outer\n  - alpha beta gamma delta\n", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := RenderHTML(tt.markdown); err != nil {
				t.Fatal(err)
			}
			if got := Plain(tt.markdown, tt.width); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
