---
schema_version: 1
id: head
title: head
summary: Print the first lines or bytes of files.
aliases: [ghead]
tags: [text, files, logs]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/head-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=head&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`head` prints the beginning of each file you name, ten lines by default. Use it to peek at a large file or to keep the first results of a pipeline.

### Examples
Print the first ten lines (all of a short file):

```sh
head notes.txt
```

Print the first two lines:

```sh
head -n 2 notes.txt
```

Print the first line of several files; each is introduced by a header:

```sh
head -n 1 app.log logs/server.log
```

Print the first five bytes:

```sh
head -c 5 notes.txt
```

Keep only the first two lines of a pipeline's output:

```sh
printf '1\n2\n3\n4\n' | head -n 2
```

### Option: -n
Print this many lines instead of ten. The count is the next argument, as in `-n 2`.

### Option: -c
Print this many bytes instead of lines. A multibyte character can be cut in half.

### Common mistakes
A count like `-n 2` is a number of lines, not a line number. When several files are given, `head` prints a `==> name <==` header for each; use `-q` where supported to hide it.

### Portability
Only `-n` and `-c` with a positive count are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### GNU examples
Print everything except the last two lines. A negative count is a GNU extension:

```sh
head -n -2 notes.txt
```

## Apple
### Apple notes
macOS `head` is the BSD implementation. It does not accept a negative count such as `-n -2`; use `tail` or `sed` to drop trailing lines.

## BusyBox
### BusyBox notes
BusyBox `head` supports `-n` and `-c`; other options depend on the build. A negative count is not part of the common examples.
