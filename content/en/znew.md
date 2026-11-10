---
schema_version: 1
id: znew
title: znew
summary: Recompress .Z files from the old compress tool as gzip files.
aliases: []
tags: [files, compression]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
---
## Common
### Overview
`znew` converts files made by the old `compress` tool (names ending in `.Z`) to gzip files (ending in `.gz`), which are usually smaller. It checks the new file before removing the old one. It matters only if you still have `.Z` files, which is rare today.

### Background
The macOS manual describes `znew` as converting compressed files to gzipped files. Today the name covers the GNU gzip script `znew` and the `znew` on macOS. The BusyBox 1.37.0 build used for this page has no `znew` applet. The examples on this page were reviewed against the manuals but not executed, because making a `.Z` file needs the `compress` program, which is on macOS but not on the Linux guest.

### Examples
Convert one `.Z` file to a `.gz` file; the `.Z` file is replaced:

```sh
znew old.Z
```

Convert several files and keep the `.Z` originals with `-K`:

```sh
znew -K one.Z two.Z
```

### Option: -K
Keep the `.Z` file if it is smaller than the `.gz` file.

### Common mistakes
`znew` only handles files that end in `.Z`. It needs both `gzip` and a program that can read the old format, so it can fail on a system without it. Check the result before deleting backups.

### Portability
This page has no BusyBox view: BusyBox has no `znew` applet. The examples were reviewed only. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway.

## GNU
### About GNU znew
GNU znew is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS znew
This is the znew in `/usr/bin` on macOS, part of its gzip tools.
