---
schema_version: 1
id: uniq
title: uniq
summary: Report or omit repeated adjacent lines.
aliases: [guniq]
tags: [text, files, counting]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/uniq-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=uniq&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`uniq` collapses runs of identical adjacent lines. It does not find duplicates that are not next to each other, so sort the input first when you want every duplicate treated alike.

### Background
A `uniq` command appeared in Version 3 Unix, according to the FreeBSD and macOS manuals. It only collapses adjacent duplicate lines, which is why it usually follows `sort`. Today the name covers GNU uniq, the uniq that macOS ships, and the uniq applet in BusyBox. They agree on `-c`, `-d`, and `-u`, and differ on options such as `-i` and field skipping.

### Examples
Collapse duplicates after sorting:

```sh
sort fruit.txt | uniq
```

Count how many times each line occurs:

```sh
sort fruit.txt | uniq -c
```

Show only lines that occur more than once:

```sh
sort fruit.txt | uniq -d
```

Show only lines that occur exactly once:

```sh
sort fruit.txt | uniq -u
```

Rank values by how often they occur:

```sh
sort fruit.txt | uniq -c | sort -rn
```

Without sorting, only adjacent repeats are removed, so this file is unchanged:

```sh short="Without sorting, only adjacent repeats are removed"
uniq fruit.txt
```

### Option: -c
Prefix each line with the number of times it occurred in a row. The padding differs by implementation.

### Option: -d
Print only one copy of each line that is repeated.

### Option: -u
Print only lines that are not repeated.

### Common mistakes
`uniq` only compares neighbours, so unsorted input keeps non-adjacent duplicates. Use `sort -u` when you only need unique lines and do not need counts. A second file argument names an output file, which it overwrites, so pass input on standard input in scripts to be safe.

### Portability
`-c`, `-d`, and `-u` are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). Options such as `-i` (ignore case) and field skipping differ between builds.

## GNU
### About GNU uniq
GNU uniq is part of GNU coreutils, and this page's checks ran coreutils 9.11. It adds long options such as `--count` and options such as `-i` and `-f`.

### GNU notes
GNU `uniq` adds long options such as `--count` and options such as `-i` (ignore case) and `-f` (skip fields). Confirm that a macOS or BusyBox build provides an option before relying on it.

## Apple
### About macOS uniq
This is the uniq in `/usr/bin` on macOS. GNU long options such as `--count` are not available, so use the short spellings.

### Apple notes
macOS `uniq` is the BSD implementation. GNU long options such as `--count` are not available; use the short spellings.

## BusyBox
### About BusyBox uniq
BusyBox bundles uniq as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `uniq` supports `-c`, `-d`, and `-u`; other options depend on the build.
