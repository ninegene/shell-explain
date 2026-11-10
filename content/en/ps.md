---
schema_version: 1
id: ps
title: ps
summary: List the processes that are running.
aliases: []
tags: [process, system]
implementations: [apple, procps, busybox]
updated: 2026-11-10
references:
  - https://gitlab.com/procps-ng/procps
  - https://man.freebsd.org/cgi/man.cgi?query=ps&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`ps` prints a snapshot of the running processes: their numbers (PIDs), who runs them, and what they are. A process is a running program. `ps` shows only the moment you ask, and it does not update; `top` is the live view. The options are the oldest source of confusion in the process tools, because BSD systems such as macOS and Linux tools (procps) accept different styles.

### Background
The `ps` command is part of every Unix-like system, and three implementations are covered here: the BSD ps that macOS ships, procps-ng ps (the usual one on Linux) and the ps applet in BusyBox. The procps version accepts both the BSD style (`ps aux`) and the Unix style (`ps -ef`), macOS accepts mostly the BSD style, and BusyBox accepts very few options. They agree on `-o` for choosing columns and on `-e` for all processes in the versions tested here.

### Examples
Print the heading that `ps` always prints; the column names are `PID` and others:

```sh
ps | head -n 1 | grep -c PID
```

Choose a single column with `-o`; here the heading shows `PID`:

```sh
ps -o pid | head -n 1
```

List all processes, one number per line, and check that your own shell is in the list (`$$` is its number):

```sh
ps -e -o pid= | grep -c "^ *$$\$"
```

Start a background process and find it in the list by its number:

```sh
sleep 7$$ & p=$!; ps -e -o pid= | grep -c "^ *$p\$"; kill $p
```

### Option: -e
Show every process, not only those of your terminal.

### Option: -o
Choose the columns to print, as in `-o pid,comm`; a name followed by `=` removes its heading.

### Common mistakes
Without `-e` or similar, `ps` shows only the processes of the current terminal, so in a script or a container it can print only itself. The style of options matters: `ps aux` is the BSD form and `ps -ef` the Unix form, and on macOS `ps -ef` works but shows different columns than on Linux. Looking for a program with `ps | grep name` also finds the `grep` itself; use `pgrep` for that. The command name column is cut to a limit that differs between systems, so use `-o args` or `pgrep -f` for full command lines.

### Portability
Only the heading, `-o` and `-e` are in the Common examples. The columns printed without `-o`, `-ef` and `aux` differ widely, and BusyBox prints a very short list. The Common examples were executed on macOS system tools, procps-ng and BusyBox 1.37.0 (see the development guide for the exact check).

## Apple
### About macOS ps
This is the ps in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Print the command name of your own shell with `-p` to select a process by number:

```sh
ps -p $$ -o comm= | grep -c .
```

### Apple notes
macOS `ps` accepts `aux` and `-ef`, and has `-p PID` to select a process.

## BusyBox
### About BusyBox ps
BusyBox bundles ps as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `ps` prints a short list and accepts `-o` for columns. In the 1.37.0 check it also accepted `-e`, which it ignores, but it has no `-p`.

## Procps
### About procps-ng ps
This is ps from procps-ng, the package that supplies the process tools of most Linux distributions, and this page's checks ran the version in Alpine 3.24. It accepts the BSD style (`ps aux`) and the Unix style (`ps -ef`), and has long options such as `--pid`.

### Procps examples
Print the command name of your own shell with `-p` to select a process by number:

```sh
ps -p $$ -o comm= | grep -c .
```
