---
schema_version: 1
id: du
title: du
summary: Show how much disk space files and folders use.
aliases: [gdu]
tags: [files, system, size]
implementations: [gnu, apple, busybox]
updated: 2026-11-09
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/du-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=du&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`du` (disk usage) adds up the space used by files and by every folder below the ones you name, and prints one line per folder: a size, a tab and the path. The size is disk space in blocks, which is usually larger than the length of the data, because a file uses whole blocks. Use `ls -l` or `wc -c` for the exact number of bytes.

### Background
The `du` utility and its `-a` and `-s` options first appeared in Version 1 AT&T UNIX, according to the macOS manual. Today the name covers GNU du, the du that macOS ships, and the du applet in BusyBox. They agree on `-s`, `-k`, `-h`, `-c` and `-d`; the default unit and a few options differ.

### Examples
Show the total for a folder with `-s` (summary), in kilobyte units with `-k`; `cut` keeps only the path to show the shape of the line:

```sh short="Total for a folder (-s) in kilobytes (-k), path only"
du -sk logs | cut -f2
```

Show sizes in a readable form, such as `4.0K` or `1.2M`, with `-h`:

```sh
du -sh logs | grep -c '^[0-9.]*[BKMGT]*[[:space:]]logs$'
```

Show every file as well as every folder with `-a`:

```sh
du -a logs | cut -f2
```

List only the first level of folders below here with `-d 1`:

```sh short="List only the first level of folders (-d 1)"
du -d 1 . | cut -f2
```

Add a grand total with `-c`:

```sh
du -c notes.txt app.log | tail -n 1 | cut -f2
```

Check that the size and the path are separated by a tab and that the size is a number:

```sh
du -sk logs | grep -c '^[0-9][0-9]*[[:space:]]logs$'
```

### Option: -s
Print only a total for each argument, not a line for every folder below it.

### Option: -h
Print sizes with a unit suffix such as K, M or G, chosen for each line.

### Option: -k
Print sizes in units of 1024 bytes. Without it the unit differs between systems.

### Option: -c
Print a grand total as the last line.

### Option: -d
Show lines only down to this many levels below the argument, as in `-d 1`. `-d 0` is the same as `-s`.

### Common mistakes
`du` reports space used on disk, so a 34-byte file shows as 4K on most systems and the number can change with the file system. The default unit is not the same everywhere (see Portability), so use `-k` or `-h` when you compare or script. A hard-linked file is counted once, in the first place `du` meets it. `du -sh *` skips names that start with a dot, and sorts by name, not by size; pipe `du -sk *` to `sort -n` to find the largest.

### Portability
The Common examples use `-s`, `-k`, `-h`, `-a`, `-c` and `-d`. GNU also accepts `--max-depth` and `-b` for bytes; macOS has no `-b`. The default unit is 512-byte blocks on macOS and 1024-byte blocks on GNU, unless `POSIXLY_CORRECT` or `BLOCKSIZE` changes it. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). Their output only shows the shape of a line, because the sizes depend on the file system.

## GNU
### About GNU du
GNU du is part of GNU coreutils, and this page's checks ran coreutils 9.11. By default it counts 1024-byte blocks, and it can show exact sizes and apparent length.

### GNU examples
Print the exact length of the data in bytes, not disk blocks, with `-b`:

```sh
du -b notes.txt
```

### Option: -b, --bytes
Print the apparent size in bytes, the length of the data rather than the disk space used.

### Option: --max-depth
The long spelling of `-d`: `--max-depth=1` shows lines for the first level only.

## Apple
### About macOS du
This is the du in `/usr/bin` on macOS. It follows the BSD implementation, and without `-k` or `-h` it counts 512-byte blocks.

### Apple notes
macOS `du` uses 512-byte blocks by default; the `BLOCKSIZE` variable changes this. It has no `-b` for bytes. Use `-k` or `-h` for portable numbers.

## BusyBox
### About BusyBox du
BusyBox bundles du as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `du` supports `-s`, `-k`, `-h`, `-a`, `-c` and `-d`. The GNU long options depend on the build.
