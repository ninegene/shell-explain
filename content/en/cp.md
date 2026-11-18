---
schema_version: 1
id: cp
title: cp
summary: Copy files and directories.
aliases: [gcp]
tags: [files, directories]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/cp-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=cp&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`cp` copies files. With two names it copies the first to the second; with more, the last must be an existing directory that receives the copies. It overwrites an existing target file without asking.

### Background
A `cp` command appeared in Version 1 Unix, according to the FreeBSD and macOS manuals. Today the name covers GNU cp, the cp that macOS ships, and the cp applet in BusyBox. They agree on `-R` and `-p`, and differ on archive mode `-a`, no-clobber `-n`, and update `-u`.

### Examples
Copy a file to a new name and show the copy:

```sh
cp notes.txt copy.txt && cat copy.txt
```

Copy a file into an existing directory, keeping its name:

```sh
cp app.log logs/ && ls logs
```

Copy several files into a directory; the last name is the destination:

```sh short="Copy several files into the last-named directory"
cp notes.txt app.log logs && ls logs
```

Copy a directory and everything inside it. `-R` is the portable spelling of recursive copy:

```sh short="Copy a directory and its contents (-R, portable)"
cp -R logs logs-copy && ls logs-copy
```

Keep the modification time and permissions of the original:

```sh short="Keep modification time and permissions (-p)"
cp -p notes.txt kept.txt && ls kept.txt
```

### Option: -R
Copy directories recursively. `-r` is accepted by all three implementations too, but `-R` is the one POSIX defines.

### Option: -p
Preserve the modification time, access time and permission bits where possible.

### Common mistakes
`cp dir other` fails without `-R`, and `cp -R logs logs-copy` creates `logs-copy/server.log` the first time but `logs-copy/logs/server.log` if `logs-copy` already exists. A trailing slash on a source directory means different things on GNU and BSD, so avoid it when copying directories. `cp` silently overwrites; `-i` asks first on a terminal.

### Portability
Only `-R` and `-p` are in the Common examples. Archive mode `-a`, no-clobber `-n` and update `-u` differ in meaning or availability between GNU, BSD and BusyBox. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU cp
GNU cp is part of GNU coreutils, and this page's checks ran coreutils 9.11. It adds `-t` to name the destination directory first.

### GNU examples
Name the destination directory first with `-t`, which suits `xargs` pipelines. This option is a GNU extension:

```sh
cp -t logs notes.txt app.log && ls logs
```

## Apple
### About macOS cp
This is the cp in `/bin` on macOS. It has no `-t`, so the destination goes last. It adds `-c` to clone files on APFS, which GNU and BusyBox do not have.

### Apple notes
macOS `cp` is the BSD implementation. It has no `-t`; put the destination last. It adds `-c` to clone files on APFS, which GNU and BusyBox do not have.

## BusyBox
### About BusyBox cp
BusyBox bundles cp as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `cp` supports `-R`, `-p` and `-a`; other options and GNU long spellings depend on the build.
