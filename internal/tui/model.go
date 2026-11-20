// Package tui implements `explain browse`, an interactive terminal browser over
// the same catalog and renderer as the command line. It never executes the
// examples it shows.
package tui

import (
	"fmt"
	"regexp"
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
	exLine   []int // content line of each example
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

// listWidth fits the longest command name plus the two-column marker.
func (m *Model) listWidth() int {
	w := 0
	for _, e := range m.env.Catalog.Entries {
		w = max(w, len([]rune(e.ID)))
	}
	return min(max(w+3, 12), max(m.width/3, 12))
}

func (m *Model) layout() {
	m.input.SetWidth(max(m.width-len(m.input.Prompt)-1, 10))
	m.preview.SetWidth(max(m.width-m.listWidth()-3, 10))
	m.preview.SetHeight(max(m.height-4, 3))
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
	// Command names come first: an exact name, then names that start with the
	// query, then names that contain it, then everything else (still in
	// catalog order).
	q := strings.ToLower(query)
	rank := func(e *catalog.Entry) int {
		id := strings.ToLower(e.ID)
		switch {
		case q == "":
			return 0
		case id == q:
			return 0
		case strings.HasPrefix(id, q):
			return 1
		case strings.Contains(id, q):
			return 2
		}
		return 3
	}
	slices.SortStableFunc(m.entries, func(a, b *catalog.Entry) int { return rank(a) - rank(b) })
	m.cursor = 0
	if keep != nil && query == "" {
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
	matched := len(m.hits[e.ID]) > 0
	for _, h := range m.hits[e.ID] {
		m.examples = append(m.examples, catalog.Example{Description: h.Description, Short: h.Short, Code: h.Code})
	}
	if !matched {
		for _, s := range sectionsOf(e, sel) {
			m.examples = append(m.examples, catalog.ParseExamples(s.Markdown)...)
		}
		m.examples = slices.DeleteFunc(m.examples, func(x catalog.Example) bool { return x.Code == "" })
	}
	if !reset {
		m.exIdx = min(idx, max(len(m.examples)-1, 0))
	}
	// Both parts are indented two columns so the marker of the selected
	// example has room.
	text, err := cli.RenderLookup(m.env.Catalog, e, sel, cli.LookupOptions{Details: m.details, Width: m.preview.Width() - 2, Color: m.env.Color, Heading: m.dim})
	if err != nil {
		m.preview.SetContent(err.Error())
		return
	}
	var lines []string
	m.exLine = make([]int, len(m.examples))
	if matched {
		for i, ex := range m.examples {
			m.exLine[i] = len(lines)
			block := strings.Split(strings.TrimRight(cli.RenderExamples([]catalog.Example{ex}, m.preview.Width()-2, m.env.Color, true), "\n"), "\n")
			lines = append(lines, block...)
		}
		lines = append(lines, "")
	}
	body := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if matched && !m.details {
		// The matches are enough; the rest is one keypress away.
		body = []string{m.dim("enter: the full reference for " + e.ID)}
	}
	if !matched {
		from := 0
		for i, ex := range m.examples {
			m.exLine[i] = -1
			want := strings.TrimSpace(firstLine(ex.Code))
			for j := from; j < len(body); j++ {
				if strings.Contains(ansiRE.ReplaceAllString(body[j], ""), want) {
					m.exLine[i], from = len(lines)+j, j+1
					break
				}
			}
		}
	}
	lines = append(lines, body...)
	for i, l := range lines {
		marker := "  "
		if m.focus == panePreview && slices.Contains(m.exLine, i) && m.exLine[m.exIdx] == i {
			marker = "▸ "
		}
		lines[i] = marker + l
	}
	m.preview.SetContent(strings.Join(lines, "\n") + "\n")
	if reset {
		m.preview.GotoTop()
	} else {
		m.preview.SetYOffset(off)
	}
}

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*m")

// showExample scrolls the preview just enough to bring the selected example's
// first line into view.
func (m *Model) showExample() {
	if m.exIdx >= len(m.exLine) || m.exLine[m.exIdx] < 0 {
		return
	}
	line, top, h := m.exLine[m.exIdx], m.preview.YOffset(), m.preview.Height()
	switch {
	case line < top:
		m.preview.SetYOffset(line)
	case line >= top+h:
		m.preview.SetYOffset(line - h + 1)
	}
}

func sectionsOf(e *catalog.Entry, sel detect.Selection) []catalog.Section {
	sections, err := e.Select(sel.Implementation, "")
	if err != nil {
		return nil
	}
	return slices.DeleteFunc(sections, func(s catalog.Section) bool { return !s.IsExamples() })
}

func (m *Model) dim(s string) string {
	if m.env.Color {
		return lipgloss.NewStyle().Faint(true).Render(s)
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
	m.showExample()
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
		case "right", "n":
			m.step(1)
		case "left", "p":
			m.step(-1)
		case "j":
			m.preview.ScrollDown(1)
		case "k":
			m.preview.ScrollUp(1)
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
	body := max(m.height-4, 3)
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
			line = "  " + e.ID
			if idx == m.cursor {
				line = "▸ " + e.ID
			}
		}
		line = fit(line, lw)
		if idx == m.cursor && idx < len(m.entries) && m.env.Color {
			line = bold.Render(line)
		}
		list[i] = line
	}
	right := strings.Split(m.preview.View(), "\n")
	sep := faint.Render(" │ ")
	cmdTitle, exTitle := fmt.Sprintf("Commands %d/%d", min(m.cursor+1, len(m.entries)), len(m.entries)), "Examples 0/0"
	if len(m.examples) > 0 {
		exTitle = fmt.Sprintf("Examples %d/%d", m.exIdx+1, len(m.examples))
	}
	if m.details {
		exTitle = "Details · " + exTitle
	}
	// The focused pane's title is bold; the other is dim.
	active := bold
	if m.focus == panePreview {
		cmdTitle, exTitle = faint.Render(fit(cmdTitle, lw)), active.Render(exTitle)
	} else {
		cmdTitle, exTitle = active.Render(fit(cmdTitle, lw)), faint.Render(exTitle)
	}
	var rows []string
	rows = append(rows, m.input.View(), cmdTitle+sep+exTitle)
	for i := 0; i < body; i++ {
		r := ""
		if i < len(right) {
			r = right[i]
		}
		rows = append(rows, list[i]+sep+r)
	}
	// The footer takes two rows: key hints (or a status message), then the
	// selected example's command in color.
	footer := faint.Render("↑/↓ list · tab commands · enter details · ctrl+y copy · esc clear/quit")
	if m.focus == panePreview {
		footer = faint.Render("↑/↓ example · j/k scroll · y copy · pgup/pgdn page · enter details · tab/esc back · q quit")
	}
	if m.details {
		footer = strings.ReplaceAll(footer, "enter details", "enter hide details")
	}
	if m.status != "" {
		footer = bold.Render(m.status)
	}
	second := ""
	if m.exIdx < len(m.examples) {
		second = strings.TrimRight(cli.RenderExamples([]catalog.Example{{Code: firstLine(m.examples[m.exIdx].Code)}}, m.width, m.env.Color, false), "\n")
	}
	rows = append(rows, lipgloss.NewStyle().MaxWidth(m.width).Render(footer), lipgloss.NewStyle().MaxWidth(m.width).Render(second))
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
