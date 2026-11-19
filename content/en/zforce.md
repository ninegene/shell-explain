---
schema_version: 1
id: zforce
title: zforce
summary: Add the .gz suffix to gzip files that lack it.
aliases: []
tags: [files, compression]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/zforce.1.html
---
## Common
### Overview
`zforce` renames every gzip file that does not end in `.gz` so that it does. The point is to stop `gzip` from compressing a file twice: `gzip` skips a name that ends in `.gz`, but compresses a gzip file with another name again. It looks inside the file, so only real gzip data is renamed.

### Background
The macOS manual describes `zforce` as forcing gzip files to have a `.gz` suffix. Today the name covers the GNU gzip script `zforce` and the `zforce` on macOS. The BusyBox 1.37.0 build used for this page has no `zforce` applet. They agree on the rename; the message printed differs.

### Examples
Rename a gzip file that has no suffix:

```sh
printf data | gzip > file && zforce file > /dev/null && ls file.gz
```

Leave a file that is not gzip data alone:

```sh
printf plain > file && zforce file > /dev/null && ls file
```

### Common mistakes
`zforce` changes names in place, so scripts that expect the old name break; list what it will touch first. A name that already ends in `.gz` is left alone. If the new name already exists, the file is not renamed. The message it prints ("renamed" on macOS, "replaced" on GNU in the check) differs, so do not read it in a script.

### Portability
This page has no BusyBox view: BusyBox has no `zforce` applet. The examples hide the message. The GNU versions are shell scripts, which the implementation probe never runs, so they are not identified automatically; the Linux guest ran them anyway. The Common examples were executed on GNU gzip scripts in the Linux guest and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU zforce
GNU zforce is part of GNU gzip, and it is a shell script that calls `gzip`.

## Apple
### About macOS zforce
This is the zforce in `/usr/bin` on macOS, part of its gzip tools.
