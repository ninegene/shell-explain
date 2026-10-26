package main

import (
	"context"
	"fmt"
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
	width := 80
	if n, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && n >= 20 && n <= 1000 {
		width = n
	}
	os.Exit(cli.Run(context.Background(), c, os.Args[1:], os.Stdout, os.Stderr, cli.Config{Version: version, Width: width, Color: cli.TerminalColor(os.Stdout), Detector: detect.Detector{OS: runtime.GOOS, Path: os.Getenv("PATH")}}))
}
