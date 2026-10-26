//go:build darwin || linux

package detect

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// A separate process group cleans up ordinary children. It is not a sandbox:
// descendants can escape the group, so Cmd.WaitDelay also bounds pipe waiting.
func configureProbe(cmd *exec.Cmd) func() {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	kill := func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return err
	}
	cmd.Cancel = kill
	return func() { _ = kill() }
}
