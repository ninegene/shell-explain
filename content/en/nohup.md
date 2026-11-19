---
schema_version: 1
id: nohup
title: nohup
summary: Run a command that keeps going after you log out.
aliases: [gnohup]
tags: [process, scripts]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/nohup-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=nohup&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnunohup.1.html
---
## Common
### Overview
`nohup` runs a command so that it ignores the hangup signal, the one a terminal sends when you close it or log out. The command keeps running in the background. If its output would go to the terminal, `nohup` redirects it to a file named `nohup.out`.

### Background
The macOS manual says `nohup` is expected to be POSIX.2 compatible, and gives no date. Today the name covers GNU nohup, the nohup that macOS ships, and the nohup applet in BusyBox. They agree on ignoring the hangup signal and on the `nohup.out` default.

### Examples
Run a command with `nohup`; here its output is not a terminal, so it is passed through:

```sh short="Run a command; non-terminal output is passed through"
nohup echo hi
```

Send the output to a file of your choice and read it afterwards:

```sh
nohup sh -c 'echo hi' > out.log 2>&1 && cat out.log
```

Start a long command in the background so it survives the terminal closing:

```sh
nohup sleep 1 > /dev/null 2>&1 & wait; echo started
```

### Common mistakes
`nohup` does not put the command in the background; add `&` yourself. When both output streams go to the terminal, `nohup.out` appears in the current folder, which can surprise you. A command that is already stopped by a signal other than hangup is not helped by `nohup`. For long-lived services, a process supervisor or a terminal multiplexer such as `tmux` is a better choice.

### Portability
Only running a command and redirecting its output are in the Common examples. The message `nohup: ignoring input` and its wording differ between implementations and go to standard error. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU nohup
GNU nohup is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS nohup
This is the nohup in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox nohup
BusyBox bundles nohup as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `nohup` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
