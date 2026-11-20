package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/cli"
	"github.com/ninegene/shell-explain/internal/detect"
)

var version = "dev"

func main() {
	c, err := catalog.Load(content.Files)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	width, guideWidth := catalog.DescriptionWidth, 0
	if n, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && n >= 20 && n <= 1000 {
		width, guideWidth = n, n
	}
	// Long text output goes through the pager on a terminal (see cli.WantsPager).
	var out io.Writer = os.Stdout
	finish := func() {}
	if cli.WantsPager(os.Args[1:]) {
		out, finish = cli.StartPager(os.Stdout, false)
	}
	status := cli.Run(context.Background(), c, os.Args[1:], out, os.Stderr, cli.Config{Version: version, Width: width, GuideWidth: guideWidth, Color: cli.TerminalColor(os.Stdout), Detector: detect.Detector{OS: runtime.GOOS, Path: os.Getenv("PATH")}})
	finish()
	os.Exit(status)
}
