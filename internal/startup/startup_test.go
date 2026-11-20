package startup

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func flow(t *testing.T, path, answer string) (Flow, *int, *int, *bytes.Buffer) {
	t.Helper()
	var helps, browses int
	var errOut bytes.Buffer
	return Flow{
		ConfigPath: path,
		In:         strings.NewReader(answer),
		Err:        &errOut,
		ShowHelp:   func() error { helps++; return nil },
		Browse:     func() error { browses++; return nil },
	}, &helps, &browses, &errOut
}

func TestFirstRunAsksThenRemembersNo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "explain", "config")
	f, helps, browses, errOut := flow(t, path, "n\n")
	if err := f.Run(); err != nil || *helps != 1 || *browses != 1 || !strings.Contains(errOut.String(), "[Y/n]") {
		t.Fatalf("first run: helps=%d browses=%d err=%v out=%q", *helps, *browses, err, errOut.String())
	}
	if show, known := Load(path); show || !known {
		t.Fatalf("answer not saved as no: show=%v known=%v", show, known)
	}
	f, helps, browses, errOut = flow(t, path, "")
	if err := f.Run(); err != nil || *helps != 0 || *browses != 1 || errOut.Len() != 0 {
		t.Fatalf("after no: helps=%d browses=%d out=%q", *helps, *browses, errOut.String())
	}
}

func TestDefaultAnswerKeepsShowingUsage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	f, _, _, _ := flow(t, path, "\n")
	if err := f.Run(); err != nil {
		t.Fatal(err)
	}
	f, helps, browses, errOut := flow(t, path, "")
	if err := f.Run(); err != nil || *helps != 1 || *browses != 1 || errOut.Len() != 0 {
		t.Fatalf("after yes: helps=%d browses=%d out=%q", *helps, *browses, errOut.String())
	}
}

func TestUnsavablePreferenceIsReportedAndIgnored(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "file")
	if err := os.WriteFile(blocker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	f, _, browses, errOut := flow(t, filepath.Join(blocker, "config"), "n\n")
	if err := f.Run(); err != nil || *browses != 1 || !strings.Contains(errOut.String(), "could not save") {
		t.Fatalf("browses=%d err=%v out=%q", *browses, err, errOut.String())
	}
	f, _, browses, errOut = flow(t, "", "\n")
	if err := f.Run(); err != nil || *browses != 1 || !strings.Contains(errOut.String(), "not saved") {
		t.Fatalf("no config dir: browses=%d err=%v out=%q", *browses, err, errOut.String())
	}
}

func TestLoadIgnoresGarbage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	for _, text := range []string{"", "show_help_on_start = maybe\n", "other = true\n", "\x00\x01"} {
		if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, known := Load(path); known {
			t.Errorf("%q must not count as an answer", text)
		}
	}
}
