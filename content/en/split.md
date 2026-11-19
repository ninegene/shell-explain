---
schema_version: 1
id: split
title: split
summary: Cut a file into pieces.
aliases: [gsplit]
tags: [files, text]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/split-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=split&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnusplit.1.html
---
## Common
### Overview
`split` cuts a file into smaller files, by number of lines (`-l`) or by size (`-b`). The pieces are named with a prefix you choose and a suffix that counts up: `part_aa`, `part_ab`, and so on. To put them back together, `cat part_*` joins them in order.

### Background
A `split` command appeared in Version 3 AT&T UNIX, according to the macOS manual, which also says it conforms to POSIX.1-2001. Today the name covers GNU split, the split that macOS ships, and the split applet in BusyBox. They agree on `-l`, `-b`, the prefix argument and the `aa` suffixes.

### Examples
Split a file into pieces of two lines each; `part_` is the prefix of the new files:

```sh
seq 6 > n6.txt && split -l 2 n6.txt part_ && ls part_*
```

Read one piece:

```sh
seq 6 > n6.txt && split -l 2 n6.txt part_ && cat part_aa
```

Join the pieces back together in order:

```sh
seq 6 > n6.txt && split -l 2 n6.txt part_ && cat part_* | paste -sd, -
```

Split by size with `-b`; here every piece has four bytes, and `-` reads standard input:

```sh
printf abcdefgh | split -b 4 - chunk_ && cat chunk_aa
```

### Option: -l
Put this many lines in each piece.

### Option: -b
Put this many bytes in each piece.

### Common mistakes
If you leave out the prefix, the pieces are named `xaa`, `xab` and so on, in the current folder, which is easy to miss. `split` never changes the original file. The suffix has two letters by default, so a file that needs more than 676 pieces needs a longer suffix (`-a` in GNU) or it stops with an error.

### Portability
Only `-l`, `-b` and a prefix are in the Common examples. `-d` for numeric suffixes works on GNU and macOS but was rejected by BusyBox 1.37.0. GNU adds `--filter` and `--additional-suffix`, which macOS does not have. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU split
GNU split is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-n` splits into a given number of pieces, and `--additional-suffix` adds an extension; macOS has `-n` too.

### GNU examples
Use numbers instead of letters in the names with `-d`:

```sh
seq 6 > n6.txt && split -d -l 2 n6.txt p_ && ls p_*
```

Split into exactly three pieces with `-n`:

```sh
seq 6 > n6.txt && split -n 3 n6.txt part_ && ls part_* | grep -c .
```

### Option: -d
Use numeric suffixes `00`, `01` instead of `aa`, `ab`.

### Option: -n, --number
Split into this many pieces of about equal size.

## Apple
### About macOS split
This is the split in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Use numbers instead of letters in the names with `-d`:

```sh
seq 6 > n6.txt && split -d -l 2 n6.txt p_ && ls p_*
```

### Option: -d
Use numeric suffixes `00`, `01` instead of `aa`, `ab`.

### Apple notes
macOS `split` has `-l`, `-b`, `-d`, `-a` and `-n`, with the same meaning as in GNU.

## BusyBox
### About BusyBox split
BusyBox bundles split as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `split` supports `-l`, `-b` and a prefix. In the 1.37.0 check it rejected `-d`, so numeric suffixes are not available; other options depend on the build.
