---
schema_version: 1
id: mknod
title: mknod
summary: Create special files such as named pipes and device nodes.
aliases: [gmknod]
tags: [files, devices, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/mknod-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=mknod&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`mknod` creates a special file: a named pipe, or a device node that stands for a piece of hardware. Device nodes need administrator rights, and most systems now create them automatically, so in everyday use the only thing people make with `mknod` is a named pipe, which `mkfifo` does more simply. The examples make only a pipe.

### Background
Today the name covers GNU mknod, the mknod that macOS ships (in `/sbin`), and the mknod applet in BusyBox. The macOS manual gives no date for it. They agree on `mknod NAME p` for a named pipe.

### Examples
Create a named pipe with the type letter `p`, and check that it is one:

```sh
mknod p p && test -p p && echo fifo
```

See the `p` at the start of the long listing:

```sh
mknod p p && ls -l p | cut -c1
```

### Common mistakes
Device nodes (`c` for character, `b` for block) need a major and a minor number and administrator rights; getting the numbers wrong gives a node for the wrong device. Do not make nodes in a shared folder. For a named pipe, `mkfifo` is shorter and has the same effect.

### Portability
Only `NAME p` for a named pipe is in the Common examples. The syntax for devices differs between systems; read the manual of your system and prefer `mkfifo` for pipes. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU mknod
GNU mknod is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS mknod
This is the mknod in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox mknod
BusyBox bundles mknod as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `mknod` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
