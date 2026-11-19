---
schema_version: 1
id: install
title: install
summary: Copy files and set their permissions, or create directories.
aliases: [ginstall]
tags: [files, permissions, scripts]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/install-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=install&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuinstall.1.html
---
## Common
### Overview
`install` copies a file to a place and sets its permission bits in the same step, and can create directories. It is what build scripts use to put programs into place. Unlike `cp`, it can set the mode with `-m` in one command, and with `-d` it makes directories.

### Background
The `install` utility appeared in 4.2BSD, according to the macOS manual. Today the name covers GNU install, the install that macOS ships, and the install applet in BusyBox. They agree on `-m` for the mode and `-d` for directories.

### Examples
Copy a file and give it the mode 644, which is read and write for the owner and read for others:

```sh
install -m 644 notes.txt copy.txt && ls -l copy.txt | cut -c1-10
```

Copy a program and make it executable with the mode 755:

```sh
install -m 755 notes.txt run.sh && ls -l run.sh | cut -c1-10
```

Create a directory and its missing parents with `-d`:

```sh
install -d x/y && ls -d x/y
```

Create a directory that only its owner can use:

```sh
install -d -m 700 pv && ls -ld pv | cut -c1-10
```

### Option: -m
Set the permission bits of the installed file or directory, as in `-m 755`.

### Option: -d
Create the named directories, with any missing parents, instead of copying files.

### Common mistakes
`install` replaces an existing destination. The default mode is 755 even for a plain file, so give `-m 644` for data. Owner and group options such as `-o` and `-g` usually need administrator rights. Unlike `cp -p`, `install` does not keep the original times unless you ask with `-p`.

### Portability
Only `-m` with a file and `-d` are in the Common examples. GNU `-D` creates missing parent folders of the destination, which macOS does not do with the same option. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU install
GNU install is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-D` creates every missing parent directory of the destination.

### GNU examples
Copy a file into a folder that does not exist yet with `-D`:

```sh
install -D notes.txt q/r/n.txt && ls q/r
```

### Option: -D
Create all missing parent directories of the destination before copying.

## Apple
### About macOS install
This is the install in the system folders of macOS. It follows the BSD implementation.

### Apple notes
On macOS, `-D` is not the GNU option that creates parent folders: the same command printed a usage message in the check. Use `install -d` or `mkdir -p` first.

## BusyBox
### About BusyBox install
BusyBox bundles install as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `install` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
