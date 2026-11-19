---
schema_version: 1
id: timeout
title: timeout
summary: Run a command and stop it if it takes too long.
aliases: [gtimeout]
tags: [scripts, process]
implementations: [gnu, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/timeout-invocation.html
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnutimeout.1.html
---
## Common
### Overview
`timeout` runs a command and stops it after a time limit, so a stuck command cannot hang a script forever. It sends the command a signal, `TERM` by default, when the time is up. On GNU the exit status is 124 when the limit was hit, and otherwise the command's own status.

### Background
`timeout` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it. They agree on the form `timeout SECONDS COMMAND`.

### Examples
Run a command that finishes in time; its output is passed through:

```sh short="Run a command that finishes in time"
timeout 5 echo hi
```

Stop a command that runs too long; the shell shows that it did not succeed:

```sh short="Stop a command that runs too long"
timeout 1 sleep 5 || echo stopped
```

Use the command's own exit status when it finishes in time:

```sh
timeout 5 sh -c 'exit 3'; echo $?
```

Send a different signal with `-s`; here the stronger `KILL`, and the shell's own message about it is hidden:

```sh
{ timeout -s KILL 1 sleep 5; } 2> /dev/null || echo killed
```

### Option: -s
Send this signal instead of `TERM` when the time is up.

### Common mistakes
A command can ignore `TERM`, in which case it keeps running; use `-s KILL` or GNU's `-k` to follow up. `timeout` stops the command it started, but the command's children may survive. The exit status for a timeout is 124 on GNU, and may be different on other implementations, so test for "not 0" unless you know which you have. The time can have a suffix such as `m` for minutes in GNU.

### Portability
This page has no Apple view: macOS has no `timeout`; a script there can install `gtimeout` from GNU coreutils. Only whole seconds and `-s` are in the Common examples. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU timeout
GNU timeout is part of GNU coreutils, and this page's checks ran coreutils 9.11. It exits with status 124 on a timeout, and `-k` sends a `KILL` after a further delay.

### GNU examples
Check the status 124 that means a timeout:

```sh
timeout 1 sleep 5; echo $?
```

### Option: -k, --kill-after
Send `KILL` this long after the first signal if the command is still running.

## BusyBox
### About BusyBox timeout
BusyBox bundles timeout as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `timeout` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
