package tui

import (
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/detect"
)

func newModel(t *testing.T, copied *[]string) *Model {
	t.Helper()
	c, err := catalog.Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	m := New(Env{
		Catalog: c,
		Select: func(e *catalog.Entry) (detect.Selection, error) {
			return detect.Selection{Implementation: "common"}, nil
		},
		Copy: func(text string) tea.Cmd {
			*copied = append(*copied, text)
			return nil
		},
	})
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	return m
}

func typeText(m *Model, s string) {
	for _, r := range s {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
}

func press(m *Model, code rune) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: code})
	return cmd
}

func TestTypingFiltersLikeSearch(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	if len(m.Entries()) != len(m.env.Catalog.Entries) {
		t.Fatalf("empty query must list every command, got %d", len(m.Entries()))
	}
	typeText(m, "extract archive")
	want := map[string]bool{}
	for _, e := range m.env.Catalog.Search("extract archive") {
		want[e.ID] = true
	}
	for _, h := range m.env.Catalog.SearchExamples("extract archive", 0) {
		want[h.Command] = true
	}
	got := m.Entries()
	if len(got) != len(want) || !slices.Contains(got, "tar") {
		t.Fatalf("list %v does not match search results %v", got, want)
	}
	for _, id := range got {
		if !want[id] {
			t.Errorf("%s listed but not found by search", id)
		}
	}
}

func TestPreviewShowsMatchingExamplesFirst(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "extract archive")
	for m.Selected().ID != "tar" {
		press(m, tea.KeyDown)
	}
	v := m.View().Content
	if !strings.Contains(v, "Examples 1/") || !strings.Contains(v, "tar -xf") {
		t.Fatalf("matching example missing from preview:\n%s", v)
	}
}

func TestCopyReceivesExactExampleCode(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "tar")
	for m.Selected().ID != "tar" {
		press(m, tea.KeyDown)
	}
	var tarHits []catalog.ExampleHit
	for _, h := range m.env.Catalog.SearchExamples("tar", 0) {
		if h.Command == "tar" {
			tarHits = append(tarHits, h)
		}
	}
	if len(tarHits) < 2 {
		t.Fatal("need two tar examples")
	}
	m.Update(tea.KeyPressMsg{Code: 'y', Mod: tea.ModCtrl})
	if len(copied) != 1 || copied[0] != tarHits[0].Code {
		t.Fatalf("copied %q, want %q", copied, tarHits[0].Code)
	}
	// Preview focus: the next example can be selected and copied with y.
	press(m, tea.KeyTab)
	press(m, tea.KeyRight)
	m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	if len(copied) != 2 || copied[1] != tarHits[1].Code {
		t.Fatalf("copied %q, want the second example %q", copied, tarHits[1].Code)
	}
}

func TestEscapeAndQuit(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "tar")
	if cmd := press(m, tea.KeyEscape); cmd != nil && isQuit(cmd) {
		t.Fatal("esc with a query must clear it, not quit")
	}
	if len(m.Entries()) != len(m.env.Catalog.Entries) {
		t.Fatal("esc did not clear the query")
	}
	if cmd := press(m, tea.KeyEscape); !isQuit(cmd) {
		t.Fatal("esc with an empty query must quit")
	}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if !isQuit(cmd) {
		t.Fatal("ctrl+c must quit")
	}
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func TestDetailsToggleAndEmptyResult(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "tar")
	before := m.preview.GetContent()
	press(m, tea.KeyEnter)
	if m.preview.GetContent() == before || !strings.Contains(m.preview.GetContent(), "References") {
		t.Fatal("enter did not switch to the full reference")
	}
	typeText(m, "zzzqqq")
	if len(m.Entries()) != 0 || !strings.Contains(m.View().Content, "Nothing matched") {
		t.Fatal("empty result not reported")
	}
	m.Update(tea.KeyPressMsg{Code: 'y', Mod: tea.ModCtrl})
	if len(copied) != 0 {
		t.Fatal("copied with nothing selected")
	}
}

func TestNarrowWindowDoesNotPanic(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	m.Update(tea.WindowSizeMsg{Width: 40, Height: 8})
	typeText(m, "grep")
	if m.View().Content == "" {
		t.Fatal("empty view")
	}
}

func TestPreviewArrowsStepExamplesAndFollow(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "ssh")
	for len(m.examples) < 2 && m.cursor < len(m.entries)-1 {
		press(m, tea.KeyDown) // find a command with several matching examples
	}
	press(m, tea.KeyTab)
	if len(m.examples) < 2 {
		t.Fatalf("need two matching examples, got %d", len(m.examples))
	}
	view := m.View().Content
	if !strings.Contains(view, "▸ ") || !strings.Contains(view, "Examples 1/") {
		t.Fatalf("first example should be marked:\n%s", view)
	}
	press(m, tea.KeyDown)
	view = m.View().Content
	if !strings.Contains(view, "Examples 2/") {
		t.Fatalf("footer should show the second example:\n%s", view)
	}
	for i, l := range strings.Split(view, "\n") {
		_, right, _ := strings.Cut(l, "│")
		if strings.Contains(right, "▸ ") && !strings.Contains(right, prefix(firstLine(m.examples[1].Code), 20)) && !strings.Contains(right, prefix(m.examples[1].Description, 15)) {
			t.Fatalf("marker line %d is not the second example: %q", i, l)
		}
	}
	// stepping far enough scrolls the selected example into view
	for range m.examples {
		press(m, tea.KeyDown)
	}
	if _, right, _ := strings.Cut(m.View().Content, "│"); !strings.Contains(right, "▸ ") {
		t.Fatalf("selected example must stay visible:\n%s", m.View().Content)
	}
}

func TestJKScrollsPreview(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "tar")
	press(m, tea.KeyEnter) // full reference is taller than the pane
	press(m, tea.KeyTab)
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if got := m.preview.YOffset(); got != 2 {
		t.Fatalf("j twice should scroll two lines, offset %d", got)
	}
	m.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	if got := m.preview.YOffset(); got != 1 {
		t.Fatalf("k should scroll back a line, offset %d", got)
	}
}

func TestMatchesOnlyUntilDetails(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "ssh")
	for m.Selected().ID != "rmt" { // few matches, so the hint line is on screen
		press(m, tea.KeyDown)
	}
	e := m.Selected()
	heading := e.ID + " — "
	view := m.View().Content
	if !strings.Contains(view, "Examples 1/") || strings.Contains(view, heading) || !strings.Contains(view, "enter: the full reference") {
		t.Fatalf("a search with matches should show only the matches:\n%s", view)
	}
	press(m, tea.KeyEnter)
	if !strings.Contains(m.View().Content, "Examples 1/") {
		t.Fatal("matches should stay on top of the full reference")
	}
}

func TestCommandNameMatchesComeFirst(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	typeText(m, "ssh")
	ids := m.Entries()
	if len(ids) < 2 || ids[0] != "ssh" {
		t.Fatalf("an exact command name must be first: %v", ids)
	}
	if m.Selected().ID != "ssh" {
		t.Fatalf("the best name match must be highlighted, got %s", m.Selected().ID)
	}
	prefix := 1
	for prefix < len(ids) && strings.HasPrefix(ids[prefix], "ssh") {
		prefix++
	}
	for _, id := range ids[prefix:] {
		if strings.HasPrefix(id, "ssh") {
			t.Fatalf("names starting with the query must precede the rest: %v", ids)
		}
	}
}

func prefix(s string, n int) string { return s[:min(n, len(s))] }

func TestHeaderShowsDetailsMode(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	if strings.Contains(m.View().Content, "Details") {
		t.Fatal("details mode is off at start")
	}
	press(m, tea.KeyEnter)
	v := m.View().Content
	if !strings.Contains(v, "Details · Examples") || !strings.Contains(v, "enter hide details") {
		t.Fatalf("details mode should be visible:\n%s", v)
	}
	press(m, tea.KeyEnter)
	if strings.Contains(m.View().Content, "Details") {
		t.Fatal("details mode should turn off again")
	}
}

func TestFocusedPaneTitleIsMarked(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	header := func() string { return strings.Split(m.View().Content, "\n")[1] }
	left := header()
	press(m, tea.KeyTab)
	if right := header(); left == right {
		t.Fatal("moving focus must change the header styling")
	}
}

func TestHeaderStaysAlignedInNarrowWindow(t *testing.T) {
	var copied []string
	m := newModel(t, &copied)
	m.Update(tea.WindowSizeMsg{Width: 36, Height: 20})
	for range 100 {
		press(m, tea.KeyDown)
	}
	rows := strings.Split(m.View().Content, "\n")
	col := func(l string) int {
		plain := ansiRE.ReplaceAllString(l, "")
		return len([]rune(plain[:strings.Index(plain, "│")]))
	}
	if col(rows[1]) != col(rows[2]) {
		t.Fatalf("header separator at %d, list separator at %d", col(rows[1]), col(rows[2]))
	}
}
