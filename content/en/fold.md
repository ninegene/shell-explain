---
schema_version: 1
id: fold
title: fold
summary: Wrap long lines to a fixed width.
aliases: [gfold]
tags: [text, formatting]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/fold-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=fold&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`fold` breaks lines that are longer than a width into shorter lines. It cuts at exactly that width, even in the middle of a word, unless you add `-s`, which breaks at spaces. It is a quick way to make text fit a narrow screen or a fixed-width format.

### Background
The `fold` utility first appeared in 1BSD and was rewritten for 4.3BSD-Reno, and `-b` and `-s` were added to NetBSD 1.0 for POSIX.2 compliance, according to the macOS manual. Today the name covers GNU fold, the fold that macOS ships, and the fold applet in BusyBox. They agree on `-w` and `-s`.

### Examples
Wrap lines at 10 characters, even in the middle of a word:

```sh
echo 'the quick brown fox' | fold -w 10
```

Break at spaces instead, so words stay whole, with `-s`:

```sh
echo 'the quick brown fox' | fold -s -w 12
```

Compare with the same width without `-s`, which cuts a word:

```sh
echo 'the quick brown fox' | fold -w 12
```

Count the lines a long text becomes:

```sh
echo 'the quick brown fox' | fold -s -w 10 | grep -c .
```

### Option: -w
Set the width. The default is 80.

### Option: -s
Break after the last space that fits, so a word is not split when it can be avoided.

### Common mistakes
Without `-s`, `fold` splits words, which makes the text hard to read. The width counts characters in the common case, and tabs and some wide characters are counted differently on different systems. `fold -s` keeps the space at the end of the line, so you may see a trailing space.

### Portability
Only `-w` and `-s` are in the Common examples. GNU and macOS also have `-b` to count bytes instead of columns. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU fold
GNU fold is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### GNU examples
Count bytes instead of columns with `-b`, which matters for text with tabs or backspaces:

```sh
printf 'abcdefgh\n' | fold -b -w 4
```

### Option: -b, --bytes
Count bytes instead of display columns.

## Apple
### About macOS fold
This is the fold in `/usr/bin` on macOS. It follows the BSD implementation.

### Apple notes
macOS `fold` has `-b`, `-s` and `-w`, with the same meaning as GNU's.

## BusyBox
### About BusyBox fold
BusyBox bundles fold as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `fold` supports `-w` and `-s`. `-b` depends on the build.
