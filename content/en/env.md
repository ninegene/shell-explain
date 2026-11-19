---
schema_version: 1
id: env
title: env
summary: Run a command with a changed environment, or print the environment.
aliases: [genv]
tags: [shell, environment, scripts]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/env-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=env&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuenv.1.html
---
## Common
### Overview
`env` runs a program with changes to its environment variables: add or override some with `NAME=value`, remove one with `-u`, or start from nothing with `-i`. With no command it prints every variable. It is also the portable way to find an interpreter in a script's first line, as in `#!/usr/bin/env python3`.

### Background
The `env` command appeared in 4.4BSD, according to the macOS manual, which also says it conforms to POSIX.1-2001 and that `-u` is a non-standard extension. Today the name covers GNU env, the env that macOS ships, and the env applet in BusyBox. They agree on `NAME=value`, `-i` and `-u`.

### Examples
Run a command with one extra variable; the change lasts only for that command:

```sh short="Run a command with one extra variable, only for it"
env FOO=bar sh -c 'echo "$FOO"'
```

Start with an empty environment except for what you give, with `-i`; `env` then prints just that one variable:

```sh short="Start with an empty environment plus what you give (-i)"
env -i FOO=bar env
```

Remove a variable for one command with `-u`:

```sh
HOME=/x env -u HOME sh -c 'echo ${HOME:-unset}'
```

Print the environment and count the line that holds `PATH`:

```sh
env | grep -c '^PATH='
```

### Option: -i
Start with an empty environment, ignoring the one the shell passed down.

### Option: -u
Remove this variable from the environment of the command.

### Common mistakes
`env FOO=bar` with no command only prints the environment with that change, it does not set `FOO` in your shell; use `export FOO=bar` for that. `-i` also removes `PATH`, so a command may no longer be found; give a full path or set `PATH` again. A variable set on the command line of your shell before a program, as in `FOO=bar prog`, does the same thing as `env FOO=bar prog` for external programs.

### Portability
Only `NAME=value`, `-i` and `-u` are in the Common examples. GNU adds `-S` and `-C`, and macOS adds `-P` and `-S` in its own form, so avoid them in portable scripts. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU env
GNU env is part of GNU coreutils, and this page's checks ran coreutils 9.11. Its `-S` option splits one string into several arguments, which is useful on the first line of a script.

### GNU examples
Print the variables with a NUL byte after each, so values that contain newlines stay unambiguous, then count them:

```sh
env -0 | tr '\0' '\n' | grep -c '^PATH='
```

### Option: -0, --null
End each printed variable with a NUL byte instead of a newline.

## Apple
### About macOS env
This is the env in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `env` has `-i`, `-u`, `-P` and `-S`; its `-0` is available on current releases. Check the manual on your system before relying on them.

## BusyBox
### About BusyBox env
BusyBox bundles env as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `env` supports `NAME=value`, `-i` and `-u`. `-0` and `-S` depend on the build.
