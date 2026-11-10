---
schema_version: 1
id: dirname
title: dirname
summary: Print the folder part of a path.
aliases: [gdirname]
tags: [files, paths]
implementations: [gnu, apple, busybox]
updated: 2026-11-09
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/dirname-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=dirname&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`dirname` removes the last part of a path and prints the folder that is left. It is the partner of `basename`. It works on the text of the path only, so the file does not have to exist.

### Background
The `basename` and `dirname` utilities first appeared in 4.4BSD, according to the macOS manual. Today the name covers GNU dirname, the dirname that macOS ships, and the dirname applet in BusyBox. They agree on how to treat a single path.

### Examples
Print the folder of a path:

```sh
dirname logs/server.log
```

A name with no folder gives the current folder, written as a dot:

```sh
dirname server.log
```

An absolute path:

```sh
dirname /usr/local/bin/tool
```

Go up one level from a path that ends in a slash:

```sh
dirname /usr/local/bin/
```

Use the folder of a file, for example to create it first:

```sh
mkdir -p "$(dirname cache/data/out.txt)" && ls cache
```

### Common mistakes
`dirname` removes the last part even when it is a folder, so `dirname logs` prints `.` and `dirname /usr/local/bin/` prints `/usr/local`. It never checks the file system, so a path that does not exist still gives an answer. Quote the argument and the command substitution around a path with spaces. In a shell script, `${path%/*}` does the same without a process, but it does not give a dot for a name with no folder.

### Portability
Only the single-path form is in the Common examples. POSIX describes one path; several paths at once are a GNU and BSD feature that is only shown in the GNU example below. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU dirname
GNU dirname is part of GNU coreutils, and this page's checks ran coreutils 9.11. It accepts several paths and prints one folder per path.

### GNU examples
Print the folder of several paths at once:

```sh
dirname logs/server.log projects/notes.txt
```

## Apple
### About macOS dirname
This is the dirname in `/usr/bin` on macOS. It follows the BSD implementation.

### Apple notes
macOS `dirname` follows the BSD behavior and prints one folder per path.

## BusyBox
### About BusyBox dirname
BusyBox bundles dirname as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `dirname` prints the folder part of a path.
