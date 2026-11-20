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
	if !strings.Contains(v, "Matching examples") || !strings.Contains(v, "tar -xf") {
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
	press(m, tea.KeyDown)
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
