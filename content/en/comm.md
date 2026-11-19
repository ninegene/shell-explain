---
schema_version: 1
id: comm
title: comm
summary: Compare two sorted files and show lines unique to each or shared.
aliases: [gcomm]
tags: [text, files, compare]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/comm-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=comm&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnucomm.1.html
---
## Common
### Overview
`comm` compares two files that are already sorted and prints three columns: lines only in the first file, lines only in the second, and lines in both. You hide a column by adding its number as an option, so `comm -12` shows only what is in both, and `comm -23` only what is in the first file alone.

### Background
A `comm` command appeared in Version 4 AT&T UNIX, according to the macOS manual. Today the name covers GNU comm, the comm that macOS ships, and the comm applet in BusyBox. They agree on `-1`, `-2` and `-3`. `comm` needs sorted input in the same order on both sides, or the result is wrong.

### Examples
Show the lines that are in both files, which is the `-1 -2` hiding of the first two columns:

```sh
printf 'a\nb\nc\n' > one.txt && printf 'b\nc\nd\n' > two.txt && comm -12 one.txt two.txt
```

Show the lines only in the first file:

```sh
printf 'a\nb\nc\n' > one.txt && printf 'b\nc\nd\n' > two.txt && comm -23 one.txt two.txt
```

Show the lines only in the second file:

```sh
printf 'a\nb\nc\n' > one.txt && printf 'b\nc\nd\n' > two.txt && comm -13 one.txt two.txt
```

Sort unsorted input first, with `sort`, by writing the sorted copies to files:

```sh
printf 'c\na\nb\n' | sort > one.txt && printf 'd\nb\nc\n' | sort > two.txt && comm -12 one.txt two.txt
```

### Option: -1
Hide the first column, the lines only in the first file.

### Option: -2
Hide the second column, the lines only in the second file.

### Option: -3
Hide the third column, the lines in both files.

### Common mistakes
Unsorted input gives wrong answers without an error, and GNU may print a warning about it. Both files must also be sorted the same way, for example with the same `LC_ALL`, so sort them with the same settings you use for `comm`. When you show all three columns, the second and third are indented with tabs, which can confuse a script; hide the columns you do not need.

### Portability
Only the three column options are in the Common examples. GNU adds `--check-order` and `--output-delimiter`, and macOS has `-i` to ignore case. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU comm
GNU comm is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can check that the input is sorted and change the column separator.

### GNU examples
Stop with an error when the input is not sorted, instead of giving a wrong answer, with `--check-order`:

```sh
printf 'b\na\n' > one.txt && printf 'a\nb\n' > two.txt && { comm --check-order one.txt two.txt > /dev/null 2>&1 || echo not sorted; }
```

### Option: --check-order
Fail when an input file is not sorted, even when it would not matter for the comparison.

### Option: --output-delimiter
Separate the columns with this string instead of a tab.

## Apple
### About macOS comm
This is the comm in `/usr/bin` on macOS. It follows the BSD implementation, and has `-i` to ignore case.

### Apple notes
macOS `comm` has `-i` to ignore case when comparing. It has no `--check-order` or `--output-delimiter`.

## BusyBox
### About BusyBox comm
BusyBox bundles comm as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `comm` supports `-1`, `-2` and `-3`. GNU long options depend on the build.
