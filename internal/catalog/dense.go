package catalog

import (
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
	Code        string
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
			examples = append(examples, Example{
				Description: CleanDescription(strings.Join(prose, " ")),
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
