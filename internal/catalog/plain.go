package catalog

import (
	"bufio"
	"fmt"
	"html"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Plain renders the accepted Markdown subset without ANSI codes. Code is never wrapped.
func Plain(markdown string, width int) string {
	source := []byte(markdown)
	doc := goldmark.New(goldmark.WithExtensions(extension.Table)).Parser().Parse(text.NewReader(source))
	var out strings.Builder
	var inline func(ast.Node) []inlineSpan
	inline = func(n ast.Node) []inlineSpan {
		var spans []inlineSpan
		appendText := func(s string) { spans = append(spans, inlineSpan{text: s}) }
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch v := c.(type) {
			case *ast.Text:
				appendText(plainText(v.Segment.Value(source), v.IsRaw()))
				if v.SoftLineBreak() || v.HardLineBreak() {
					appendText(" ")
				}
			case *ast.String:
				appendText(plainText(v.Value, v.IsRaw() || v.IsCode()))
			case *ast.CodeSpan:
				var code strings.Builder
				for child := c.FirstChild(); child != nil; child = child.NextSibling() {
					value := string(child.(*ast.Text).Segment.Value(source))
					if strings.HasSuffix(value, "\n") {
						value = strings.TrimSuffix(value, "\n") + " "
					}
					code.WriteString(value)
				}
				spans = append(spans, inlineSpan{text: code.String(), code: true})
			case *ast.Link:
				spans = append(spans, inline(c)...)
				appendText(fmt.Sprintf(" (%s)", v.Destination))
			default:
				spans = append(spans, inline(c)...)
			}
		}
		return spans
	}
	var block func(ast.Node, string)
	block = func(n ast.Node, prefix string) {
		switch v := n.(type) {
		case *ast.Paragraph, *ast.TextBlock:
			out.WriteString(wrap(inline(n), width, prefix))
			out.WriteString("\n\n")
		case *ast.Heading:
			out.WriteString(prefix + strings.TrimSpace(flatText(n, source)) + "\n\n")
		case *east.Table:
			// Tables and diagrams are printed as written, with no added indent,
			// so they keep their shape in a terminal.
			for _, line := range strings.Split(layoutTable(tableRows(v, source)), "\n") {
				out.WriteString(prefix + line + "\n")
			}
			out.WriteByte('\n')
		case *ast.FencedCodeBlock:
			indent := "    "
			if string(v.Language(source)) == "diagram" {
				indent = ""
			}
			for i := 0; i < v.Lines().Len(); i++ {
				seg := v.Lines().At(i)
				out.WriteString(prefix + indent)
				out.Write(seg.Value(source))
			}
			out.WriteByte('\n')
		case *ast.List:
			idx := v.Start
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				marker := "- "
				if v.IsOrdered() {
					marker = fmt.Sprintf("%d. ", idx)
					idx++
				}
				out.WriteString(prefix + marker)
				padding := prefix + strings.Repeat(" ", len(marker))
				first := true
				for item := c.FirstChild(); item != nil; item = item.NextSibling() {
					switch item.(type) {
					case *ast.TextBlock, *ast.Paragraph:
						rendered := wrap(inline(item), width, padding)
						if first {
							rendered = strings.TrimPrefix(rendered, padding)
						}
						out.WriteString(rendered)
						out.WriteByte('\n')
						if _, loose := item.(*ast.Paragraph); loose {
							out.WriteByte('\n')
						}
					default:
						if first {
							out.WriteByte('\n')
						}
						block(item, padding)
					}
					first = false
				}
			}
			out.WriteByte('\n')
		default:
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				block(c, prefix)
			}
		}
	}
	block(doc, "")
	return strings.TrimRight(out.String(), "\n") + "\n"
}

// Use the HTML renderer's single-pass CommonMark decoding, then remove only
// the HTML escaping it adds. Code and raw text bypass Markdown decoding.
func plainText(source []byte, raw bool) string {
	if raw {
		return string(source)
	}
	var rendered strings.Builder
	w := bufio.NewWriter(&rendered)
	goldmarkhtml.DefaultWriter.Write(w, source)
	_ = w.Flush()
	return strings.Map(func(r rune) rune {
		// Entity references must not introduce terminal control sequences.
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return '\uFFFD'
		}
		return r
	}, html.UnescapeString(rendered.String()))
}

type inlineSpan struct {
	text string
	code bool
}

func wrap(spans []inlineSpan, width int, prefix string) string {
	if width <= 0 {
		var raw strings.Builder
		raw.WriteString(prefix)
		for _, span := range spans {
			raw.WriteString(span.text)
		}
		return raw.String()
	}
	var out strings.Builder
	length := utf8.RuneCountInString(prefix)
	out.WriteString(prefix)
	// Split only prose whitespace; inline code and adjacent punctuation stay intact.
	var words []string
	var word strings.Builder
	flush := func() {
		if word.Len() > 0 {
			words = append(words, word.String())
			word.Reset()
		}
	}
	for _, span := range spans {
		if span.code {
			word.WriteString(span.text)
			continue
		}
		for _, r := range span.text {
			if unicode.IsSpace(r) {
				flush()
			} else {
				word.WriteRune(r)
			}
		}
	}
	flush()
	for i, w := range words {
		n := utf8.RuneCountInString(w)
		if i > 0 {
			if length+1+n > width {
				out.WriteString("\n" + prefix)
				length = utf8.RuneCountInString(prefix)
			} else {
				out.WriteByte(' ')
				length++
			}
		}
		out.WriteString(w)
		length += n
	}
	return out.String()
}

// flatText returns the visible text of a node's inline content on one line.
func flatText(n ast.Node, source []byte) string {
	var b strings.Builder
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch v := c.(type) {
		case *ast.Text:
			b.WriteString(plainText(v.Segment.Value(source), v.IsRaw()))
			if v.SoftLineBreak() || v.HardLineBreak() {
				b.WriteByte(' ')
			}
		case *ast.String:
			b.WriteString(plainText(v.Value, v.IsRaw() || v.IsCode()))
		case *ast.CodeSpan:
			for child := c.FirstChild(); child != nil; child = child.NextSibling() {
				value := string(child.(*ast.Text).Segment.Value(source))
				b.WriteString(strings.TrimSuffix(value, "\n"))
				if strings.HasSuffix(value, "\n") {
					b.WriteByte(' ')
				}
			}
		case *ast.Link:
			b.WriteString(flatText(c, source))
			fmt.Fprintf(&b, " (%s)", v.Destination)
		default:
			b.WriteString(flatText(c, source))
		}
	}
	return b.String()
}

// tableRows returns a table's cell texts, header row first.
func tableRows(t *east.Table, source []byte) [][]string {
	var rows [][]string
	for r := t.FirstChild(); r != nil; r = r.NextSibling() {
		var row []string
		for c := r.FirstChild(); c != nil; c = c.NextSibling() {
			row = append(row, strings.TrimSpace(flatText(c, source)))
		}
		rows = append(rows, row)
	}
	return rows
}

// layoutTable aligns rows into columns separated by two spaces, with a dashed
// line under the header. Cells are never wrapped.
func layoutTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				widths[i] = max(widths[i], utf8.RuneCountInString(cell))
			}
		}
	}
	line := func(row []string) string {
		cells := make([]string, len(widths))
		for i := range widths {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			cells[i] = cell + strings.Repeat(" ", widths[i]-utf8.RuneCountInString(cell))
		}
		return strings.TrimRight(strings.Join(cells, "  "), " ")
	}
	lines := []string{line(rows[0])}
	rule := make([]string, len(widths))
	for i, w := range widths {
		rule[i] = strings.Repeat("-", w)
	}
	lines = append(lines, strings.Join(rule, "  "))
	for _, row := range rows[1:] {
		lines = append(lines, line(row))
	}
	return strings.Join(lines, "\n")
}
