---
schema_version: 1
id: printenv
title: printenv
summary: Print the value of environment variables.
aliases: [gprintenv]
tags: [shell, environment]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/printenv-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=printenv&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`printenv` prints the value of the environment variable you name, or all variables when you name none. It exits with status 1 when a named variable is not set, which makes it useful in `if` tests. It shows only exported variables, those that programs started by the shell can see.

### Background
The `printenv` command appeared in 3.0BSD, according to the macOS manual, which also says it is not specified by any standard. Today the name covers GNU printenv, the printenv that macOS ships, and the printenv applet in BusyBox. They agree on printing one variable's value and on the exit status for a missing one.

### Examples
Print the value of one variable:

```sh
FOO=bar printenv FOO
```

Check whether a variable is set; the exit status is 1 when it is not:

```sh
printenv NOPE || echo unset
```

Print every variable and count the line that holds `PATH`:

```sh
printenv | grep -c '^PATH='
```

Print the value in a script without the name:

```sh
FOO=bar sh -c 'echo "value is $(printenv FOO)"'
```

### Common mistakes
A shell variable that was never exported is not in the environment, so `printenv NAME` fails for it even though `echo "$NAME"` works; use `export`. `printenv NAME` prints only the value, while `printenv` alone prints `NAME=value` lines. `printenv` is not in POSIX; `env` and `echo "$NAME"` are the portable choices.

### Portability
Only printing one variable and listing all are in the Common examples. GNU has `-0` for NUL-terminated output. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU printenv
GNU printenv is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### Option: -0, --null
End each output line with a NUL byte instead of a newline.

## Apple
### About macOS printenv
This is the printenv in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `printenv` has no `-0`.

## BusyBox
### About BusyBox printenv
BusyBox bundles printenv as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `printenv` prints named variables or all of them. GNU options depend on the build.
