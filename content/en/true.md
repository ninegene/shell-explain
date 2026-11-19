---
schema_version: 1
id: true
title: true
summary: Do nothing and succeed.
aliases: [gtrue]
tags: [shell, scripts]
implementations: [bash, gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-true
  - https://www.gnu.org/software/coreutils/manual/html_node/true-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=true&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/true.1.html
---
## Common
### Overview
`true` does nothing and exits with status 0, which means success. It is a placeholder for a command that must succeed, for example as the condition of an endless loop, or to keep a script going when a command may fail (`cmd || true`). Its partner `false` always fails.

### Background
The `true` command exists so that a shell script can say "success" without doing anything. Today the name covers the Bash builtin, GNU true, the true that macOS ships, and the true applet in BusyBox. They agree on exiting with status 0 and printing nothing.

### Examples
Show that the exit status is 0:

```sh
true; echo $?
```

Run the second command because the first succeeded:

```sh
true && echo yes
```

Keep a script going when a command may fail:

```sh
false || true; echo continued
```

Use it as the condition of a loop that stops by itself:

```sh
n=0; while true; do n=$((n + 1)); test "$n" -ge 3 && break; done; echo "$n"
```

Use it as an empty `then` branch, so the `if` has a body:

```sh
if true; then true; fi; echo ok
```

### Common mistakes
`true` ignores its arguments, so `true anything` still succeeds. In a script with `set -e`, a failing command stops the script unless it is followed by `|| true`. The shell also has the `:` builtin, which does the same job; use whichever your team knows. A loop `while true` needs a `break` or an `exit`, or it never ends.

### Portability
The Common examples ran in the Bash builtin, GNU true, macOS `/usr/bin/true` and BusyBox true with the same output. GNU `true --help` and `--version` print text, and other implementations ignore them. The Common examples were executed on Bash 5.2, GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash true
In Bash, `true` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it.

## GNU
### About GNU true
GNU true is part of GNU coreutils, and this page's checks ran coreutils 9.11. This view is the executable, which runs when you call it by path, as in `/bin/true`, or when the shell has no builtin of that name.

## Apple
### About macOS true
This is the true in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/true`.

## BusyBox
### About BusyBox true
BusyBox bundles true as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `true` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
