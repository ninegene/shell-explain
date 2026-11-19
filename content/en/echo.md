---
schema_version: 1
id: echo
title: echo
summary: Print text followed by a newline.
aliases: [gecho]
tags: [shell, text, scripts]
implementations: [bash, gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bash-Builtins.html#index-echo
  - https://www.gnu.org/software/coreutils/manual/html_node/echo-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=echo&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuecho.1.html
---
## Common
### Overview
`echo` prints its arguments, separated by single spaces, and ends with a newline. It is the simplest way to show a message or to put a known line into a pipe or a file. Every shell has its own `echo` builtin, and there is also an `/bin/echo` program; they agree on plain text and disagree on options and on backslashes, which is why `printf` is safer for anything but plain text.

### Background
The macOS manual for `echo` says the utility conforms to POSIX, whose specification leaves option handling and backslash escapes to the implementation. Today the name covers the Bash builtin, GNU echo, the echo that macOS ships, and the echo applet in BusyBox. They agree on printing the arguments and a newline.

### Examples
Print a word:

```sh
echo hello
```

Print several words; the arguments are joined with single spaces:

```sh
echo one two
```

Keep extra spaces by quoting the text:

```sh
echo "a   b"
```

Print the result of an arithmetic expression:

```sh
echo "$((1 + 2))"
```

Write a line to a file and read it back:

```sh
echo hi > out.txt && cat out.txt
```

Print an empty line by giving no arguments:

```sh
echo
```

### Common mistakes
Without quotes, the shell splits and expands the text first, so `echo *` prints the file names and extra spaces vanish. `-n` (no newline) and `-e` (interpret backslashes) are not portable: the builtin of a given shell, GNU `echo` and macOS `/bin/echo` treat them differently, and some print `-n` as text. Use `printf` for anything with backslashes, no newline or a leading dash. `echo` cannot reliably print text that starts with `-`.

### Portability
Only plain arguments, quoting and redirection are in the Common examples, which gave the same output in the Bash builtin, GNU echo, macOS `/bin/echo` and BusyBox echo. In a plain `sh` script the builtin may behave differently from both Bash and the executable; the macOS `/bin/sh`, for example, printed `-n` as text for `echo -n` in the check. The Common examples were executed on Bash 5.2, GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash echo
In Bash, `echo` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it. `-n` leaves out the newline and `-e` turns on backslash escapes such as `\t`; by default backslashes are printed as written.

### Bash examples
Leave out the newline with `-n`:

```sh
echo -n hi; echo
```

Interpret backslash escapes with `-e`; the tab is shown by `od`:

```sh
echo -e 'a\tb' | od -An -c
```

### Option: -n
Do not print the newline at the end.

### Option: -e
Interpret backslash escapes such as `\t` and `\n`.

## GNU
### About GNU echo
GNU echo is part of GNU coreutils, and this page's checks ran coreutils 9.11. This view is the executable, which runs when you call it by path, as in `/bin/echo`, or when the shell has no builtin of that name. It also has `-n` and `-e`, and `-E` to turn escapes off again.

### GNU examples
Interpret backslash escapes with `-e`:

```sh
echo -e 'a\tb' | od -An -c
```

### Option: -n
Do not print the newline at the end.

### Option: -e
Interpret backslash escapes such as `\t` and `\n`.

## Apple
### About macOS echo
This is the echo in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/echo`. It supports `-n`, but `-e` is not recognized and is printed as text.

### Apple examples
Leave out the newline with `-n`:

```sh
echo -n hi; echo
```

### Option: -n
Do not print the newline at the end.

### Apple notes
macOS `/bin/echo` has no `-e`: `/bin/echo -e 'a\tb'` prints `-e a\tb` as text. Use `printf` for escapes.

## BusyBox
### About BusyBox echo
BusyBox bundles echo as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `echo` prints its arguments and a newline. Its `-n`, `-e` and `-E` options depend on the build.
