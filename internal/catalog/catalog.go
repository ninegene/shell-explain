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
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.yaml.in/yaml/v3"
)

const SchemaVersion = 1
const VersionPolicy = "Implementation-level reference; examples are not certified for every installed version."

var names = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var optionName = regexp.MustCompile(`^--?[A-Za-z0-9][A-Za-z0-9-]*$`)
var variants = map[string]string{"Common": "common", "GNU": "gnu", "Apple": "apple", "BusyBox": "busybox", "Bash": "bash"}

type Metadata struct {
	SchemaVersion   int      `yaml:"schema_version" json:"schema_version"`
	ID              string   `yaml:"id" json:"id"`
	Title           string   `yaml:"title" json:"title"`
	Summary         string   `yaml:"summary" json:"summary"`
	Aliases         []string `yaml:"aliases" json:"aliases"`
	Tags            []string `yaml:"tags" json:"tags"`
	Implementations []string `yaml:"implementations" json:"implementations"`
	References      []string `yaml:"references" json:"references"`
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

func (s Section) IsExamples() bool {
	return s.Title == "Examples" || strings.HasSuffix(s.Title, " examples")
}

type Entry struct {
	Metadata
	VersionPolicy string    `json:"version_policy"`
	Sections      []Section `json:"sections"`
}
type Catalog struct {
	SchemaVersion int     `json:"schema_version"`
	Revision      string  `json:"revision"`
	Entries       []Entry `json:"entries"`
}

func Load(files fs.FS) (*Catalog, error) {
	paths, err := fs.Glob(files, "en/*.md")
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no command Markdown files found")
	}
	c := &Catalog{SchemaVersion: SchemaVersion, Entries: []Entry{}}
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
	c.Revision = hex.EncodeToString(hash.Sum(nil))
	return c, nil
}

func Parse(data []byte) (*Entry, error) {
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("content must be UTF-8")
	}
	if bytes.Contains(data, []byte("\r")) {
		return nil, fmt.Errorf("use LF line endings")
	}
	for _, b := range string(data) {
		if unicode.IsControl(b) && b != '\n' && b != '\t' {
			return nil, fmt.Errorf("control characters are not allowed")
		}
	}
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, fmt.Errorf("missing YAML front matter")
	}
	parts := bytes.SplitN(data[4:], []byte("\n---\n"), 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("unclosed front matter")
	}
	e := &Entry{VersionPolicy: VersionPolicy, Sections: []Section{}}
	dec := yaml.NewDecoder(bytes.NewReader(parts[0]))
	dec.KnownFields(true)
	if err := dec.Decode(&e.Metadata); err != nil {
		return nil, err
	}
	if err := validateMetadata(e.Metadata); err != nil {
		return nil, err
	}
	body := string(parts[1])
	headings, err := sectionHeadings(parts[1])
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
			view, ok = variants[strings.TrimPrefix(line, "## ")]
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
	if !declared["common"] || !seen["common/Overview"] || !seen["common/Examples"] || !seen["common/Common mistakes"] || !seen["common/Portability"] {
		return nil, fmt.Errorf("Common requires Overview, Examples, Common mistakes and Portability")
	}
	for _, impl := range e.Implementations {
		if !declared[impl] {
			return nil, fmt.Errorf("missing declared variant %q", impl)
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
	if len(m.Tags) == 0 || len(m.Implementations) == 0 || len(m.References) == 0 {
		return fmt.Errorf("tags, implementations and references are required")
	}
	for _, list := range [][]string{m.Aliases, m.Tags, m.Implementations} {
		seen := map[string]bool{}
		for _, n := range list {
			if !names.MatchString(n) || seen[n] {
				return fmt.Errorf("invalid or duplicate name %q", n)
			}
			seen[n] = true
		}
	}
	for _, v := range m.Implementations {
		if v == "common" || !slices.Contains([]string{"gnu", "apple", "busybox", "bash"}, v) {
			return fmt.Errorf("unsupported implementation %q", v)
		}
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
	md := goldmark.New()
	node := md.Parser().Parse(text.NewReader(source))
	code := []string{}
	err := ast.Walk(node, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Document, *ast.Paragraph, *ast.TextBlock, *ast.Text, *ast.String, *ast.CodeSpan, *ast.Emphasis, *ast.List, *ast.ListItem:
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
