---
schema_version: 1
id: killall
title: killall
summary: Send a signal to every process that has a given name.
aliases: []
tags: [process, system]
implementations: [apple, busybox]
updated: 2026-11-10
references:
  - https://man.freebsd.org/cgi/man.cgi?query=killall&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`killall` sends a signal, `TERM` by default, to every process whose name matches exactly. It is the quick way to stop all copies of a program, such as `killall Safari`. It is more blunt than `pkill`, because it works by exact name and has fewer ways to narrow the choice. On Linux systems with the psmisc package, `killall` has other options; and on Solaris, `killall` shuts down the whole machine, which is why the examples are marked reviewed only.

### Background
The macOS manual documents `killall` as sending a signal to processes by name. This page covers the killall that macOS ships and the killall applet in BusyBox. The psmisc `killall` of many Linux distributions is another program whose options differ and is not covered here. The examples on this page were reviewed against the manuals but not executed, because a name such as `sleep` would signal every process of that name on the machine, including ones that are not part of the check.

### Examples
Stop every process named `sleep` that you own, with the default `TERM`:

```sh
killall sleep
```

Name the signal with `-s` or a dash and its name:

```sh
killall -s TERM sleep
```

Use the strongest signal when a process will not stop:

```sh
killall -9 sleep
```

Check that a process with that name exists without stopping it, with signal 0:

```sh
killall -0 sleep && echo running
```

### Common mistakes
`killall` signals every process with the name, including the ones you did not mean; list them with `pgrep -x name` first. On macOS, an exact name is needed, so `killall Terminal` and not `killall terminal`. Never run `killall` on a Solaris-style system, where it means "kill all processes". A process that is stopped with `-9` cannot clean up. Programs that were started by a supervisor may be restarted at once.

### Portability
The examples were reviewed only. The options of the psmisc `killall` (`-r`, `-w`, `-y`) are not in the BSD or BusyBox ones.

## Apple
### About macOS killall
This is the killall in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox killall
BusyBox bundles killall as an applet, and anything beyond the common options depends on how it was built.

### BusyBox notes
BusyBox `killall` signals processes by exact name, like the macOS one. Its options depend on the build.
