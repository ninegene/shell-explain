package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// copiedMsg reports the outcome of a clipboard write.
type copiedMsg struct {
	how string
	err error
}

type clipboardTool struct {
	program string
	args    []string
}

// clipboardTools lists the programs that put text on the system clipboard, in
// the order they are tried. They only read the text from standard input.
func clipboardTools(goos string) []clipboardTool {
	if goos == "darwin" {
		return []clipboardTool{{"pbcopy", nil}}
	}
	return []clipboardTool{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}
}

// CopyToClipboard returns a command that copies text to the system clipboard
// with the first available helper program, or with the terminal's OSC 52
// sequence when there is none. It never executes the text.
func CopyToClipboard(text string) tea.Cmd {
	for _, tool := range clipboardTools(runtime.GOOS) {
		path, err := exec.LookPath(tool.program)
		if err != nil {
			continue
		}
		return func() tea.Msg {
			cmd := exec.Command(path, tool.args...)
			cmd.Stdin = strings.NewReader(text)
			if err := cmd.Run(); err != nil {
				return copiedMsg{err: fmt.Errorf("%s: %w", tool.program, err)}
			}
			return copiedMsg{how: tool.program}
		}
	}
	return tea.Batch(tea.SetClipboard(text), func() tea.Msg { return copiedMsg{how: "terminal clipboard (OSC 52)"} })
}
