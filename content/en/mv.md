---
schema_version: 1
id: mv
title: mv
summary: Move or rename files and directories.
aliases: [gmv]
tags: [files, directories]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/mv-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=mv&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`mv` renames a file or directory, or moves it into another directory. Within one file system it is a rename, so it is fast and keeps the contents. It overwrites an existing target file without asking.

### Background
A `mv` command appeared in Version 1 Unix, according to the FreeBSD and macOS manuals. Today the name covers GNU mv, the mv that macOS ships, and the mv applet in BusyBox. They agree on `-f`, and differ on no-clobber `-n`, update `-u`, and the target-directory option `-t`.

### Examples
Rename a file and show that the contents moved with it:

```sh
mv notes.txt renamed.txt && head -n 1 renamed.txt
```

Move a file into an existing directory:

```sh
mv app.log logs/ && ls logs
```

Move several files into a directory; the last name is the destination:

```sh short="Move several files into the last-named directory"
mv app.log notes.txt logs && ls logs
```

Rename a directory:

```sh
mv logs old-logs && ls old-logs
```

Replace an existing file without asking; `-f` suppresses a prompt:

```sh short="Overwrite without asking (-f)"
mv -f notes.txt app.log && head -n 1 app.log
```

### Option: -f
Overwrite the target without prompting, even if `-i` appeared earlier.

### Common mistakes
`mv a b` overwrites `b` silently, and when `b` is an existing directory it moves `a` inside it instead of renaming. Moving across file systems copies and then deletes, so an interrupted move can leave both. Use `-i` to be asked before overwriting.

### Portability
Only `-f` is in the Common examples. No-clobber `-n`, update `-u` and the target-directory option `-t` differ or are missing between implementations. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU mv
GNU mv is part of GNU coreutils, and this page's checks ran coreutils 9.11. It adds `-t` to name the destination directory first.

### GNU examples
Name the destination directory first with `-t`. This option is a GNU extension:

```sh
mv -t logs notes.txt app.log && ls logs
```

## Apple
### About macOS mv
This is the mv in `/bin` on macOS. It has no `-t`, so the destination goes last.

### Apple notes
macOS `mv` is the BSD implementation. It has no `-t`; put the destination last.

## BusyBox
### About BusyBox mv
BusyBox bundles mv as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `mv` supports `-f` and `-i`; GNU long options depend on the build.
