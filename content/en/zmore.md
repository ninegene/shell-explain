---
schema_version: 1
id: zmore
title: zmore
summary: View a compressed file one screen at a time, with the simple more pager.
aliases: []
tags: [text, compression, pager]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/zmore.1.html
---
## Common
### Overview
`zmore` shows a compressed file one screen at a time with the `more` pager. It is the older and simpler partner of `zless`: it moves forward by a screen with the space bar, and most versions of `more` cannot scroll back. When output is not a terminal it prints everything.

### Background
The macOS manual describes `zmore` and `zless` together, as tools to view compressed files. Today the name covers the GNU gzip script `zmore` and the `zmore` on macOS. The BusyBox 1.37.0 build used for this page has no `zmore` applet.

### Examples
View a compressed file at a terminal; the space bar shows the next screen and `q` quits:

```sh short="View a compressed file (space: next screen, q quits)"
zmore access.log.gz
```

Put it into a pipe, where it prints everything:

```sh
printf 'a\nb\n' | gzip > f.gz && zmore f.gz | cat
```

### Common mistakes
`zmore` is rarely the better choice: `zless` can go back and search. It needs the `more` program. As with the other `z` tools, a file that is not compressed is shown as plain text on GNU.

### Portability
This page has no BusyBox view: BusyBox has no `zmore` applet. The first example needs a terminal and was reviewed only; the pipe example was executed on GNU gzip scripts in the Linux guest and macOS system tools. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU coreutils 9.11, and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU zmore
GNU zmore is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS zmore
This is the zmore in `/usr/bin` on macOS, part of its gzip tools.
