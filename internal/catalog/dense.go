package catalog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Example is one authored example: the prose that introduces a code block and
// the exact code in it. The website shows each as its own block; the terminal
// lays the same pairs out more densely (see internal/cli).
type Example struct {
	Description string
	// Short is an optional label of about 50 characters for the default terminal
	// layout, authored on the code fence as ```sh short="...". The website and
	// --details keep the full Description.
	Short string
	Code  string
}

// ParseExamples returns the examples in an examples section, in authored
// order. Prose before a code block becomes its description; prose with no code
// after it is kept as an example without code so no explanation is lost.
func ParseExamples(markdown string) []Example {
	source := []byte(markdown)
	doc := goldmark.New().Parser().Parse(text.NewReader(source))
	var examples []Example
	var prose []string
	flush := func() {
		if len(prose) > 0 {
			examples = append(examples, Example{Description: CleanDescription(strings.Join(prose, " "))})
			prose = nil
		}
	}
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.(type) {
		case *ast.Paragraph, *ast.TextBlock:
			var raw strings.Builder
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				raw.Write(segment.Value(source))
			}
			prose = append(prose, strings.TrimSpace(Plain(raw.String(), 1<<20)))
			return ast.WalkSkipChildren, nil
		case *ast.FencedCodeBlock:
			var code strings.Builder
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				code.Write(segment.Value(source))
			}
			short, _ := ShortLabel(fenceInfo(n.(*ast.FencedCodeBlock), source))
			examples = append(examples, Example{
				Description: CleanDescription(strings.Join(prose, " ")),
				Short:       short,
				Code:        strings.TrimRight(code.String(), "\n"),
			})
			prose = nil
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	flush()
	return examples
}

var reviewedOnly = regexp.MustCompile(`\s*[;,]?\s*reviewed only(,? because [^:.]*)?`)

// CleanDescription makes an authored lead-in short enough for a terminal
// comment: it joins lines, drops the "reviewed only" remark (a note for web
// readers about how the example is checked) and the trailing colon.
func CleanDescription(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	s = reviewedOnly.ReplaceAllString(s, "")
	s = strings.TrimSpace(strings.TrimRight(strings.TrimSpace(s), ":"))
	return strings.TrimSpace(strings.TrimRight(s, ";,"))
}

// MaxShortLabel is the longest short label, in characters.
const MaxShortLabel = 70

var shortLabel = regexp.MustCompile(`^short="([^"]*)"$`)

// fenceInfo returns what follows the language on a code fence's opening line.
func fenceInfo(n *ast.FencedCodeBlock, source []byte) string {
	if n.Info == nil {
		return ""
	}
	info := strings.TrimSpace(string(n.Info.Segment.Value(source)))
	lang := string(n.Language(source))
	return strings.TrimSpace(strings.TrimPrefix(info, lang))
}

// ShortLabel reads the `short="..."` attribute from the text after the
// language on a code fence. The empty string means there is none. The label is
// plain text on one line, so a terminal comment can use it as it is.
func ShortLabel(info string) (string, error) {
	if info == "" {
		return "", nil
	}
	m := shortLabel.FindStringSubmatch(info)
	if m == nil {
		return "", fmt.Errorf("unsupported code fence attribute %q; only short=\"...\" is allowed", info)
	}
	label := m[1]
	if label == "" || len([]rune(label)) > MaxShortLabel || strings.TrimSpace(label) != label {
		return "", fmt.Errorf("short label must be 1-%d characters without surrounding spaces: %q", MaxShortLabel, label)
	}
	for _, r := range label {
		if r < ' ' || r == 0x7f || r == '`' {
			return "", fmt.Errorf("short label contains a control character or backtick: %q", label)
		}
	}
	return label, nil
}
