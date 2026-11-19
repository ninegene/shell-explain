---
schema_version: 1
id: logname
title: logname
summary: Print the name you logged in with.
aliases: [glogname]
tags: [users, system]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/logname-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=logname&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnulogname.1.html
---
## Common
### Overview
`logname` prints the user name of the login session, which can differ from the current user after `sudo` or `su`. After `sudo`, `whoami` says `root`, while `logname` still says who logged in. It fails when the process has no login session, such as in a cron job or a container.

### Background
The `logname` command appeared in 4.4BSD, according to the macOS manual, which also says it is expected to conform to POSIX.2. Today the name covers GNU logname and the logname that macOS ships. The BusyBox 1.37.0 build used for this page has no `logname` applet.

### Examples
Print the login name, or say none when there is no login session; either way one line is printed:

```sh
{ logname || echo none; } 2> /dev/null | grep -c .
```

Compare it with the current user in a script:

```sh
{ test "$(logname 2> /dev/null)" = "$(id -un)" && echo same || echo different; } | grep -c .
```

### Common mistakes
`logname` can fail with "no login name" in a program started by cron, a service or a container, so a script should handle the failure. It reports the first login of the session, not the current effective user. For the user you are running as, use `id -un`.

### Portability
This page has no BusyBox view: BusyBox has no `logname` applet. The output depends on the session, so the Common examples count lines. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU logname
GNU logname is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS logname
This is the logname in the system folders of macOS. It follows the BSD implementation.
