---
schema_version: 1
id: coreutils
title: coreutils
summary: The GNU coreutils package, and its single multi-call program.
aliases: [gcoreutils]
tags: [package, system]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/index.html
---
## Common
### Overview
GNU coreutils is the package that supplies many of the commands on this site: `ls`, `cp`, `mv`, `cat`, `sort`, `date` and about a hundred more. The package can be built as one program, named `coreutils`, that behaves as any of its commands: you name the command as the first argument, or give it a link whose name is the command. Distributions that build this way save disk space, because the commands share code.

### Background
The GNU manual describes coreutils as the basic file, shell and text manipulation utilities of the GNU operating system. The multi-call `coreutils` program is built only on request, and many systems do not install it. macOS ships BSD tools instead, and BusyBox has its own multi-call program named `busybox`, so this page has only a GNU view.

### Examples
Run a command through the multi-call program:

```sh
coreutils --coreutils-prog=echo hello
```

Use another command the same way:

```sh
coreutils --coreutils-prog=basename logs/server.log .log
```

Run it with a name the program knows, as in a link:

```sh
coreutils --coreutils-prog=true; echo $?
```

### Common mistakes
There is no standard command for "the coreutils package": the package is a set of separate commands, and an individual command can come from another package, as on Alpine Linux, where `groups` and `uptime` are not part of the coreutils package. A system with the multi-call program installed may still have the commands as separate files, so `which ls` tells you more than `coreutils` does. Do not confuse GNU coreutils with the BSD tools on macOS, whose options differ, and the `g`-prefixed Homebrew names point to the GNU ones.

### Portability
This page covers GNU only. The Homebrew coreutils used for the host checks does not install the multi-call program, so these examples ran in the Linux guest only. The Common examples were executed on GNU coreutils 9.11 in the Linux guest only (see the development guide for the exact check).

## GNU
### About GNU coreutils
GNU coreutils is part of GNU coreutils, and this page's checks ran coreutils 9.11. This is the multi-call program, which the Alpine package installs as `coreutils`.
