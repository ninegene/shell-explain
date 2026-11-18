---
schema_version: 1
id: mkdir
title: mkdir
summary: Create directories.
aliases: [gmkdir]
tags: [files, directories]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/mkdir-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=mkdir&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`mkdir` creates the directories you name. It fails if a directory already exists or if a parent directory is missing, unless you use `-p`.

### Background
A `mkdir` command appeared in Version 1 Unix, according to the FreeBSD and macOS manuals. Today the name covers GNU mkdir, the mkdir that macOS ships, and the mkdir applet in BusyBox. They agree on `-p` and `-m` with an octal mode, and differ on the wording and availability of `-v`.

### Examples
Create one directory and confirm it exists:

```sh
mkdir reports && ls -d reports
```

Create a directory and any missing parents; `-p` is also quiet when the directory already exists:

```sh short="Create a directory and missing parents (-p)"
mkdir -p reports/2026/q1 && ls reports/2026
```

Create several directories in one command:

```sh
mkdir one two && ls -d one two
```

Create a directory that only its owner can use; `-m` sets the permission bits and `cut` shows the first ten characters of the long listing:

```sh
mkdir -m 700 private && ls -ld private | cut -c1-10
```

Make sure a directory exists without an error if it already does:

```sh
mkdir -p cache && echo ok
```

### Option: -p
Create missing parent directories and do not fail when the target already exists.

### Option: -m
Set the permission bits of the new directory, as in `-m 700`. Parents created by `-p` use the default permissions.

### Common mistakes
Without `-p`, `mkdir a/b/c` fails when `a/b` does not exist, and `mkdir logs` fails when `logs` exists. Quote names that contain spaces, or `mkdir my docs` creates two directories.

### Portability
Only `-p` and `-m` with an octal mode are in the Common examples. The verbose option `-v` prints different text on GNU and BSD, and not every BusyBox build has it. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU mkdir
GNU mkdir is part of GNU coreutils, and this page's checks ran coreutils 9.11. With `-v` it prints a sentence for every directory created.

### GNU examples
Print a line for every directory created; the wording is specific to GNU:

```sh
mkdir -pv reports/2026/q2
```

## Apple
### About macOS mkdir
This is the mkdir in `/bin` on macOS. With `-v` it prints just the path of each directory created.

### Apple notes
macOS `mkdir` is the BSD implementation. `-v` prints just the path of each directory created rather than GNU's sentence.

## BusyBox
### About BusyBox mkdir
BusyBox bundles mkdir as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `mkdir` supports `-p` and `-m`; `-v` and GNU long options depend on the build.
