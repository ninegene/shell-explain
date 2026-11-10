package catalog

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/ninegene/shell-explain/content"
)

var updatedLine = regexp.MustCompile(`(?m)^updated: .*\n`)

// withUpdated replaces the first updated line of a file, or removes it when
// value is empty. Tests use it so they do not depend on a page's current date.
func withUpdated(s, value string) string {
	line := ""
	if value != "" {
		line = "updated: " + value + "\n"
	}
	loc := updatedLine.FindStringIndex(s)
	if loc == nil {
		return s
	}
	return s[:loc[0]] + line + s[loc[1]:]
}

// afterUpdated inserts extra front matter on the line after updated.
func afterUpdated(s, extra string) string {
	loc := updatedLine.FindStringIndex(s)
	if loc == nil {
		return s
	}
	return s[:loc[1]] + extra + s[loc[1]:]
}

// isolatedGuide removes the sample guide's links to other guides, so a test
// can load it on its own. The shipped guide links to its neighbours.
func isolatedGuide(guide string) string {
	return regexp.MustCompile(`(?m)^related_guides: .*$`).ReplaceAllString(guide, "related_guides: []")
}

// isolatedGlossary points every term at one guide, for the same reason.
func isolatedGlossary(terms, guide string) string {
	return regexp.MustCompile(`(?m)^See: .*$`).ReplaceAllString(terms, "See: "+guide)
}

func sampleGuide(t *testing.T) string {
	t.Helper()
	raw, err := content.Files.ReadFile("en/guides/big-picture.md")
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestSampleGuideParses(t *testing.T) {
	g, err := ParseGuide([]byte(sampleGuide(t)))
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Sections) != len(GuideSections) {
		t.Fatalf("%d sections", len(g.Sections))
	}
	for i, s := range g.Sections {
		if s.Title != GuideSections[i] || s.Implementation != "common" || s.ID != "big-picture-"+string(rune('0'+(i+1)/10))+string(rune('0'+(i+1)%10)) {
			t.Fatalf("section %d: %+v", i, s)
		}
	}
	picture := g.Sections[2]
	if !strings.Contains(picture.HTML, `class="language-diagram"`) || len(picture.Code) != 1 {
		t.Fatalf("diagram not rendered: %s", picture.HTML)
	}
	if table := g.Sections[3]; !strings.Contains(table.HTML, "<table>") || !strings.Contains(table.HTML, "<th>Layer</th>") {
		t.Fatalf("table not rendered: %s", table.HTML)
	}
}

func TestInvalidGuides(t *testing.T) {
	valid := sampleGuide(t)
	wide := "```diagram\n" + strings.Repeat("x", MaxDiagramWidth+1) + "\n```\n\nDiagram: too wide."
	tests := map[string]string{
		"unknown field":          strings.Replace(valid, "level: beginner", "level: beginner\nsurprise: true", 1),
		"bad schema":             strings.Replace(valid, "schema_version: 1", "schema_version: 2", 1),
		"bad id":                 strings.Replace(valid, "id: big-picture", "id: ../big", 1),
		"unknown series":         strings.Replace(valid, "series: how-a-computer-works", "series: nope", 1),
		"zero order":             strings.Replace(valid, "order: 1", "order: 0", 1),
		"bad level":              strings.Replace(valid, "level: beginner", "level: expert", 1),
		"no reading time":        strings.Replace(valid, "reading_minutes: 5", "reading_minutes: 0", 1),
		"bad date":               withUpdated(valid, "2026-13-45"),
		"missing date":           withUpdated(valid, ""),
		"unsafe ref":             strings.Replace(valid, "https://docs.kernel.org/", "javascript:alert(1)", 1),
		"no references":          strings.Replace(valid, "references:\n  - https://www.gnu.org/software/bash/manual/\n  - https://docs.kernel.org/\n", "references: []\n", 1),
		"missing section":        strings.Replace(valid, "## Check yourself", "## Quiz", 1),
		"sections out of order":  strings.Replace(strings.Replace(valid, "## Goal", "## TMP", 1), "## Where next", "## Goal", 1),
		"extra section":          valid + "\n## Bonus\nMore.\n",
		"text before sections":   strings.Replace(valid, "---\n## Goal", "---\nStray text.\n\n## Goal", 1),
		"level one heading":      strings.Replace(valid, "## Goal\n", "## Goal\n# Big\n", 1),
		"level four heading":     strings.Replace(valid, "## Goal\n", "## Goal\n#### Deep\n", 1),
		"HTML":                   strings.Replace(valid, "## Goal\n", "## Goal\n<script>alert(1)</script>\n\n", 1),
		"image":                  strings.Replace(valid, "## Goal\n", "## Goal\n![x](https://example.org/a.png)\n\n", 1),
		"unsafe link":            strings.Replace(valid, "## Goal\n", "## Goal\n[x](javascript:alert)\n\n", 1),
		"unknown fence language": strings.Replace(valid, "```sh\nls\n```", "```python\nls\n```", 1),
		"diagram too wide":       strings.Replace(valid, "## Goal\n", "## Goal\n"+wide+"\n\n", 1),
		"diagram not ascii":      strings.Replace(valid, " YOU                      type:  ls", " YOU   é                  type:  ls", 1),
		"diagram no description": strings.Replace(valid, "Diagram: seven layers", "Seven layers", 1),
		"five column table":      strings.Replace(valid, "| Layer | In a restaurant |\n| --- | --- |", "| a | b | c | d | e |\n| - | - | - | - | - |\n| 1 | 2 | 3 | 4 | 5 |\n\n| Layer | In a restaurant |\n| --- | --- |", 1),
		"table too wide":         strings.Replace(valid, "| You | The customer with an order |", "| You | "+strings.Repeat("w", MaxTableWidth)+" |", 1),
		"empty section":          strings.Replace(valid, "## Goal\nBy the end", "## Goal\n\n## The idea in plain words\nBy the end", 1),
		"control character":      valid + "\x1b[0m",
	}
	for name, s := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseGuide([]byte(s)); err == nil {
				t.Fatal("invalid guide accepted")
			}
		})
	}
}

func TestPlainRendersTablesAndDiagrams(t *testing.T) {
	md := "Intro text.\n\n| Layer | Job |\n| --- | --- |\n| Shell | Starts programs |\n| Kernel | Shares the CPU |\n\n```diagram\n A\n |\n v\n B\n```\n\nDiagram: A then B.\n\n```sh\nls\n```\n"
	got := Plain(md, 80)
	want := "Intro text.\n\nLayer   Job\n------  ---------------\nShell   Starts programs\nKernel  Shares the CPU\n\n A\n |\n v\n B\n\nDiagram: A then B.\n\n    ls\n"
	if got != want {
		t.Fatalf("got:\n%q\nwant:\n%q", got, want)
	}
	// The sample guide renders in every section without losing a diagram line.
	g, err := ParseGuide([]byte(sampleGuide(t)))
	if err != nil {
		t.Fatal(err)
	}
	plain := Plain(g.Sections[2].Markdown, 80)
	for _, line := range strings.Split(strings.TrimSuffix(g.Sections[2].Code[0], "\n"), "\n") {
		if !strings.Contains(plain, "\n"+line+"\n") {
			t.Fatalf("diagram line missing or indented: %q", line)
		}
	}
	for _, line := range strings.Split(Plain(g.Sections[3].Markdown, 80), "\n") {
		if len(line) > MaxTableWidth {
			t.Fatalf("table line too wide: %q", line)
		}
	}
}

func TestHeadingsInsideGuideSections(t *testing.T) {
	md := strings.Replace(sampleGuide(t), "## Goal\n", "## Goal\n### A sub-point\nText under it.\n\n", 1)
	g, err := ParseGuide([]byte(md))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(g.Sections[0].HTML, "<h3>A sub-point</h3>") || !strings.Contains(Plain(g.Sections[0].Markdown, 80), "A sub-point\n") {
		t.Fatalf("level-three heading lost: %s", g.Sections[0].HTML)
	}
}

func TestGlossaryParsing(t *testing.T) {
	raw, err := content.Files.ReadFile("en/glossary/terms.md")
	if err != nil {
		t.Fatal(err)
	}
	g, err := ParseGlossary(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Terms) == 0 {
		t.Fatal("no terms")
	}
	for _, term := range g.Terms {
		if term.ID == "" || term.Markdown == "" || strings.Contains(term.Markdown, "See:") || len(term.Guides) < 1 {
			t.Fatalf("term %+v", term)
		}
	}
	var os GlossaryTerm
	for _, term := range g.Terms {
		if term.ID == "operating-system" {
			os = term
		}
	}
	if os.Term != "Operating system" || len(os.Guides) < 2 {
		t.Fatalf("anchor id or guide links: %+v", os)
	}
	// Terms are listed alphabetically, ignoring case, so readers can scan them.
	for i := 1; i < len(g.Terms); i++ {
		if strings.ToLower(g.Terms[i-1].Term) > strings.ToLower(g.Terms[i].Term) {
			t.Errorf("%q should come after %q", g.Terms[i-1].Term, g.Terms[i].Term)
		}
	}
	valid := string(raw)
	for name, s := range map[string]string{
		"wrong id":       strings.Replace(valid, "id: glossary", "id: words", 1),
		"bad date":       withUpdated(valid, "soon"),
		"duplicate term": valid + "\n## Shell\nAgain.\n",
		"no terms":       strings.SplitN(valid, "\n## ", 2)[0] + "\n",
		"only See":       strings.Replace(valid, "## Command\nA short instruction", "## Command\nSee: big-picture\n\n## Extra\nA short instruction", 1),
		"bad See id":     strings.Replace(valid, "See: big-picture", "See: Big Picture!", 1),
		"HTML":           valid + "\n## Bad\n<b>x</b>\n",
	} {
		if _, err := ParseGlossary([]byte(s)); err == nil {
			t.Errorf("%s: invalid glossary accepted", name)
		}
	}
}

func TestLoadChecksReferencesAndRevision(t *testing.T) {
	grep, err := content.Files.ReadFile("en/grep.md")
	if err != nil {
		t.Fatal(err)
	}
	guide := isolatedGuide(sampleGuide(t))
	rawTerms, err := content.Files.ReadFile("en/glossary/terms.md")
	if err != nil {
		t.Fatal(err)
	}
	terms := []byte(isolatedGlossary(string(rawTerms), "big-picture"))
	load := func(command, guide, glossary string) (*Catalog, error) {
		files := fstest.MapFS{
			"en/grep.md":                {Data: []byte(command)},
			"en/guides/big-picture.md":  {Data: []byte(guide)},
			"en/glossary/terms.md":      {Data: []byte(glossary)},
			"en/guides/not-markdown.go": {Data: []byte("ignored")},
		}
		return Load(files)
	}
	// A guide that names a command with no page needs the page to exist.
	withLs := strings.Replace(guide, "related_commands: [ls, cat]", "related_commands: [grep]", 1)
	base, err := load(string(grep), withLs, string(terms))
	if err != nil {
		t.Fatal(err)
	}
	if len(base.Guides) != 1 || len(base.Glossary.Terms) == 0 || len(base.Series) != len(Series) {
		t.Fatalf("catalog: %d guides, %d terms", len(base.Guides), len(base.Glossary.Terms))
	}
	// Any change to a guide or the glossary changes the content revision.
	edited, err := load(string(grep), strings.Replace(withLs, "By the end", "By the end,", 1), string(terms))
	if err != nil || edited.Revision == base.Revision {
		t.Fatalf("guide edit left the revision unchanged: %v", err)
	}
	edited, err = load(string(grep), withLs, strings.Replace(string(terms), "A short instruction", "A tiny instruction", 1))
	if err != nil || edited.Revision == base.Revision {
		t.Fatalf("glossary edit left the revision unchanged: %v", err)
	}
	for name, tc := range map[string]struct{ command, guide, glossary string }{
		"related command missing": {string(grep), strings.Replace(withLs, "related_commands: [grep]", "related_commands: [nosuch]", 1), string(terms)},
		"related guide missing":   {string(grep), strings.Replace(withLs, "related_guides: []", "related_guides: [nosuch]", 1), string(terms)},
		"related guide is itself": {string(grep), strings.Replace(withLs, "related_guides: []", "related_guides: [big-picture]", 1), string(terms)},
		"command related guide":   {afterUpdated(string(grep), "related_guides: [nosuch]\n"), withLs, string(terms)},
		"glossary guide missing":  {string(grep), withLs, strings.Replace(string(terms), "See: big-picture", "See: nosuch", 1)},
		"filename mismatch":       {string(grep), strings.Replace(withLs, "id: big-picture", "id: other", 1), string(terms)},
	} {
		if _, err := load(tc.command, tc.guide, tc.glossary); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	// A command page may point at an existing guide.
	linked := afterUpdated(string(grep), "related_guides: [big-picture]\n")
	c, err := load(linked, withLs, string(terms))
	if err != nil || len(c.Entries[0].RelatedGuides) != 1 {
		t.Fatalf("command related_guides: %v", err)
	}
	// Two guides cannot share a position in a series.
	two := fstest.MapFS{
		"en/grep.md":       {Data: grep},
		"en/guides/one.md": {Data: []byte(strings.Replace(withLs, "id: big-picture", "id: one", 1))},
		"en/guides/two.md": {Data: []byte(strings.Replace(withLs, "id: big-picture", "id: two", 1))},
	}
	if _, err := Load(two); err == nil {
		t.Error("duplicate series position accepted")
	}
}

func TestExportIncludesGuidesAndGlossary(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"entries", "series", "guides", "glossary", "revision"} {
		if _, ok := out[key]; !ok {
			t.Errorf("export is missing %q", key)
		}
	}
	if len(c.Guides) == 0 || c.Guides[0].ID != "big-picture" || len(c.Glossary.Terms) == 0 {
		t.Fatalf("shipped corpus: %d guides, %d terms", len(c.Guides), len(c.Glossary.Terms))
	}
}

func TestGuideLookupOrderAndSearch(t *testing.T) {
	grep, err := content.Files.ReadFile("en/grep.md")
	if err != nil {
		t.Fatal(err)
	}
	terms, err := content.Files.ReadFile("en/glossary/terms.md")
	if err != nil {
		t.Fatal(err)
	}
	base := strings.Replace(isolatedGuide(sampleGuide(t)), "related_commands: [ls, cat]", "related_commands: [grep]", 1)
	guide := func(id, series string, order int, title string) []byte {
		g := strings.Replace(base, "id: big-picture", "id: "+id, 1)
		g = strings.Replace(g, "series: how-a-computer-works", "series: "+series, 1)
		g = strings.Replace(g, "order: 1", "order: "+string(rune('0'+order)), 1)
		return []byte(strings.Replace(g, "title: From Keypress to Screen", "title: "+title, 1))
	}
	// Files are given out of order: reading order is series, then position.
	c, err := Load(fstest.MapFS{
		"en/grep.md":                {Data: grep},
		"en/glossary/terms.md":      {Data: []byte(isolatedGlossary(string(terms), "alpha"))},
		"en/guides/zeta.md":         {Data: guide("zeta", "terminal-shell-commands", 1, "Zeta Guide")},
		"en/guides/beta.md":         {Data: guide("beta", "how-a-computer-works", 2, "Beta Guide")},
		"en/guides/alpha.md":        {Data: guide("alpha", "how-a-computer-works", 1, "Alpha Guide")},
		"en/glossary/ignored.txt":   {Data: []byte("x")},
		"en/guides/readme-notes.go": {Data: []byte("x")},
	})
	if err != nil {
		t.Fatal(err)
	}
	var order []string
	for _, g := range c.Guides {
		order = append(order, g.ID)
	}
	if strings.Join(order, ",") != "alpha,beta,zeta" {
		t.Fatalf("reading order %v", order)
	}
	for key, want := range map[string]string{"alpha": "alpha", "1": "alpha", "2": "beta", "3": "zeta", "zeta": "zeta"} {
		if g, ok := c.FindGuide(key); !ok || g.ID != want {
			t.Errorf("FindGuide(%q) = %v, %v", key, g, ok)
		}
	}
	for _, key := range []string{"0", "4", "nope", ""} {
		if _, ok := c.FindGuide(key); ok {
			t.Errorf("FindGuide(%q) found a guide", key)
		}
	}
	prev, next := c.GuideNeighbors("alpha")
	if prev != nil || next == nil || next.ID != "beta" {
		t.Errorf("alpha neighbors: %v %v", prev, next)
	}
	prev, next = c.GuideNeighbors("beta")
	if prev == nil || prev.ID != "alpha" || next == nil || next.ID != "zeta" {
		t.Errorf("beta neighbors: %v %v", prev, next)
	}
	prev, next = c.GuideNeighbors("zeta")
	if prev == nil || prev.ID != "beta" || next != nil {
		t.Errorf("zeta neighbors: %v %v", prev, next)
	}
	if p, n := c.GuideNeighbors("nope"); p != nil || n != nil || c.GuideNumber("nope") != 0 || c.GuideNumber("zeta") != 3 {
		t.Error("unknown guide has neighbors or a number")
	}
	if got := c.SearchGuides("alpha guide"); len(got) != 1 || got[0].ID != "alpha" {
		t.Errorf("SearchGuides: %v", got)
	}
	if got := c.SearchGuides("keypress"); len(got) != 0 {
		t.Errorf("the titles were replaced, keypress should not match: %v", got)
	}
	if got := c.SearchGuides("big picture of what"); len(got) != 3 {
		t.Errorf("subtitle search: %d", len(got))
	}
	if got := c.SearchGuides("nonsense words"); got == nil || len(got) != 0 {
		t.Errorf("an empty search must be an empty list, not nil: %#v", got)
	}
}

func TestGuideSectionAndTermLookup(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, _ := c.FindGuide("big-picture")
	for _, name := range []string{"Try it", "try it", "TRY IT", "try-it", "The idea in plain words", "the-idea-in-plain-words"} {
		if _, ok := g.Section(name); !ok {
			t.Errorf("Section(%q) not found", name)
		}
	}
	for _, name := range []string{"", "nope", "try", "Bonus"} {
		if _, ok := g.Section(name); ok {
			t.Errorf("Section(%q) found", name)
		}
	}
	for _, name := range []string{"Kernel", "kernel", "KERNEL", "operating system", "operating-system", "Operating-System"} {
		if term, ok := c.FindTerm(name); !ok || (term.ID != "kernel" && term.ID != "operating-system") {
			t.Errorf("FindTerm(%q) = %v %v", name, term, ok)
		}
	}
	for _, name := range []string{"", "!!!", "kern", "operating"} {
		if _, ok := c.FindTerm(name); ok {
			t.Errorf("FindTerm(%q) found a term", name)
		}
	}
	if got := c.SearchTerms("kernel"); len(got) < 2 {
		t.Errorf("SearchTerms(kernel) = %d, want the kernel and operating system terms", len(got))
	}
	if got := c.SearchTerms("zzzz"); got == nil || len(got) != 0 {
		t.Errorf("empty SearchTerms must be an empty list: %#v", got)
	}
	term, _ := c.FindTerm("kernel")
	if got := term.FirstSentence(); got != "The core of the operating system." {
		t.Errorf("FirstSentence = %q", got)
	}
	single := GlossaryTerm{Markdown: "Only one sentence here\n"}
	if got := single.FirstSentence(); got != "Only one sentence here" {
		t.Errorf("FirstSentence without a full stop = %q", got)
	}
}
