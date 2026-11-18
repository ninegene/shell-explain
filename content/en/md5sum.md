---
schema_version: 1
id: md5sum
title: md5sum
summary: Compute or verify MD5 checksums of files.
aliases: [gmd5sum]
tags: [files, checksum, security]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/md5sum-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=md5sum&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`md5sum` reads a file, or standard input, and prints a 128-bit digest of its contents as hexadecimal, followed by the file name. Two files with the same digest are, for practical purposes, identical; a single changed byte changes the digest completely. With `-c` it reads a list of digests made earlier and checks each file against it. It is fast, and good for catching accidental damage or comparing against a published checksum, but it is broken for security: two different files can be made with the same MD5. Use `sha256sum` or stronger when someone could tamper with the file.

### Background
The macOS manual for its checksum tools says the `-sum` forms, including `md5sum`, have command-line options and an output format that match those of the similarly named GNU utilities, which is why the three implementations here are interchangeable for the common uses. Today the name covers GNU md5sum, the one macOS ships, and the applet in BusyBox.

### Examples
Print the digest of some text; `printf` is used because `echo` adds a newline that changes the digest:

```sh short="Print the digest of some text (printf, not echo)"
printf abc | md5sum
```

Print only the first eight characters of the digest:

```sh
printf abc | md5sum | cut -c1-8
```

Print the digest of a file; the file name follows the digest:

```sh
printf abc > f.txt && md5sum f.txt
```

Save digests to a file and check them later with `-c`; each file reports `OK`:

```sh
printf abc > f.txt && md5sum f.txt > SUMS && md5sum -c SUMS
```

Detect a changed file; `-c` prints `FAILED` and exits with a non-zero status, so the `echo` runs:

```sh
printf abc > f.txt && md5sum f.txt > SUMS && printf x >> f.txt && md5sum -c SUMS 2> /dev/null || echo changed
```

### Option: -c
Read digests and file names from a file, and check that each file still matches.

### Common mistakes
Hashing `echo text` and `printf text` gives different digests, because `echo` adds a newline; make sure you hash what you mean to. A digest proves a file is unchanged only if the digest itself came from a source you trust, for example over a different channel than the file. The text mode and binary mode markers (`-t` and `-b`) make no difference on Linux and macOS. Checking with `-c` needs the file names in the list to be reachable from where you run it.

### Portability
Only the plain form and `-c` are in the Common examples; the output line `digest  file` (two spaces) is the same on all the implementations here. GNU adds `--tag`, `--quiet` and `--ignore-missing`, which are also in the macOS tools, and `--zero`. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU md5sum
GNU md5sum is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### Option: --tag
Print BSD-style lines such as `SHA256 (file) = digest` instead of the default.

### Option: --quiet
With `-c`, do not print `OK` for files that match.

## Apple
### About macOS md5sum
This is the md5sum that macOS ships, in `/sbin`. Its options and output format are written to match GNU's.

### Apple notes
On macOS the same digests are also available through `shasum -a` and the BSD-style `md5` and `sha256` commands, which print lines like `SHA256 (file) = digest`; they are not covered on this page.

## BusyBox
### About BusyBox md5sum
BusyBox bundles md5sum as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `md5sum` supports printing digests and `-c`. GNU long options depend on the build.
