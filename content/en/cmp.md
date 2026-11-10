---
schema_version: 1
id: cmp
title: cmp
summary: Compare two files byte by byte.
aliases: [gcmp]
tags: [files, compare]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/diffutils/manual/html_node/Invoking-cmp.html
  - https://man.freebsd.org/cgi/man.cgi?query=cmp&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`cmp` compares two files byte by byte and, if they differ, tells where the first difference is. It works on any kind of file, including binary ones, where `diff` would not help. Its exit status is 0 when the files are the same, 1 when they differ and 2 on an error, so `cmp -s` is the quiet way to ask "are these the same?".

### Background
Today the name covers GNU cmp, the cmp that macOS ships, and the cmp applet in BusyBox. They agree on `-s` for silent and on the exit statuses.

### Examples
Ask quietly whether two files are the same; the answer is in the exit status:

```sh
printf 'a\nb\n' > f1 && cp f1 f2 && cmp -s f1 f2 && echo same
```

Run something when they differ:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && { cmp -s f1 f2 || echo differ; }
```

Compare a file with itself:

```sh
printf 'a\n' > f1 && cmp -s f1 f1 && echo same
```

Compare files of different lengths; one is a prefix of the other, so they differ:

```sh
printf 'a\n' > f1 && printf 'a\nb\n' > f2 && { cmp -s f1 f2 || echo differ; }
```

### Option: -s
Print nothing; only set the exit status.

### Common mistakes
`cmp` reports the first difference only, and without `-s` it prints a message that differs between implementations, so scripts should use `-s` and the exit status. A file that is a prefix of the other differs. `cmp` reads both files completely when they are the same, so it can be slow on big files; for that case a checksum may be easier to store. For text files, `diff` shows all the differences.

### Portability
Only `-s` and the exit status are in the Common examples. The message text and the `-l` list of differing bytes are described in the GNU and Apple sections. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU cmp
GNU cmp is part of GNU diffutils, and this page's checks ran version 3.12. `-l` lists every differing byte, and `-n LIMIT` stops after that many bytes.

### GNU examples
Without `-s`, `cmp` says where the first difference is:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && cmp f1 f2; test $? -eq 1 && echo differ
```

### Option: -l, --verbose
List the position and value of every differing byte.

## Apple
### About macOS cmp
This is the cmp in the system folders of macOS.

### Apple examples
Without `-s`, `cmp` says where the first difference is:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && cmp f1 f2; test $? -eq 1 && echo differ
```

### Option: -l
List the position and value of every differing byte.

## BusyBox
### About BusyBox cmp
BusyBox bundles cmp as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `cmp` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
