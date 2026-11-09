---
schema_version: 1
id: gzip
title: gzip
summary: Compress and decompress files with gzip.
aliases: [ggzip]
tags: [files, compression]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`gzip` compresses each file you name and replaces it with a file ending in `.gz`. `-d` goes the other way. It compresses one file at a time; use `tar` to bundle several files first.

### Background
The macOS manual says `gzip` was originally written by Jean-loup Gailly, and that Matthew R. Green wrote a simple front end for NetBSD, later made mostly feature-compatible with the original GNU gzip. It also describes the macOS implementation as a port of the NetBSD gzip that first appeared in FreeBSD 7.0. Today the name covers GNU gzip, the Apple gzip that macOS ships, and the gzip applet in BusyBox. They agree on `-k`, `-d`, and `-c`, and differ on options such as `-r`, `-l`, and `--rsyncable`.

### Examples
Compress a file; the original is replaced by `app.log.gz`:

```sh
gzip app.log && ls
```

Compress and keep the original file:

```sh
gzip -k app.log && ls app.log*
```

Print the uncompressed contents of a compressed file without writing a file:

```sh
gzip -k app.log && gzip -dc app.log.gz
```

Decompress a file in place:

```sh
gzip notes.txt && gzip -d notes.txt.gz && head -n 1 notes.txt
```

Compress into a chosen file, leaving the input alone, by writing to standard output:

```sh
gzip -c app.log > copy.log.gz && ls *.gz
```

### Option: -k
Keep the input file instead of deleting it after compressing or decompressing.

### Option: -d
Decompress. `gunzip` is the same as `gzip -d`.

### Option: -c
Write to standard output and leave the input file unchanged.

### Common mistakes
`gzip file` deletes `file` unless you pass `-k` or `-c`, so a failed copy is easy to miss. It refuses to overwrite an existing `.gz` file without `-f`. Compression levels and the stored file name differ between versions, so compare contents, not archive bytes.

### Portability
The Common examples use `-k`, `-d` and `-c`. The `-r` option for recursive compression, `-l` listing and `--rsyncable` differ or are missing between implementations. The Common examples were executed on GNU gzip, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU gzip
GNU gzip is the GNU Project's gzip, and this page's checks ran GNU gzip 1.14 on macOS (Homebrew) and the Alpine GNU gzip. It accepts `-S` to choose the suffix.

### GNU examples
Compress with a different suffix. GNU accepts `-S`, and BusyBox builds may not:

```sh
gzip -k -S .z app.log && ls app.log*
```

## Apple
### About macOS gzip
This is the gzip in `/usr/bin` on macOS, and it reports itself as Apple gzip 479. Its version output differs from GNU's, so do not rely on parsing `--version`.

### Apple notes
macOS ships Apple's own `gzip` (`Apple gzip 479`), which accepts the options in the common examples. Version output differs from GNU, so do not rely on parsing `--version`.

## BusyBox
### About BusyBox gzip
BusyBox bundles gzip as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `gzip` supports `-k`, `-d` and `-c` in recent builds; `-r` and other GNU options are not part of the common examples.
