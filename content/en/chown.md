---
schema_version: 1
id: chown
title: chown
summary: Change the owner and group of files.
aliases: [gchown]
tags: [files, permissions]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/chown-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=chown&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`chown` changes who owns a file, and with a colon, its group too. Ownership decides which permission bits apply to you. Normal users cannot give files away, so changing the owner to someone else usually needs administrator rights; the examples change a file to its current owner so that they are safe to run.

### Background
Today the name covers GNU chown, the chown that macOS ships (in `/usr/sbin`), and the chown applet in BusyBox. The macOS manual gives no date for it. They agree on `OWNER`, `OWNER:GROUP` and `-R`.

### Examples
Set the owner to yourself, which is always allowed:

```sh
chown "$(id -u)" notes.txt && echo ok
```

Set the owner and the group together with a colon:

```sh
chown "$(id -u):$(id -g)" notes.txt && echo ok
```

Check the owner afterwards; `ls -l` shows it in the third column:

```sh
chown "$(id -u)" notes.txt && test "$(ls -l notes.txt | awk '{print $3}')" = "$(id -un)" && echo yes
```

Apply it to a folder and everything inside it with `-R`:

```sh short="Apply it to a folder and everything inside (-R)"
chown -R "$(id -u)" projects logs && echo ok
```

### Option: -R
Change the owner of folders and everything inside them.

### Common mistakes
`chown -R` on the wrong folder can change the owner of thousands of files, and it has no undo; check the path first. Be careful with `-R` and symbolic links, because a link can lead outside the folder you meant. Giving a file to another user is a one-way move for a normal user, since only an administrator can take it back. Use `chmod` for permissions; `chown` does not change them.

### Portability
Only owner, `OWNER:GROUP` and `-R` are in the Common examples, using your own IDs. The old `OWNER.GROUP` spelling with a dot is not portable. The examples use numeric IDs from `id`, which work everywhere. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU chown
GNU chown is part of GNU coreutils, and this page's checks ran coreutils 9.11. `--reference=FILE` copies the owner and group of another file.

### Option: --reference
Use the owner and group of this file instead of naming them.

## Apple
### About macOS chown
This is the chown in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `chown` is in `/usr/sbin`. It has `-R`, `-h` and `-L`, but not GNU's `--reference`.

## BusyBox
### About BusyBox chown
BusyBox bundles chown as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `chown` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
