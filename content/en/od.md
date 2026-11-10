---
schema_version: 1
id: od
title: od
summary: Dump the bytes of a file as octal, hexadecimal, decimal or characters.
aliases: [god]
tags: [files, binary, debugging]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/od-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=od&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`od` (octal dump) shows what is really in a file: every byte, in a form you choose. It is the tool to reach for when a file contains characters you cannot see, such as tabs, carriage returns or non-text data. Each output line starts with the position of its first byte, which you can hide with `-An`.

### Background
An `od` command appeared in Version 1 AT&T UNIX, according to the macOS manual, which also says it conforms to POSIX.1-2001. Today the name covers GNU od, the od that macOS ships, and the od applet in BusyBox. They agree on `-c`, `-t`, `-A` and `-An`; the exact column spacing differs, so the examples compare output with runs of spaces collapsed.

### Examples
Show the characters of some text, one per column, without the position column:

```sh
printf abc | od -An -c
```

Show the bytes as hexadecimal numbers with `-t x1`, which is one byte per number:

```sh
printf abc | od -An -tx1
```

Show the bytes as decimal numbers:

```sh
printf abc | od -An -td1
```

Show a newline as `\\n` and keep the position column:

```sh
printf 'A\n' | od -c | head -n 1
```

Find a hidden tab in a line:

```sh
printf 'a\tb\n' | od -An -c
```

### Option: -c
Show each byte as a character, using escapes such as `\\n` and `\\t` for those that cannot be seen.

### Option: -t
Choose the type and size: `x1` is hexadecimal bytes, `d1` is decimal bytes, `c` is characters.

### Option: -An
Do not print the position column.

### Common mistakes
The spacing of `od` output is not the same on every system, so do not compare it character by character in a script; use `-An` and split on whitespace. Multi-byte formats such as `x2` depend on the byte order of the machine. `xxd` and `hexdump -C` show the same bytes next to their characters, and many people find them easier to read, but they are not on every system.

### Portability
Only `-c`, `-t x1`, `-t d1`, `-A` and `-An` are in the Common examples, with spaces collapsed in the comparison. The type `a` for named characters and the option `-v` for repeated lines also exist on all three, with small differences. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU od
GNU od is part of GNU coreutils, and this page's checks ran coreutils 9.11. It also has `-j` to skip bytes and `-N` to stop after some.

### GNU examples
Look at only the first two bytes with `-N`:

```sh
printf abc | od -An -tx1 -N 2
```

### Option: -N, --read-bytes
Read only this many bytes.

## Apple
### About macOS od
This is the od in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `od` has `-j` to skip bytes and `-N` to stop after some, like GNU.

## BusyBox
### About BusyBox od
BusyBox bundles od as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `od` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
