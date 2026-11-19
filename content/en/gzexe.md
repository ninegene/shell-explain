---
schema_version: 1
id: gzexe
title: gzexe
summary: Compress a program so that it unpacks itself each time it runs.
aliases: []
tags: [files, compression, programs]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/gzexe.1.html
---
## Common
### Overview
`gzexe` replaces an executable with a compressed copy that decompresses itself when it is run. The program still works as before, but takes less disk space and starts a little slower. It keeps the original under the same name with a `~` added. It is of use only on small disks, and it is not suited to modern systems with shared libraries.

### Background
The macOS manual says `gzexe` uses `gzip` to compress executables, producing executables that decompress on the fly. Today the name covers the GNU gzip script `gzexe` and the `gzexe` on macOS. The BusyBox 1.37.0 build used for this page has no `gzexe` applet.

### Examples
Compress a small script; it still runs:

```sh
printf '#!/bin/sh\necho hi\n' > prog && chmod +x prog && gzexe prog > /dev/null 2>&1 && ./prog
```

The original is kept as `prog~`:

```sh
printf '#!/bin/sh\necho hi\n' > prog && chmod +x prog && gzexe prog > /dev/null 2>&1 && ls prog*
```

### Common mistakes
A tiny file gets bigger, as the message `-11.1%` in the check shows, because the unpacking code is longer than what it saves. The compressed program needs `gzip` at run time, and each start unpacks it again into a temporary file, which slows it down. It does not work on programs that cannot run from a temporary file. The `-d` option that restores the original did not work on the macOS release tested.

### Portability
This page has no BusyBox view: BusyBox has no `gzexe` applet. The examples hide the messages, which differ. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU gzip scripts in the Linux guest and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU gzexe
GNU gzexe is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS gzexe
This is the gzexe in `/usr/bin` on macOS, part of its gzip tools.
