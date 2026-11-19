---
schema_version: 1
id: truncate
title: truncate
summary: Set a file to an exact size, shrinking or extending it.
aliases: [gtruncate]
tags: [files, size]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/truncate-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=truncate&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnutruncate.1.html
---
## Common
### Overview
`truncate` changes the size of a file without opening it in an editor. A smaller size cuts off the end; a larger size adds zero bytes. It is used to empty a file (`-s 0`), to make a file of a known size for testing, and to cut a log down without removing it.

### Background
The `truncate` utility first appeared in FreeBSD 4.2, according to the macOS manual, which also says it conforms to no known standards. Today the name covers GNU truncate, the truncate that macOS ships, and the truncate applet in BusyBox. They agree on `-s` with a size.

### Examples
Make a file exactly 10 bytes long; the `ls -l` size column is checked with `awk`:

```sh
truncate -s 10 f.txt && ls -l f.txt | awk '{print $5}'
```

Empty a file but keep it:

```sh
truncate -s 0 notes.txt && ls -l notes.txt | awk '{print $5}'
```

Make a new file of zero bytes, then look at it:

```sh
truncate -s 5 new.bin && od -An -tx1 new.bin
```

Cut a file down to its first bytes:

```sh
truncate -s 3 notes.txt && cat notes.txt
```

### Option: -s
Set the size in bytes, as in `-s 10`.

### Common mistakes
Cutting a file with `truncate` loses the end of it for good, so check the size first. Making a file bigger fills it with zero bytes, not spaces. If the file does not exist, `truncate` creates it. Emptying a log that a program is writing to may leave a hole of zero bytes at the start when the program writes at its old position, unless it opened the file for appending.

### Portability
Only `-s` with a plain number of bytes is in the Common examples. Sizes with suffixes such as `10K`, and relative sizes such as `+5`, differ between implementations. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU truncate
GNU truncate is part of GNU coreutils, and this page's checks ran coreutils 9.11. Sizes accept suffixes such as `K` and `M`, and a leading `+` or `-` changes the size relatively.

### GNU examples
Grow a file by five bytes with a relative size:

```sh
truncate -s 5 f.txt && truncate -s +5 f.txt && ls -l f.txt | awk '{print $5}'
```

## Apple
### About macOS truncate
This is the truncate in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `truncate` accepts `-s` with a size, including a leading `+` or `-` for a relative change.

## BusyBox
### About BusyBox truncate
BusyBox bundles truncate as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `truncate` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
