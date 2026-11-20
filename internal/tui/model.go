// Package tui implements `explain browse`, an interactive terminal browser over
// the same catalog and renderer as the command line. It never executes the
// examples it shows.
package tui

import (
	"fmt"
	"slices"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/cli"
	"github.com/ninegene/shell-explain/internal/detect"
)

// Env is what the browser needs from the rest of the program.
type Env struct {
	Catalog *catalog.Catalog
	// Select chooses the documentation for a command; it must not start processes.
	Select func(e *catalog.Entry) (detect.Selection, error)
	// Copy returns the command that puts text on the clipboard.
	Copy  func(text string) tea.Cmd
	Color bool
	Query string
}

type pane int

const (
	paneInput pane = iota
	panePreview
)

// Model is the Bubble Tea model for the browser.
type Model struct {
	env     Env
	input   textinput.Model
	preview viewport.Model

	focus    pane
	details  bool
	width    int
	height   int
	entries  []*catalog.Entry
	cursor   int
	hits     map[string][]catalog.ExampleHit
	sels     map[string]detect.Selection
	examples []catalog.Example
	exIdx    int
	status   string
}

// New returns a browser listing every command, filtered by env.Query.
func New(env Env) *Model {
	if env.Copy == nil {
		env.Copy = CopyToClipboard
	}
	in := textinput.New()
	in.Prompt = "search: "
	in.Placeholder = "a command, or what you want to do"
	in.SetValue(env.Query)
	in.CursorEnd()
	_ = in.Focus()
	vp := viewport.New()
	vp.SoftWrap = true
	m := &Model{env: env, input: in, preview: vp, sels: map[string]detect.Selection{}, width: 80, height: 24}
	m.layout()
	m.refilter()
	return m
}

// Selected returns the highlighted command, if any.
func (m *Model) Selected() *catalog.Entry {
	if m.cursor < 0 || m.cursor >= len(m.entries) {
		return nil
	}
	return m.entries[m.cursor]
}

// Entries returns the commands in the list, in order.
func (m *Model) Entries() []string {
	ids := make([]string, len(m.entries))
	for i, e := range m.entries {
		ids[i] = e.ID
	}
	return ids
}

func (m *Model) Init() tea.Cmd { return textinput.Blink }

func (m *Model) listWidth() int {
	w := m.width / 3
	return min(max(w, 16), 30)
}

func (m *Model) layout() {
	m.input.SetWidth(max(m.width-len(m.input.Prompt)-1, 10))
	m.preview.SetWidth(max(m.width-m.listWidth()-3, 10))
	m.preview.SetHeight(max(m.height-2, 3))
}

// refilter lists the commands matching the query the way `explain search`
// does: by command metadata, or because one of their examples matches.
func (m *Model) refilter() {
	query := strings.TrimSpace(m.input.Value())
	keep := m.Selected()
	m.hits = map[string][]catalog.ExampleHit{}
	m.entries = m.entries[:0]
	if query == "" {
		for i := range m.env.Catalog.Entries {
			m.entries = append(m.entries, &m.env.Catalog.Entries[i])
		}
	} else {
		match := map[string]bool{}
		for _, e := range m.env.Catalog.Search(query) {
			match[e.ID] = true
		}
		for _, h := range m.env.Catalog.SearchExamples(query, 0) {
			match[h.Command] = true
			m.hits[h.Command] = append(m.hits[h.Command], h)
		}
		for i := range m.env.Catalog.Entries {
			if e := &m.env.Catalog.Entries[i]; match[e.ID] {
				m.entries = append(m.entries, e)
			}
		}
	}
	m.cursor = 0
	if keep != nil {
		if i := slices.IndexFunc(m.entries, func(e *catalog.Entry) bool { return e.ID == keep.ID }); i >= 0 {
			m.cursor = i
		}
	}
	m.refreshPreview()
}

func (m *Model) selection(e *catalog.Entry) (detect.Selection, error) {
	if s, ok := m.sels[e.ID]; ok {
		return s, nil
	}
	s, err := m.env.Select(e)
	if err == nil {
		m.sels[e.ID] = s
	}
	return s, err
}

// refreshPreview rebuilds the preview for the highlighted command and goes
// back to its first example and the top.
func (m *Model) refreshPreview() { m.rebuildPreview(true) }

// rebuildPreview shows examples that match the query first, then the same text
// `explain COMMAND` prints (or the full reference with --details). Without
// reset it keeps the selected example and the scroll position.
func (m *Model) rebuildPreview(reset bool) {
	idx, off := m.exIdx, m.preview.YOffset()
	m.exIdx, m.examples = 0, nil
	e := m.Selected()
	if e == nil {
		m.preview.SetContent("Nothing matched. Clear the search to list every command.")
		return
	}
	sel, err := m.selection(e)
	if err != nil {
		m.preview.SetContent("Cannot select documentation for " + e.ID + ": " + err.Error())
		return
	}
	text, err := cli.RenderLookup(m.env.Catalog, e, sel, cli.LookupOptions{Details: m.details, Width: m.preview.Width(), Color: m.env.Color})
	if err != nil {
		m.preview.SetContent(err.Error())
		return
	}
	for _, h := range m.hits[e.ID] {
		m.examples = append(m.examples, catalog.Example{Description: h.Description, Short: h.Short, Code: h.Code})
	}
	matched := len(m.examples) > 0
	if !matched {
		for _, s := range sectionsOf(e, sel) {
			m.examples = append(m.examples, catalog.ParseExamples(s.Markdown)...)
		}
		m.examples = slices.DeleteFunc(m.examples, func(x catalog.Example) bool { return x.Code == "" })
	}
	if !reset {
		m.exIdx = min(idx, max(len(m.examples)-1, 0))
	}
	var b strings.Builder
	if matched {
		b.WriteString(m.heading("Matching examples") + "\n")
		for i, ex := range m.examples {
			marker := "  "
			if i == m.exIdx && m.focus == panePreview {
				marker = "▸ "
			}
			b.WriteString(indent(cli.RenderExamples([]catalog.Example{ex}, m.preview.Width()-2, m.env.Color, true), marker))
		}
		b.WriteString("\n")
	}
	b.WriteString(text)
	m.preview.SetContent(b.String())
	if reset {
		m.preview.GotoTop()
	} else {
		m.preview.SetYOffset(off)
	}
}

func sectionsOf(e *catalog.Entry, sel detect.Selection) []catalog.Section {
	sections, err := e.Select(sel.Implementation, "")
	if err != nil {
		return nil
	}
	return slices.DeleteFunc(sections, func(s catalog.Section) bool { return !s.IsExamples() })
}

func indent(s, marker string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		if i == 0 {
			lines[i] = marker + l
		} else {
			lines[i] = strings.Repeat(" ", len([]rune(marker))) + l
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func (m *Model) heading(s string) string {
	if m.env.Color {
		return lipgloss.NewStyle().Bold(true).Render(s)
	}
	return s
}

func (m *Model) move(delta int) {
	if len(m.entries) == 0 {
		return
	}
	m.cursor = min(max(m.cursor+delta, 0), len(m.entries)-1)
	m.status = ""
	m.refreshPreview()
}

func (m *Model) step(delta int) {
	if len(m.examples) == 0 {
		return
	}
	m.exIdx = min(max(m.exIdx+delta, 0), len(m.examples)-1)
	m.status = ""
	m.rebuildPreview(false)
}

func (m *Model) copySelected() tea.Cmd {
	if m.exIdx >= len(m.examples) {
		m.status = "No example to copy"
		return nil
	}
	code := m.examples[m.exIdx].Code
	m.status = "Copying…"
	return m.env.Copy(code)
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		m.refreshPreview()
		return m, nil
	case copiedMsg:
		if msg.err != nil {
			m.status = "Copy failed: " + msg.err.Error()
		} else {
			m.status = "Copied to " + msg.how
		}
		return m, nil
	case tea.KeyPressMsg:
		return m.key(msg)
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) key(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	k := msg.String()
	switch k {
	case "ctrl+c", "ctrl+d":
		return m, tea.Quit
	case "ctrl+y":
		return m, m.copySelected()
	case "ctrl+n", "down":
		if m.focus == panePreview {
			m.step(1)
		} else {
			m.move(1)
		}
		return m, nil
	case "ctrl+p", "up":
		if m.focus == panePreview {
			m.step(-1)
		} else {
			m.move(-1)
		}
		return m, nil
	case "enter":
		m.details = !m.details
		m.refreshPreview()
		return m, nil
	case "tab":
		if m.focus == paneInput {
			m.focus = panePreview
			m.input.Blur()
		} else {
			m.focus = paneInput
			m.rebuildPreview(false)
			return m, m.input.Focus()
		}
		m.rebuildPreview(false)
		return m, nil
	case "pgdown":
		m.preview.PageDown()
		return m, nil
	case "pgup":
		m.preview.PageUp()
		return m, nil
	case "esc":
		switch {
		case m.focus == panePreview:
			m.focus = paneInput
			m.rebuildPreview(false)
			return m, m.input.Focus()
		case m.input.Value() != "":
			m.input.SetValue("")
			m.refilter()
			return m, nil
		}
		return m, tea.Quit
	}
	if m.focus == panePreview {
		switch k {
		case "y":
			return m, m.copySelected()
		case "q":
			return m, tea.Quit
		case "space", " ":
			m.preview.PageDown()
		}
		return m, nil
	}
	before := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() != before {
		m.status = ""
		m.refilter()
	}
	return m, cmd
}

// View implements tea.Model.
func (m *Model) View() tea.View {
	faint := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)
	lw := m.listWidth()
	body := max(m.height-2, 3)
	list := make([]string, body)
	first := 0
	if m.cursor >= body {
		first = m.cursor - body + 1
	}
	for i := range list {
		idx := first + i
		line := ""
		if idx < len(m.entries) {
			e := m.entries[idx]
			line = "  " + e.ID + " " + e.Summary
			if idx == m.cursor {
				line = "▸ " + e.ID + " " + e.Summary
			}
		}
		line = fit(line, lw)
		if idx == m.cursor && idx < len(m.entries) && m.env.Color {
			line = bold.Render(line)
		}
		list[i] = line
	}
	right := strings.Split(m.preview.View(), "\n")
	var rows []string
	rows = append(rows, m.input.View())
	for i := 0; i < body; i++ {
		r := ""
		if i < len(right) {
			r = right[i]
		}
		rows = append(rows, list[i]+faint.Render(" │ ")+r)
	}
	footer := "↑/↓ list · tab examples · enter details · ctrl+y copy · esc clear/quit"
	if m.focus == panePreview {
		footer = "↑/↓ example · y copy · pgup/pgdn scroll · enter details · tab/esc back · q quit"
	}
	if m.status != "" {
		footer = bold.Render(m.status)
	} else if m.focus == panePreview && m.exIdx < len(m.examples) {
		footer = fmt.Sprintf("example %d/%d: %s", m.exIdx+1, len(m.examples), firstLine(m.examples[m.exIdx].Code))
	}
	rows = append(rows, faint.Render(lipgloss.NewStyle().MaxWidth(m.width).Render(footer)))
	v := tea.NewView(strings.Join(rows, "\n"))
	v.AltScreen = true
	return v
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

// Run starts the browser on the terminal and returns when the user quits.
func Run(env Env) error {
	_, err := tea.NewProgram(New(env)).Run()
	return err
}

// fit truncates s to width columns (with an ellipsis) and pads it with spaces,
// so list rows never wrap into the preview column.
func fit(s string, width int) string {
	r := []rune(s)
	if len(r) > width {
		r = append(r[:max(width-1, 0)], '…')
	}
	return string(r) + strings.Repeat(" ", max(width-len(r), 0))
}
