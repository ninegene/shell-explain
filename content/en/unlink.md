---
schema_version: 1
id: unlink
title: unlink
summary: Remove one file with the plain one-argument form.
aliases: [gunlink]
tags: [files]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/unlink-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=unlink&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`unlink` removes a single file or link. It does what `rm FILE` does for one file, but it takes no options and never asks. It mirrors the system function of the same name, which removes one name of a file.

### Background
Today the name covers GNU unlink, the unlink that macOS ships, and the unlink applet in BusyBox. The macOS manual has no separate page for `unlink`, so no date is given. All three take exactly one argument.

### Examples
Remove a file and check that it is gone:

```sh
unlink app.log && { ls app.log 2> /dev/null || echo gone; }
```

Remove one name of a hard link; the other name still works:

```sh
ln notes.txt hard.txt && unlink hard.txt && head -n 1 notes.txt
```

Remove a symbolic link without touching its target:

```sh
ln -s notes.txt link && unlink link && head -n 1 notes.txt
```

### Common mistakes
`unlink` accepts exactly one name; it has no `-r` or `-f`, and it does not remove directories on most systems, so use `rmdir` or `rm -r` for those. It does not ask for confirmation, and the file is gone at once, so check the name before you press Enter. When a file has several hard links, `unlink` removes only the name you gave.

### Portability
Only the one-argument form exists, and it is the same on all three implementations. `rm` is more common and safer in scripts because of its options. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU unlink
GNU unlink is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### GNU notes
GNU `unlink` takes exactly one argument and no other options.

## Apple
### About macOS unlink
This is the unlink in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `unlink` takes exactly one argument.

## BusyBox
### About BusyBox unlink
BusyBox bundles unlink as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `unlink` takes exactly one argument.
