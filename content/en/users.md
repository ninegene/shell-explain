---
schema_version: 1
id: users
title: users
summary: Print the names of users who are logged in.
aliases: [gusers]
tags: [users, system]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/users-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=users&sektion=1
---
## Common
### Overview
`users` prints, on one line, the names of the users logged in now, sorted, with a name repeated for each session. With nobody in the login records it prints nothing, or an empty line on some systems. `who` shows more detail.

### Background
The `users` command appeared in 3.0BSD, according to the macOS manual. Today the name covers GNU users and the users that macOS ships. The BusyBox 1.37.0 build used for this page has no `users` applet.

### Examples
Run it; the output depends on who is logged in, so only success is shown:

```sh
users > /dev/null && echo ok
```

Send the output to a pipe and check that the command still succeeds when nobody is logged in; the output may be empty:

```sh
users | cat > /dev/null && echo ok
```

### Common mistakes
`users` reads the system login records, so a session that did not register there, such as some containers or remote tools, does not appear. A user with two terminals is listed twice. It shows names only; use `who` to see terminals and times.

### Portability
This page has no BusyBox view: BusyBox has no `users` applet. The output depends on the system, so the Common examples check only that the command succeeds. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU users
GNU users is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS users
This is the users in the system folders of macOS. It follows the BSD implementation.
