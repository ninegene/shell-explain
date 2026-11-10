---
schema_version: 1
id: realpath
title: realpath
summary: Print the absolute, fully resolved form of a path.
aliases: [grealpath]
tags: [files, paths]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/realpath-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=realpath&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`realpath` turns any path into its absolute form: it adds the current folder to a relative path, removes `.` and `..`, and follows symbolic links. Use it to find out where a path really points, or to give another program a path that does not depend on the current folder.

### Background
The `realpath` utility first appeared in FreeBSD 4.3, according to the macOS manual. Today the name covers GNU realpath, the realpath that macOS ships, and the realpath applet in BusyBox. They agree on resolving a path that exists; they differ on a path that does not exist yet.

### Examples
Print the absolute path of a file reached through `..`. The output starts with a slash, so `grep -c` counts one match:

```sh
realpath logs/../notes.txt | grep -c '^/'
```

Show only the final part of the resolved path:

```sh
basename "$(realpath logs/../notes.txt)"
```

Resolve a symbolic link to the file it points at:

```sh
ln -s notes.txt link && basename "$(realpath link)"
```

Resolve the current folder:

```sh
cd logs && basename "$(realpath .)"
```

### Common mistakes
`realpath` gives a different answer from the text you typed whenever a link or `..` is involved, which is the point, but it can surprise you when a folder is a link into another disk. A path that does not exist is an error on macOS and is accepted by GNU, so a script that depends on either needs a check first with `test -e`. Quote the path and the command substitution if the name may contain spaces.

### Portability
Only a path that exists is in the Common examples. GNU's `--relative-to`, `-s`, `-m` and `-e` are not on every system. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU realpath
GNU realpath is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can express a path relative to another folder and can accept missing paths.

### GNU examples
Print a path relative to another folder:

```sh
realpath --relative-to=logs notes.txt
```

Accept a path whose last part does not exist yet; `-m` allows missing parts, and `basename` keeps the final name:

```sh
basename "$(realpath -m logs/new.txt)"
```

### Option: --relative-to
Print the result relative to this folder instead of as an absolute path.

### Option: -m, --canonicalize-missing
Do not fail when parts of the path do not exist.

## Apple
### About macOS realpath
This is the realpath in `/bin` on macOS. It follows the BSD implementation and requires the path to exist.

### Apple notes
macOS `realpath` fails for a path that does not exist, and it has no `--relative-to` or `-m`. Use `-q` to hide the error message.

## BusyBox
### About BusyBox realpath
BusyBox bundles realpath as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `realpath` prints the resolved path of each argument; the GNU options depend on the build.
