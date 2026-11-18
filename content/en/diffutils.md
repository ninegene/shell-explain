---
schema_version: 1
id: diffutils
title: diffutils
summary: The GNU package that provides diff, cmp, diff3 and sdiff.
aliases: []
tags: [package, files, compare]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/diffutils/manual/html_node/index.html
---
## Common
### Overview
GNU diffutils is a package, not a single command. It supplies four programs that compare files: `diff` (show the differences between two files or folders), `cmp` (compare two files byte by byte), `diff3` (compare three files) and `sdiff` (compare two files side by side, or merge them). The `patch` command that applies a difference is a separate GNU package. Each program has its own page on this site; this page explains how the package relates to what you have installed.

### Background
The GNU manual documents the four programs together as the diffutils manual. Linux distributions install the package by default or on request. macOS ships its own `diff` (based on FreeBSD diff), `cmp`, `diff3` and `sdiff`; BusyBox has small `diff` and `cmp` applets and no `diff3` or `sdiff`. The version checked here is GNU diffutils 3.12.

### Examples
Ask a program of the package for its version; the first line names the package. Reviewed only, because this page has no program of its own to run:

```sh short="Ask a program of the package for its version"
diff --version | head -n 1
```

The other programs report the same package; reviewed only:

```sh
cmp --version | head -n 1
```

### Common mistakes
The words `diff` and `cmp` on your system do not tell you which implementation you have: macOS ships the BSD-derived ones, and BusyBox has smaller ones, so options such as `--color` and `-y` may be missing. Check with `diff --version`, which GNU supports. `diff3` and `sdiff` run `diff` by name, so on a system where `diff` is another implementation they can fail. Each program has its own page for options.

### Portability
This is a package overview. The programs are `diff`, `cmp`, `diff3` and `sdiff`, each with its own page. The examples are reviewed only, because there is no executable named `diffutils` to run.

## GNU
### About GNU diffutils
GNU diffutils is a GNU project package, and this page was checked against version 3.12.
