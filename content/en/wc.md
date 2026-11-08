---
schema_version: 1
id: wc
title: wc
summary: Count lines, words, and bytes.
aliases: [gwc]
tags: [text, files, count]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=wc&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`wc` counts the lines, words, and bytes in each file, or in standard input. A line is counted at each newline, so a last line without a trailing newline is not counted.

### Examples
Show lines, words, and bytes for a file:

```sh
wc notes.txt
```

Count lines only:

```sh
wc -l notes.txt
```

Count words only:

```sh
wc -w notes.txt
```

Count bytes only:

```sh
wc -c notes.txt
```

Count the lines of a command's output; reading standard input prints no file name:

```sh
cat notes.txt | wc -l
```

Count several files; a total is added:

```sh
wc -l app.log logs/server.log
```

### Option: -l
Count newline characters, which is the number of lines.

### Option: -w
Count words, meaning runs of characters separated by whitespace.

### Option: -c
Count bytes. For characters in a multibyte encoding use `-m` where it is supported.

### Common mistakes
The column padding of the output differs between implementations, so do not compare it as text; split on whitespace instead. Counting a file by name prints the name, while counting standard input does not, which matters when you capture the number in a script.

### Portability
`-l`, `-w`, and `-c` are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). Output padding and the default `-m` behavior differ.

## GNU
### GNU examples
Print the length of the longest line. `-L` is a GNU option:

```sh
wc -L notes.txt
```

## Apple
### Apple notes
macOS `wc` is the BSD implementation. It pads each count to a fixed width, so `wc -l` output has leading spaces.

## BusyBox
### BusyBox notes
BusyBox `wc` supports `-l`, `-w`, and `-c`; other options depend on the build.
