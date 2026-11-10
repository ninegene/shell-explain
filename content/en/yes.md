---
schema_version: 1
id: yes
title: yes
summary: Repeat a line of text forever, until the reader stops.
aliases: [gyes]
tags: [text, pipes, scripts]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/yes-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=yes&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`yes` prints `y` followed by a newline, again and again, until the program reading its output stops. It exists to answer "y" to prompts that ask for confirmation, and it is a handy source of endless test input. Give it text and it repeats that text instead.

### Background
The `yes` command appeared in Version 7 AT&T UNIX, according to the macOS manual. Today the name covers GNU yes, the yes that macOS ships, and the yes applet in BusyBox. They agree on repeating a line without end and on repeating a single word given as an argument.

### Examples
Print three lines of `y`; `head` stops the endless output after three lines:

```sh
yes | head -n 3
```

Repeat a word instead of `y`:

```sh
yes ok | head -n 2
```

Make 1000 lines of the same text and count them:

```sh
yes data | head -n 1000 | grep -c data
```

Answer a confirmation prompt automatically, with a program that asks twice:

```sh
yes | (read a; read b; echo "$a$b")
```

### Common mistakes
`yes` never ends by itself, so on its own it floods the screen until you press Ctrl+C; pipe it into something that stops, such as `head`. Answering every prompt with `y` is risky for a command that deletes or overwrites, so use it only when you have already checked what the command will do. Many programs have their own option for this, such as `-y` or `--yes`, and that is safer than a pipe.

### Portability
Printing `y` or the given words is the same on all three implementations. The Common examples avoid long runs, because the output is unlimited. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU yes
GNU yes is part of GNU coreutils, and this page's checks ran coreutils 9.11. It joins several arguments with spaces and repeats them as one line, and it stops quietly when the reader closes the pipe.

### GNU examples
Repeat several words as one line:

```sh
yes one two | head -n 2
```

## Apple
### About macOS yes
This is the yes in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `yes` repeats only the first argument: `yes one two` prints `one`. Use `printf` or `echo` in a loop, or `yes "one two"`, for a line with spaces.

## BusyBox
### About BusyBox yes
BusyBox bundles yes as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `yes` repeats `y` or the text you give. Pass a single quoted argument for text with spaces.
