---
schema_version: 1
id: false
title: false
summary: Do nothing and fail.
aliases: [gfalse]
tags: [shell, scripts]
implementations: [bash, gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-false
  - https://www.gnu.org/software/coreutils/manual/html_node/false-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=false&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`false` does nothing and exits with status 1, which means failure. It is useful for testing how a script behaves when a step fails, or as a condition that is never true. Its partner `true` always succeeds.

### Background
The `false` command exists so that a shell script can say "failure" without doing anything. Today the name covers the Bash builtin, GNU false, the false that macOS ships, and the false applet in BusyBox. They agree on exiting with a non-zero status and printing nothing.

### Examples
Show that the exit status is 1:

```sh
false; echo $?
```

Run the `||` branch because the first command failed:

```sh
false || echo failed
```

Choose the `else` branch of an `if`:

```sh
if false; then echo a; else echo b; fi
```

The `&&` branch is skipped after a failure:

```sh
false && echo skipped; echo end
```

### Common mistakes
`false` ignores its arguments. Do not treat the exact number 1 as part of its contract in a portable script: test for "not zero", as `if ! cmd` does. In a script with `set -e`, a bare `false` stops the script, which is a handy way to see whether `set -e` is on. `false` never prints anything, so a failure is silent.

### Portability
The Common examples ran in the Bash builtin, GNU false, macOS `/usr/bin/false` and BusyBox false with the same output. The Common examples were executed on Bash 5.2, GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash false
In Bash, `false` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it.

## GNU
### About GNU false
GNU false is part of GNU coreutils, and this page's checks ran coreutils 9.11. This view is the executable, which runs when you call it by path, as in `/bin/false`, or when the shell has no builtin of that name.

## Apple
### About macOS false
This is the false in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/false`.

## BusyBox
### About BusyBox false
BusyBox bundles false as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `false` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
