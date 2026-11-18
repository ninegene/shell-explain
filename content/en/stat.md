---
schema_version: 1
id: stat
title: stat
summary: Show file size, permissions, type and timestamps.
aliases: [gstat]
tags: [files, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/stat-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=stat&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`stat` reports details about a file that `ls` hides or formats loosely: exact size, permission bits, type, owner and timestamps. Its output format is different on each implementation, and so is the option that chooses the fields, so the portable use is the exit status. Pick the format option in the GNU, Apple or BusyBox section for what you need.

### Background
The `stat` utility appeared in NetBSD 1.6 and FreeBSD 4.10, according to the FreeBSD and macOS manuals. Today the name covers GNU stat, the stat that macOS ships, and the stat applet in BusyBox. They agree only on whether a file exists, and differ on the option that chooses the fields: GNU and BusyBox use `-c`, macOS uses `-f`, and the format letters do not overlap.

### Examples
Test that a file exists. The output is discarded, and `stat` exits with a non-zero status for a missing file:

```sh short="Test that a file exists"
stat notes.txt > /dev/null && echo exists
```

Handle a missing file:

```sh
stat missing.txt 2> /dev/null || echo missing
```

### Common mistakes
Parsing the default output of `stat` in a script breaks when you move between Linux and macOS. Choose the format option for the implementation in use, or use `test -e`, `test -d` and `ls` where they are enough. By default `stat` describes a symbolic link itself; `-L` describes its target.

### Portability
GNU and BusyBox use `-c FORMAT` with letters such as `%s` for size and `%a` for the octal mode. macOS uses `-f FORMAT` with `%z` for size and `%Lp` for the octal mode. The two format languages do not overlap, so only the existence check is common. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU stat
GNU stat is part of GNU coreutils, and this page's checks ran coreutils 9.11. It prints chosen fields with `-c FORMAT`.

### GNU examples
Print only the size in bytes:

```sh
stat -c '%s' notes.txt
```

Print the octal permission bits and the file name; `%a` is the mode and `%n` the name:

```sh short="Print the octal permission bits and the name (%a, %n)"
stat -c '%a %n' notes.txt
```

Print the file type as words, such as `regular file` or `directory`:

```sh
stat -c '%F' logs
```

### Option: -c
Print using a format string instead of the default report. `%s` is the size in bytes, `%a` the octal permissions, `%n` the name, `%F` the file type, `%U` the owner and `%y` the modification time.

## Apple
### About macOS stat
This is the stat in `/usr/bin` on macOS. It prints chosen fields with `-f FORMAT`, using format letters that differ from GNU's.

### Apple examples
Print only the size in bytes:

```sh
stat -f '%z' notes.txt
```

Print the octal permission bits and the file name; `%Lp` is the low permission bits and `%N` the name:

```sh short="Print the octal permission bits and the name (%Lp, %N)"
stat -f '%Lp %N' notes.txt
```

Print the file type as words, such as `Regular File` or `Directory`:

```sh
stat -f '%HT' logs
```

### Option: -f
Print using a format string instead of the default report. `%z` is the size in bytes, `%Lp` the octal permissions, `%N` the name, `%HT` the file type, `%Su` the owner and `%Sm` the modification time.

## BusyBox
### About BusyBox stat
BusyBox bundles stat as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox examples
BusyBox `stat` takes the same `-c` format as GNU for the fields below:

```sh
stat -c '%s' notes.txt
```

```sh
stat -c '%a %n' notes.txt
```

### Option: -c
Print using a format string instead of the default report. BusyBox supports the common GNU letters; the set depends on the build.
