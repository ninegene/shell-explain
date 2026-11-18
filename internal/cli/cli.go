// Package cli implements explain's offline interface.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/detect"
)

var Help = `explain — practical, offline command help

Usage:
  explain COMMAND [flags]
  explain list [--format json]
  explain search WORDS [--format json]
  explain learn [GUIDE] [--section NAME]       Beginner guides; GUIDE is a name or number
  explain glossary [WORD]                      Short definitions of words you will see

Flags (may appear before or after COMMAND):
  --implementation ID                          Select documentation (see below)
  --platform macos|linux|bash                   Target platform, not detection
  --option=-n                                  Show a documented option
  --details                                    Show full explanations, options and references
  --format text|json                           Output format (default text)
  --width N                                    Prose width; compact descriptions capped at 100
  --section NAME                               Show one section of a guide (learn only)
  --no-pager                                   Print guides and the glossary without a pager
  --color auto|always|never                     Heading color (default auto)
  --no-probe                                   Disable native --version probes
  --help                                      Show this help
  --version                                   Show binary and content versions

Implementations: ` + implementationList + `.
Text output shows example usages by default. --option selects an option directly.
JSON always includes the full selected reference, independently of --details.
Examples are never executed. PATH detection cannot see parent-shell aliases/functions.
`

var implementationList = catalog.ImplementationIDs(", ")

type Config struct {
	Detector detect.Detector
	Version  string
	Width    int
	// GuideWidth is the prose width for guides: the terminal width when known.
	// Zero means use the default of 80 columns.
	GuideWidth int
	Color      bool
}
type options struct {
	implementation, platform, option, format, color, section string
	width                                                    int
	widthSet                                                 bool
	noProbe, help, version, details, noPager                 bool
	args                                                     []string
}
type response struct {
	SchemaVersion   int               `json:"schema_version"`
	ContentRevision string            `json:"content_revision"`
	Command         *catalog.Entry    `json:"command"`
	Selection       detect.Selection  `json:"selection"`
	Sections        []catalog.Section `json:"sections"`
}

func Run(ctx context.Context, c *catalog.Catalog, args []string, out, errOut io.Writer, cfg Config) int {
	op, err := parse(args, cfg.Width)
	if err != nil {
		return failure(out, errOut, op.format, err)
	}
	if op.help || len(args) == 0 {
		fmt.Fprint(out, Help)
		return 0
	}
	if op.version {
		fmt.Fprintf(out, "explain %s (schema %d, content %s)\n", cfg.Version, c.SchemaVersion, c.Revision)
		return 0
	}
	if len(op.args) == 0 {
		return failure(out, errOut, op.format, fmt.Errorf("specify a command, list, or search"))
	}
	switch op.args[0] {
	case "learn":
		return runLearn(c, op, out, errOut, cfg)
	case "glossary":
		return runGlossary(c, op, out, errOut, cfg)
	case "list", "search":
		if op.implementation != "" || op.platform != "" || op.option != "" || op.details {
			return failure(out, errOut, op.format, fmt.Errorf("selection flags apply to command lookup only"))
		}
		entries := c.Entries
		var guides []catalog.Guide
		var terms []catalog.GlossaryTerm
		if op.args[0] == "search" {
			if len(op.args) < 2 {
				return failure(out, errOut, op.format, fmt.Errorf("search needs a query"))
			}
			query := strings.Join(op.args[1:], " ")
			entries = c.Search(query)
			guides = c.SearchGuides(query)
			terms = c.SearchTerms(query)
		} else if len(op.args) > 1 {
			return failure(out, errOut, op.format, fmt.Errorf("list takes no arguments"))
		}
		if op.format == "json" {
			result := map[string]any{"schema_version": c.SchemaVersion, "content_revision": c.Revision, "entries": entries}
			if op.args[0] == "search" {
				summaries := []map[string]any{}
				for _, g := range guides {
					summaries = append(summaries, guideSummary(g, c.GuideNumber(g.ID)))
				}
				result["guides"] = summaries
				result["terms"] = terms
			}
			return encode(out, result)
		}
		for _, e := range entries {
			fmt.Fprintf(out, "%-8s %s\n", e.ID, e.Summary)
		}
		for _, g := range guides {
			fmt.Fprintf(out, "guide    explain learn %s  (%s)\n", g.ID, g.Title)
		}
		for _, t := range terms {
			fmt.Fprintf(out, "word     explain glossary %q  (%s)\n", t.Term, t.FirstSentence())
		}
		if len(entries)+len(guides)+len(terms) == 0 {
			fmt.Fprintln(out, "Nothing matched. Try: explain list, explain learn, or explain glossary")
		}
		return 0
	}
	if len(op.args) != 1 {
		return failure(out, errOut, op.format, fmt.Errorf("explain accepts one command name, not a command line"))
	}
	name := op.args[0]
	e, ok := c.Find(filepath.Base(name))
	if !ok {
		return failure(out, errOut, op.format, fmt.Errorf("no reference for %q; try explain list or explain search", name))
	}
	sel, err := cfg.Detector.Select(ctx, detect.Request{Name: name, Command: e.ID, Implementation: op.implementation, Platform: op.platform, NoProbe: op.noProbe, ShellBuiltin: e.ShellOnly(), Shells: e.ShellImplementations(), Single: singleImplementation(e)})
	if err != nil {
		return failure(out, errOut, op.format, err)
	}
	sections, err := e.Select(sel.Implementation, op.option)
	if err != nil {
		return failure(out, errOut, op.format, err)
	}
	if op.format == "json" {
		meta := *e
		meta.Sections = nil
		return encode(out, response{c.SchemaVersion, c.Revision, &meta, sel, sections})
	}
	color := op.color == "always" || (op.color == "auto" && cfg.Color)
	heading := func(s string) {
		if color {
			fmt.Fprintln(out, paint(true, ansiHeading, s))
		} else {
			fmt.Fprintln(out, s)
		}
	}
	single := singleImplementation(e) != ""
	if !op.details {
		if op.option != "" {
			heading(fmt.Sprintf("%s — %s (%s)", e.Title, sel.Implementation, sel.Confidence))
			fmt.Fprintln(out)
			heading("Options")
			fmt.Fprint(out, RenderOptions(sections, op.width, color))
			return 0
		}
		heading(e.Title + " — " + e.Summary)
		if !single {
			// Which implementation was chosen only matters when there is a choice.
			fmt.Fprintln(out, paint(color, ansiLabel, fmt.Sprintf("showing %s (%s)", sel.Implementation, sel.Confidence)))
		}
		fmt.Fprintln(out)
		for _, s := range sections {
			if s.IsExamples() {
				fmt.Fprint(out, RenderExamples(catalog.ParseExamples(s.Markdown), op.width, color))
			}
		}
		writeLinks(out, e, sel.Implementation, single, color)
		return 0
	}
	heading(e.Title + " — " + e.Summary)
	meta := func(label, value string) {
		fmt.Fprintln(out, paint(color, ansiLabel, label)+" "+value)
	}
	fmt.Fprintln(out)
	meta("Documentation:", fmt.Sprintf("%s (%s)", sel.Implementation, sel.Confidence))
	if sel.Executable != "" {
		meta("Executable:", sel.Executable)
	}
	if sel.Version != "" {
		meta("Version:", sel.Version)
	}
	meta("Selection:", sel.Evidence)
	meta("Last updated:", e.Updated)
	fmt.Fprintln(out, catalog.Plain(e.VersionPolicy, op.width))
	for _, warning := range sel.Warnings {
		fmt.Fprint(out, catalog.Plain("Note: "+warning, op.width))
	}
	optionsShown := false
	for _, s := range sections {
		if len(s.Options) > 0 {
			if !optionsShown {
				heading("Options")
				fmt.Fprintln(out, RenderOptions(sections, op.width, color))
				optionsShown = true
			}
			continue
		}
		heading(s.Title)
		if s.IsExamples() {
			fmt.Fprintln(out, RenderExamples(catalog.ParseExamples(s.Markdown), op.width, color))
		} else {
			fmt.Fprintln(out, catalog.Plain(s.Markdown, op.width))
		}
	}
	learnMore(c, e, out, color)
	heading("References")
	for _, ref := range e.References {
		fmt.Fprintln(out, paint(color, ansiLink, ref))
	}
	fmt.Fprintln(out)
	writeMore(out, e, sel.Implementation, single, color)
	return 0
}

// singleImplementation returns the only non-shell implementation an entry
// documents when that implementation is "upstream".
func singleImplementation(e *catalog.Entry) string {
	if len(e.Implementations) == 1 && e.Implementations[0] == "upstream" {
		return "upstream"
	}
	return ""
}
func encode(w io.Writer, v any) int {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return 1
	}
	return 0
}
func failure(out, errOut io.Writer, format string, err error) int {
	if format == "json" {
		_ = encode(out, map[string]any{"schema_version": catalog.SchemaVersion, "error": map[string]string{"code": "invalid_request", "message": err.Error()}})
	} else {
		fmt.Fprintln(errOut, "explain:", err)
	}
	return 2
}
func parse(args []string, width int) (options, error) {
	o := options{format: "text", color: "auto", width: width, args: []string{}}
	literal := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		if literal {
			o.args = append(o.args, a)
			continue
		}
		if a == "--" {
			literal = true
			continue
		}
		if !strings.HasPrefix(a, "-") {
			o.args = append(o.args, a)
			continue
		}
		switch a {
		case "--help", "-h":
			o.help = true
			continue
		case "--version":
			o.version = true
			continue
		case "--no-probe":
			o.noProbe = true
			continue
		case "--details":
			o.details = true
			continue
		case "--no-pager":
			o.noPager = true
			continue
		}
		key, val, has := strings.Cut(a, "=")
		if !has {
			if i+1 == len(args) {
				return o, fmt.Errorf("flag %s needs a value", key)
			}
			i++
			val = args[i]
		}
		switch key {
		case "--implementation":
			o.implementation = val
		case "--platform":
			o.platform = val
		case "--option":
			o.option = val
		case "--section":
			o.section = val
		case "--format":
			o.format = val
		case "--color":
			o.color = val
		case "--width":
			n, err := strconv.Atoi(val)
			if err != nil || n < 0 || n > 1000 {
				return o, fmt.Errorf("width must be 0–1000")
			}
			o.width, o.widthSet = n, true
		default:
			return o, fmt.Errorf("unknown flag %s", key)
		}
	}
	if o.format != "text" && o.format != "json" {
		return o, fmt.Errorf("format must be text or json")
	}
	if o.color != "auto" && o.color != "always" && o.color != "never" {
		return o, fmt.Errorf("color must be auto, always, or never")
	}
	return o, nil
}
func TerminalColor(out *os.File) bool {
	info, err := out.Stat()
	_, noColor := os.LookupEnv("NO_COLOR")
	return err == nil && info.Mode()&os.ModeCharDevice != 0 && !noColor && os.Getenv("TERM") != "dumb"
}
