---
schema_version: 1
id: kill
title: kill
summary: Send a signal to a process, usually to stop it.
aliases: []
tags: [process, shell, system]
implementations: [bash, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/bash/manual/html_node/Job-Control-Builtins.html#index-kill
  - https://man.freebsd.org/cgi/man.cgi?query=kill&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`kill` sends a signal to a process, named by its number (PID). By default the signal is `TERM`, which asks the process to stop and lets it clean up. `-9` (`KILL`) stops it at once and cannot be caught. Signal 0 does not stop anything; it only checks that the process exists and that you may signal it.

### Background
The `kill` command sends signals to processes; the Bash builtin also understands job numbers such as `%1`. Today the name covers the Bash builtin, the kill that macOS ships, and the kill applet in BusyBox. GNU coreutils can build a `kill` program too, but the Alpine and Homebrew coreutils used for these checks do not include it, so this page has no GNU view. They agree on `-s SIGNAL`, a numeric signal and a process number.

### Examples
Start a background process, stop it with the default signal and wait for it:

```sh
sleep 5 & kill $!; wait $! 2> /dev/null; echo done
```

Check that a process exists without sending a real signal, using signal 0; `$$` is the shell itself:

```sh
kill -0 $$ && echo alive
```

Name the signal with `-s`:

```sh
sleep 5 & kill -s TERM $!; wait $! 2> /dev/null; echo done
```

Use the strongest signal for a process that will not stop; the shell's own message is hidden:

```sh
sleep 5 & kill -9 $!; wait $! 2> /dev/null; echo done
```

### Option: -s
Send the named signal, such as `-s TERM`, instead of the default.

### Option: -9
Send `KILL`, which cannot be caught or ignored.

### Common mistakes
Try the default signal first: `kill -9` gives the process no chance to save its work or remove its files. Check the number twice, because a wrong PID stops a different program, and a PID can be reused once a process has ended. You can signal only your own processes, unless you are an administrator. Many processes ignore `TERM` while busy, so a process that remains is not necessarily broken. To stop programs by name, see `pkill`.

### Portability
Only the default signal, `-0`, `-s NAME` and `-9` are in the Common examples. The list of signals printed by `kill -l` differs: the Bash builtin prints numbers with names, and the executables print names only. The Common examples were executed on Bash 5.2, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## Bash
### About Bash kill
In Bash, `kill` is a shell builtin, so the shell runs it without starting a program. This view documents Bash 5.2, and the Common examples were also executed with it. `kill %1` signals a job by its job number, and `kill -l` lists the signals with their numbers.

## Apple
### About macOS kill
This is the kill in the system folders of macOS. It follows the BSD implementation, and is the executable that runs when you call it by path, as in `/bin/kill`.

## BusyBox
### About BusyBox kill
BusyBox bundles kill as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `kill` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
