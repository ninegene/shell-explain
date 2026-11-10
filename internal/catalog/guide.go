package catalog

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.yaml.in/yaml/v3"
)

// Limits that keep guides readable in a terminal and on a phone.
const (
	// MaxDiagramWidth is the widest line of a diagram block, in columns.
	MaxDiagramWidth = 72
	// MaxTableWidth is the widest rendered line of a table, in columns.
	MaxTableWidth = 78
	// MaxTableColumns is the most columns a table may have.
	MaxTableColumns = 4
)

// SeriesInfo describes a group of guides that are read in order.
type SeriesInfo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Promise string `json:"promise"`
}

// Series is the registry of guide series, in display order. A guide's series
// field must name one of them.
var Series = []SeriesInfo{
	{ID: "how-a-computer-works", Title: "How a Computer Works, Just Enough", Promise: "What happens between your keyboard and the screen."},
	{ID: "terminal-shell-commands", Title: "Terminal, Shell and Commands", Promise: "What you are typing into, and how it becomes a result."},
	{ID: "why-so-many-versions", Title: "Why Are There So Many Versions?", Promise: "Where Unix, BSD, GNU and Linux came from, and why one command differs from system to system."},
	{ID: "getting-things-done", Title: "Getting Things Done", Promise: "Start from the result you want and pick a route to it."},
}

// LookupSeries returns a registered series by ID.
func LookupSeries(id string) (SeriesInfo, bool) {
	for _, s := range Series {
		if s.ID == id {
			return s, true
		}
	}
	return SeriesInfo{}, false
}

// GuideSections are the required level-two sections of every guide, in order.
var GuideSections = []string{
	"Goal",
	"The idea in plain words",
	"A picture",
	"Another way to see it",
	"Walk through it",
	"Try it",
	"Common confusions",
	"Check yourself",
	"Where next",
}

// GuideMetadata is the front matter of a guide.
type GuideMetadata struct {
	SchemaVersion   int      `yaml:"schema_version" json:"schema_version"`
	ID              string   `yaml:"id" json:"id"`
	Title           string   `yaml:"title" json:"title"`
	Subtitle        string   `yaml:"subtitle" json:"subtitle"`
	Summary         string   `yaml:"summary" json:"summary"`
	Series          string   `yaml:"series" json:"series"`
	Order           int      `yaml:"order" json:"order"`
	Level           string   `yaml:"level" json:"level"`
	ReadingMinutes  int      `yaml:"reading_minutes" json:"reading_minutes"`
	Updated         string   `yaml:"updated" json:"updated"`
	Tags            []string `yaml:"tags" json:"tags"`
	RelatedCommands []string `yaml:"related_commands" json:"related_commands"`
	RelatedGuides   []string `yaml:"related_guides" json:"related_guides"`
	References      []string `yaml:"references" json:"references"`
}

// Guide is one beginner guide: metadata plus its nine sections.
type Guide struct {
	GuideMetadata
	Sections []Section `json:"sections"`
}

// GlossaryTerm is one short definition.
type GlossaryTerm struct {
	ID       string   `json:"id"`
	Term     string   `json:"term"`
	Markdown string   `json:"markdown"`
	HTML     string   `json:"html"`
	Code     []string `json:"code"`
	Guides   []string `json:"guides"`
}

// Glossary is the list of terms, with the front matter of terms.md.
type Glossary struct {
	SchemaVersion int            `json:"schema_version"`
	Title         string         `json:"title"`
	Summary       string         `json:"summary"`
	Updated       string         `json:"updated"`
	Terms         []GlossaryTerm `json:"terms"`
}

// frontMatter checks the file-level rules shared by every authored file and
// splits it into front matter and body.
func frontMatter(data []byte) (front, body []byte, err error) {
	if !utf8.Valid(data) {
		return nil, nil, fmt.Errorf("content must be UTF-8")
	}
	if bytes.Contains(data, []byte("\r")) {
		return nil, nil, fmt.Errorf("use LF line endings")
	}
	for _, b := range string(data) {
		if unicode.IsControl(b) && b != '\n' && b != '\t' {
			return nil, nil, fmt.Errorf("control characters are not allowed")
		}
	}
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, nil, fmt.Errorf("missing YAML front matter")
	}
	parts := bytes.SplitN(data[4:], []byte("\n---\n"), 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("unclosed front matter")
	}
	return parts[0], parts[1], nil
}

func validDate(s string) error {
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return fmt.Errorf("updated must be a calendar date written YYYY-MM-DD, got %q", s)
	}
	return nil
}

func validateGuideMetadata(m GuideMetadata) error {
	if m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema version %d", m.SchemaVersion)
	}
	if !names.MatchString(m.ID) {
		return fmt.Errorf("invalid guide id %q", m.ID)
	}
	for _, s := range []string{m.Title, m.Subtitle, m.Summary} {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("title, subtitle and summary are required")
		}
		for _, r := range s {
			if unicode.IsControl(r) {
				return fmt.Errorf("metadata must not contain control characters")
			}
		}
	}
	if _, ok := LookupSeries(m.Series); !ok {
		return fmt.Errorf("unknown series %q", m.Series)
	}
	if m.Order < 1 {
		return fmt.Errorf("order must be 1 or more")
	}
	if m.Level != "beginner" {
		return fmt.Errorf("level must be beginner, got %q", m.Level)
	}
	if m.ReadingMinutes < 1 || m.ReadingMinutes > 60 {
		return fmt.Errorf("reading_minutes must be between 1 and 60")
	}
	if err := validDate(m.Updated); err != nil {
		return err
	}
	if len(m.References) == 0 {
		return fmt.Errorf("references are required")
	}
	for _, ref := range m.References {
		if !safeURL(ref) {
			return fmt.Errorf("unsafe reference URL %q", ref)
		}
	}
	for _, list := range [][]string{m.Tags, m.RelatedCommands, m.RelatedGuides} {
		seen := map[string]bool{}
		for _, n := range list {
			if !IsLookupName(n) || seen[n] {
				return fmt.Errorf("invalid or duplicate name %q", n)
			}
			seen[n] = true
		}
	}
	return nil
}

// ParseGuide parses and validates one guide file.
func ParseGuide(data []byte) (*Guide, error) {
	front, body, err := frontMatter(data)
	if err != nil {
		return nil, err
	}
	g := &Guide{Sections: []Section{}}
	dec := yaml.NewDecoder(bytes.NewReader(front))
	dec.KnownFields(true)
	if err := dec.Decode(&g.GuideMetadata); err != nil {
		return nil, err
	}
	if err := validateGuideMetadata(g.GuideMetadata); err != nil {
		return nil, err
	}
	sections, err := splitSections(body)
	if err != nil {
		return nil, err
	}
	if len(sections) != len(GuideSections) {
		return nil, fmt.Errorf("a guide needs exactly these sections in order: %s", strings.Join(GuideSections, ", "))
	}
	for i, s := range sections {
		if s.title != GuideSections[i] {
			return nil, fmt.Errorf("section %d is %q, want %q", i+1, s.title, GuideSections[i])
		}
		html, code, err := RenderGuideHTML(s.markdown)
		if err != nil {
			return nil, fmt.Errorf("section %q: %w", s.title, err)
		}
		g.Sections = append(g.Sections, Section{
			ID:             fmt.Sprintf("%s-%02d", g.ID, i+1),
			Implementation: "common",
			Title:          s.title,
			Options:        []string{},
			Markdown:       s.markdown,
			HTML:           html,
			Code:           code,
		})
	}
	return g, nil
}

type rawSection struct{ title, markdown string }

// splitSections cuts a body into level-two sections. Level-three headings stay
// inside their section; any other heading level, or text before the first
// section, is rejected. Headings inside code fences are literal text.
func splitSections(body []byte) ([]rawSection, error) {
	headings, err := sectionHeadings(body)
	if err != nil {
		return nil, err
	}
	var out []rawSection
	var lines []string
	title := ""
	flush := func() error {
		md := strings.Trim(strings.Join(lines, "\n"), "\n") + "\n"
		if title == "" {
			if strings.TrimSpace(md) != "" {
				return fmt.Errorf("text must be inside a level-two section")
			}
		} else {
			if strings.TrimSpace(md) == "" {
				return fmt.Errorf("empty section %q", title)
			}
			out = append(out, rawSection{title: title, markdown: md})
		}
		lines = nil
		return nil
	}
	offset := 0
	for _, line := range strings.Split(string(body), "\n") {
		level := headings[offset]
		offset += len(line) + 1
		switch level {
		case 2:
			if err := flush(); err != nil {
				return nil, err
			}
			title = strings.TrimSpace(strings.TrimPrefix(line, "## "))
		case 0, 3:
			lines = append(lines, line)
		default:
			return nil, fmt.Errorf("heading level %d is not allowed; use ## for sections and ### inside them", level)
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	return out, nil
}

var diagramLine = regexp.MustCompile(`^[\x20-\x7e]*$`)

// RenderGuideHTML renders the guide subset of Markdown: everything commands
// allow, plus pipe tables, level-three headings and diagram blocks. A diagram
// is a fenced block with the language diagram, at most MaxDiagramWidth plain
// ASCII columns wide, and it must be followed by a paragraph starting with
// "Diagram:" that describes it in words. A table has at most MaxTableColumns
// columns and renders within MaxTableWidth columns.
func RenderGuideHTML(markdown string) (string, []string, error) {
	source := []byte(markdown)
	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	node := md.Parser().Parse(text.NewReader(source))
	code := []string{}
	err := ast.Walk(node, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Document, *ast.Paragraph, *ast.TextBlock, *ast.Text, *ast.String, *ast.CodeSpan, *ast.Emphasis, *ast.List, *ast.ListItem,
			*east.TableHeader, *east.TableRow, *east.TableCell:
		case *ast.Heading:
			if v.Level != 3 {
				return ast.WalkStop, fmt.Errorf("heading level %d is not allowed", v.Level)
			}
		case *east.Table:
			rows := tableRows(v, source)
			if len(rows) == 0 || len(rows[0]) > MaxTableColumns {
				return ast.WalkStop, fmt.Errorf("a table may have at most %d columns", MaxTableColumns)
			}
			for _, line := range strings.Split(layoutTable(rows), "\n") {
				if utf8.RuneCountInString(line) > MaxTableWidth {
					return ast.WalkStop, fmt.Errorf("table line is wider than %d columns: %q", MaxTableWidth, line)
				}
			}
		case *ast.Link:
			if !safeURL(string(v.Destination)) {
				return ast.WalkStop, fmt.Errorf("unsafe or unsupported link %q", v.Destination)
			}
		case *ast.FencedCodeBlock:
			lang := string(v.Language(source))
			if lang != "sh" && lang != "bash" && lang != "console" && lang != "text" && lang != "diagram" {
				return ast.WalkStop, fmt.Errorf("unsupported code language %q", lang)
			}
			var b strings.Builder
			for i := 0; i < v.Lines().Len(); i++ {
				seg := v.Lines().At(i)
				b.Write(seg.Value(source))
			}
			if lang == "diagram" {
				if err := checkDiagram(v, b.String(), source); err != nil {
					return ast.WalkStop, err
				}
			}
			code = append(code, b.String())
		default:
			return ast.WalkStop, fmt.Errorf("unsupported Markdown node %s", n.Kind().String())
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", nil, err
	}
	var out bytes.Buffer
	if err := md.Renderer().Render(&out, source, node); err != nil {
		return "", nil, err
	}
	return out.String(), code, nil
}

func checkDiagram(block *ast.FencedCodeBlock, content string, source []byte) error {
	for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
		if !diagramLine.MatchString(line) {
			return fmt.Errorf("diagram lines may contain only plain ASCII characters: %q", line)
		}
		if len(line) > MaxDiagramWidth {
			return fmt.Errorf("diagram line is wider than %d columns: %q", MaxDiagramWidth, line)
		}
	}
	next, ok := block.NextSibling().(*ast.Paragraph)
	if !ok || !strings.HasPrefix(strings.TrimSpace(string(next.Lines().Value(source))), "Diagram:") {
		return fmt.Errorf("a diagram must be followed by a paragraph that starts with %q and describes it in words", "Diagram:")
	}
	return nil
}

// ParseGlossary parses and validates the glossary file. Each level-two
// section is one term: its heading is the term, its body the definition. A
// final paragraph starting with "See:" lists guide IDs separated by commas.
func ParseGlossary(data []byte) (*Glossary, error) {
	front, body, err := frontMatter(data)
	if err != nil {
		return nil, err
	}
	var m struct {
		SchemaVersion int    `yaml:"schema_version"`
		ID            string `yaml:"id"`
		Title         string `yaml:"title"`
		Summary       string `yaml:"summary"`
		Updated       string `yaml:"updated"`
	}
	dec := yaml.NewDecoder(bytes.NewReader(front))
	dec.KnownFields(true)
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	if m.SchemaVersion != SchemaVersion || m.ID != "glossary" || strings.TrimSpace(m.Title) == "" || strings.TrimSpace(m.Summary) == "" {
		return nil, fmt.Errorf("glossary needs schema_version %d, id glossary, a title and a summary", SchemaVersion)
	}
	if err := validDate(m.Updated); err != nil {
		return nil, err
	}
	g := &Glossary{SchemaVersion: m.SchemaVersion, Title: m.Title, Summary: m.Summary, Updated: m.Updated, Terms: []GlossaryTerm{}}
	sections, err := splitSections(body)
	if err != nil {
		return nil, err
	}
	if len(sections) == 0 {
		return nil, fmt.Errorf("the glossary has no terms")
	}
	seen := map[string]bool{}
	for _, s := range sections {
		id := termID(s.title)
		if id == "" || seen[id] {
			return nil, fmt.Errorf("invalid or duplicate term %q", s.title)
		}
		seen[id] = true
		markdown, guides, err := splitSee(s.markdown)
		if err != nil {
			return nil, fmt.Errorf("term %q: %w", s.title, err)
		}
		html, code, err := RenderGuideHTML(markdown)
		if err != nil {
			return nil, fmt.Errorf("term %q: %w", s.title, err)
		}
		g.Terms = append(g.Terms, GlossaryTerm{ID: id, Term: s.title, Markdown: markdown, HTML: html, Code: code, Guides: guides})
	}
	return g, nil
}

var notTermID = regexp.MustCompile(`[^a-z0-9]+`)

// termID turns a term heading into a lowercase hyphenated anchor.
func termID(term string) string {
	id := strings.Trim(notTermID.ReplaceAllString(strings.ToLower(term), "-"), "-")
	if !names.MatchString(id) {
		return ""
	}
	return id
}

// splitSee removes a trailing "See: a, b" paragraph and returns the guide IDs.
func splitSee(markdown string) (string, []string, error) {
	trimmed := strings.TrimRight(markdown, "\n")
	i := strings.LastIndex(trimmed, "\n\n")
	last := trimmed
	rest := ""
	if i >= 0 {
		last, rest = trimmed[i+2:], trimmed[:i]
	}
	if !strings.HasPrefix(last, "See: ") {
		return markdown, []string{}, nil
	}
	if rest == "" {
		return "", nil, fmt.Errorf("a term needs a definition before its See: line")
	}
	guides := []string{}
	for _, id := range strings.Split(strings.TrimPrefix(last, "See: "), ",") {
		id = strings.TrimSpace(id)
		if !names.MatchString(id) {
			return "", nil, fmt.Errorf("invalid guide id %q in See:", id)
		}
		guides = append(guides, id)
	}
	return rest + "\n", guides, nil
}

const (
	guideDir    = "en/guides"
	glossaryDir = "en/glossary"
)

// loadGuides reads en/guides/*.md and en/glossary/terms.md when present, and
// adds their files to the content revision hash after the command pages.
func (c *Catalog) loadGuides(files fs.FS, hash io.Writer) error {
	paths, err := fs.Glob(files, guideDir+"/*.md")
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	slots := map[string]bool{}
	for _, p := range paths {
		data, err := fs.ReadFile(files, p)
		if err != nil {
			return err
		}
		g, err := ParseGuide(data)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		if path.Base(p) != g.ID+".md" {
			return fmt.Errorf("%s: filename must match id", p)
		}
		slot := fmt.Sprintf("%s/%d", g.Series, g.Order)
		if seen[g.ID] || slots[slot] {
			return fmt.Errorf("%s: duplicate guide id or position %s", p, slot)
		}
		seen[g.ID], slots[slot] = true, true
		hash.Write([]byte(p))
		hash.Write([]byte{0})
		hash.Write(data)
		hash.Write([]byte{0})
		c.Guides = append(c.Guides, *g)
	}
	// Reading order: by series position in the registry, then by order.
	slices.SortStableFunc(c.Guides, func(a, b Guide) int {
		return cmp.Or(seriesIndex(a.Series)-seriesIndex(b.Series), a.Order-b.Order)
	})
	terms := glossaryDir + "/terms.md"
	data, err := fs.ReadFile(files, terms)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	g, err := ParseGlossary(data)
	if err != nil {
		return fmt.Errorf("%s: %w", terms, err)
	}
	c.Glossary = *g
	hash.Write([]byte(terms))
	hash.Write([]byte{0})
	hash.Write(data)
	hash.Write([]byte{0})
	return nil
}

func seriesIndex(id string) int {
	for i, s := range Series {
		if s.ID == id {
			return i
		}
	}
	return len(Series)
}

// checkReferences makes sure every cross-reference between commands, guides
// and glossary terms names something that exists.
func (c *Catalog) checkReferences() error {
	guides := map[string]bool{}
	for _, g := range c.Guides {
		guides[g.ID] = true
	}
	for _, g := range c.Guides {
		for _, name := range g.RelatedCommands {
			if _, ok := c.Find(name); !ok {
				return fmt.Errorf("guide %s: related command %q has no page", g.ID, name)
			}
		}
		for _, id := range g.RelatedGuides {
			if !guides[id] || id == g.ID {
				return fmt.Errorf("guide %s: related guide %q does not exist or is itself", g.ID, id)
			}
		}
	}
	for _, e := range c.Entries {
		for _, id := range e.RelatedGuides {
			if !guides[id] {
				return fmt.Errorf("command %s: related guide %q does not exist", e.ID, id)
			}
		}
	}
	for _, t := range c.Glossary.Terms {
		for _, id := range t.Guides {
			if !guides[id] {
				return fmt.Errorf("glossary term %q: guide %q does not exist", t.Term, id)
			}
		}
	}
	return nil
}

// FindGuide returns a guide by ID, or by its 1-based number in reading order.
func (c *Catalog) FindGuide(key string) (*Guide, bool) {
	for i := range c.Guides {
		if c.Guides[i].ID == key || fmt.Sprint(i+1) == key {
			return &c.Guides[i], true
		}
	}
	return nil, false
}

// GuideNumber is the guide's 1-based position in reading order, or 0.
func (c *Catalog) GuideNumber(id string) int {
	for i := range c.Guides {
		if c.Guides[i].ID == id {
			return i + 1
		}
	}
	return 0
}

// GuideNeighbors returns the guides before and after id in reading order; a
// missing neighbor is nil.
func (c *Catalog) GuideNeighbors(id string) (prev, next *Guide) {
	n := c.GuideNumber(id)
	if n == 0 {
		return nil, nil
	}
	if n > 1 {
		prev = &c.Guides[n-2]
	}
	if n < len(c.Guides) {
		next = &c.Guides[n]
	}
	return prev, next
}

// Section returns the guide section whose title or slug (for example
// "try-it") matches name, ignoring case.
func (g *Guide) Section(name string) (Section, bool) {
	want := termID(name)
	for _, s := range g.Sections {
		if strings.EqualFold(s.Title, name) || termID(s.Title) == want {
			return s, true
		}
	}
	return Section{}, false
}

// FindTerm returns a glossary term by its words or anchor, ignoring case.
func (c *Catalog) FindTerm(name string) (*GlossaryTerm, bool) {
	want := termID(name)
	for i := range c.Glossary.Terms {
		t := &c.Glossary.Terms[i]
		if strings.EqualFold(t.Term, name) || (want != "" && t.ID == want) {
			return t, true
		}
	}
	return nil, false
}

// matchesAll reports whether every query word occurs in text, ignoring case.
func matchesAll(text string, words []string) bool {
	text = strings.ToLower(text)
	for _, w := range words {
		if !strings.Contains(text, w) {
			return false
		}
	}
	return true
}

// SearchGuides finds guides whose ID, titles, summary or tags contain every word.
func (c *Catalog) SearchGuides(query string) []Guide {
	result := []Guide{}
	words := strings.Fields(strings.ToLower(query))
	for _, g := range c.Guides {
		if matchesAll(g.ID+" "+g.Title+" "+g.Subtitle+" "+g.Summary+" "+strings.Join(g.Tags, " "), words) {
			result = append(result, g)
		}
	}
	return result
}

// SearchTerms finds glossary terms whose words or definition contain every word.
func (c *Catalog) SearchTerms(query string) []GlossaryTerm {
	result := []GlossaryTerm{}
	words := strings.Fields(strings.ToLower(query))
	for _, t := range c.Glossary.Terms {
		if matchesAll(t.Term+" "+t.ID+" "+t.Markdown, words) {
			result = append(result, t)
		}
	}
	return result
}

// FirstSentence returns the first sentence of a term's definition as plain text.
func (t GlossaryTerm) FirstSentence() string {
	plain := strings.Join(strings.Fields(Plain(t.Markdown, 1000)), " ")
	if i := strings.Index(plain, ". "); i >= 0 {
		return plain[:i+1]
	}
	return plain
}
