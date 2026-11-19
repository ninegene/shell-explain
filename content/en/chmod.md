---
schema_version: 1
id: chmod
title: chmod
summary: Change file and directory permissions.
aliases: [gchmod]
tags: [files, permissions]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/chmod-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=chmod&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuchmod.1.html
---
## Common
### Overview
`chmod` changes who can read, write or run a file. A mode is either an octal number such as `644` or a symbolic change such as `u+x`. Only the owner of a file, or root, can change its mode.

### Background
A `chmod` command appeared in Version 1 Unix, according to the FreeBSD and macOS manuals. Today the name covers GNU chmod, the chmod that macOS ships, and the chmod applet in BusyBox. They agree on octal modes, the classes `u`, `g`, `o`, and `a`, the operators `+`, `-`, and `=`, and `-R`, and differ on GNU-only options such as `--reference`.

### Examples
Give the owner read and write access and nobody else; the digits are owner, group and others, where read is 4, write is 2 and execute is 1. `ls -l` and `cut` show the resulting permission string:

```sh
chmod 600 notes.txt && ls -l notes.txt | cut -c1-10
```

Add execute permission for the owner only:

```sh
chmod u+x notes.txt && ls -l notes.txt | cut -c1-10
```

Make a file executable for everyone:

```sh
chmod a+x notes.txt && ls -l notes.txt | cut -c1-10
```

Remove read access from the group and others:

```sh
chmod go-r app.log && ls -l app.log | cut -c1-10
```

Set exact permissions for each class at once:

```sh
chmod u=rw,go=r notes.txt && ls -l notes.txt | cut -c1-10
```

Change a directory and everything below it:

```sh
chmod -R 700 logs && ls -ld logs logs/server.log | cut -c1-10
```

### Option: -R
Apply the change to a directory and everything inside it. Modes that make sense for directories, such as the execute bit, are rarely right for plain files.

### Common mistakes
`chmod -R 777` makes everything writable by everyone and is almost never the fix. Use `chmod -R u+rwX,go+rX` to give directories the execute bit without making every file executable; GNU and BSD `chmod` accept the capital `X`. A symbolic `+x` without a class depends on the `umask`, so name the class (`u+x`, `a+x`) when the result matters.

### Portability
The Common examples use octal modes, the classes `u`, `g`, `o` and `a`, the operators `+`, `-` and `=`, and `-R`. GNU `--reference`, `--changes` and `--preserve-root` are not available in the other implementations. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU chmod
GNU chmod is part of GNU coreutils, and this page's checks ran coreutils 9.11. It adds long options such as `--reference`, `--changes`, and `--preserve-root`.

### GNU examples
Copy the mode of another file with `--reference`. This option is a GNU extension:

```sh
chmod --reference=app.log notes.txt && ls -l notes.txt | cut -c1-10
```

## Apple
### About macOS chmod
This is the chmod in `/bin` on macOS. It has no `--reference`, so read a mode with `stat` or `ls -l` and set it explicitly. It supports ACL entries with `+a` and `-a`, which GNU and BusyBox do not.

### Apple notes
macOS `chmod` is the BSD implementation. It has no `--reference`; read the mode with `stat` or `ls -l` and set it explicitly. It also supports ACL entries with `+a` and `-a`, which GNU and BusyBox do not.

## BusyBox
### About BusyBox chmod
BusyBox bundles chmod as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `chmod` supports octal and symbolic modes and `-R`; GNU long options depend on the build.
