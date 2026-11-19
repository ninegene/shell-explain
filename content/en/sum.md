---
schema_version: 1
id: sum
title: sum
summary: Print an old-style checksum and block count of files.
aliases: [gsum]
tags: [files, checksum]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/sum-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=sum&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnusum.1.html
---
## Common
### Overview
`sum` prints a 16-bit checksum of a file and the number of 1024-byte blocks it takes. It is an old command, kept for compatibility. For new uses, `cksum` or `sha256sum` are better, because a 16-bit sum misses many changes.

### Background
Today the name covers GNU sum, the sum that macOS ships, and the sum applet in BusyBox. The macOS manual has no dated history for it. All three use the BSD checksum algorithm by default, so they agree on the number.

### Examples
Print the checksum and block count for some text:

```sh
printf abc | sum
```

Print them for a file; `awk` keeps just the two numbers, because the file name is not printed on every system:

```sh
printf abc > f.txt && sum f.txt | awk '{print $1, $2}'
```

Compare two files by checksum:

```sh
printf abc > a.txt && printf abc > b.txt && test "$(sum a.txt | cut -d' ' -f1)" = "$(sum b.txt | cut -d' ' -f1)" && echo same
```

### Common mistakes
A 16-bit checksum has only 65,536 values, so two different files will often match by chance; do not use `sum` to prove files are equal. The spacing of the output is not the same on all systems, so split on whitespace. The block count is rounded up, so a 3-byte file counts as 1.

### Portability
Only the default BSD algorithm is in the Common examples. GNU also has `-s` for the System V algorithm, which gives different numbers. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU sum
GNU sum is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-s` uses the System V algorithm instead.

### Option: -s, --sysv
Use the System V checksum algorithm.

## Apple
### About macOS sum
This is the sum in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `sum` prints the BSD checksum and block count, and does not have `--sysv`.

## BusyBox
### About BusyBox sum
BusyBox bundles sum as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `sum` prints the same two numbers. In the 1.37.0 check it did not print the file name after them for a single file, which GNU and macOS do.
