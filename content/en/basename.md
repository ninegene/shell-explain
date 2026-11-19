---
schema_version: 1
id: basename
title: basename
summary: Print the file name part of a path, with or without its extension.
aliases: [gbasename]
tags: [files, paths]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/basename-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=basename&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnubasename.1.html
---
## Common
### Overview
`basename` removes the folders from a path and prints the last part. It can also remove a suffix, such as an extension. It works on the text of the path only, so the file does not have to exist.

### Background
The `basename` and `dirname` utilities first appeared in 4.4BSD, according to the macOS manual. Today the name covers GNU basename, the basename that macOS ships, and the basename applet in BusyBox. They agree on a path, and on a suffix given as a second argument; they differ in the options for several names at once.

### Examples
Print the file name of a path:

```sh
basename logs/server.log
```

Remove a suffix as well:

```sh
basename logs/server.log .log
```

A trailing slash is ignored:

```sh
basename /usr/local/bin/
```

Use it in a loop to build new names from old ones:

```sh
for f in logs/*.log; do echo "$(basename "$f" .log).txt"; done
```

### Common mistakes
`basename` never looks at the file system: `basename logs/missing.log` prints `missing.log`. Quote a path with spaces, or it is split into several arguments. A suffix only matches at the very end, and a name that equals the suffix is left whole, so `basename .log .log` prints `.log`. In a shell script, `${path##*/}` does the same without starting a process, but it is not what a trailing slash gives you.

### Portability
Only a path with an optional suffix is in the Common examples. `-a` and `-s` for several names differ between implementations and are not in POSIX. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU basename
GNU basename is part of GNU coreutils, and this page's checks ran coreutils 9.11. With `-a` it handles many paths in one call.

### GNU examples
Print the file name of every path given; `-a` treats every argument as a path:

```sh
basename -a logs/server.log projects/notes.txt
```

### Option: -a, --multiple
Treat every argument as a path and print each name on its own line.

### Option: -s, --suffix
Remove this suffix from every name; it implies `-a`.

## Apple
### About macOS basename
This is the basename in `/usr/bin` on macOS. It follows the BSD implementation and also accepts `-a` and `-s`.

### Apple notes
macOS `basename` supports `-a` and `-s SUFFIX` with the same meaning as GNU's.

## BusyBox
### About BusyBox basename
BusyBox bundles basename as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `basename` prints a name with an optional suffix; `-a` and `-s` depend on the build.
