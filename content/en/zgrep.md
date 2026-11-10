---
schema_version: 1
id: zgrep
title: zgrep
summary: Search compressed files for a pattern without unpacking them.
aliases: [zegrep, zfgrep]
tags: [text, search, compression]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
---
## Common
### Overview
`zgrep` runs `grep` on gzip-compressed files, decompressing them on the fly, so you can search a rotated log such as `app.log.2.gz` without unpacking it on disk. It takes the options of `grep`. `zegrep` and `zfgrep` are the same with `grep -E` and `grep -F`, and this page answers for those names too.

### Background
Today the name covers the GNU gzip script `zgrep`, and the `zgrep` on macOS, which the macOS manual lists in the same page as `grep`, `egrep` and `fgrep`. The BusyBox 1.37.0 build used for this page has no `zgrep` applet. They agree on passing grep options through and on reading compressed files by name.

### Examples
Search a compressed file for a word:

```sh
printf 'a\nb\n' | gzip > f.gz && zgrep b f.gz
```

Count the matching lines with `-c`:

```sh
printf 'a\nb\n' | gzip > f.gz && zgrep -c a f.gz
```

Ignore case with `-i`:

```sh
printf 'a\nb\n' | gzip > f.gz && zgrep -i B f.gz
```

Use an extended regular expression through the `zegrep` name:

```sh
printf 'a\nb\n' | gzip > f.gz && zegrep 'a|b' f.gz | paste -sd, -
```

Search for a fixed string through the `zfgrep` name:

```sh
printf 'a\nb\n' | gzip > f.gz && zfgrep a f.gz
```

### Common mistakes
`zgrep` reads the whole file every time and is slower than `grep` on the same text, so unpack a file you will search many times. With several files it prefixes each match with the file name, as `grep` does. It handles gzip files, and compressed files of other kinds (`.bz2`, `.xz`) need their own tools, such as `bzgrep` and `xzgrep`. A file that is not compressed is searched as plain text on GNU.

### Portability
This page has no BusyBox view: BusyBox has no `zgrep` applet. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU gzip scripts in the Linux guest and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU zgrep
GNU zgrep is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS zgrep
This is the zgrep in `/usr/bin` on macOS, part of its gzip tools.
