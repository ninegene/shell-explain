---
schema_version: 1
id: rm
title: rm
summary: Remove files and directories.
aliases: [grm]
tags: [files, directories]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/rm-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=rm&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`rm` deletes the files you name. It does not move them to a trash, so a deleted file is gone. Directories need `-r`.

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

```sh
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
### GNU examples
Ask once before removing more than three files or recursing; `-I` is a safer habit than `-f`. GNU and macOS accept it, but BusyBox may not. The question is read from the terminal, so only the quiet case is shown:

```sh
rm -I app.log && ls
```

## Apple
### Apple notes
macOS `rm` is the BSD implementation. It also accepts `-I`, and adds `-P` to overwrite files before deleting and `-W` to undelete on file systems that support it, which GNU does not have.

## BusyBox
### BusyBox notes
BusyBox `rm` supports `-r`, `-f` and `-i`; GNU long options depend on the build.
