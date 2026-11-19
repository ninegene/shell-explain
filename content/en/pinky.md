---
schema_version: 1
id: pinky
title: pinky
summary: Show a short report on logged-in users, in the style of the finger command.
aliases: [gpinky]
tags: [users, system]
implementations: [gnu]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/pinky-invocation.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnupinky.1.html
---
## Common
### Overview
`pinky` is a lightweight `finger`: it lists the users logged in now with their login name, real name, terminal, idle time and login time. It does not contact other machines, which `finger` can. It is rarely needed today; `who` and `w` cover most uses.

### Background
`pinky` is part of GNU coreutils; the GNU manual describes it as a lightweight finger. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Run it; the lines depend on who is logged in, so only success is shown:

```sh short="Run it; only success is shown"
pinky > /dev/null && echo ok
```

Use the short form, which leaves out the home directory and shell fields:

```sh short="Use the short form (no home directory or shell fields)"
pinky -s > /dev/null && echo ok
```

### Common mistakes
`pinky` reads the login records, so it shows nothing for sessions that are not registered there. The real name comes from the password database and may be empty. It never talks to other computers, so `pinky user@host` is not supported.

### Portability
This page covers GNU only. On macOS, use `who` or `finger`; in a portable script, use `who`. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU pinky
GNU pinky is part of GNU coreutils, and this page's checks ran coreutils 9.11.
