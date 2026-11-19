---
schema_version: 1
id: readlink
title: readlink
summary: Show where a symbolic link points, or resolve a path completely.
aliases: [greadlink]
tags: [files, links, paths]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/readlink-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=readlink&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnureadlink.1.html
---
## Common
### Overview
`readlink` prints the target stored in a symbolic link, which is the text you gave to `ln -s`. It prints nothing and fails for a file that is not a link. With `-f` it instead follows every link in a path and prints the final absolute path.

### Background
The macOS manual gives no date for the first appearance of `readlink`. Today the name covers GNU readlink, the readlink that macOS ships, and the readlink applet in BusyBox. They agree on printing a link's target and on `-f`; the other options differ.

### Examples
Show the target of a link, exactly as stored:

```sh
ln -s notes.txt link && readlink link
```

Resolve a path through every link to its absolute name, then show only the final part:

```sh
ln -s notes.txt link && basename "$(readlink -f link)"
```

Use the failure on a file that is not a link:

```sh
readlink notes.txt || echo not a link
```

Show a relative link as stored, which can differ from where it leads:

```sh
ln -s ../notes.txt projects/up && readlink projects/up
```

### Option: -f
Follow every link in the path, as far as it goes, and print the absolute path of the result. All but the last part of the path must exist.

### Common mistakes
Plain `readlink` shows only one step: if `a` points to `b` and `b` points to `c`, it prints `b`. Use `-f` for the final target. A file that is not a link produces no output and exit status 1, which is easy to miss in a script.

### Portability
Only the plain form and `-f` are in the Common examples. GNU's `-e`, `-m`, `-n` and `-z` are not shared by every implementation. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU readlink
GNU readlink is part of GNU coreutils, and this page's checks ran coreutils 9.11. It has options that decide how strict the resolution is.

### GNU examples
Resolve a path only if every part of it exists, with `-e`; a missing part is an error:

```sh
readlink -e missing.txt || echo missing
```

### Option: -e, --canonicalize-existing
Like `-f`, but every part of the path, including the last, must exist.

### Option: -m, --canonicalize-missing
Like `-f`, but no part of the path has to exist.

## Apple
### About macOS readlink
This is the readlink in `/usr/bin` on macOS. It follows the BSD implementation, with `-f` and `-n` available on current releases.

### Apple notes
macOS `readlink` has `-f` and `-n`; it does not have GNU's `-e` and `-m`.

## BusyBox
### About BusyBox readlink
BusyBox bundles readlink as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `readlink` supports `-f`. `-e` and `-m` depend on the build.
