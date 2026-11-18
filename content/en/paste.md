---
schema_version: 1
id: paste
title: paste
summary: Join lines from several files side by side, or join one file's lines into one.
aliases: [gpaste]
tags: [text, columns, files]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/paste-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=paste&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`paste` glues files together line by line: line 1 of each file goes on one output line, separated by a tab, then line 2, and so on. It is the opposite of `cut`, which pulls columns apart. With `-s` it turns the lines of one file into a single line, and a `-` reads standard input.

### Background
A `paste` command appeared in Version 7 AT&T UNIX/32V, according to the macOS manual. Today the name covers GNU paste, the paste that macOS ships, and the paste applet in BusyBox. They agree on the tab default, `-d` and `-s`.

### Examples
Put two files side by side, separated by a tab:

```sh
printf 'a\nb\n' > names.txt && printf '1\n2\n' > nums2.txt && paste names.txt nums2.txt
```

Use a comma as the separator with `-d`:

```sh
printf 'a\nb\n' > names.txt && printf '1\n2\n' > nums2.txt && paste -d, names.txt nums2.txt
```

Join all the lines of a file into one with `-s`:

```sh
printf 'a\nb\nc\n' > names.txt && paste -s names.txt
```

Join them with commas instead of tabs:

```sh
printf 'a\nb\nc\n' > names.txt && paste -sd, names.txt
```

Arrange a stream into columns; each `-` takes the next line of standard input, so three of them make three columns:

```sh short="Arrange a stream into columns, one - per column"
seq 6 | paste - - -
```

Join numbers into an expression for a calculator, with a plus sign as the separator:

```sh short="Join numbers into an expression, separated by +"
seq 3 | paste -sd+ -
```

### Option: -d
Use these characters, one after the other, instead of a tab between the joined lines.

### Option: -s
Paste one file at a time: all its lines are joined into one output line.

### Common mistakes
`paste` does not match lines by content; it just pairs the nth line of each file, so files of different lengths leave empty columns. To join by a shared key, use `join`. The default separator is a tab, which looks like spaces on a screen but is a different character when you search or sort the output.

### Portability
Only `-d`, `-s` and `-` for standard input are in the Common examples. Behavior with a multi-character `-d` list is the same idea on every system, but check it on yours before relying on it. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU paste
GNU paste is part of GNU coreutils, and this page's checks ran coreutils 9.11. It adds `-z` for NUL-separated lines.

### Option: -z, --zero-terminated
Treat lines as ending in a NUL byte instead of a newline.

## Apple
### About macOS paste
This is the paste in `/usr/bin` on macOS. It follows the BSD implementation.

### Apple notes
macOS `paste` has `-d` and `-s` as above; GNU's `-z` is not available.

## BusyBox
### About BusyBox paste
BusyBox bundles paste as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `paste` supports `-d` and `-s`. GNU long options depend on the build.
