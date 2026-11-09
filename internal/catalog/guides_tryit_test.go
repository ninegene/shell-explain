package catalog

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/examplecheck"
)

// Every command a guide asks a reader to try must be safe, must succeed on
// this machine, and must print what the guide says it prints. Commands run
// only in a temporary folder, and only as an explicit test operation.
func TestGuideTryItCommands(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	ran := 0
	for _, g := range c.Guides {
		for _, s := range g.Sections {
			if s.Title != "Try it" {
				continue
			}
			blocks := examplecheck.TryItBlocks(s.Markdown)
			if len(blocks) == 0 {
				t.Errorf("guide %s has a Try it section with no command", g.ID)
			}
			for i, b := range blocks {
				t.Run(g.ID+"/"+string(rune('a'+i)), func(t *testing.T) {
					if why := examplecheck.UnsafeTryIt(b.Code); why != "" {
						t.Fatalf("%q is not allowed in a Try it block: %s", why, b.Code)
					}
					dir := t.TempDir()
					if err := examplecheck.Setup(dir, g.ID); err != nil {
						t.Fatal(err)
					}
					cmd := exec.Command("/bin/sh", "-c", b.Code)
					cmd.Dir = dir
					cmd.Env = []string{"PATH=/usr/bin:/bin:/usr/sbin:/sbin", "HOME=" + dir, "LC_ALL=C", "TERM=dumb"}
					timer := time.AfterFunc(20*time.Second, func() { _ = cmd.Process.Kill() })
					defer timer.Stop()
					out, err := cmd.Output()
					if err != nil {
						t.Fatalf("%s: %v", b.Code, err)
					}
					if err := examplecheck.CheckTryIt(b, string(out)); err != nil {
						t.Fatalf("%s: %v", b.Code, err)
					}
					ran++
				})
			}
		}
	}
	if ran == 0 {
		t.Fatal("no Try it commands ran")
	}
	// The runner must leave the working folder alone.
	if _, err := os.Stat("practice"); err == nil {
		t.Error("a Try it command wrote into the test's own folder")
	}
}
