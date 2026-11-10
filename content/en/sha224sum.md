---
schema_version: 1
id: sha224sum
title: sha224sum
summary: Compute or verify SHA-224 checksums of files.
aliases: [gsha224sum]
tags: [files, checksum, security]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/sha224sum-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=sha224sum&sektion=1
---
## Common
### Overview
`sha224sum` reads a file, or standard input, and prints a 224-bit digest of its contents as hexadecimal, followed by the file name. Two files with the same digest are, for practical purposes, identical; a single changed byte changes the digest completely. With `-c` it reads a list of digests made earlier and checks each file against it. It is a shorter variant of the SHA-2 family, used less often than SHA-256.

### Background
The macOS manual for its checksum tools says the `-sum` forms, including `sha224sum`, have command-line options and an output format that match those of the similarly named GNU utilities, which is why the three implementations here are interchangeable for the common uses. Today the name covers GNU sha224sum, the one macOS ships.

### Examples
Print the digest of some text; `printf` is used because `echo` adds a newline that changes the digest:

```sh
printf abc | sha224sum
```

Print only the first eight characters of the digest:

```sh
printf abc | sha224sum | cut -c1-8
```

Print the digest of a file; the file name follows the digest:

```sh
printf abc > f.txt && sha224sum f.txt
```

Save digests to a file and check them later with `-c`; each file reports `OK`:

```sh
printf abc > f.txt && sha224sum f.txt > SUMS && sha224sum -c SUMS
```

Detect a changed file; `-c` prints `FAILED` and exits with a non-zero status, so the `echo` runs:

```sh
printf abc > f.txt && sha224sum f.txt > SUMS && printf x >> f.txt && sha224sum -c SUMS 2> /dev/null || echo changed
```

### Option: -c
Read digests and file names from a file, and check that each file still matches.

### Common mistakes
Hashing `echo text` and `printf text` gives different digests, because `echo` adds a newline; make sure you hash what you mean to. A digest proves a file is unchanged only if the digest itself came from a source you trust, for example over a different channel than the file. The text mode and binary mode markers (`-t` and `-b`) make no difference on Linux and macOS. Checking with `-c` needs the file names in the list to be reachable from where you run it.

### Portability
Only the plain form and `-c` are in the Common examples; the output line `digest  file` (two spaces) is the same on all the implementations here. GNU adds `--tag`, `--quiet` and `--ignore-missing`, which are also in the macOS tools, and `--zero`. The Common examples were executed on GNU coreutils 9.11, and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU sha224sum
GNU sha224sum is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### Option: --tag
Print BSD-style lines such as `SHA256 (file) = digest` instead of the default.

### Option: --quiet
With `-c`, do not print `OK` for files that match.

## Apple
### About macOS sha224sum
This is the sha224sum that macOS ships, in `/sbin`. Its options and output format are written to match GNU's.

### Apple notes
On macOS the same digests are also available through `shasum -a` and the BSD-style `md5` and `sha256` commands, which print lines like `SHA256 (file) = digest`; they are not covered on this page.
