// Package catalog validates and normalizes the shared Markdown collection.
package catalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.yaml.in/yaml/v3"
)

const SchemaVersion = 1
const VersionPolicy = "Implementation-level reference; examples are not certified for every installed version."

var names = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var optionName = regexp.MustCompile(`^--?[A-Za-z0-9][A-Za-z0-9-]*$`)

// symbolicNames are lookup spellings that cannot be IDs, filenames or URL
// segments. Add a spelling here only for a real command name, such as test's `[`.
var symbolicNames = map[string]bool{"[": true}

// Implementation is one selectable documentation view. Implementations are
// distinct upstream projects or shells, never operating-system guesses.
type Implementation struct {
	ID      string `json:"id"`
	Heading string `json:"heading"`
	Label   string `json:"label"`
	Shell   bool   `json:"shell"`
}

// Implementations is the single registry used by the parser, CLI and website.
var Implementations = []Implementation{
	{ID: "common", Heading: "Common", Label: "Portable / common"},
	{ID: "gnu", Heading: "GNU", Label: "GNU"},
	{ID: "apple", Heading: "Apple", Label: "macOS / Apple"},
	{ID: "busybox", Heading: "BusyBox", Label: "BusyBox"},
	{ID: "mawk", Heading: "Mawk", Label: "mawk"},
	{ID: "bash", Heading: "Bash", Label: "Bash builtin", Shell: true},
	{ID: "zsh", Heading: "Zsh", Label: "Zsh builtin", Shell: true},
	{ID: "upstream", Heading: "Upstream", Label: "Upstream project"},
}

// LookupImplementation returns a registered implementation by ID.
func LookupImplementation(id string) (Implementation, bool) {
	for _, impl := range Implementations {
		if impl.ID == id {
			return impl, true
		}
	}
	return Implementation{}, false
}

// ImplementationIDs lists registered IDs in display order, separated by sep.
func ImplementationIDs(sep string) string {
	ids := make([]string, len(Implementations))
	for i, impl := range Implementations {
		ids[i] = impl.ID
	}
	return strings.Join(ids, sep)
}

func variantForHeading(heading string) (string, bool) {
	for _, impl := range Implementations {
		if impl.Heading == heading {
			return impl.ID, true
		}
	}
	return "", false
}

type Metadata struct {
	SchemaVersion   int      `yaml:"schema_version" json:"schema_version"`
	ID              string   `yaml:"id" json:"id"`
	Title           string   `yaml:"title" json:"title"`
	Summary         string   `yaml:"summary" json:"summary"`
	Aliases         []string `yaml:"aliases" json:"aliases"`
	Tags            []string `yaml:"tags" json:"tags"`
	Implementations []string `yaml:"implementations" json:"implementations"`
	// Updated is the date (YYYY-MM-DD) the page's content was last reviewed or
	// changed. Authors set it by hand when they edit a page's claims.
	Updated    string   `yaml:"updated" json:"updated"`
	References []string `yaml:"references" json:"references"`
	// RelatedGuides names beginner guides that explain ideas this page assumes.
	// Optional; each name must be an existing guide.
	RelatedGuides []string `yaml:"related_guides" json:"related_guides,omitempty"`
}

type Section struct {
	ID             string   `json:"id"`
	Implementation string   `json:"implementation"`
	Title          string   `json:"title"`
	Options        []string `json:"options"`
	Markdown       string   `json:"markdown"`
	HTML           string   `json:"html"`
	Code           []string `json:"code"`
	ExampleText    string   `json:"example_text,omitempty"`
	OptionText     string   `json:"option_text,omitempty"`
}

// IsLookupName reports whether s is a valid command ID or symbolic alias spelling.
func IsLookupName(s string) bool { return names.MatchString(s) || symbolicNames[s] }

func (s Section) IsExamples() bool {
	return s.Title == "Examples" || strings.HasSuffix(s.Title, " examples")
}

type Entry struct {
	Metadata
	VersionPolicy string    `json:"version_policy"`
	Sections      []Section `json:"sections"`
}

// ShellImplementations lists the entry's declared shell-builtin views.
func (e *Entry) ShellImplementations() []string {
	shells := []string{}
	for _, id := range e.Implementations {
		if impl, ok := LookupImplementation(id); ok && impl.Shell {
			shells = append(shells, id)
		}
	}
	return shells
}

// ShellOnly reports whether every declared implementation is a shell builtin,
// so no executable on PATH can describe what the user's shell runs.
func (e *Entry) ShellOnly() bool {
	return len(e.Implementations) > 0 && len(e.ShellImplementations()) == len(e.Implementations)
}

type Catalog struct {
	SchemaVersion   int              `json:"schema_version"`
	Revision        string           `json:"revision"`
	Implementations []Implementation `json:"implementations"`
	Entries         []Entry          `json:"entries"`
	Series          []SeriesInfo     `json:"series"`
	Guides          []Guide          `json:"guides"`
	Glossary        Glossary         `json:"glossary"`
}

func Load(files fs.FS) (*Catalog, error) {
	paths, err := fs.Glob(files, "en/*.md")
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no command Markdown files found")
	}
	c := &Catalog{SchemaVersion: SchemaVersion, Implementations: Implementations, Entries: []Entry{}, Series: Series, Guides: []Guide{}, Glossary: Glossary{Terms: []GlossaryTerm{}}}
	hash := sha256.New()
	seen := map[string]bool{}
	for _, p := range paths {
		data, err := fs.ReadFile(files, p)
		if err != nil {
			return nil, err
		}
		entry, err := Parse(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if path.Base(p) != entry.ID+".md" {
			return nil, fmt.Errorf("%s: filename must match id", p)
		}
		for _, name := range append([]string{entry.ID}, entry.Aliases...) {
			if seen[name] {
				return nil, fmt.Errorf("duplicate command or alias %q", name)
			}
			seen[name] = true
		}
		hash.Write([]byte(p))
		hash.Write([]byte{0})
		hash.Write(data)
		hash.Write([]byte{0})
		c.Entries = append(c.Entries, *entry)
	}
	if err := c.loadGuides(files, hash); err != nil {
		return nil, err
	}
	if err := c.checkReferences(); err != nil {
		return nil, err
	}
	c.Revision = hex.EncodeToString(hash.Sum(nil))
	return c, nil
}

func Parse(data []byte) (*Entry, error) {
	front, rest, err := frontMatter(data)
	if err != nil {
		return nil, err
	}
	e := &Entry{VersionPolicy: VersionPolicy, Sections: []Section{}}
	dec := yaml.NewDecoder(bytes.NewReader(front))
	dec.KnownFields(true)
	if err := dec.Decode(&e.Metadata); err != nil {
		return nil, err
	}
	if err := validateMetadata(e.Metadata); err != nil {
		return nil, err
	}
	body := string(rest)
	headings, err := sectionHeadings(rest)
	if err != nil {
		return nil, err
	}
	view := ""
	title := ""
	var lines []string
	seen := map[string]bool{}
	declared := map[string]bool{}
	flush := func() error {
		if title == "" {
			if strings.TrimSpace(strings.Join(lines, "\n")) != "" {
				return fmt.Errorf("text must be inside a level-three section")
			}
			lines = nil
			return nil
		}
		md := strings.Trim(strings.Join(lines, "\n"), "\n") + "\n"
		if strings.TrimSpace(md) == "" {
			return fmt.Errorf("empty section %q", title)
		}
		section := Section{ID: fmt.Sprintf("%s-%02d", view, len(e.Sections)+1), Implementation: view, Title: title, Options: []string{}, Markdown: md, Code: []string{}}
		if strings.HasPrefix(title, "Option: ") {
			for _, op := range strings.Split(strings.TrimPrefix(title, "Option: "), ",") {
				op = strings.TrimSpace(op)
				if !optionName.MatchString(op) {
					return fmt.Errorf("invalid option %q", op)
				}
				section.Options = append(section.Options, op)
			}
		}
		if title == "Background" && view != "common" {
			return fmt.Errorf("Background belongs in Common, not %q", view)
		}
		if strings.HasPrefix(title, "About ") {
			if view == "common" {
				return fmt.Errorf("%q belongs in an implementation view, not Common", title)
			}
			if seen[view+"/about"] {
				return fmt.Errorf("more than one About section in %q", view)
			}
			seen[view+"/about"] = true
		}
		key := view + "/" + title
		if seen[key] {
			return fmt.Errorf("duplicate section %q", key)
		}
		seen[key] = true
		for _, op := range section.Options {
			key = view + "/option/" + op
			if seen[key] {
				return fmt.Errorf("duplicate option %q", op)
			}
			seen[key] = true
		}
		html, code, err := RenderHTML(md)
		if err != nil {
			return err
		}
		section.HTML = html
		section.Code = code
		if section.IsExamples() {
			section.ExampleText = ExampleText(md)
		}
		if len(section.Options) > 0 {
			section.OptionText = OptionText(section, DescriptionWidth)
		}
		e.Sections = append(e.Sections, section)
		lines = nil
		return nil
	}
	offset := 0
	for _, line := range strings.Split(body, "\n") {
		level := headings[offset]
		offset += len(line) + 1
		if level == 2 {
			if err := flush(); err != nil {
				return nil, err
			}
			title = ""
			var ok bool
			view, ok = variantForHeading(strings.TrimPrefix(line, "## "))
			if !ok {
				return nil, fmt.Errorf("unknown variant heading %q", line)
			}
			if declared[view] {
				return nil, fmt.Errorf("duplicate variant %q", view)
			}
			declared[view] = true
			if view != "common" && !slices.Contains(e.Implementations, view) {
				return nil, fmt.Errorf("undeclared variant %q", view)
			}
		} else if level == 3 {
			if view == "" {
				return nil, fmt.Errorf("section requires a variant heading")
			}
			if err := flush(); err != nil {
				return nil, err
			}
			title = strings.TrimPrefix(line, "### ")
		} else {
			lines = append(lines, line)
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if !declared["common"] || !seen["common/Overview"] || !seen["common/Background"] || !seen["common/Examples"] || !seen["common/Common mistakes"] || !seen["common/Portability"] {
		return nil, fmt.Errorf("Common requires Overview, Background, Examples, Common mistakes and Portability")
	}
	for _, impl := range e.Implementations {
		if !declared[impl] {
			return nil, fmt.Errorf("missing declared variant %q", impl)
		}
		if !seen[impl+"/about"] {
			return nil, fmt.Errorf("variant %q requires an About section", impl)
		}
	}
	return e, nil
}

// Goldmark decides which headings are actual top-level blocks, rather than
// literal text inside code or list containers. Keep the source heading contract.
func sectionHeadings(source []byte) (map[int]int, error) {
	fences := &closedFenceParser{BlockParser: parser.NewFencedCodeBlockParser(), open: map[ast.Node]bool{}}
	// Wrap fences before Goldmark's default fenced-code parser (priority 700).
	md := goldmark.New(goldmark.WithParserOptions(parser.WithBlockParsers(util.Prioritized(fences, 699))))
	doc := md.Parser().Parse(text.NewReader(source))
	if len(fences.open) != 0 {
		return nil, fmt.Errorf("unclosed code fence")
	}
	headings := map[int]int{}
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok || (h.Level != 2 && h.Level != 3) || h.Lines().Len() == 0 {
			continue
		}
		start := h.Lines().At(0).Start
		lineStart := bytes.LastIndexByte(source[:start], '\n') + 1
		if bytes.HasPrefix(source[lineStart:], []byte(strings.Repeat("#", h.Level)+" ")) {
			headings[lineStart] = h.Level
		}
	}
	return headings, nil
}

// CommonMark permits implicit closing at EOF or a container boundary; our
// authoring format requires an explicit closer. Delegate all syntax to Goldmark
// and remember fences whose Continue never consumed a closing delimiter.
type closedFenceParser struct {
	parser.BlockParser
	open map[ast.Node]bool
}

func (p *closedFenceParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	n, state := p.BlockParser.Open(parent, reader, pc)
	if n != nil {
		p.open[n] = true
	}
	return n, state
}
func (p *closedFenceParser) Continue(n ast.Node, reader text.Reader, pc parser.Context) parser.State {
	state := p.BlockParser.Continue(n, reader, pc)
	if state&parser.Close != 0 {
		delete(p.open, n)
	}
	return state
}

func validateMetadata(m Metadata) error {
	if m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema_version %d", m.SchemaVersion)
	}
	if !names.MatchString(m.ID) || m.Title != m.ID || strings.TrimSpace(m.Summary) == "" {
		return fmt.Errorf("invalid id, title or summary")
	}
	for _, s := range []string{m.ID, m.Title, m.Summary} {
		for _, r := range s {
			if unicode.IsControl(r) {
				return fmt.Errorf("metadata must not contain control characters")
			}
		}
	}
	if err := validDate(m.Updated); err != nil {
		return err
	}
	if len(m.Tags) == 0 || len(m.Implementations) == 0 || len(m.References) == 0 {
		return fmt.Errorf("tags, implementations and references are required")
	}
	for i, list := range [][]string{m.Aliases, m.Tags, m.Implementations, m.RelatedGuides} {
		valid := names.MatchString
		if i == 0 {
			valid = IsLookupName
		}
		seen := map[string]bool{}
		for _, n := range list {
			if !valid(n) || seen[n] {
				return fmt.Errorf("invalid or duplicate name %q", n)
			}
			seen[n] = true
		}
	}
	for _, v := range m.Implementations {
		if _, ok := LookupImplementation(v); v == "common" || !ok {
			return fmt.Errorf("unsupported implementation %q", v)
		}
	}
	if slices.Contains(m.Implementations, "upstream") && len(m.Implementations) > 1 {
		return fmt.Errorf("upstream documents a single-implementation tool and cannot be combined with other implementations")
	}
	for _, ref := range m.References {
		if !safeURL(ref) {
			return fmt.Errorf("unsafe reference URL %q", ref)
		}
	}
	return nil
}
func safeURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" && u.User == nil
}

// RenderHTML allows only a small, non-executable CommonMark subset.
func RenderHTML(markdown string) (string, []string, error) {
	source := []byte(markdown)
	// Parse pipe tables only to reject them: the terminal renderer understands
	// tables, so a command page must not contain one that the web would show
	// as plain text.
	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	node := md.Parser().Parse(text.NewReader(source))
	code := []string{}
	err := ast.Walk(node, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Document, *ast.Paragraph, *ast.TextBlock, *ast.Text, *ast.String, *ast.CodeSpan, *ast.Emphasis, *ast.List, *ast.ListItem:
		case *east.Table:
			return ast.WalkStop, fmt.Errorf("tables are not supported in command pages")
		case *ast.Link:
			if !safeURL(string(v.Destination)) {
				return ast.WalkStop, fmt.Errorf("unsafe or unsupported link %q", v.Destination)
			}
		case *ast.FencedCodeBlock:
			lang := string(v.Language(source))
			if lang != "sh" && lang != "bash" && lang != "console" && lang != "text" {
				return ast.WalkStop, fmt.Errorf("unsupported code language %q", lang)
			}
			var b strings.Builder
			for i := 0; i < v.Lines().Len(); i++ {
				seg := v.Lines().At(i)
				b.Write(seg.Value(source))
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

func (c *Catalog) Find(name string) (*Entry, bool) {
	for i := range c.Entries {
		e := &c.Entries[i]
		if e.ID == name || slices.Contains(e.Aliases, name) {
			return e, true
		}
	}
	return nil, false
}
func (c *Catalog) Search(query string) []Entry {
	result := []Entry{}
	terms := strings.Fields(strings.ToLower(query))
	for _, e := range c.Entries {
		hay := strings.ToLower(e.ID + " " + e.Summary + " " + strings.Join(e.Tags, " ") + " " + strings.Join(e.Aliases, " "))
		match := true
		for _, t := range terms {
			if !strings.Contains(hay, t) {
				match = false
				break
			}
		}
		if match {
			result = append(result, e)
		}
	}
	return result
}
func (e *Entry) Select(implementation, option string) ([]Section, error) {
	if implementation != "common" && !slices.Contains(e.Implementations, implementation) {
		return nil, fmt.Errorf("%s has no %s documentation", e.ID, implementation)
	}
	sections := []Section{}
	for _, s := range e.Sections {
		if s.Implementation != "common" && s.Implementation != implementation {
			continue
		}
		if option != "" && !slices.Contains(s.Options, option) {
			continue
		}
		sections = append(sections, s)
	}
	if len(sections) == 0 {
		return nil, fmt.Errorf("no documentation for option %q in %s", option, implementation)
	}
	return sections, nil
}
