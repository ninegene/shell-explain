package cli

import (
	"fmt"
	"io"
	"strings"
)

// runBrowse starts the interactive browser. It needs a terminal and takes no
// selection flags, because the browser chooses documentation itself.
func runBrowse(op options, out, errOut io.Writer, cfg Config) int {
	if op.format != "text" {
		return failure(out, errOut, op.format, fmt.Errorf("browse is interactive; use explain search --format json"))
	}
	if op.implementation != "" || op.platform != "" || op.option != "" || op.details || op.noProbe || op.section != "" || op.limitSet {
		return failure(out, errOut, op.format, fmt.Errorf("selection flags apply to command lookup only"))
	}
	if cfg.Browse == nil {
		return failure(out, errOut, op.format, fmt.Errorf("browse needs a terminal; use explain search"))
	}
	if err := cfg.Browse(strings.Join(op.args[1:], " ")); err != nil {
		fmt.Fprintln(errOut, "explain:", err)
		return 1
	}
	return 0
}
