---
schema_version: 1
id: cksum
title: cksum
summary: Print a CRC checksum and the byte count of files.
aliases: [gcksum]
tags: [files, checksum]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/cksum-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=cksum&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnucksum.1.html
---
## Common
### Overview
`cksum` prints a number computed from the contents of a file, then its length in bytes, then its name. The number is a cyclic redundancy check (CRC), good for catching accidental damage in transfer, but not for security. The POSIX algorithm is the same everywhere, so the output matches across systems.

### Background
The `cksum` utility appeared in 4.4BSD, according to the macOS manual, which also says it is expected to conform to POSIX.2. Today the name covers GNU cksum, the cksum that macOS ships, and the cksum applet in BusyBox. They agree on the default algorithm and output.

### Examples
Print the checksum and length of some text:

```sh
printf abc | cksum
```

Print the checksum, length and name of a file:

```sh
printf abc > f.txt && cksum f.txt
```

Compare two files by checksum; the same content gives the same line apart from the name:

```sh
printf abc > a.txt && printf abc > b.txt && cksum a.txt b.txt | cut -d' ' -f1,2 | uniq | grep -c .
```

Show that a one-byte change gives a different checksum:

```sh
printf abc | cksum | cut -d' ' -f1; printf abd | cksum | cut -d' ' -f1
```

### Common mistakes
A CRC is easy to forge on purpose, so do not use `cksum` where someone may tamper with the file; use `sha256sum`. Two different files can have the same CRC, though that is rare by accident. The number is not a hash in the `md5sum` sense and it is printed in decimal.

### Portability
Only the default POSIX algorithm is in the Common examples. GNU 9 adds `-a` to choose other algorithms; macOS has `-o` for historic algorithms instead. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU cksum
GNU cksum is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-a` selects another algorithm, such as `sha256`.

### GNU examples
Print a SHA-256 digest with `-a`; GNU `cksum` prints it in a tagged line:

```sh
printf abc | cksum -a sha256 | cut -c1-30
```

### Option: -a, --algorithm
Choose the checksum algorithm, such as `crc`, `md5` or `sha256`.

## Apple
### About macOS cksum
This is the cksum that macOS ships. It follows the BSD implementation.

### Apple notes
macOS `cksum` has `-o 1`, `-o 2` and `-o 3` for historic algorithms, and no `-a`.

## BusyBox
### About BusyBox cksum
BusyBox bundles cksum as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `cksum` prints the CRC, the byte count and the name. Other options depend on the build.
