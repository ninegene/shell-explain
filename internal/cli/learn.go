package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ninegene/shell-explain/internal/catalog"
)

// defaultGuideWidth is the prose width for guides when neither --width nor the
// terminal width (COLUMNS) is given. Diagrams (72 columns) and tables (78
// columns) fit inside it, and an 80-column terminal does not wrap it.
const defaultGuideWidth = 80

func (c Config) guideWidth(op options) int {
	if op.widthSet {
		return op.width
	}
	if c.GuideWidth > 0 {
		return c.GuideWidth
	}
	return defaultGuideWidth
}

func writeHeading(out io.Writer, color bool, s string) {
	if color {
		fmt.Fprintln(out, paint(true, ansiHeading, s))
	} else {
		fmt.Fprintln(out, s)
	}
}

// guideFlags rejects flags that only make sense for command lookup.
func guideFlags(op options, allowSection bool) error {
	if op.implementation != "" || op.platform != "" || op.option != "" || op.details || op.noProbe {
		return fmt.Errorf("selection flags apply to command lookup only")
	}
	if op.section != "" && !allowSection {
		return fmt.Errorf("--section applies to explain learn GUIDE only")
	}
	return nil
}

// NoPagerEnv is the environment variable that turns paging off everywhere,
// like --no-pager. Any value except empty or "0" disables it.
const NoPagerEnv = "EXPLAIN_NO_PAGER"

// PagerDisabled reports whether EXPLAIN_NO_PAGER turns paging off.
func PagerDisabled() bool {
	v := os.Getenv(NoPagerEnv)
	return v != "" && v != "0"
}

// wantsPager reports whether the request prints text that should go through a
// pager on a terminal: all text output except --version, and never JSON, when
// neither --no-pager nor EXPLAIN_NO_PAGER is set. Short output still prints
// directly, because the default pager (less -FRX) exits when it fits on one
// screen.
func wantsPager(op options) bool {
	return !op.noPager && !PagerDisabled() && op.format == "text" && !op.version && (len(op.args) == 0 || op.args[0] != "browse")
}

// WantsPager reports whether args ask for text output that should go through a
// pager when standard output is a terminal.
func WantsPager(args []string) bool {
	op, err := parse(args, 0)
	return err == nil && wantsPager(op)
}

func guideSummary(g catalog.Guide, number int) map[string]any {
	return map[string]any{
		"id": g.ID, "number": number, "title": g.Title, "subtitle": g.Subtitle, "summary": g.Summary,
		"series": g.Series, "order": g.Order, "reading_minutes": g.ReadingMinutes, "updated": g.Updated,
	}
}

func runLearn(c *catalog.Catalog, op options, out, errOut io.Writer, cfg Config) int {
	if err := guideFlags(op, true); err != nil {
		return failure(out, errOut, op.format, err)
	}
	if len(op.args) > 2 {
		return failure(out, errOut, op.format, fmt.Errorf("learn takes one guide name or number"))
	}
	color := op.color == "always" || (op.color == "auto" && cfg.Color)
	if len(op.args) == 1 {
		if op.section != "" {
			return failure(out, errOut, op.format, fmt.Errorf("--section applies to explain learn GUIDE only"))
		}
		return listGuides(c, op, out, color)
	}
	g, ok := c.FindGuide(op.args[1])
	if !ok {
		return failure(out, errOut, op.format, fmt.Errorf("no guide %q; run explain learn to list them", op.args[1]))
	}
	shown := *g
	if op.section != "" {
		s, ok := g.Section(op.section)
		if !ok {
			titles := make([]string, len(g.Sections))
			for i, s := range g.Sections {
				titles[i] = s.Title
			}
			return failure(out, errOut, op.format, fmt.Errorf("guide %s has no section %q; sections: %s", g.ID, op.section, strings.Join(titles, ", ")))
		}
		shown.Sections = []catalog.Section{s}
	}
	if op.format == "json" {
		return encode(out, map[string]any{"schema_version": c.SchemaVersion, "content_revision": c.Revision, "guide": shown})
	}
	width := cfg.guideWidth(op)
	if op.section == "" {
		series, _ := catalog.LookupSeries(g.Series)
		writeHeading(out, color, g.Title)
		fmt.Fprintln(out, g.Subtitle)
		fmt.Fprintf(out, "Series: %s\n", series.Title)
		fmt.Fprintf(out, "Guide %d of %d · about %d minutes · updated %s\n\n", c.GuideNumber(g.ID), len(c.Guides), g.ReadingMinutes, g.Updated)
	}
	for _, s := range shown.Sections {
		writeHeading(out, color, s.Title)
		fmt.Fprintln(out, catalog.Plain(s.Markdown, width))
	}
	if op.section != "" {
		return 0
	}
	writeHeading(out, color, "Keep going")
	prev, next := c.GuideNeighbors(g.ID)
	// A title and its command go on two lines, so long titles stay within 80 columns.
	if next != nil {
		fmt.Fprintf(out, "Next:      %s\n           explain learn %s\n", next.Title, next.ID)
	}
	if prev != nil {
		fmt.Fprintf(out, "Previous:  %s\n           explain learn %s\n", prev.Title, prev.ID)
	}
	fmt.Fprintln(out, "All guides: explain learn    Words you will see: explain glossary")
	if len(g.RelatedCommands) > 0 {
		cmds := make([]string, len(g.RelatedCommands))
		for i, name := range g.RelatedCommands {
			cmds[i] = "explain " + name
		}
		fmt.Fprintf(out, "Commands:  %s\n", strings.Join(cmds, ", "))
	}
	for _, id := range g.RelatedGuides {
		if related, ok := c.FindGuide(id); ok {
			fmt.Fprintf(out, "Related:   %s\n           explain learn %s\n", related.Title, related.ID)
		}
	}
	fmt.Fprintln(out)
	writeHeading(out, color, "References")
	for _, ref := range g.References {
		fmt.Fprintln(out, ref)
	}
	return 0
}

func listGuides(c *catalog.Catalog, op options, out io.Writer, color bool) int {
	if op.format == "json" {
		type seriesOut struct {
			catalog.SeriesInfo
			Guides []map[string]any `json:"guides"`
		}
		list := []seriesOut{}
		for _, s := range catalog.Series {
			entry := seriesOut{SeriesInfo: s, Guides: []map[string]any{}}
			for i, g := range c.Guides {
				if g.Series == s.ID {
					entry.Guides = append(entry.Guides, guideSummary(g, i+1))
				}
			}
			list = append(list, entry)
		}
		return encode(out, map[string]any{"schema_version": c.SchemaVersion, "content_revision": c.Revision, "series": list})
	}
	titleWidth := 0
	for _, g := range c.Guides {
		titleWidth = max(titleWidth, len([]rune(g.Title)))
	}
	if len(c.Guides) == 0 {
		fmt.Fprintln(out, "No guides yet.")
		return 0
	}
	for _, s := range catalog.Series {
		first := true
		for i, g := range c.Guides {
			if g.Series != s.ID {
				continue
			}
			if first {
				writeHeading(out, color, s.Title)
				fmt.Fprintln(out, s.Promise)
				first = false
			}
			fmt.Fprintf(out, "  %2d  %-*s  %2d min  explain learn %s\n", i+1, titleWidth, g.Title, g.ReadingMinutes, g.ID)
		}
		if !first {
			fmt.Fprintln(out)
		}
	}
	fmt.Fprintln(out, "Read one with its number or name, for example: explain learn 1")
	fmt.Fprintln(out, "Look up a word with: explain glossary WORD")
	return 0
}

func runGlossary(c *catalog.Catalog, op options, out, errOut io.Writer, cfg Config) int {
	if err := guideFlags(op, false); err != nil {
		return failure(out, errOut, op.format, err)
	}
	if len(op.args) > 2 {
		return failure(out, errOut, op.format, fmt.Errorf("glossary takes one term; quote words that contain spaces"))
	}
	color := op.color == "always" || (op.color == "auto" && cfg.Color)
	width := cfg.guideWidth(op)
	if len(op.args) == 1 {
		if op.format == "json" {
			return encode(out, map[string]any{"schema_version": c.SchemaVersion, "content_revision": c.Revision, "glossary": c.Glossary})
		}
		if len(c.Glossary.Terms) == 0 {
			fmt.Fprintln(out, "No glossary terms yet.")
			return 0
		}
		termWidth := 0
		for _, t := range c.Glossary.Terms {
			termWidth = max(termWidth, len([]rune(t.Term)))
		}
		writeHeading(out, color, c.Glossary.Title)
		fmt.Fprintln(out, c.Glossary.Summary)
		fmt.Fprintln(out)
		listWidth := width
		if width == 0 { // --width=0 turns wrapping off
			listWidth = 1000
		}
		for _, t := range c.Glossary.Terms {
			lines := strings.Split(strings.TrimRight(catalog.Plain(t.FirstSentence(), max(listWidth-termWidth-2, 20)), "\n"), "\n")
			for i, line := range lines {
				label := ""
				if i == 0 {
					label = t.Term
				}
				fmt.Fprintf(out, "%-*s  %s\n", termWidth, label, line)
			}
		}
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Read one in full with: explain glossary WORD")
		return 0
	}
	t, ok := c.FindTerm(op.args[1])
	if !ok {
		// Suggest terms whose name contains the typed text; only when none do,
		// fall back to terms whose definition mentions it.
		var names []string
		query := strings.ToLower(op.args[1])
		for _, t := range c.Glossary.Terms {
			if strings.Contains(strings.ToLower(t.Term), query) {
				names = append(names, t.Term)
			}
		}
		if len(names) == 0 {
			for _, t := range c.SearchTerms(op.args[1]) {
				names = append(names, t.Term)
			}
		}
		hint := ""
		if len(names) > 0 {
			hint = "; similar: " + strings.Join(names[:min(len(names), 5)], ", ")
		}
		return failure(out, errOut, op.format, fmt.Errorf("no glossary entry for %q%s; run explain glossary to list them", op.args[1], hint))
	}
	if op.format == "json" {
		return encode(out, map[string]any{"schema_version": c.SchemaVersion, "content_revision": c.Revision, "term": t})
	}
	writeHeading(out, color, t.Term)
	fmt.Fprintln(out, catalog.Plain(t.Markdown, width))
	if len(t.Guides) > 0 {
		fmt.Fprintln(out, "Learn more:")
	}
	for _, id := range t.Guides {
		if g, ok := c.FindGuide(id); ok {
			fmt.Fprintf(out, "  explain learn %s  (%s)\n", g.ID, g.Title)
		}
	}
	return 0
}

// learnMore prints the guides a command page points to.
func learnMore(c *catalog.Catalog, e *catalog.Entry, out io.Writer, color bool) {
	if len(e.RelatedGuides) == 0 {
		return
	}
	writeHeading(out, color, "Learn more")
	for _, id := range e.RelatedGuides {
		if g, ok := c.FindGuide(id); ok {
			fmt.Fprintf(out, "explain learn %s  (%s)\n", g.ID, g.Title)
		}
	}
	fmt.Fprintln(out)
}
