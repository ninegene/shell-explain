---
schema_version: 1
id: unzip
title: unzip
summary: List, test and unpack the contents of a .zip archive.
aliases: []
tags: [files, archive, compression]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://infozip.sourceforge.net/UnZip.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/unzip.1.html
---
## Common
### Overview
`unzip` takes files out of a `.zip` archive. Without options it unpacks everything into the current folder, so use `-d` to choose a folder. `-l` lists the contents, `-p` prints a file to the screen and `-t` tests the archive.

### Background
`unzip` is the Info-ZIP program, and the macOS build is UnZip 6.00 with changes by Apple. The tool tested here is that build.

### Examples
Print one file from an archive to the screen, without unpacking it:

```sh
printf data > f.txt && zip -q a.zip f.txt && unzip -p a.zip f.txt
```

Unpack into a folder with `-d`:

```sh
printf data > f.txt && zip -q a.zip f.txt && unzip -q a.zip -d out && cat out/f.txt
```

List only the names with `-Z1`:

```sh
printf data > f.txt && zip -q a.zip f.txt && unzip -Z1 a.zip
```

Test that the archive is not damaged:

```sh
printf data > f.txt && zip -q a.zip f.txt && unzip -t a.zip | tail -n 1
```

List the contents with sizes, ending with a summary line:

```sh
printf data > f.txt && zip -q a.zip f.txt && unzip -l a.zip | tail -n 1
```

### Option: -d
Unpack into this folder.

### Option: -l
List the contents.

### Option: -p
Print the files to standard output.

### Option: -t
Test the archive.

### Option: -o
Replace existing files without asking.

### Option: -q
Quiet.

### Common mistakes
Unpacking an archive from someone else can write files over yours, or outside the folder with `../` in a name on old versions, so list it with `-l` first and unpack with `-d` into a new folder. When a file exists, `unzip` asks whether to replace it; `-o` replaces without asking and `-n` never replaces. A `.zip` made on another system can have names in another character set. Passwords given with `-P` are visible in the process list.

### Portability
The options above are the same on macOS and Linux with Info-ZIP. The table of the list and the text of messages differ a little between versions, so the example that shows the summary line is compared with blanks collapsed. The examples were executed on UnZip 6.00 (macOS build). The Common examples were executed on UnZip 6.00 (macOS build) (see the development guide for the exact check).

## Upstream
### About upstream unzip
This is the upstream Info-ZIP `unzip` program. It is one program everywhere; some Linux distributions need the package installed.
