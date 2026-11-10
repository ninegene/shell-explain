---
schema_version: 1
id: zcat
title: zcat
summary: Print the contents of compressed files to the screen.
aliases: [gzcat]
tags: [files, compression]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/gzip/manual/gzip.html
  - https://man.freebsd.org/cgi/man.cgi?query=gzip&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`zcat` decompresses a compressed file and writes the result to standard output, leaving the file alone. It is `gzip -dc` under a shorter name, so you can read a `.gz` file, for example a log, without unpacking it on disk. The name goes back to the older `compress` tool, whose `.Z` files `zcat` first read, and that history is why implementations differ.

### Background
Today the name covers GNU zcat (part of gzip), the zcat that macOS ships (part of its gzip, with the name `gzcat` as an alias) and the zcat applet in BusyBox. They agree on reading a compressed file from standard input and writing it to standard output.

### Examples
Read a compressed file from standard input; this form works on every implementation:

```sh
printf data | gzip > f.gz && zcat < f.gz
```

Pipe it into another command to look at a compressed log:

```sh
printf 'a\nb\n' | gzip > f.gz && zcat < f.gz | wc -l | tr -d ' '
```

Do the same with `gunzip -c`, which names the file and works everywhere:

```sh
printf data | gzip > f.gz && gunzip -c f.gz
```

### Common mistakes
On macOS, `zcat f.gz` fails with "can\'t stat: f.gz (f.gz.Z)", because that `zcat` looks for a `.Z` file by default; use `zcat < f.gz`, `gzcat f.gz` or `gunzip -c f.gz` instead. `zcat` reads gzip data, and on GNU also the old `.Z` format; it is not for `.zip`, `.bz2` or `.xz` files, which have their own tools. A file that is not compressed gives an error, not its own text.

### Portability
Only the standard-input form and `gunzip -c` are in the Common examples. The file-name form `zcat f.gz` works on GNU and BusyBox but not on macOS. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU zcat
GNU zcat is part of GNU gzip, and this page's checks ran the gzip in the Linux guest and on the host.

### GNU examples
Name the file directly; GNU reads it:

```sh
printf data | gzip > f.gz && zcat f.gz
```

## Apple
### About macOS zcat
This is the zcat in `/usr/bin` on macOS, which is Apple's gzip under another name. `gzcat` is the same program for reading `.gz` files by name.

### Apple examples
Use `gzcat` to name the file on macOS:

```sh
printf data | gzip > f.gz && gzcat f.gz
```

### Apple notes
On macOS, name the file with `gzcat f.gz` or use standard input, as in `zcat < f.gz`.

## BusyBox
### About BusyBox zcat
BusyBox bundles zcat as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `zcat` reads a gzip file, by name or from standard input.
