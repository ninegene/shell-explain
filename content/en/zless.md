---
schema_version: 1
id: zless
title: zless
summary: Page through a compressed file.
aliases: []
tags: [text, compression, pager]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
---
## Common
### Overview
`zless` shows a compressed file one screen at a time, like `less`, decompressing it as you read. You can scroll up and down and search with `/`. When its output is not a terminal, for example in a pipe, it prints the whole text, like `zcat`. It is meant for a person at a terminal, so its examples show the pipe form.

### Background
The macOS manual describes `zmore` and `zless` together, as tools to view compressed files. Today the name covers the GNU gzip script `zless` and the `zless` on macOS. The BusyBox 1.37.0 build used for this page has no `zless` applet.

### Examples
Look at a compressed file in the pager; at a terminal you can scroll and press `q` to quit:

```sh short="Look at a compressed file in the pager (q quits)"
zless access.log.gz
```

Put it into a pipe, where it prints everything, to show what the pager would show:

```sh
printf 'a\nb\n' | gzip > f.gz && zless f.gz | cat
```

### Common mistakes
`zless` needs the `less` program, which most systems have. A file that is not compressed is shown as plain text on GNU. It cannot show a file that is still being written to the end, as `less +F` does for plain files. Binary data in the file can mess up the terminal; press `q` and run `reset` if that happens.

### Portability
This page has no BusyBox view: BusyBox has no `zless` applet. The first example needs a terminal and was reviewed only; the pipe example was executed on GNU gzip scripts in the Linux guest and macOS system tools. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU coreutils 9.11, and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU zless
GNU zless is part of GNU gzip, and it is a shell script that calls `gzip`. It runs `less` on the decompressed text.

## Apple
### About macOS zless
This is the zless in `/usr/bin` on macOS, part of its gzip tools.
