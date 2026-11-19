---
schema_version: 1
id: ln
title: ln
summary: Make hard links and symbolic links to files.
aliases: [gln]
tags: [files, links]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/ln-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=ln&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuln.1.html
---
## Common
### Overview
`ln` gives a file a second name. Without options it makes a hard link, another name for the same file. With `-s` it makes a symbolic link, a small file that holds a path to somewhere else. The first argument is the existing target and the second is the new name.

### Background
An `ln` command appeared in Version 1 AT&T UNIX, according to the macOS manual. Today the name covers GNU ln, the ln that macOS ships, and the ln applet in BusyBox. They agree on `-s` and `-f`; the options for relative links and for replacing links differ.

### Examples
Make a symbolic link and read the file through it:

```sh
ln -s notes.txt link.txt && head -n 1 link.txt
```

Check that the new name is a symbolic link:

```sh
ln -s notes.txt link.txt && test -L link.txt && echo symlink
```

Make a hard link; a change made through one name shows in the other, because both names are the same file:

```sh
ln notes.txt copy.txt && echo extra >> copy.txt && tail -n 1 notes.txt
```

Link to a file in the parent directory. A relative target is read from the folder that holds the link, not from where you ran `ln`:

```sh
ln -s ../notes.txt projects/notes.txt && head -n 1 projects/notes.txt
```

Replace an existing link with `-f`:

```sh
ln -s notes.txt current && ln -sf app.log current && head -n 1 current
```

Link to a directory:

```sh
ln -s logs latest && ls latest
```

Find out that a link points at nothing, because `test -e` follows the link:

```sh
ln -s missing.txt broken && { test -e broken || echo broken link; }
```

### Option: -s
Make a symbolic link instead of a hard link.

### Option: -f
Replace the new name if it already exists, instead of failing.

### Common mistakes
A relative symbolic link is stored as written. `ln -s notes.txt projects/notes.txt` points at `projects/notes.txt` itself, not at `notes.txt` in the parent folder, so the link is broken; write `../notes.txt` or an absolute path. Hard links cannot cross file systems, and most systems do not let you hard-link a directory. Deleting the target of a symbolic link leaves a broken link behind; deleting one name of a hard link leaves the file reachable through the others.

### Portability
Only `-s` and `-f` are in the Common examples. Options such as `-r` (relative links), `-T`, `-n` and `-v` differ in spelling or meaning between implementations. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU ln
GNU ln is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can make a link relative to the link's own folder with `-r`.

### GNU examples
Make a symbolic link whose target is rewritten relative to the link's folder, so it keeps working if the whole tree moves:

```sh
ln -sr notes.txt projects/notes.txt && head -n 1 projects/notes.txt
```

### Option: -r, --relative
With `-s`, store the target as a path relative to the folder that contains the link.

## Apple
### About macOS ln
This is the ln in `/bin` on macOS. It follows the BSD implementation, which has no `-r` option for relative links.

### Apple notes
macOS `ln` has no `-r`; write the relative target yourself, as in `ln -s ../notes.txt projects/notes.txt`. The `-F` option, which can replace an empty directory, is a BSD extension.

## BusyBox
### About BusyBox ln
BusyBox bundles ln as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `ln` supports `-s` and `-f`. GNU long options and `-r` depend on the build.
