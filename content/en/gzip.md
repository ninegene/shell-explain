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
### GNU examples
Compress with a different suffix. GNU accepts `-S`, and BusyBox builds may not:

```sh
gzip -k -S .z app.log && ls app.log*
```

## Apple
### Apple notes
macOS ships Apple's own `gzip` (`Apple gzip 479`), which accepts the options in the common examples. Version output differs from GNU, so do not rely on parsing `--version`.

## BusyBox
### BusyBox notes
BusyBox `gzip` supports `-k`, `-d` and `-c` in recent builds; `-r` and other GNU options are not part of the common examples.
