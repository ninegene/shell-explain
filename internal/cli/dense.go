package cli

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/ninegene/shell-explain/internal/catalog"
)

// SiteURL is the public website that has the full page for every command.
const SiteURL = "https://shell-explain.ninegene.com"

const (
	maxInlineCode = 46  // longest command that gets a trailing comment
	defaultWidth  = 100 // line width when none is given or the width is 0
)

// Colors use only the 8 standard ANSI colors and the bold, dim and underline
// attributes, so every terminal theme maps them to something readable on its own
// background. No 256-color or truecolor values, and no cyan, yellow or bright
// colors, which wash out on light backgrounds; bold is never combined with a
// color because many terminals draw that as bright.
const (
	ansiReset    = "\x1b[0m"
	ansiHeading  = "\x1b[1m"
	ansiCommand  = "\x1b[1m"
	ansiFlag     = "\x1b[34m"
	ansiString   = "\x1b[35m"
	ansiOperator = "\x1b[2m"
	ansiComment  = "\x1b[32m"
	ansiLabel    = "\x1b[2m"
	ansiLink     = "\x1b[4;34m"
)

func paint(on bool, code, s string) string {
	if !on || s == "" {
		return s
	}
	return code + s + ansiReset
}

func width(s string) int { return utf8.RuneCountInString(s) }

// block is one laid-out example: either a command with a trailing comment on
// one line, or a comment line (or lines) above the command.
type block struct {
	catalog.Example
	inline bool
}

// RenderExamples lays examples out for a terminal. A short command with a
// short description becomes one line, `command  # description`, with the
// comments of neighbouring lines aligned. Everything else keeps the authored
// shape: the description as comment lines above the command. The commands are
// never changed or wrapped, so they can be copied as they are shown.
func RenderExamples(examples []catalog.Example, lineWidth int, color bool) string {
	if lineWidth <= 0 {
		lineWidth = defaultWidth
	}
	blocks := make([]block, len(examples))
	for i, e := range examples {
		blocks[i] = block{Example: e, inline: e.Description != "" && e.Code != "" &&
			!strings.Contains(e.Code, "\n") && width(e.Code) <= maxInlineCode &&
			width(e.Code)+4+width(e.Description) <= lineWidth}
	}
	// Aligning a run to its longest command can push one line past the width;
	// move such lines out of the run until everything fits.
	for changed := true; changed; {
		changed = false
		for start := 0; start < len(blocks); {
			if !blocks[start].inline {
				start++
				continue
			}
			end, column := start, 0
			for end < len(blocks) && blocks[end].inline {
				column = max(column, width(blocks[end].Code))
				end++
			}
			for i := start; i < end; i++ {
				if column+4+width(blocks[i].Description) > lineWidth {
					blocks[i].inline = false
					changed = true
				}
			}
			start = end
		}
	}

	var out strings.Builder
	spaced := true // nothing written yet, or a blank line was just written
	blank := func() {
		if !spaced {
			out.WriteByte('\n')
			spaced = true
		}
	}
	for start := 0; start < len(blocks); {
		b := blocks[start]
		if b.inline {
			end, column := start, 0
			for end < len(blocks) && blocks[end].inline {
				column = max(column, width(blocks[end].Code))
				end++
			}
			for _, r := range blocks[start:end] {
				pad := strings.Repeat(" ", column-width(r.Code))
				fmt.Fprintf(&out, "%s%s  %s\n", Highlight(r.Code, color), pad, paint(color, ansiComment, "# "+r.Description))
			}
			spaced = false
			start = end
			continue
		}
		blank()
		if b.Description != "" {
			for _, line := range strings.Split(strings.TrimRight(catalog.Plain(b.Description, max(20, lineWidth-2)), "\n"), "\n") {
				out.WriteString(paint(color, ansiComment, "# "+line) + "\n")
			}
		}
		if b.Code != "" {
			for _, line := range strings.Split(b.Code, "\n") {
				out.WriteString(Highlight(line, color) + "\n")
			}
		}
		out.WriteByte('\n')
		spaced = true
		start++
	}
	return strings.TrimRight(out.String(), "\n") + "\n"
}

// Highlight colors a shell command line: the command names, flags, quoted
// strings and operators. It only adds color codes; the text is unchanged, and
// with color off it returns the line as it is.
func Highlight(line string, color bool) string {
	if !color {
		return line
	}
	var out strings.Builder
	atCommand := true
	i := 0
	for i < len(line) {
		if line[i] == ' ' || line[i] == '\t' {
			out.WriteByte(line[i])
			i++
			continue
		}
		start := i
		for i < len(line) && line[i] != ' ' && line[i] != '\t' {
			switch line[i] {
			case '\'', '"':
				quote := line[i]
				i++
				for i < len(line) && line[i] != quote {
					if quote == '"' && line[i] == '\\' && i+1 < len(line) {
						i++
					}
					i++
				}
				if i < len(line) {
					i++
				}
			case '\\':
				i += 2
			default:
				i++
			}
		}
		if i > len(line) {
			i = len(line)
		}
		token := line[start:i]
		switch {
		case isOperator(token):
			out.WriteString(paint(true, ansiOperator, token))
			atCommand = token == "|" || token == "||" || token == "&&" || token == ";"
		case atCommand:
			out.WriteString(paint(true, ansiCommand, token))
			atCommand = false
		case token[0] == '-' && len(token) > 1:
			out.WriteString(paint(true, ansiFlag, token))
		case token[0] == '\'' || token[0] == '"':
			out.WriteString(paint(true, ansiString, token))
		default:
			out.WriteString(token)
		}
	}
	return out.String()
}

func isOperator(token string) bool {
	if token == "" {
		return false
	}
	for i := 0; i < len(token); i++ {
		if !strings.ContainsRune("|&;<>0123456789", rune(token[i])) {
			return false
		}
	}
	return strings.ContainsAny(token, "|&;<>")
}

// writeLinks prints where to read more: the page on the website and the
// upstream manuals. The URLs are written out in full so a terminal can open
// them with a click.
func writeLinks(out io.Writer, e *catalog.Entry, implementation string, single bool, color bool) {
	fmt.Fprintln(out)
	writeMore(out, e, implementation, single, color)
	for i, ref := range e.References {
		name := linkLabel("", color)
		if i == 0 {
			name = linkLabel("Docs:", color)
		}
		fmt.Fprintf(out, "%s %s\n", name, paint(color, ansiLink, ref))
	}
	fmt.Fprintf(out, "%s %s\n", linkLabel("Full:", color), paint(color, ansiLabel, "explain "+e.ID+" --details"))
}

func linkLabel(s string, color bool) string { return paint(color, ansiLabel, fmt.Sprintf("%-6s", s)) }

// writeMore prints the command's page on the website.
func writeMore(out io.Writer, e *catalog.Entry, implementation string, single bool, color bool) {
	page := SiteURL + "/commands/" + e.ID + "/"
	if !single && implementation != "" && implementation != "common" {
		page += implementation + "/"
	}
	fmt.Fprintf(out, "%s %s\n", linkLabel("More:", color), paint(color, ansiLink, page))
}

const maxOptionColumn = 22 // widest option spelling that shares a line with its description

// RenderOptions lays option explanations out as two aligned columns, the
// spellings and then the description, wrapped under the description column.
// A spelling too long for the column gets its own line. The website keeps one
// block with each spelling above its description; the words are the same.
func RenderOptions(sections []catalog.Section, lineWidth int, color bool) string {
	if lineWidth <= 0 {
		lineWidth = defaultWidth
	}
	limit := min(maxOptionColumn, max(8, lineWidth/3))
	column := 0
	for _, s := range sections {
		if len(s.Options) > 0 {
			if w := width(strings.Join(s.Options, ", ")); w <= limit {
				column = max(column, w)
			}
		}
	}
	indent := strings.Repeat(" ", column+2)
	var out strings.Builder
	for _, s := range sections {
		if len(s.Options) == 0 {
			continue
		}
		spelling := strings.Join(s.Options, ", ")
		text := strings.TrimRight(catalog.Plain(s.Markdown, max(20, lineWidth-column-2)), "\n")
		lines := strings.Split(text, "\n")
		// Descriptions are explanations, so they get the comment color.
		describe := func(l string) string { return paint(color, ansiComment, l) }
		if width(spelling) > limit {
			out.WriteString(paint(color, ansiFlag, spelling) + "\n")
			for _, l := range lines {
				out.WriteString(indent + describe(l) + "\n")
			}
			continue
		}
		pad := strings.Repeat(" ", column-width(spelling))
		out.WriteString(paint(color, ansiFlag, spelling) + pad + "  " + describe(lines[0]) + "\n")
		for _, l := range lines[1:] {
			out.WriteString(strings.TrimRight(indent+describe(l), " ") + "\n")
		}
	}
	return out.String()
}
