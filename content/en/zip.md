---
schema_version: 1
id: zip
title: zip
summary: Pack files and folders into a .zip archive.
aliases: []
tags: [files, archive, compression]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://infozip.sourceforge.net/Zip.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/zip.1.html
---
## Common
### Overview
`zip` collects files into one `.zip` archive and compresses them. The format is understood on every operating system, so it is the usual way to send a folder to someone on Windows. `zip` adds to an existing archive, and `unzip` takes things out.

### Background
`zip` is the Info-ZIP program, and the macOS build is Info-ZIP Zip 3.0 with changes by Apple. The tool tested here is that build.

### Examples
Make an archive of one file; `-q` keeps it quiet:

```sh
printf data > f.txt && zip -q a.zip f.txt && unzip -Z1 a.zip
```

Add a folder and everything inside it with `-r`:

```sh
mkdir d && printf x > d/g && zip -qr b.zip d && unzip -Z1 b.zip
```

Store files without their folder names, with `-j` (junk the paths):

```sh
mkdir d && printf x > d/g && zip -q -j c.zip d/g && unzip -Z1 c.zip
```

Add a file to an archive that already exists; the list shows both:

```sh
printf a > one.txt && printf b > two.txt && zip -q a.zip one.txt && zip -q a.zip two.txt && unzip -Z1 a.zip
```

### Option: -r
Include the contents of folders.

### Option: -j
Do not store the folder part of names.

### Option: -q
Quiet: no list of files added.

### Option: -d
Delete entries from an archive.

### Common mistakes
Without `-r`, `zip` stores a folder as an empty entry and none of its files. `zip` adds to an existing archive and does not replace it, so old entries stay; delete the archive first for a fresh one, or use `-d` to remove entries. The archive does not keep Unix owners, and permissions only as far as the version supports. Names with `*` or spaces need quotes. A password given with `-P` is visible to other users in the process list, and the old zip encryption is weak: do not use it for secrets.

### Portability
The options above are the same in the Info-ZIP versions on macOS and Linux. Other programs that make `.zip` files, such as `bsdtar` and Python's `zipfile`, have different options. The examples were executed on Info-ZIP Zip 3.0 (macOS build). The Common examples were executed on Info-ZIP Zip 3.0 (macOS build) (see the development guide for the exact check).

## Upstream
### About upstream zip
This is the upstream Info-ZIP `zip` program. It is one program everywhere; some Linux distributions need the package installed.
