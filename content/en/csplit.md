---
schema_version: 1
id: csplit
title: csplit
summary: Cut a file into pieces at lines that match a pattern or a line number.
aliases: [gcsplit]
tags: [files, text]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/csplit-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=csplit&sektion=1
---
## Common
### Overview
`csplit` (context split) cuts a file where you tell it to: at a line number or at the first line that matches a pattern. The pieces are named `xx00`, `xx01`, and so on. Where `split` cuts by size, `csplit` cuts by content, which suits files made of sections.

### Background
A `csplit` command appeared in PWB UNIX, according to the macOS manual, which also says it conforms to POSIX.1-2001. Today the name covers GNU csplit and the csplit that macOS ships. The BusyBox 1.37.0 build used for this page has no `csplit` applet. They agree on line-number and `/pattern/` arguments and on `-s` and `-f`.

### Examples
Cut before line 4: lines 1 to 3 go in `xx00` and the rest in `xx01`; `-s` hides the byte counts:

```sh
seq 6 > n6.txt && csplit -s n6.txt 4 && paste -sd, xx00 && paste -sd, xx01
```

Choose the prefix of the file names with `-f`:

```sh
seq 6 > n6.txt && csplit -s -f part_ n6.txt 4 && ls part_*
```

Cut before the first line that matches a pattern:

```sh
printf 'a\nb\nc\nb\n' > words.txt && csplit -s words.txt '/c/' && paste -sd, xx00 && paste -sd, xx01
```

### Option: -s
Do not print the size of each piece.

### Option: -f
Use this prefix instead of `xx` for the names of the pieces.

### Common mistakes
Without `-s`, `csplit` prints the size in bytes of every piece, which can look like noise in a script. If the pattern is not found, `csplit` fails and may leave pieces behind. The pattern is a regular expression and a line number counts from 1. The first piece ends just before the line you name, so the matching line starts the next piece.

### Portability
This page has no BusyBox view: BusyBox has no `csplit` applet. The `{*}` repeat count and `--suppress-matched` are GNU features; macOS rejects `{*}` and needs a number such as `{1}`. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check). The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU csplit
GNU csplit is part of GNU coreutils, and this page's checks ran coreutils 9.11. `{*}` repeats the pattern for as long as it matches.

### GNU examples
Cut at every line that matches, with a repeat count of `{*}`:

```sh
printf 'a\nb\na\nb\n' > ab.txt && csplit -s ab.txt '/a/' '{*}' && ls xx* | grep -c .
```

### Option: --suppress-matched
Do not include the matching line in either piece.

## Apple
### About macOS csplit
This is the csplit in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `csplit` has `-s`, `-f`, `-k` and `-n`. It rejects the repeat count `{*}`; give a number such as `{1}` instead. It has no `--suppress-matched`.
