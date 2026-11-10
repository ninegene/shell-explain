---
schema_version: 1
id: mktemp
title: mktemp
summary: Create a temporary file or directory with a unique name.
aliases: [gmktemp]
tags: [files, directories, scripts]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/mktemp-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=mktemp&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`mktemp` creates a new, empty file or directory with a name nobody else is using and prints the name. Scripts use it for scratch space, because a fixed name such as `/tmp/data` can collide with another run or be guessed by someone else. You give a template ending in `XXXXXX`, and `mktemp` replaces the Xs with random characters.

### Background
A `mktemp` utility appeared in OpenBSD 2.1, and the FreeBSD one was written independently from the OpenBSD manual and first appeared in FreeBSD 2.2.7, according to the macOS manual. Today the name covers GNU mktemp, the mktemp that macOS ships, and the mktemp applet in BusyBox. They agree on a template and on `-d`; the default location and the extra options differ.

### Examples
Create a temporary file and keep its name in a variable:

```sh
f=$(mktemp) && test -f "$f" && echo created; rm -f "$f"
```

Create a temporary directory with `-d`:

```sh
d=$(mktemp -d) && test -d "$d" && echo created; rmdir "$d"
```

Choose where the file goes and how it starts, with a template that ends in six Xs:

```sh
f=$(mktemp cache/item.XXXXXX) && ls cache | grep -c .
```

Build a name from a template in the current folder:

```sh
f=$(mktemp work.XXXXXX) && echo "${f%.*}"
```

### Option: -d
Create a directory instead of a file.

### Common mistakes
Do not make up a name such as `/tmp/myapp.$$`: it can be guessed, and `mktemp` exists to avoid that. A file made by `mktemp` is not removed for you, so delete it when you are done, or use `trap 'rm -f "$f"' EXIT` in a script. On macOS the template's Xs must be at the very end, or they are not replaced. `mktemp` makes the file with owner-only permissions, so another user cannot read it.

### Portability
Only the template and `-d` are in the Common examples, and the default folder is left to the system: usually `/tmp` or the folder in `TMPDIR`, with details that differ. GNU's `-p` and `--suffix`, and BSD's `-t`, are not portable. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU mktemp
GNU mktemp is part of GNU coreutils, and this page's checks ran coreutils 9.11. It reads the `TMPDIR` variable and can add a suffix.

### GNU examples
Put a suffix after the random part, as in `item.XXXXXX.txt`:

```sh
f=$(mktemp item.XXXXXX --suffix=.txt) && echo "${f##*.}"
```

### Option: -p, --tmpdir
Create the file in this folder, or in `TMPDIR` when no folder is given.

### Option: --suffix
Append this text after the Xs of the template.

## Apple
### About macOS mktemp
This is the mktemp in `/usr/bin` on macOS. It follows the BSD implementation, where `-t` takes a prefix and creates the file in the user's temporary folder.

### Apple notes
macOS `mktemp` has `-t PREFIX`, which makes a name in the folder from `TMPDIR`. Plain `mktemp` ignores `TMPDIR`, and a template with text after the Xs is not filled in and makes a file with the literal name; use a template that ends in the Xs.

## BusyBox
### About BusyBox mktemp
BusyBox bundles mktemp as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `mktemp` supports a template and `-d`. `-p` and `--suffix` depend on the build.
