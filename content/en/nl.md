---
schema_version: 1
id: nl
title: nl
summary: Number the lines of a file.
aliases: [gnl]
tags: [text, files]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/nl-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=nl&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`nl` prints a file with a line number in front of each line. By default it numbers only lines that are not empty, and it puts a tab between the number and the text. For numbering every line, use `-ba`; `cat -n` is the shorter way when that is all you need.

### Background
The `nl` utility first appeared in AT&T System V Release 2 UNIX, according to the macOS manual. Today the name covers GNU nl, the nl that macOS ships, and the nl applet in BusyBox. They agree on the default numbering and on `-b`, `-s` and `-w`.

### Examples
Number the lines of a file:

```sh
nl notes.txt
```

Number every line, including empty ones, with `-b a`:

```sh
printf 'one\n\ntwo\n' | nl -b a
```

Number only the lines that are not empty, which is the default, written out as `-b t`:

```sh
printf 'one\n\ntwo\n' | nl -b t
```

Use a dot and a space between the number and the text with `-s`:

```sh
nl -s'. ' notes.txt
```

Use a narrower number column with `-w`:

```sh
nl -w 2 notes.txt
```

### Option: -b
Choose which lines to number: `a` for all, `t` for non-empty ones, `n` for none.

### Option: -s
Put this text between the number and the line.

### Option: -w
Set the width of the number column.

### Common mistakes
By default blank lines are not numbered and are not counted, so the numbers can look as if they skip. Add `-b a` to number them. Line numbers are padded with spaces, so a script that reads them should split on whitespace rather than counting characters.

### Portability
Only `-b`, `-s` and `-w` are in the Common examples. Section and page handling options differ, and the padding of the number column is not the same on every system, so the examples compare output after collapsing blanks. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU nl
GNU nl is part of GNU coreutils, and this page's checks ran coreutils 9.11. It has long option names and number formats.

### GNU examples
Pad the numbers with zeros instead of spaces, with the number format `rz` (right-justified, zero-filled):

```sh
nl -n rz -w 3 notes.txt
```

### Option: -n, --number-format
Choose how the number is written: `ln` left-justified, `rn` right-justified, or `rz` right-justified with zeros.

## Apple
### About macOS nl
This is the nl in `/usr/bin` on macOS. It follows the BSD implementation.

### Apple notes
macOS `nl` has `-b`, `-s` and `-w`, plus `-n` to pick the number format.

## BusyBox
### About BusyBox nl
BusyBox bundles nl as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `nl` supports `-b`, `-s` and `-w`. Other options depend on the build.
