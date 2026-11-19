---
schema_version: 1
id: rm
title: rm
summary: Remove files and directories.
aliases: [grm]
tags: [files, directories]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/rm-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=rm&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/rm.1.html
---
## Common
### Overview
`rm` deletes the files you name. It does not move them to a trash, so a deleted file is gone. Directories need `-r`.

### Background
An `rm` command appeared in Version 1 Unix, according to the FreeBSD and macOS manuals. Today the name covers GNU rm, the rm that macOS ships, and the rm applet in BusyBox. They agree on `-r` and `-f`, and differ on interactive `-i`, verbose `-v`, and the empty-directory option `-d`.

### Examples
Delete one file and list what is left:

```sh
rm app.log && ls
```

Delete several files:

```sh
rm app.log notes.txt && ls
```

Delete a directory and everything inside it:

```sh
rm -r logs && ls
```

Ignore a file that does not exist, with no error and a zero exit status:

```sh short="Ignore a file that does not exist (-f)"
rm -f missing.txt && echo done
```

### Option: -r
Remove directories and their contents recursively. `-R` is accepted too.

### Option: -f
Never prompt, and do not complain about files that do not exist.

### Common mistakes
There is no undo. Check a pattern with `ls` before you pass it to `rm`, quote or escape names that contain spaces, and be careful with variables that may be empty: `rm -rf $DIR/` with an empty `DIR` aims at `/`. Use `--` before names that start with a dash, as in `rm -- -file`.

### Portability
Only `-r` and `-f` are in the Common examples. Interactive `-i`, verbose `-v` and the empty-directory option `-d` behave differently or are missing between implementations. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU rm
GNU rm is part of GNU coreutils, and this page's checks ran coreutils 9.11. It accepts `-I`, which asks once before removing many files or recursing.

### GNU examples
Ask once before removing more than three files or recursing; `-I` is a safer habit than `-f`. GNU and macOS accept it, but BusyBox may not. The question is read from the terminal, so only the quiet case is shown:

```sh
rm -I app.log && ls
```

## Apple
### About macOS rm
This is the rm in `/bin` on macOS. It also accepts `-I`, and adds `-P` to overwrite files before deleting and `-W` to undelete on file systems that support it, which GNU does not have.

### Apple notes
macOS `rm` is the BSD implementation. It also accepts `-I`, and adds `-P` to overwrite files before deleting and `-W` to undelete on file systems that support it, which GNU does not have.

## BusyBox
### About BusyBox rm
BusyBox bundles rm as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `rm` supports `-r`, `-f` and `-i`; GNU long options depend on the build.
