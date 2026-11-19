---
schema_version: 1
id: nice
title: nice
summary: Run a command with a lower or higher scheduling priority.
aliases: [gnice]
tags: [process, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/nice-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=nice&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnunice.1.html
---
## Common
### Overview
`nice` starts a command with a changed "niceness", a number that tells the system how considerate the command should be to other programs. A higher number is nicer, which means a lower priority: the command gets the processor when nothing more important wants it. Use it to run a heavy job without slowing the machine down.

### Background
A `nice` utility appeared in Version 4 AT&T UNIX, according to the macOS manual, which also says it conforms to POSIX.1-2001. Today the name covers GNU nice, the nice that macOS ships, and the nice applet in BusyBox. They agree on `-n` with a number.

### Examples
Run a command with a niceness of 5; its output is unchanged:

```sh
nice -n 5 echo hi
```

Start a heavy job with a lower priority:

```sh
nice -n 10 sh -c 'echo working'
```

Check that the shell found `nice` and passed the number on, by running a program that exits with a status:

```sh short="Check that nice passes on the program's status"
nice -n 5 sh -c 'exit 3'; echo $?
```

### Option: -n
Add this number to the current niceness; positive is nicer, and only an administrator can use a negative one.

### Common mistakes
A normal user can only raise the number, not lower it again, so a command started with `nice -n 10` cannot get its priority back. `nice` changes the start of a command; `renice` changes one that is already running. Niceness matters only when the machine is busy: on an idle computer a nice command runs at full speed.

### Portability
Only `nice -n NUMBER command` is in the Common examples. GNU `nice` with no command prints the current niceness, which macOS does not: it prints a usage message instead. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU nice
GNU nice is part of GNU coreutils, and this page's checks ran coreutils 9.11. With no command, it prints the current niceness.

### GNU examples
Check that `-n 5` raises the niceness by five:

```sh
a=$(nice); b=$(nice -n 5 nice); echo $((b - a))
```

## Apple
### About macOS nice
This is the nice in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `nice` needs a command: with none it prints a usage message and exits.

## BusyBox
### About BusyBox nice
BusyBox bundles nice as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `nice` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
