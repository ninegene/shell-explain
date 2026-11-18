---
schema_version: 1
id: link
title: link
summary: Make a hard link with the plain two-argument form.
aliases: [glink]
tags: [files, links]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/link-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=link&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`link` makes a hard link: a second name for an existing file. It does exactly what `ln FILE LINKNAME` does, with no options. Most people use `ln`, which can also make symbolic links; `link` is the minimal call that mirrors the system function of the same name.

### Background
Today the name covers GNU link, the link that macOS ships, and the link applet in BusyBox. The macOS manual has no separate page for `link`, so no date is given. All three take exactly two arguments, the existing file and the new name.

### Examples
Make a second name for a file; both names are the same file:

```sh short="Make a second name for a file"
link notes.txt hard.txt && head -n 1 hard.txt
```

Show that a change through one name appears in the other:

```sh
link notes.txt hard.txt && echo extra >> hard.txt && tail -n 1 notes.txt
```

Check that both names have the same size:

```sh
link notes.txt hard.txt && test notes.txt -ef hard.txt && echo same file
```

### Common mistakes
`link` has no `-s`, so it cannot make a symbolic link; use `ln -s`. Both arguments must be on the same file system, and the target must not exist yet. It refuses directories on most systems. If you need options, use `ln`.

### Portability
Only the two-argument form exists, and it is the same on all three implementations. `ln` is the portable choice in scripts. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU link
GNU link is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### GNU notes
GNU `link` takes exactly two arguments and no other options.

## Apple
### About macOS link
This is the link in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `link` takes exactly two arguments.

## BusyBox
### About BusyBox link
BusyBox bundles link as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `link` takes exactly two arguments.
