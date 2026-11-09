package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/ninegene/shell-explain/internal/catalog"
	"github.com/ninegene/shell-explain/internal/detect"
)

const upstreamPage = `---
schema_version: 1
id: test
title: test
summary: Evaluate a condition.
aliases: ['[']
tags: [shell]
implementations: [upstream]
references:
  - https://example.org/test
---

## Common

### Overview

Checks a condition.

### Background

Test is a shell command.

### Examples

Check that a file exists.

` + "```sh\ntest -e example.txt\n```" + `

### Common mistakes

Forgetting quotes.

### Portability

Same everywhere here.

## Upstream

### About test

One implementation.

### Upstream examples

Compare numbers.

` + "```sh\n[ \"$n\" -gt 3 ]\n```\n"

func TestSymbolicAliasAndUpstreamSelection(t *testing.T) {
	c, err := catalog.Load(fstest.MapFS{"en/test.md": {Data: []byte(upstreamPage)}})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"[", "/usr/bin/[", "test"} {
		var out, stderr bytes.Buffer
		status := Run(context.Background(), c, []string{name}, &out, &stderr, Config{Version: "t", Width: 40, Detector: detect.Detector{OS: "darwin", Path: t.TempDir()}})
		if status != 0 || !strings.Contains(out.String(), "test — upstream (inferred)") || !strings.Contains(out.String(), `[ "$n" -gt 3 ]`) {
			t.Errorf("%q: %d %q %q", name, status, out.String(), stderr.String())
		}
	}
	if !strings.Contains(Help, "upstream") || !strings.Contains(Help, "zsh") {
		t.Fatal("help does not list registered implementations")
	}
}
