package catalog

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/examplecheck"
)

// Run the copyable examples in zsh with a simulated, delayed process.
func TestPgrepExamplesInZsh(t *testing.T) {
	zsh, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh is not installed")
	}
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := c.Find("pgrep")
	if !ok {
		t.Fatal("pgrep page missing")
	}
	count := 0
	for _, section := range entry.Sections {
		for _, code := range section.Code {
			if strings.Contains(code, "lookup_status=") {
				count++
				checkPgrepZshExample(t, zsh, code)
			}
		}
	}
	if count != 4 {
		t.Fatalf("got %d pgrep examples, want 4", count)
	}
}

func checkPgrepZshExample(t *testing.T, zsh, code string) {
	t.Helper()
	dir := t.TempDir()
	tools := filepath.Join(dir, "tools")
	if err := os.Mkdir(tools, 0o700); err != nil {
		t.Fatal(err)
	}
	stubs := map[string]string{
		"sleep": "#!/bin/sh\nif [ \"$1\" = 1 ]; then exec /bin/sleep 0.05; fi\n/bin/sleep 0.1\nprintf '%s\\n' \"$$\" > \"$READY_FILE\"\n",
		"pgrep": "#!/bin/sh\n[ -s \"$READY_FILE\" ] || exit 1\ncase \"$1\" in -c) echo 1;; -l) printf '%s sleep\\n' \"$(cat \"$READY_FILE\")\";; *) cat \"$READY_FILE\";; esac\n",
	}
	for name, script := range stubs {
		if err := os.WriteFile(filepath.Join(tools, name), []byte(script), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, zsh, "-f", "-c", code)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PATH="+tools+":/usr/bin:/bin", "READY_FILE="+filepath.Join(dir, "ready"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("zsh example failed: %v: %s", err, output)
	}
	if err := examplecheck.Verify(code, dir, string(output)); err != nil {
		t.Fatal(err)
	}
}
