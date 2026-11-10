---
schema_version: 1
id: mkfifo
title: mkfifo
summary: Create a named pipe, a file that connects two programs.
aliases: [gmkfifo]
tags: [files, pipes, ipc]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/mkfifo-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=mkfifo&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`mkfifo` creates a named pipe, also called a FIFO: a special file that is a pipe. One program writes into it and another reads from it, like `|` but with a name on the file system, so the two programs do not have to be started together in one line. The writer waits until a reader opens the file, and the reader waits for a writer.

### Background
The `mkfifo` command appeared in 4.4BSD, according to the macOS manual, which also says it is expected to be POSIX.2 compliant. Today the name covers GNU mkfifo, the mkfifo that macOS ships, and the mkfifo applet in BusyBox. They agree on creating the pipe and on `-m` for the mode.

### Examples
Create a named pipe and check that its type is a pipe:

```sh
mkfifo p && test -p p && echo fifo
```

See that it shows as a pipe in a long listing; the first character is `p`:

```sh
mkfifo p && ls -l p | cut -c1
```

Send text through it: a background writer and a reader in the foreground:

```sh
mkfifo p && (echo hi > p &) && cat p
```

### Option: -m
Set the permission bits of the pipe, as in `-m 600`.

### Common mistakes
Opening a named pipe blocks until the other side opens it, so `cat p` with no writer waits forever; run the writer in the background or in another terminal. Data in a pipe is not stored: when the last reader closes it, unread data is lost. A named pipe stays on disk until you remove it with `rm`, and an empty one looks like a file of size 0.

### Portability
Only creating a pipe is in the Common examples. GNU adds `-Z` for a security context. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU mkfifo
GNU mkfifo is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS mkfifo
This is the mkfifo in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox mkfifo
BusyBox bundles mkfifo as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `mkfifo` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
