---
schema_version: 1
id: zdiff
title: zdiff
summary: Show the differences between compressed files.
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
`zdiff` runs `diff` on gzip-compressed files, decompressing them first, so you see the differences in the text. It takes the options of `diff`, such as `-u`, and gives the same exit status.

### Background
The macOS manual describes `zcmp` and `zdiff` together, as tools to compare compressed files. Today the name covers the GNU gzip script `zdiff` and the `zdiff` on macOS. The BusyBox 1.37.0 build used for this page has no `zdiff` applet.

### Examples
Show the differences in the default format; the status is 1 when they differ:

```sh
printf 'a\nb\n' | gzip > f.gz && printf 'a\nc\n' | gzip > g.gz && zdiff f.gz g.gz; test $? -eq 1 && echo differ
```

Say only whether they differ, with `-q`:

```sh
printf 'a\nb\n' | gzip > f.gz && printf 'a\nb\n' | gzip > h.gz && zdiff -q f.gz h.gz && echo same
```

### Common mistakes
A status of 1 means "different", as with `diff`, so a script with `set -e` stops. The file names in a unified diff (`-u`) may show temporary names for the decompressed text on some systems. Compressed files that hold the same text are the same for `zdiff`, even if the compressed bytes differ.

### Portability
This page has no BusyBox view: BusyBox has no `zdiff` applet. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU gzip scripts in the Linux guest and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU zdiff
GNU zdiff is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS zdiff
This is the zdiff in `/usr/bin` on macOS, part of its gzip tools.
