//go:build !darwin && !linux

package detect

import "os/exec"

// Only macOS and Linux currently support descendant cleanup.
func configureProbe(cmd *exec.Cmd) func() {
	return func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
}
