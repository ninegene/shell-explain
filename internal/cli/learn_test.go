package cli

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestLearnListAndReading(t *testing.T) {
	status, list, stderr := run(t, "learn")
	if status != 0 || stderr != "" {
		t.Fatalf("learn: %d %s", status, stderr)
	}
	for _, want := range []string{"How a Computer Works, Just Enough", "  1  From Keypress to Screen", "  2  The Computer's Body: CPU, Memory, Disk", "  5  Programs That Are Awake: Processes", "5 min", "explain learn big-picture", "explain learn processes", "explain glossary WORD"} {
		if !strings.Contains(list, want) {
			t.Errorf("list is missing %q:\n%s", want, list)
		}
	}
	// A series with no guides yet is not listed.
	if strings.Contains(list, "Getting Things Done") {
		t.Error("empty series listed")
	}
	_, byName, _ := run(t, "learn", "big-picture", "--width=70")
	_, byNumber, _ := run(t, "learn", "1", "--width=70")
	if byName != byNumber || !strings.HasPrefix(byName, "From Keypress to Screen\n") {
		t.Fatalf("name and number differ or the header is wrong:\n%s", byName)
	}
	for _, want := range []string{"Series: How a Computer Works, Just Enough", "Guide 1 of 5 · about 5 minutes", "Next:      The Computer's Body: CPU, Memory, Disk\n           explain learn computer-body\n", "Related:   The Computer's Body: CPU, Memory, Disk\n           explain learn computer-body\n", "           explain learn who-is-boss\n", "Goal\n", "A picture\n", "Try it\n", "Check yourself\n", "Where next\n", "Keep going\n", "Commands:  explain ls, explain cat", "References\nhttps://www.gnu.org/software/bash/manual/\nhttps://docs.kernel.org/\n"} {
		if !strings.Contains(byName, want) {
			t.Errorf("guide is missing %q", want)
		}
	}
	// Diagrams and tables are never wrapped or indented; prose obeys the width.
	if !strings.Contains(byName, "\n YOU                      type:  ls\n   |  keys\n") {
		t.Error("diagram changed")
	}
	if !strings.Contains(byName, "\nLayer     In a restaurant\n--------  ----") {
		t.Error("table changed")
	}
	for _, line := range strings.Split(byName, "\n") {
		if strings.HasPrefix(line, " ") || strings.Contains(line, "  ") || strings.HasPrefix(line, "Layer") || strings.HasPrefix(line, "-") {
			continue // diagram, table or list continuation lines
		}
		if utf8.RuneCountInString(line) > 70 && !strings.HasPrefix(line, "https://") && !strings.HasPrefix(line, "Next") {
			t.Errorf("prose line wider than 70: %q", line)
		}
	}
}

// Every guide, with its header and footer, fits an 80-column terminal.
func TestEveryGuideFitsEightyColumns(t *testing.T) {
	_, list, _ := run(t, "learn", "--format=json")
	var ids struct {
		Series []struct{ Guides []struct{ ID string } }
	}
	if err := json.Unmarshal([]byte(list), &ids); err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, s := range ids.Series {
		for _, g := range s.Guides {
			_, text, stderr := run(t, "learn", g.ID, "--width=80")
			if stderr != "" || text == "" {
				t.Fatalf("%s: %q", g.ID, stderr)
			}
			for _, line := range strings.Split(text, "\n") {
				if n := utf8.RuneCountInString(line); n > 80 && !strings.HasPrefix(line, "https://") {
					t.Errorf("%s: line of %d characters: %q", g.ID, n, line)
				}
			}
			checked++
		}
	}
	if checked < 5 {
		t.Errorf("only %d guides checked", checked)
	}
}

func TestLearnSectionAndJSON(t *testing.T) {
	_, section, _ := run(t, "learn", "big-picture", "--section", "try it")
	if !strings.HasPrefix(section, "Try it\n") || strings.Contains(section, "Goal") || strings.Contains(section, "Keep going") || strings.Contains(section, "References") {
		t.Fatalf("section output:\n%s", section)
	}
	_, slug, _ := run(t, "learn", "big-picture", "--section=try-it")
	if slug != section {
		t.Fatal("section slug and title differ")
	}
	var guide struct {
		SchemaVersion   int    `json:"schema_version"`
		ContentRevision string `json:"content_revision"`
		Guide           struct {
			ID       string `json:"id"`
			Sections []struct {
				Title string `json:"title"`
				HTML  string `json:"html"`
			} `json:"sections"`
		} `json:"guide"`
	}
	_, data, _ := run(t, "learn", "big-picture", "--format=json")
	if err := json.Unmarshal([]byte(data), &guide); err != nil || guide.Guide.ID != "big-picture" || len(guide.Guide.Sections) != 9 || guide.ContentRevision == "" || strings.Contains(data, "\x1b") {
		t.Fatalf("guide json: %v %s", err, data[:min(len(data), 200)])
	}
	_, data, _ = run(t, "learn", "big-picture", "--format=json", "--section=Goal")
	if err := json.Unmarshal([]byte(data), &guide); err != nil || len(guide.Guide.Sections) != 1 || guide.Guide.Sections[0].Title != "Goal" {
		t.Fatalf("section json: %v", err)
	}
	var list struct {
		Series []struct {
			ID     string `json:"id"`
			Guides []struct {
				ID     string `json:"id"`
				Number int    `json:"number"`
			} `json:"guides"`
		} `json:"series"`
	}
	_, data, _ = run(t, "learn", "--format=json")
	if err := json.Unmarshal([]byte(data), &list); err != nil || len(list.Series) != 4 || list.Series[0].Guides[0].ID != "big-picture" || list.Series[0].Guides[0].Number != 1 || len(list.Series[0].Guides) != 5 || list.Series[0].Guides[4].ID != "processes" || list.Series[0].Guides[4].Number != 5 || len(list.Series[3].Guides) != 0 {
		t.Fatalf("list json: %v %s", err, data)
	}
}

func TestLearnAndGlossaryErrors(t *testing.T) {
	for _, args := range [][]string{
		{"learn", "nosuch"}, {"learn", "9"}, {"learn", "0"},
		{"learn", "big-picture", "extra"},
		{"learn", "big-picture", "--section", "nope"},
		{"learn", "--section", "Goal"},
		{"learn", "--implementation=gnu"}, {"learn", "--details"}, {"learn", "big-picture", "--option=-n"}, {"learn", "--platform=linux"},
		{"glossary", "nosuch"}, {"glossary", "a", "b"}, {"glossary", "--section=Goal"}, {"glossary", "--details"},
	} {
		status, out, stderr := run(t, args...)
		if status == 0 || stderr == "" || out != "" {
			t.Errorf("%v: status %d, stdout %q, stderr %q", args, status, out, stderr)
		}
	}
	status, out, _ := run(t, "learn", "nosuch", "--format=json")
	if status == 0 || !strings.Contains(out, `"code": "invalid_request"`) {
		t.Errorf("json error: %d %s", status, out)
	}
}

func TestGlossary(t *testing.T) {
	_, list, stderr := run(t, "glossary", "--width=60")
	if stderr != "" || !strings.HasPrefix(list, "Words You Will See\n") || !strings.Contains(list, "Kernel  ") || !strings.Contains(list, "explain glossary WORD") {
		t.Fatalf("glossary list (%s):\n%s", stderr, list)
	}
	_, term, _ := run(t, "glossary", "kernel")
	_, same, _ := run(t, "glossary", "KERNEL")
	_, anchor, _ := run(t, "glossary", "operating-system")
	_, spaced, _ := run(t, "glossary", "Operating system")
	if term != same || anchor != spaced || !strings.HasPrefix(term, "Kernel\nThe core of the operating system.") || !strings.Contains(term, "Learn more:\n  explain learn big-picture  (From Keypress to Screen)\n  explain learn what-is-an-os  (What an Operating System Really Does)\n  explain learn who-is-boss  (") || strings.Contains(term, "See:") {
		t.Fatalf("term output:\n%s", term)
	}
	var one struct {
		Term struct {
			ID     string   `json:"id"`
			Guides []string `json:"guides"`
		} `json:"term"`
	}
	_, data, _ := run(t, "glossary", "kernel", "--format=json")
	if err := json.Unmarshal([]byte(data), &one); err != nil || one.Term.ID != "kernel" || !slices.Equal(one.Term.Guides, []string{"big-picture", "what-is-an-os", "who-is-boss"}) {
		t.Fatalf("term json: %v %s", err, data)
	}
	var all struct {
		Glossary struct {
			Terms []struct{ ID string } `json:"terms"`
		} `json:"glossary"`
	}
	_, data, _ = run(t, "glossary", "--format=json")
	if err := json.Unmarshal([]byte(data), &all); err != nil || len(all.Glossary.Terms) != 21 {
		t.Fatalf("glossary json: %v", err)
	}
	_, _, stderr = run(t, "glossary", "kerne")
	if !strings.Contains(stderr, "similar: Kernel;") {
		t.Errorf("name matches should be the only suggestions: %s", stderr)
	}
	// With no name match, fall back to definitions that mention the word.
	_, _, stderr = run(t, "glossary", "processor")
	if !strings.Contains(stderr, "similar: CPU,") {
		t.Errorf("no definition-based suggestion: %s", stderr)
	}
}

func TestWidthZeroTurnsWrappingOff(t *testing.T) {
	_, guide, _ := run(t, "learn", "big-picture", "--width=0")
	if !strings.Contains(guide, "By the end you can name the layers between your fingers and the text on your screen, and say what each layer does.") {
		t.Error("guide prose wrapped with --width=0")
	}
	_, list, _ := run(t, "glossary", "--width=0")
	if !strings.Contains(list, "Disk              Long-term storage that keeps your files when the computer is turned off, such as a hard disk or a solid-state drive (SSD).\n") {
		t.Errorf("glossary list wrapped with --width=0:\n%s", list)
	}
}

func TestSearchFindsGuidesAndTerms(t *testing.T) {
	_, out, _ := run(t, "search", "keypress")
	if !strings.Contains(out, "guide    explain learn big-picture  (From Keypress to Screen)") {
		t.Errorf("guide not found:\n%s", out)
	}
	_, out, _ = run(t, "search", "kernel")
	if !strings.Contains(out, `word     explain glossary "Kernel"`) {
		t.Errorf("term not found:\n%s", out)
	}
	_, out, _ = run(t, "search", "logs")
	if !strings.Contains(out, "\ngrep     Find lines of text that match a pattern.\n") || strings.Contains(out, "guide ") || strings.Contains(out, "word ") {
		t.Errorf("command search changed:\n%s", out)
	}
	_, out, _ = run(t, "search", "zzzznothing")
	if !strings.Contains(out, "Nothing matched") {
		t.Errorf("empty result: %s", out)
	}
	var result struct {
		Entries []json.RawMessage `json:"entries"`
		Guides  []struct{ ID string }
		Terms   []struct{ ID string }
	}
	_, data, _ := run(t, "search", "kernel", "--format=json")
	if err := json.Unmarshal([]byte(data), &result); err != nil || result.Entries == nil || len(result.Terms) == 0 || result.Guides == nil {
		t.Fatalf("search json: %v %s", err, data)
	}
	// list is unchanged: commands only.
	_, data, _ = run(t, "list", "--format=json")
	if strings.Contains(data, `"guides"`) || strings.Contains(data, `"terms"`) {
		t.Error("list now includes guides or terms")
	}
}

func TestDetailsPointToRelatedGuides(t *testing.T) {
	_, ls, _ := run(t, "ls", "--implementation=gnu", "--details")
	if !strings.Contains(ls, "Learn more\nexplain learn big-picture  (From Keypress to Screen)\n\nReferences\n") {
		t.Errorf("ls details:\n%s", ls[max(0, len(ls)-400):])
	}
	_, grep, _ := run(t, "grep", "--implementation=gnu", "--details")
	if strings.Contains(grep, "Learn more") {
		t.Error("grep points to a guide it does not list")
	}
	_, short, _ := run(t, "ls", "--implementation=gnu")
	if strings.Contains(short, "Learn more") {
		t.Error("default output shows Learn more")
	}
}

func TestWantsPager(t *testing.T) {
	for args, want := range map[string]bool{
		"learn big-picture":               true,
		"learn 1":                         true,
		"learn":                           false,
		"learn big-picture --no-pager":    false,
		"learn big-picture --format json": false,
		"glossary":                        true,
		"glossary kernel":                 false,
		"glossary --no-pager":             false,
		"grep":                            false,
		"list":                            false,
		"--help":                          false,
		"learn big-picture --bogus":       false,
	} {
		if got := WantsPager(strings.Fields(args)); got != want {
			t.Errorf("WantsPager(%q) = %v, want %v", args, got, want)
		}
	}
}

func TestStartPager(t *testing.T) {
	t.Setenv("LESS", "")
	pipe := func() (*os.File, *os.File) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		return r, w
	}
	// Not a terminal: output goes straight to stdout.
	r, w := pipe()
	defer r.Close()
	if out, finish := StartPager(w, false); out != w {
		t.Error("a pipe got a pager")
	} else {
		finish()
	}
	// A pager that is not installed falls back to direct output.
	t.Setenv("PAGER", "definitely-not-a-pager-xyz")
	if out, _ := StartPager(w, true); out != w {
		t.Error("a missing pager was used")
	}
	w.Close()
	// A real pager receives the text and writes it to stdout.
	r, w = pipe()
	defer r.Close()
	t.Setenv("PAGER", "cat")
	out, finish := StartPager(w, true)
	if out == w {
		t.Fatal("cat was not started")
	}
	if _, err := out.Write([]byte("through the pager\n")); err != nil {
		t.Fatal(err)
	}
	finish()
	w.Close()
	buf := make([]byte, 64)
	n, _ := r.Read(buf)
	if string(buf[:n]) != "through the pager\n" {
		t.Fatalf("pager output %q", buf[:n])
	}
}
