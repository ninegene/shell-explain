---
schema_version: 1
id: zcmp
title: zcmp
summary: Compare compressed files byte by byte.
aliases: []
tags: [files, compare, compression]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
---
## Common
### Overview
`zcmp` runs `cmp` on gzip-compressed files, decompressing them first, so it compares what is inside and not the compressed bytes. It takes the options of `cmp`, and its exit status is the same: 0 for the same, 1 for different.

### Background
The macOS manual describes `zcmp` and `zdiff` together, as tools to compare compressed files. Today the name covers the GNU gzip script `zcmp` and the `zcmp` on macOS. The BusyBox 1.37.0 build used for this page has no `zcmp` applet.

### Examples
Ask whether the contents of two compressed files are the same:

```sh
printf 'a\nb\n' | gzip > f.gz && printf 'a\nb\n' | gzip > h.gz && zcmp f.gz h.gz && echo same
```

Ask quietly with `-s`, as for `cmp`:

```sh
printf 'a\nb\n' | gzip > f.gz && printf 'a\nc\n' | gzip > g.gz && { zcmp -s f.gz g.gz || echo differ; }
```

Compare a compressed file with a plain one; the plain file is read as it is:

```sh
printf 'a\nb\n' | gzip > f.gz && printf 'a\nb\n' > plain.txt && zcmp f.gz plain.txt && echo same
```

### Common mistakes
Two compressed files of the same text can differ as files, because of the time stored in the header and the level used, so use `zcmp` and not `cmp` for this question. Without `-s`, the message names a temporary file, as in `- /tmp/zcmp.XXXX differ`, which differs between systems. `zcmp` with one file argument compares it with the file of the same name without the `.gz`.

### Portability
This page has no BusyBox view: BusyBox has no `zcmp` applet. The examples use `-s` and the exit status, since the message differs. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU gzip scripts in the Linux guest and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU zcmp
GNU zcmp is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS zcmp
This is the zcmp in `/usr/bin` on macOS, part of its gzip tools.
