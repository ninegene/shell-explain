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
	"github.com/ninegene/shell-explain/internal/startup"
	"github.com/ninegene/shell-explain/internal/tui"
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
	cfg := cli.Config{Version: version, Width: width, GuideWidth: guideWidth, Color: cli.TerminalColor(os.Stdout), Detector: detect.Detector{OS: runtime.GOOS, Path: os.Getenv("PATH")}}
	interactive := cli.IsTerminal(os.Stdin) && cli.IsTerminal(os.Stdout)
	if interactive {
		cfg.Browse = func(query string) error { return browse(c, cfg, query) }
		if len(os.Args) == 1 {
			os.Exit(startInteractive(cfg))
		}
	}
	// Long text output goes through the pager on a terminal (see cli.WantsPager).
	var out io.Writer = os.Stdout
	finish := func() {}
	if cli.WantsPager(os.Args[1:]) {
		out, finish = cli.StartPager(os.Stdout, false)
	}
	status := cli.Run(context.Background(), c, os.Args[1:], out, os.Stderr, cfg)
	finish()
	os.Exit(status)
}

// browse starts the interactive browser. It selects documentation without
// launching any process, like `--no-probe`.
func browse(c *catalog.Catalog, cfg cli.Config, query string) error {
	return tui.Run(tui.Env{
		Catalog: c,
		Select: func(e *catalog.Entry) (detect.Selection, error) {
			return cli.SelectImplementation(context.Background(), cfg, e, true)
		},
		Color: cfg.Color,
		Query: query,
	})
}

// startInteractive handles `explain` with no arguments on a terminal: show the
// usage when wanted (asking the first time), then open the browser.
func startInteractive(cfg cli.Config) int {
	path, _ := startup.Path()
	flow := startup.Flow{
		ConfigPath: path,
		In:         os.Stdin,
		Err:        os.Stderr,
		ShowHelp: func() error {
			out, finish := io.Writer(os.Stdout), func() {}
			if !cli.PagerDisabled() {
				out, finish = cli.StartPager(os.Stdout, false)
			}
			_, err := fmt.Fprint(out, cli.Help)
			finish()
			return err
		},
		Browse: func() error { return cfg.Browse("") },
	}
	if err := flow.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "explain:", err)
		return 1
	}
	return 0
}
