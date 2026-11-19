---
schema_version: 1
id: expand
title: expand
summary: Convert tabs to spaces.
aliases: [gexpand]
tags: [text, formatting]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/expand-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=expand&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuexpand.1.html
---
## Common
### Overview
`expand` replaces each tab character with enough spaces to reach the next tab stop, which is every 8 columns by default. Use it when text with tabs must look the same everywhere, for example in a plain-text report or when pasting into a program that treats tabs strangely. The reverse is `unexpand`.

### Background
The `expand` utility first appeared in 1BSD, according to the macOS manual, which also says `expand` and `unexpand` conform to POSIX.1-2001. Today the name covers GNU expand, the expand that macOS ships, and the expand applet in BusyBox. They agree on `-t` for the tab width.

### Examples
Convert a tab to spaces with the default width of 8 columns:

```sh
printf 'a\tb\n' | expand
```

Use a tab stop every 4 columns with `-t`:

```sh
printf 'a\tb\n' | expand -t 4
```

Check that no tab is left, by looking for the `\t` that `od -c` prints for a tab:

```sh
printf 'a\tb\n' | expand | od -c | grep -q '\\t' || echo no tabs left
```

Convert a file, writing the result to a new one:

```sh
printf 'x\ty\n' > in.txt && expand in.txt > out.txt && cat out.txt
```

### Option: -t
Set the distance between tab stops, in columns, as in `-t 4`.

### Common mistakes
A tab moves to the next tab stop, not a fixed number of spaces, so `a<TAB>b` and `abcdefg<TAB>b` produce different amounts of space. `expand` reads a file and writes to standard output; it does not change the file in place. In a Makefile a tab is required, so never run `expand` on one.

### Portability
Only `-t` with one number is in the Common examples. GNU also accepts a list of stops and `-i` to change only leading tabs, which is not in every implementation. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU expand
GNU expand is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### GNU examples
Change only the tabs at the start of each line with `-i`, leaving later tabs alone:

```sh
printf '\ta\tb\n' | expand -i -t 4 | od -c | grep -c '\\t'
```

### Option: -i, --initial
Convert only the tabs at the start of a line.

## Apple
### About macOS expand
This is the expand in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `expand` has `-t`. It does not have GNU's `-i`.

## BusyBox
### About BusyBox expand
BusyBox bundles expand as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `expand` supports `-t`. `-i` and the list form of `-t` depend on the build.
