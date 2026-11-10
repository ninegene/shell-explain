---
schema_version: 1
id: test
title: test
summary: Check files, strings and numbers, and report with the exit status.
aliases: [gtest, "["]
tags: [shell, scripts, conditions]
implementations: [bash, gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-test
  - https://www.gnu.org/software/coreutils/manual/html_node/test-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=test&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`test` evaluates a condition and reports the result with its exit status: 0 for true, 1 for false. It checks whether a file exists or is a folder, whether a string is empty or equal to another, and how two numbers compare. The same command is also written `[ condition ]`, with the closing bracket as the last argument; that is what `if [ -f file ]` uses.

### Background
The `test` command, also written `[`, is how shell scripts have always asked questions about files and values. Today the name covers the Bash builtin, GNU test, the test that macOS ships, and the test applet in BusyBox. They agree on the file tests `-e`, `-f`, `-d`, `-s`, `-r`, the string tests `-z`, `-n`, `=` and the number tests `-eq`, `-lt`, `-gt`.

### Examples
Check that a folder exists:

```sh
test -d logs && echo dir
```

Check that a regular file exists:

```sh
test -f notes.txt && echo file
```

Check that something does not exist:

```sh
test -e missing.txt || echo missing
```

Check that a file is not empty:

```sh
test -s notes.txt && echo nonempty
```

Check that a string is empty:

```sh
test -z "" && echo empty
```

Compare strings with a single equals sign:

```sh
test "a" = "a" && echo equal
```

Compare numbers with `-lt`:

```sh
test 3 -lt 5 && echo yes
```

Negate a test with `!`:

```sh
test ! -e missing.txt && echo not there
```

Write the same test with brackets; the spaces inside them are required:

```sh
[ -n "x" ] && echo set
```

### Option: -e
True if the path exists, whatever its type.

### Option: -f
True if the path is a regular file.

### Option: -d
True if the path is a folder.

### Option: -s
True if the file exists and is not empty.

### Option: -z
True if the string is empty.

### Option: -n
True if the string is not empty.

### Common mistakes
Spaces matter: `[ -f file ]` works, and `[-f file]` is an error. Quote variables, as in `test "$a" = "$b"`, or an empty value shifts the arguments. Use `=` for strings and `-eq` for numbers: `[ 10 = 10.0 ]` is false and `[ 10 -eq 10 ]` is true. Do not use `-a` and `-o` to combine tests, because they are ambiguous; chain several `test` commands with `&&` and `||`. The Bash `[[ ... ]]` is a different, shell-only keyword and is not covered here.

### Portability
Only the tests listed above are in the Common examples. `==` for strings, `-a` and `-o`, `-nt`, `-ot` and `-ef` are common but not in POSIX. The examples use `test`; the `[` form runs the shell builtin in the executable checks. The Common examples were executed on Bash 5.2, GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash test
In Bash, `test` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it. Bash also accepts `==` for strings, and has the separate `[[ ]]` keyword.

## GNU
### About GNU test
GNU test is part of GNU coreutils, and this page's checks ran coreutils 9.11. This view is the executable, which runs when you call it by path, as in `/bin/test`, or when the shell has no builtin of that name. `/bin/[` is the same program under another name.

## Apple
### About macOS test
This is the test in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/test`.

## BusyBox
### About BusyBox test
BusyBox bundles test as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `test` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
