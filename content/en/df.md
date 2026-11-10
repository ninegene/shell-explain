---
schema_version: 1
id: df
title: df
summary: Show free and used space on mounted file systems.
aliases: [gdf]
tags: [files, system, size]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/df-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=df&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`df` (disk free) shows, for each mounted file system, its size, how much is used, how much is available and where it is mounted. Give a path and it reports only the file system that holds that path. It answers "is the disk full?"; use `du` to find out which folder is using the space.

### Background
A `df` command appeared in Version 1 AT&T UNIX, according to the macOS manual, which also says POSIX defines only `-k`, `-P` and `-t`. Today the name covers GNU df, the df that macOS ships, and the df applet in BusyBox. They agree on `-k`, `-P` and on `-h` for readable sizes.

### Examples
Show the file system that holds the current folder, in 1024-byte units with `-k`; the first column heading is `Filesystem`:

```sh
df -k . | head -n 1 | cut -d' ' -f1
```

Use the POSIX format with `-P`, which keeps every file system on one line and uses the same columns everywhere:

```sh
df -Pk . | head -n 1 | cut -d' ' -f1
```

Print sizes in a readable form, such as `20G`, with `-h`:

```sh
df -h . | head -n 1 | cut -d' ' -f1
```

Count the lines for one path: the heading and one file system:

```sh
df -Pk . | grep -c .
```

Get the percentage used of the current file system for a script:

```sh
df -Pk . | awk 'NR == 2 {print $5}' | grep -c '%$'
```

### Option: -k
Show sizes in units of 1024 bytes.

### Option: -P
Use the portable POSIX output format, with one line per file system.

### Option: -h
Show sizes with a unit suffix such as K, M or G.

### Common mistakes
The default unit is 512-byte blocks on macOS and 1024-byte blocks on GNU, so numbers from two systems do not compare; use `-k` or `-h`. Without `-P`, a long file system name wraps the line on some systems and breaks a script that reads columns. `df` reports file systems, not folders, so a folder on the same disk shows the whole disk. A full disk is not the only way to run out of space: a file system can also run out of inodes (the entries that identify files), which `df -i` shows.

### Portability
Only `-k`, `-P` and `-h` are in the Common examples, and the examples check only the heading and the shape of the output because the sizes depend on the machine. GNU also has `--output`, `-T` and `-x`, and the column set differs on macOS. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU df
GNU df is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can choose the columns to show and can exclude file system types.

### Option: -T, --print-type
Add a column with the file system type.

### Option: --output
Choose the columns to show, as in `--output=source,pcent`.

## Apple
### About macOS df
This is the df in the system folders of macOS. It follows the BSD implementation. Without `-k` it counts 512-byte blocks, and the `BLOCKSIZE` variable changes this.

### Apple notes
macOS `df` shows extra inode columns by default and has no `--output`. Use `-P` to get the POSIX layout.

## BusyBox
### About BusyBox df
BusyBox bundles df as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `df` supports `-k`, `-P` and `-h`. GNU long options depend on the build.
