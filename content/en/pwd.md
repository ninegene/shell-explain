---
schema_version: 1
id: pwd
title: pwd
summary: Print the folder you are in.
aliases: [gpwd]
tags: [shell, navigation, directories]
implementations: [bash, gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html#index-pwd
  - https://www.gnu.org/software/coreutils/manual/html_node/pwd-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=pwd&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`pwd` (print working directory) prints the full path of the current folder. With `-L` it prints the path as you reached it, including symbolic links; with `-P` it prints the physical path, with every link resolved. Without an option the builtin of most shells uses `-L`.

### Background
The `pwd` command prints the working directory, and the Bash builtin keeps the logical path in the `PWD` variable. Today the name covers the Bash builtin, GNU pwd, the pwd that macOS ships, and the pwd applet in BusyBox. They agree on `-L` and `-P`.

### Examples
Print the current folder and check that it is an absolute path:

```sh
pwd | grep -c '^/'
```

Move into a folder and print where you are:

```sh
cd logs && pwd | grep -c 'logs$'
```

After going through a symbolic link, `-L` keeps the link name:

```sh
ln -s logs lnk && cd lnk && pwd -L | grep -c 'lnk$'
```

`-P` shows the real folder behind the link:

```sh
ln -s logs lnk && cd lnk && pwd -P | grep -c 'logs$'
```

### Option: -L
Print the logical path, keeping symbolic links as you used them.

### Option: -P
Print the physical path, with symbolic links resolved.

### Common mistakes
The two options give different answers inside a linked folder, so a script that compares paths should choose one on purpose. On many systems `/tmp` or `/var` is itself a link, so `pwd -P` can print a longer path than you expect. `pwd` prints the folder of the shell that runs it; a script that changes folders in a subshell does not change yours.

### Portability
Only `-L` and `-P` with the paths above are in the Common examples. The default when no option is given differs: the Bash builtin and the macOS executable use `-L`, while GNU `/bin/pwd` uses `-P` unless the `POSIXLY_CORRECT` variable is set. Give the option you mean. The Common examples were executed on Bash 5.2, GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash pwd
In Bash, `pwd` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it. The `PWD` variable holds the same logical path, so `echo "$PWD"` needs no process.

## GNU
### About GNU pwd
GNU pwd is part of GNU coreutils, and this page's checks ran coreutils 9.11. This view is the executable, which runs when you call it by path, as in `/bin/pwd`, or when the shell has no builtin of that name. Without an option it prints the physical path, as `-P` does.

## Apple
### About macOS pwd
This is the pwd in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/pwd`. Without an option it prints the logical path, as `-L` does.

## BusyBox
### About BusyBox pwd
BusyBox bundles pwd as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `pwd` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
