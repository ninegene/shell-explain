---
schema_version: 1
id: printf
title: printf
summary: Print text using a format string, with exact control over newlines and escapes.
aliases: [gprintf]
tags: [shell, text, scripts]
implementations: [bash, gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-printf
  - https://www.gnu.org/software/coreutils/manual/html_node/printf-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=printf&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuprintf.1.html
---
## Common
### Overview
`printf` prints text built from a format string and arguments, as the C function of the same name does. Unlike `echo`, it does not add a newline, understands backslash escapes the same way everywhere, and can pad, round and convert numbers. When the arguments outnumber the format, the format is used again, which makes a one-line list into many lines.

### Background
The `printf` command takes its format language from the C library function. Today the name covers the Bash builtin, GNU printf, the printf that macOS ships, and the printf applet in BusyBox. They agree on `%s`, `%d`, `%f`, `%b` and the usual backslash escapes.

### Examples
Print two values with a format; `\n` ends the line:

```sh
printf '%s-%s\n' a b
```

Round a number to two decimals, padded to five characters:

```sh
printf '%5.2f|\n' 3.14159
```

Convert a hexadecimal number to decimal:

```sh
printf '%d\n' 0x1f
```

Reuse the format for every extra argument, one line each:

```sh
printf '%s\n' a b c
```

Print a tab and a newline from escapes:

```sh
printf 'a\tb\n'
```

Expand escapes inside an argument with `%b`:

```sh
printf '%b' 'a\tb\n'
```

### Common mistakes
Never put untrusted text in the format position: `printf "$text"` treats `%` and backslashes in it as instructions. Use `printf "%s\n" "$text"` instead. The format needs `\n` for a newline; nothing is added for you. With no arguments `%s` prints nothing and `%d` prints 0, which can hide a missing variable. Floating-point results can differ in the last digits between implementations.

### Portability
Only `%s`, `%d` with a hexadecimal argument, `%5.2f`, `%b` and the tab and newline escapes are in the Common examples. `%q` for shell quoting and `-v` to store the result in a variable are Bash builtin features, not in the executables. In a plain `sh` script the builtin may behave differently from both Bash and the executable; the macOS `/bin/sh`, for example, printed `-n` as text for `echo -n` in the check. The Common examples were executed on Bash 5.2, GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash printf
In Bash, `printf` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it. `-v VAR` stores the result in a variable instead of printing it, and `%q` quotes a string for the shell.

### Bash examples
Store the result in a variable with `-v` instead of printing it:

```sh
printf -v greeting '%s' hi && echo "$greeting"
```

Quote a string with a space so the shell can read it back with `%q`:

```sh
printf '%q\n' 'a b'
```

### Option: -v
Assign the output to this variable instead of printing it.

## GNU
### About GNU printf
GNU printf is part of GNU coreutils, and this page's checks ran coreutils 9.11. This view is the executable, which runs when you call it by path, as in `/bin/printf`, or when the shell has no builtin of that name.

## Apple
### About macOS printf
This is the printf in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/printf`.

## BusyBox
### About BusyBox printf
BusyBox bundles printf as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `printf` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
