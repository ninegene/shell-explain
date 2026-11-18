---
schema_version: 1
id: b2sum
title: b2sum
summary: Compute or verify BLAKE2 checksums of files.
aliases: [gb2sum]
tags: [files, checksum, security]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/b2sum-invocation.html
---
## Common
### Overview
`b2sum` prints a BLAKE2b digest of a file or standard input, in the same layout as `sha256sum`: the digest in hexadecimal, then the file name. BLAKE2 is fast and is considered secure. With `-c` it checks files against a list made earlier.

### Background
`b2sum` is part of GNU coreutils; the GNU manual documents it as the BLAKE2 checksum tool. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Print the first eight characters of the BLAKE2b digest of some text:

```sh short="Print the start of the BLAKE2b digest of some text"
printf abc | b2sum | cut -c1-8
```

Save a digest and check it later with `-c`:

```sh
printf abc > f.txt && b2sum f.txt > SUMS && b2sum -c SUMS
```

Use a shorter digest with `-l`, given in bits:

```sh
printf abc | b2sum -l 256 | cut -c1-8
```

### Option: -c
Read digests and file names from a file, and check that each file still matches.

### Option: -l
Set the digest length in bits; the default is 512.

### Common mistakes
The digest length changes the digest, so a list made with `-l 256` must be checked with `-l 256` as well. A digest is only as trustworthy as the place it came from. Because `b2sum` is not on macOS or in BusyBox, a script that must run everywhere should use `sha256sum`.

### Portability
This page covers GNU only. The same digests can be made on macOS with other tools, such as `openssl dgst -blake2b512`, which is not covered here. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU b2sum
GNU b2sum is part of GNU coreutils, and this page's checks ran coreutils 9.11.
