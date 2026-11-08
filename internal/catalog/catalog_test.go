package catalog

import (
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/ninegene/shell-explain/content"
)

func TestSharedCorpus(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	again, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Entries) != 12 || c.Revision != again.Revision || len(c.Revision) != 64 {
		t.Fatal("non-deterministic or incomplete corpus")
	}
	for _, e := range c.Entries {
		for _, v := range append([]string{"common"}, e.Implementations...) {
			sections, err := e.Select(v, "")
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range sections {
				if s.Implementation != "common" && s.Implementation != v {
					t.Fatal("variant leaked")
				}
				plain := Plain(s.Markdown, 40)
				for _, code := range s.Code {
					for _, line := range strings.Split(strings.TrimSuffix(code, "\n"), "\n") {
						if !strings.Contains(plain, "    "+line) {
							t.Fatalf("code changed in %s/%s: %q", e.ID, v, line)
						}
					}
				}
			}
		}
	}
	e, ok := c.Find("ggrep")
	if !ok || e.ID != "grep" {
		t.Fatal("alias not found")
	}
	sections, err := e.Select("gnu", "-n")
	if err != nil || len(sections) != 1 {
		t.Fatal("option lookup must use structured association")
	}
	if _, err = e.Select("apple", "--line-number"); err == nil {
		t.Fatal("GNU option leaked into Apple")
	}
	if _, err = e.Select("freebsd", ""); err == nil {
		t.Fatal("unreviewed implementation accepted")
	}
	if len(c.Search("search logs")) != 1 || len(c.Search("nonexistent")) != 0 {
		t.Fatal("search semantics")
	}
}

func TestInvalidContent(t *testing.T) {
	raw, err := content.Files.ReadFile("en/grep.md")
	if err != nil {
		t.Fatal(err)
	}
	valid := string(raw)
	tests := map[string]string{
		"unknown field":     strings.Replace(valid, "summary:", "surprise: true\nsummary:", 1),
		"duplicate YAML":    strings.Replace(valid, "id: grep", "id: grep\nid: grep", 1),
		"bad schema":        strings.Replace(valid, "schema_version: 1", "schema_version: 2", 1),
		"bad id":            strings.Replace(valid, "id: grep", "id: ../grep", 1),
		"empty tags":        strings.Replace(valid, "tags: [text, search, logs]", "tags: []", 1),
		"unsafe ref":        strings.Replace(valid, "https://www.gnu.org/software/grep/manual/", "javascript:alert(1)", 1),
		"unknown variant":   strings.Replace(valid, "## GNU", "## Other", 1),
		"missing coverage":  strings.Replace(valid, "implementations: [gnu, apple, busybox]", "implementations: [gnu, apple, busybox, bash]", 1),
		"duplicate variant": valid + "\n## GNU\n### Again\nText.\n",
		"duplicate section": valid + "\n### BusyBox notes\nText.\n",
		"HTML":              strings.Replace(valid, "Use `grep`", "<script>alert(1)</script>\n\nUse `grep`", 1),
		"unsafe link":       strings.Replace(valid, "Use `grep`", "[click](javascript:alert) Use `grep`", 1),
		"image":             strings.Replace(valid, "Use `grep`", "![picture](https://example.org/a.png) Use `grep`", 1),
		"unclosed fence":    valid + "\n```sh\necho hello\n",
		"escape":            valid + "\x1b[0m",
		"empty section":     valid + "\n### Empty\n",
	}
	for name, s := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse([]byte(s)); err == nil {
				t.Fatal("invalid content accepted")
			}
		})
	}
	sed, err := content.Files.ReadFile("en/sed.md")
	if err != nil {
		t.Fatal(err)
	}
	duplicate := fstest.MapFS{"en/grep.md": {Data: raw}, "en/sed.md": {Data: []byte(strings.Replace(string(sed), "aliases: [gsed]", "aliases: [ggrep]", 1))}}
	if _, err := Load(duplicate); err == nil {
		t.Fatal("duplicate alias across valid pages accepted")
	}
}

func TestMarkdownSubset(t *testing.T) {
	md := "First **important** and `literal` with [manual](https://example.org/manual).\n\n- One item\n- Two items\n\n1. First\n2. Second\n\n```sh\nprintf '%s\\n' 'a < b & c'\n```\n"
	html, code, err := RenderHTML(md)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "a &lt; b &amp; c") || len(code) != 1 {
		t.Fatal("HTML escaping or code extraction failed")
	}
	plain := Plain(md, 30)
	for _, want := range []string{"important", "literal", "https://example.org/manual", "- One item", "1. First", "    printf '%s\\n' 'a < b & c'"} {
		if !strings.Contains(plain, want) {
			t.Errorf("missing %q in %s", want, plain)
		}
	}
}

func TestPlainWrapping(t *testing.T) {
	text := Plain("A paragraph with enough words to wrap across several lines.\n\n```sh\nprintf '%s\\n' 'a very long command that must stay exactly as authored'\n```\n", 25)
	for _, line := range strings.Split(text, "\n") {
		if !strings.HasPrefix(line, "    ") && len(line) > 25 {
			t.Errorf("prose did not wrap: %q", line)
		}
	}
}

func TestWorkingTreeSources(t *testing.T) {
	if _, err := Load(os.DirFS("../../content")); err != nil {
		t.Fatal(err)
	}
}

func TestFenceDelimitersPreserveLiteralHeadings(t *testing.T) {
	raw, err := content.Files.ReadFile("en/grep.md")
	if err != nil {
		t.Fatal(err)
	}
	original, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ name, open, close, code string }{
		{"shorter backticks", "````sh", "````", "cat <<'EOF'\n```\n### literal text\n## GNU\nEOF\n"},
		{"shorter tildes", "~~~~sh", "~~~~", "cat <<'EOF'\n~~~\n### literal text\nEOF\n"},
		{"longer closer", "````sh", "`````", "echo hello\n"},
		{"closer with padding", "````sh", "```` \t", "echo hello\n"},
		{"trailing content is not a closer", "````sh", "````", "cat <<'EOF'\n````still-code\n### literal text\nEOF\n"},
		{"different character is not a closer", "````sh", "````", "cat <<'EOF'\n~~~~\n### literal text\nEOF\n"},
		{"four spaces are literal", "```sh", "```", "cat <<'EOF'\n    ```\n### literal text\n## GNU\nEOF\n"},
		{"tab is literal", "```sh", "```", "cat <<'EOF'\n\t```\n### literal text\nEOF\n"},
		{"tab after spaces is literal", "```sh", "```", "cat <<'EOF'\n  \t```\n### literal text\nEOF\n"},
		{"indented tildes are literal", "~~~sh", "~~~", "cat <<'EOF'\n    ~~~\n### literal text\nEOF\n"},
		{"one space closer", "```sh", " ```", "echo hello\n"},
		{"three space closer", "```sh", "   ```", "echo hello\n"},
		{"three space opener", "   ```sh", "```", "echo hello\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			markdown := strings.Replace(string(raw), "```sh\ngrep 'ERROR' app.log\n```", tt.open+"\n"+tt.code+tt.close, 1)
			entry, err := Parse([]byte(markdown))
			if err != nil {
				t.Fatal(err)
			}
			if len(entry.Sections) != len(original.Sections) {
				t.Fatal("literal headings changed the section structure")
			}
			if got := entry.Sections[1].Code[0]; got != tt.code {
				t.Fatalf("code changed: got %q, want %q", got, tt.code)
			}
		})
	}
	for _, closer := range []string{"```", "~~~~", "```` trailing-text", "    ````", "\t````"} {
		markdown := strings.Replace(string(raw), "```sh\ngrep 'ERROR' app.log\n```", "````sh\necho hello\n"+closer, 1)
		if _, err := Parse([]byte(markdown)); err == nil {
			t.Fatalf("accepted invalid closer %q", closer)
		}
	}
}

func TestFencesInListContainers(t *testing.T) {
	raw, err := content.Files.ReadFile("en/grep.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, block := range []string{
		"- Example:\n\n  ```sh\n  echo hello\n  ### literal text\n  ```",
		"- ```sh\n  echo hello\n  ### literal text\n  ```",
		"- Nested:\n  - Example:\n\n    ```sh\n    echo hello\n    ### literal text\n    ```",
	} {
		markdown := strings.Replace(string(raw), "```sh\ngrep 'ERROR' app.log\n```", block, 1)
		e, err := Parse([]byte(markdown))
		if err != nil {
			t.Fatalf("list fence rejected: %v\n%s", err, block)
		}
		if got := e.Sections[1].Code[0]; got != "echo hello\n### literal text\n" {
			t.Fatalf("list code changed: %q", got)
		}
	}
	block := "- ```sh\n  echo hello\n\nOutside the list."
	markdown := strings.Replace(string(raw), "```sh\ngrep 'ERROR' app.log\n```", block, 1)
	if _, err := Parse([]byte(markdown)); err == nil || !strings.Contains(err.Error(), "unclosed code fence") {
		t.Fatalf("missing list closer was not rejected: %v", err)
	}
}

func TestSectionPreservesFenceIndentation(t *testing.T) {
	raw, err := content.Files.ReadFile("en/grep.md")
	if err != nil {
		t.Fatal(err)
	}
	block := "   ```sh\n    echo hello\n   ```"
	markdown := strings.Replace(string(raw), "```sh\ngrep 'ERROR' app.log\n```", block, 1)
	e, err := Parse([]byte(markdown))
	if err != nil {
		t.Fatal(err)
	}
	// CommonMark removes the opener's three-space indent, leaving one space.
	if got := e.Sections[1].Code[0]; got != " echo hello\n" {
		t.Fatalf("section trimming changed code indentation: %q", got)
	}
}
