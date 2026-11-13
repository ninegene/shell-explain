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

// Exercise the authored snippets with a delayed process and a failed lookup,
// without signalling any real process by name.
func TestPkillExamplesWaitForReadiness(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := c.Find("pkill")
	if !ok {
		t.Fatal("pkill page missing")
	}
	sections, err := entry.Select("common", "")
	if err != nil {
		t.Fatal(err)
	}
	var examples []string
	for _, section := range sections {
		for _, code := range section.Code {
			if strings.Contains(code, "trap 'kill") {
				examples = append(examples, code)
			}
		}
	}
	if len(examples) != 2 {
		t.Fatalf("got %d guarded pkill examples, want 2", len(examples))
	}
	for _, neverReady := range []bool{false, true} {
		for _, code := range examples {
			checkPkillReadiness(t, code, neverReady)
		}
	}
}

func checkPkillReadiness(t *testing.T, code string, neverReady bool) {
	t.Helper()
	dir := t.TempDir()
	tools := filepath.Join(dir, "tools")
	if err := os.Mkdir(tools, 0o700); err != nil {
		t.Fatal(err)
	}
	stubs := map[string]string{
		"sleep": "#!/bin/sh\nif [ \"$1\" = 1 ]; then exec /bin/sleep 0.05; fi\n/bin/sleep 0.1\n: > \"$READY_FILE\"\n",
		"pgrep": "#!/bin/sh\n[ \"$NEVER_READY\" = false ] && [ -f \"$READY_FILE\" ]\n",
		"pkill": "#!/bin/sh\nrm -f \"$READY_FILE\"\n",
	}
	for name, script := range stubs {
		if err := os.WriteFile(filepath.Join(tools, name), []byte(script), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	script := code
	if neverReady {
		script = "trap 'echo outer-exit' 0\n" + code + "\necho status:$?\necho outer-alive"
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PATH="+tools+":/usr/bin:/bin", "READY_FILE="+filepath.Join(dir, "ready"), "NEVER_READY="+map[bool]string{false: "false", true: "true"}[neverReady])
	output, runErr := cmd.CombinedOutput()
	timedOut := ctx.Err() == context.DeadlineExceeded
	cancel()
	if timedOut {
		t.Fatal("example timed out")
	}
	if neverReady {
		if runErr != nil || string(output) != "status:1\nouter-alive\nouter-exit\n" {
			t.Fatalf("failed lookup did not preserve the calling shell: %v: %q", runErr, output)
		}
		return
	}
	if runErr != nil {
		t.Fatalf("example failed: %v: %s", runErr, output)
	}
	if err := examplecheck.Verify(code, dir, string(output)); err != nil {
		t.Fatal(err)
	}
}
