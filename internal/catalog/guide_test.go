package catalog

import (
	"encoding/json"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/ninegene/shell-explain/content"
)

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
		"bad date":               strings.Replace(valid, "updated: 2026-11-08", "updated: 2026-13-45", 1),
		"missing date":           strings.Replace(valid, "updated: 2026-11-08\n", "", 1),
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
		if term.ID == "" || term.Markdown == "" || strings.Contains(term.Markdown, "See:") || len(term.Guides) != 1 {
			t.Fatalf("term %+v", term)
		}
	}
	if g.Terms[7].ID != "operating-system" || g.Terms[7].Term != "Operating system" {
		t.Fatalf("anchor id: %+v", g.Terms[7])
	}
	valid := string(raw)
	for name, s := range map[string]string{
		"wrong id":       strings.Replace(valid, "id: glossary", "id: words", 1),
		"bad date":       strings.Replace(valid, "updated: 2026-11-08", "updated: soon", 1),
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
	guide := sampleGuide(t)
	terms, err := content.Files.ReadFile("en/glossary/terms.md")
	if err != nil {
		t.Fatal(err)
	}
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
		"command related guide":   {strings.Replace(string(grep), "updated: 2026-11-08\n", "updated: 2026-11-08\nrelated_guides: [nosuch]\n", 1), withLs, string(terms)},
		"glossary guide missing":  {string(grep), withLs, strings.Replace(string(terms), "See: big-picture", "See: nosuch", 1)},
		"filename mismatch":       {string(grep), strings.Replace(withLs, "id: big-picture", "id: other", 1), string(terms)},
	} {
		if _, err := load(tc.command, tc.guide, tc.glossary); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	// A command page may point at an existing guide.
	linked := strings.Replace(string(grep), "updated: 2026-11-08\n", "updated: 2026-11-08\nrelated_guides: [big-picture]\n", 1)
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
