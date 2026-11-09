---
schema_version: 1
id: cat
title: cat
summary: Print files, or standard input, one after another.
aliases: [gcat]
tags: [text, files, pipes]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
related_guides: [big-picture]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/cat-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=cat&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`cat` writes the contents of the files you name to standard output, in order. With no file, or with `-`, it copies standard input. Use it to view a short file, join files, or feed a pipeline.

### Background
The name is short for concatenate: `cat` joins files and writes them to standard output. A `cat` utility appeared in Version 1 Unix, and Dennis Ritchie designed and wrote the first manual page, according to the FreeBSD and macOS manuals. Today the name covers GNU cat, the cat that macOS ships, and the cat applet in BusyBox. They agree on plain concatenation and the numbering example below, which ran on all three, and differ on options such as `-A` and `-s`.

### Examples
Print a file:

```sh
cat notes.txt
```

Join two files, in the order given:

```sh
cat app.log notes.txt
```

Number every output line:

```sh
cat -n notes.txt
```

Mix standard input and a file; `-` stands for standard input:

```sh
echo 'from stdin' | cat - notes.txt
```

Copy a file by redirecting the output:

```sh
cat notes.txt > copy.txt
```

### Option: -n
Number all output lines starting at 1. The number and a tab or spaces precede each line; exact padding differs by implementation.

### Common mistakes
Redirecting to a file you are also reading, as in `cat notes.txt > notes.txt`, empties the file before `cat` reads it. `cat file | command` is usually unnecessary: most commands read a file name or `<` redirection directly. For large files prefer a pager or `head`/`tail`.

### Portability
Only the options listed here are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). Options such as `-s`, `-A`, `-e`, and `-t` are not uniform; check the installed manual.

## GNU
### About GNU cat
GNU cat is part of GNU coreutils, and this page's checks ran coreutils 9.11. It has the extra options named in the notes below.

### GNU notes
GNU `cat` adds options such as `-A` (show line ends and tabs), `-s` (squeeze repeated blank lines), and long spellings like `--number`. Confirm that a macOS or BusyBox build provides an option before relying on it.

## Apple
### About macOS cat
This is the cat in `/bin` on macOS. GNU cat from a package manager may be installed as `gcat`, so the executable path, not the operating system, tells you which one you are running.

### Apple notes
macOS `cat` is the BSD implementation, normally `/bin/cat`; GNU `cat` installed by a package manager may be named `gcat`. Long options such as `--number` are GNU only; use the short spellings.

## BusyBox
### About BusyBox cat
BusyBox bundles cat as an applet, and options beyond plain concatenation depend on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `cat` options such as `-n` depend on the build configuration. Run `cat --help` in the actual environment before relying on anything beyond plain concatenation.
