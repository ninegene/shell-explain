package cli

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// StartPager sends long output through the user's pager. It returns the writer
// to use and a function that waits for the pager to exit. Output goes straight
// to stdout, with no pager, unless stdout is a terminal (or force is set, for
// tests) and a pager program exists.
//
// The pager is $PAGER, or less. For less, FRX is the default in $LESS, so short
// output prints directly and colors and the screen contents are kept.
func StartPager(stdout *os.File, force bool) (io.Writer, func()) {
	direct := func() {}
	if !force && !isTerminal(stdout) {
		return stdout, direct
	}
	fields := strings.Fields(os.Getenv("PAGER"))
	if len(fields) == 0 {
		fields = []string{"less"}
	}
	program, err := exec.LookPath(fields[0])
	if err != nil {
		return stdout, direct
	}
	cmd := exec.Command(program, fields[1:]...)
	cmd.Stdout, cmd.Stderr = stdout, os.Stderr
	cmd.Env = os.Environ()
	if filepath.Base(program) == "less" && os.Getenv("LESS") == "" {
		cmd.Env = append(cmd.Env, "LESS=FRX")
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		return stdout, direct
	}
	if err := cmd.Start(); err != nil {
		return stdout, direct
	}
	return in, func() {
		_ = in.Close()
		_ = cmd.Wait()
	}
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
