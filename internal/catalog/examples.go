package catalog

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ExampleText renders descriptions as shell comments and preserves code without
// adding terminal indentation. The CLI and website share this generated text.
func ExampleText(markdown string) string {
	source := []byte(markdown)
	doc := goldmark.New().Parser().Parse(text.NewReader(source))
	var out strings.Builder
	afterCode := false
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.(type) {
		case *ast.Paragraph, *ast.TextBlock:
			if afterCode {
				out.WriteByte('\n')
			}
			var prose strings.Builder
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				prose.Write(segment.Value(source))
			}
			for _, line := range strings.Split(strings.TrimRight(Plain(prose.String(), DescriptionWidth-2), "\n"), "\n") {
				out.WriteString("# " + line + "\n")
			}
			afterCode = false
			return ast.WalkSkipChildren, nil
		case *ast.FencedCodeBlock:
			if afterCode {
				out.WriteByte('\n')
			}
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				out.Write(segment.Value(source))
			}
			afterCode = true
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimRight(out.String(), "\n") + "\n"
}
