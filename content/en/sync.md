---
schema_version: 1
id: sync
title: sync
summary: Write cached data to disk.
aliases: [gsync]
tags: [files, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/sync-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=sync&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`sync` tells the system to write any data held in memory buffers out to the disk. Operating systems delay writes to be fast, so data you saved a moment ago may not be on the disk yet. `sync` before unplugging a drive, or before a risky step, makes sure it is.

### Background
Today the name covers GNU sync, the sync that macOS ships, and the sync applet in BusyBox. The macOS manual has no separate page for `sync`, so no date is given. They agree on flushing all pending writes when run with no arguments.

### Examples
Flush all pending writes; there is no output when it works:

```sh short="Flush all pending writes; no output when it works"
sync && echo done
```

Flush after writing a file, before copying it away:

```sh
printf data > f.txt && sync && cat f.txt
```

Use it before a command that might crash the machine:

```sh
sync; echo flushed
```

### Common mistakes
`sync` makes the data safe on the disk, but does not unmount a drive; use the system's eject or unmount command before removing it. It does not report errors in a useful way. On a busy system it can take a while. The file-system sync that real programs use is `fsync` on a file, not a command.

### Portability
Plain `sync` is the same on all three. GNU adds `-d`, `-f` and file names to sync only part of the data, which macOS and BusyBox do not have in the same form. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU sync
GNU sync is part of GNU coreutils, and this page's checks ran coreutils 9.11. With file names it syncs only those files, and `-f` syncs the whole file system that holds them.

### GNU examples
Flush one file with a file name:

```sh
printf data > f.txt && sync f.txt && echo ok
```

### Option: -f, --file-system
Sync the file systems that contain the files.

## Apple
### About macOS sync
This is the sync in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `sync` takes no options and flushes all pending writes.

## BusyBox
### About BusyBox sync
BusyBox bundles sync as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `sync` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
