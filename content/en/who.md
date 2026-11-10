---
schema_version: 1
id: who
title: who
summary: Show who is logged in, and from which terminal.
aliases: [gwho]
tags: [users, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/who-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=who&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`who` lists the users logged in now, with their terminal and the time they logged in. With `am i` (or `-m`) it shows only your own login. It reads the same login records as `users`, so it shows nothing for sessions that did not register there.

### Background
A `who` command appeared in Version 1 AT&T UNIX, according to the macOS manual, which also says it conforms to POSIX.1-2001. Today the name covers GNU who, the who that macOS ships, and the who applet in BusyBox. They agree on listing logged-in users.

### Examples
Run it; the lines depend on who is logged in, so only success is shown:

```sh
who > /dev/null && echo ok
```

Show the login records with the column headings, with `-H`:

```sh
who -H > /dev/null && echo ok
```

### Common mistakes
`who` can print nothing in a container, over some remote tools or under `sudo`, even though you are working. The columns and the time format differ between implementations, so do not parse them. `who am i` reports the terminal of the session, not the effective user; use `id -un` for that.

### Portability
Only plain `who` and `-H` are in the Common examples; `who am i` works on GNU and macOS but BusyBox 1.37.0 printed its usage and failed. Options such as `-q` for a count and `-b` for the boot time print in different forms on each system. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU who
GNU who is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-q` prints a count of users, and `-b` the time of the last boot.

### GNU examples
Show only the line of your own login; the command succeeds even if the system has no record:

```sh
who am i > /dev/null && echo ok
```

## Apple
### About macOS who
This is the who in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Show only the line of your own login; the command succeeds even if the system has no record:

```sh
who am i > /dev/null && echo ok
```

## BusyBox
### About BusyBox who
BusyBox bundles who as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `who` supports `-H`. It rejected `am i` in the 1.37.0 check and printed its usage, so use `id -un` or `tty` there.
