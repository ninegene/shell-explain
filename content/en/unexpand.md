---
schema_version: 1
id: unexpand
title: unexpand
summary: Convert runs of spaces to tabs.
aliases: [gunexpand]
tags: [text, formatting]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/unexpand-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=unexpand&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuunexpand.1.html
---
## Common
### Overview
`unexpand` replaces runs of spaces with tabs where they reach a tab stop. By default only spaces at the start of a line are converted; `-a` converts them everywhere. It is the reverse of `expand`, and is used to compress aligned text or to restore tabs in a file where they are required.

### Background
The macOS manual describes `expand` and `unexpand` together, saying `expand` first appeared in 1BSD and that both conform to POSIX.1-2001. Today the name covers GNU unexpand, the unexpand that macOS ships, and the unexpand applet in BusyBox. They agree on `-a` and `-t`.

### Examples
Convert 8 leading spaces to a tab; the default converts only leading blanks:

```sh
printf '        a\n' | unexpand | od -c | grep -c '\\t'
```

Convert runs of spaces everywhere in the line with `-a`:

```sh
printf 'a       b\n' | unexpand -a | od -c | grep -c '\\t'
```

Use a tab stop every 4 columns with `-t`; this also turns on `-a`:

```sh
printf 'a   b\n' | unexpand -t 4 | od -c | grep -c '\\t'
```

### Option: -a
Convert runs of spaces everywhere in a line, not only at the start.

### Option: -t
Set the distance between tab stops, in columns, as in `-t 4`.

### Common mistakes
Only runs of spaces that end exactly on a tab stop become a tab, and single spaces are left alone, so the result can look the same as the input. Because `-t` implies `-a` in the common case, a file with spaces inside text, such as in a string, will be changed. Do not run `unexpand` on source code where alignment spaces are inside strings.

### Portability
The examples read the result with `od` and count lines that hold a tab, because a tab is hard to see. Without `-a`, GNU and macOS convert only leading spaces, but BusyBox 1.37.0 also converted spaces inside the line in the check, so use `-a` or `-t` explicitly when you want that behavior on every system. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU unexpand
GNU unexpand is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### GNU examples
Leave spaces inside the line alone without `-a`:

```sh
printf 'a       b\n' | unexpand | od -c | grep -q '\\t' || echo no tab
```

### Option: --first-only
Convert only the first run of spaces on each line.

## Apple
### About macOS unexpand
This is the unexpand in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Leave spaces inside the line alone without `-a`:

```sh
printf 'a       b\n' | unexpand | od -c | grep -q '\\t' || echo no tab
```

### Apple notes
macOS `unexpand` has `-a` and `-t`. It does not have GNU's `--first-only`.

## BusyBox
### About BusyBox unexpand
BusyBox bundles unexpand as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `unexpand` supports `-a` and `-t`. In the 1.37.0 check, plain `unexpand` also converted spaces inside the line, which GNU and macOS leave alone without `-a`. Other options and GNU long options depend on the build.
