---
schema_version: 1
id: tty
title: tty
summary: Print the name of the terminal connected to standard input.
aliases: [gtty]
tags: [terminal, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/tty-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=tty&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`tty` prints the file name of the terminal your input comes from, such as `/dev/ttys003` or `/dev/pts/0`. When input is not a terminal, for example from a pipe or a file, it prints `not a tty` and exits with status 1. Scripts use it to decide whether they are running interactively.

### Background
A `tty` command appeared in Version 1 AT&T UNIX, according to the macOS manual, which also says it is POSIX.2 compatible. Today the name covers GNU tty, the tty that macOS ships, and the tty applet in BusyBox. They agree on printing the terminal name or `not a tty`, and on the exit status.

### Examples
Ask when input is not a terminal; `/dev/null` is not one, so the answer is `not a tty` and the status is 1:

```sh short="Ask when input is not a terminal: status 1"
tty < /dev/null; echo $?
```

Use the exit status in a test:

```sh
tty < /dev/null || echo not a terminal
```

Branch on it in a script:

```sh
if tty > /dev/null < /dev/null; then echo interactive; else echo batch; fi
```

### Common mistakes
`tty` looks at standard input, not output, so `tty | cat` in a terminal still reports the terminal, while `cat file | tty` reports `not a tty`. To test standard output use `test -t 1`, and to test input in a script `test -t 0`, which needs no process. Inside a program run from cron, a container or an editor, there may be no terminal at all.

### Portability
Only plain `tty` and its exit status are in the Common examples, run with input from `/dev/null` so the answer is the same everywhere. The silent option `-s` is described in the sections below. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU tty
GNU tty is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-s` prints nothing and only sets the exit status.

### GNU examples
Use the silent form `-s` in a test:

```sh
tty -s < /dev/null || echo not a terminal
```

### Option: -s
Print nothing and exit with status 0 only when input is a terminal.

## Apple
### About macOS tty
This is the tty in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Use the silent form `-s` in a test:

```sh
tty -s < /dev/null || echo not a terminal
```

### Option: -s
Print nothing and exit with status 0 only when input is a terminal.

## BusyBox
### About BusyBox tty
BusyBox bundles tty as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
In the 1.37.0 check, BusyBox `tty -s` was not silent: it printed `tty: ignoring all arguments` and `not a tty`. Redirect the output or use `test -t 0` instead.
