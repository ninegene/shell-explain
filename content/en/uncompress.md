---
schema_version: 1
id: uncompress
title: uncompress
summary: Restore files made by the old compress tool.
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
`uncompress` restores files that the old `compress` tool made, which end in `.Z`. The `.Z` file is replaced by the restored one. It is an old format, kept for compatibility, and `gzip` replaced it for new files. On GNU systems, `gzip -d` can also read `.Z` files.

### Background
Today the name covers the `uncompress` of GNU gzip and the `uncompress` on macOS, which is the BSD `compress` tool under another name. The BusyBox 1.37.0 build used for this page has no `uncompress` applet. They agree on restoring a `.Z` file. The examples on this page were reviewed against the manuals but not executed, because making a `.Z` file needs the `compress` program, which is on macOS but not on the Linux guest.

### Examples
Restore a file; `old.txt.Z` becomes `old.txt`:

```sh
uncompress old.txt.Z
```

Print the restored text and keep the `.Z` file, with `-c`:

```sh
uncompress -c old.txt.Z
```

### Option: -c
Write the restored text to standard output and keep the file.

### Common mistakes
`uncompress` handles only the old `.Z` format, not `.gz`; use `gunzip` for those. The `.Z` file is removed after success unless you give `-c`. Few tools still produce `.Z` files, so most readers will meet them only in old archives.

### Portability
This page has no BusyBox view: BusyBox has no `uncompress` applet. The examples were reviewed only.

## GNU
### About GNU uncompress
GNU uncompress comes with GNU gzip and reads the old `.Z` format.

## Apple
### About macOS uncompress
This is the `uncompress` in `/usr/bin` on macOS, part of the BSD `compress` tools.
