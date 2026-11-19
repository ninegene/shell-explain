---
schema_version: 1
id: rmdir
title: rmdir
summary: Remove empty directories.
aliases: [grmdir]
tags: [files, directories]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/rmdir-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=rmdir&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnurmdir.1.html
---
## Common
### Overview
`rmdir` removes directories, but only empty ones. That makes it the safe way to delete a folder: it refuses if anything is inside. To remove a folder and its contents, use `rm -r` and read its warnings first.

### Background
An `rmdir` command appeared in Version 1 AT&T UNIX, according to the macOS manual. Today the name covers GNU rmdir, the rmdir that macOS ships, and the rmdir applet in BusyBox. They agree on removing empty directories and on `-p`; the options for ignoring a non-empty directory differ.

### Examples
Remove an empty directory:

```sh
mkdir empty && rmdir empty && echo removed
```

Remove several empty directories at once:

```sh
mkdir one two && rmdir one two && echo ok
```

See that a directory with something inside is refused:

```sh
rmdir logs 2> /dev/null || echo not empty
```

Remove a directory and then its empty parents with `-p`:

```sh
mkdir -p a/b/c && rmdir -p a/b/c && { ls -d a 2> /dev/null || echo gone; }
```

### Option: -p
Remove the directory, then each empty parent in the path, so `rmdir -p a/b/c` removes `a/b/c`, `a/b` and `a`.

### Common mistakes
`rmdir` fails on a directory that has any file in it, including hidden files whose names start with a dot. It does not remove files, and it does not follow a symbolic link to a directory. Do not reach for `rm -rf` just to silence the error: the error is telling you the folder still has something in it.

### Portability
Only the plain form and `-p` are in the Common examples. GNU adds `--ignore-fail-on-non-empty` and `-v`, which macOS and BusyBox do not have in the same form. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU rmdir
GNU rmdir is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can skip a directory that is not empty without reporting an error.

### GNU examples
Do not complain about a directory that still has files in it; it is left in place and the exit status stays 0:

```sh
rmdir --ignore-fail-on-non-empty logs && ls logs
```

### Option: --ignore-fail-on-non-empty
Do not treat a non-empty directory as an error. The directory is not removed.

## Apple
### About macOS rmdir
This is the rmdir in `/bin` on macOS. It follows the BSD implementation.

### Apple notes
macOS `rmdir` has `-p` and `-v`; it has no `--ignore-fail-on-non-empty`.

## BusyBox
### About BusyBox rmdir
BusyBox bundles rmdir as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `rmdir` supports `-p`. GNU long options depend on the build.
