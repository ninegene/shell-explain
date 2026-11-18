package cli

import (
	"strings"
	"testing"

	"github.com/ninegene/shell-explain/internal/catalog"
)

func TestRenderOptionsAlignsColumnsAndWrapsUnderThem(t *testing.T) {
	sections := []catalog.Section{
		{Options: []string{"-p"}, Markdown: "The port to connect to, if not 22.\n"},
		{Options: []string{"-i", "--identity-file"}, Markdown: "Use this private key file for authentication when logging in to the host.\n"},
		{Options: []string{"--a-very-long-option-spelling"}, Markdown: "Own line.\n"},
	}
	got := RenderOptions(sections, 70, false)
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if lines[0] != "-p                   The port to connect to, if not 22." {
		t.Errorf("first line not aligned to the widest spelling: %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "-i, --identity-file  Use this private key") || !strings.HasPrefix(lines[2], strings.Repeat(" ", 21)) {
		t.Errorf("description did not wrap under its column:\n%s", got)
	}
	last := lines[len(lines)-2]
	if last != "--a-very-long-option-spelling" || lines[len(lines)-1] != strings.Repeat(" ", 21)+"Own line." {
		t.Errorf("a long spelling should get its own line:\n%s", got)
	}
}
