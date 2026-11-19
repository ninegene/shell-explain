---
schema_version: 1
id: id
title: id
summary: Print the user and group IDs and names of a user.
aliases: [gid]
tags: [users, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/id-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=id&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuid.1.html
---
## Common
### Overview
`id` prints who a user is to the system: the numeric user ID, the group ID and the groups the user belongs to. With a flag it prints just one of these, which makes it handy in scripts, for example to test whether you are the root user (ID 0). With a user name as the last argument it describes that user instead of you.

### Background
The macOS manual says `id` is expected to conform to POSIX.2, and that the historic `groups` command is equivalent to `id -Gn` and `whoami` to `id -un`. Today the name covers GNU id, the id that macOS ships, and the id applet in BusyBox. They agree on `-u`, `-g`, `-G` and `-n`.

### Examples
Print your numeric user ID:

```sh
id -u | grep -c '^[0-9][0-9]*$'
```

Print your user name; `-u` selects the user and `-n` asks for the name:

```sh short="Compare id -un with whoami: same name"
test "$(id -un)" = "$(whoami)" && echo same
```

Print your primary group ID:

```sh
id -g | grep -c '^[0-9][0-9]*$'
```

List the names of all your groups on one line:

```sh
id -Gn | grep -c .
```

Check in a script whether you are the root user:

```sh
{ test "$(id -u)" -eq 0 && echo root || echo not root; } | grep -c .
```

### Option: -u
Print only the user ID.

### Option: -g
Print only the primary group ID.

### Option: -G
Print the IDs of all groups.

### Option: -n
Print names instead of numbers; used with `-u`, `-g` or `-G`.

### Common mistakes
`-n` does nothing on its own; combine it with `-u`, `-g` or `-G`. The group list shows the groups a session started with, so a group added later may not appear until you log in again. Test the number, as in `id -u`, rather than the name `root`, since a system can have another account with ID 0 or none named root.

### Portability
Only `-u`, `-g`, `-G` and `-n` are in the Common examples, and the example outputs show only their shape because they depend on who runs them. GNU `-Z` for security context and `-r` for the real ID are not on every system. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU id
GNU id is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### Option: -r, --real
Print the real ID instead of the effective ID; used with `-u`, `-g` or `-G`.

## Apple
### About macOS id
This is the id in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `id` has `-u`, `-g`, `-G` and `-n`; it also has `-P` for a password-file style line.

## BusyBox
### About BusyBox id
BusyBox bundles id as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `id` supports `-u`, `-g`, `-G` and `-n`. `-r` and `-Z` depend on the build.
