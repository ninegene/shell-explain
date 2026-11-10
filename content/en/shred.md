---
schema_version: 1
id: shred
title: shred
summary: Overwrite a file so its contents are hard to recover.
aliases: [gshred]
tags: [files, security]
implementations: [gnu, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/shred-invocation.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`shred` overwrites a file several times with random data, so the original contents are harder to recover, and with `-u` it then removes the file. It is meant for sensitive data. It works only when the file system overwrites data in place, which many modern ones do not, so it is not a guarantee.

### Background
`shred` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it. They agree on overwriting a file, and on `-u` to remove it afterwards.

### Examples
Overwrite a file in place; its size stays the same but the contents change:

```sh
printf secret > s.txt && cp s.txt orig.txt && shred s.txt && { cmp -s s.txt orig.txt || echo overwritten; }
```

Overwrite and then delete the file with `-u`:

```sh
printf secret > s.txt && shred -u s.txt && { ls s.txt 2> /dev/null || echo removed; }
```

Use one pass instead of the default three with `-n`:

```sh
printf secret > s.txt && cp s.txt orig.txt && shred -n 1 s.txt && { cmp -s s.txt orig.txt || echo overwritten; }
```

### Option: -u
Remove the file after overwriting it.

### Option: -n
Overwrite this many times; the default is 3.

### Common mistakes
`shred` is not a guarantee: on journaling, copy-on-write and flash storage (such as most SSDs and APFS), old copies of the data can survive elsewhere, and backups and snapshots are not touched. For sensitive disks use full-disk encryption. `shred` destroys the data for good, so check the file name. Without `-u` the file stays, filled with garbage.

### Portability
This page has no Apple view: macOS has no `shred`; `rm -P` overwrote files on older systems and is not a safe substitute today. Only `-u` and `-n` are in the Common examples. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU shred
GNU shred is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-z` adds a last pass of zeros to hide that the file was shredded.

### Option: -z, --zero
Add a final overwrite with zeros.

## BusyBox
### About BusyBox shred
BusyBox bundles shred as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `shred` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
