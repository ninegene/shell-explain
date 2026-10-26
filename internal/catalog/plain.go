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
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Plain renders the accepted Markdown subset without ANSI codes. Code is never wrapped.
func Plain(markdown string, width int) string {
	source := []byte(markdown)
	doc := goldmark.New().Parser().Parse(text.NewReader(source))
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
		case *ast.FencedCodeBlock:
			for i := 0; i < v.Lines().Len(); i++ {
				seg := v.Lines().At(i)
				out.WriteString(prefix + "    ")
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
