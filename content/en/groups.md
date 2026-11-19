---
schema_version: 1
id: groups
title: groups
summary: Print the groups a user belongs to.
aliases: [ggroups]
tags: [users, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/groups-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=groups&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnugroups.1.html
---
## Common
### Overview
`groups` prints the names of the groups you belong to, or those of the user you name. Group membership decides which files you may read or write besides your own. It is the same list as `id -Gn`.

### Background
The macOS manual lists the historic `groups` command as equivalent to `id -Gn`, and gives no date. Today the name covers GNU groups, the groups that macOS ships, and the groups applet in BusyBox. They agree on printing group names for the current user.

### Examples
Print your groups; the names depend on the system, so only that something is printed is shown:

```sh short="Print your groups (the names depend on the system)"
groups | grep -c .
```

Count the lines of output; the list is on one line:

```sh
groups | wc -l | tr -d ' '
```

### Common mistakes
`groups` shows names, and a group that has no name on the system appears as a number. The list reflects the groups your session started with, so a group you were added to a moment ago may be missing until you log in again. GNU `groups user` prints `user : groups`, which has a different layout from `groups` alone.

### Portability
Only the no-argument form is in the Common examples, and the output differs between systems, so the examples check its shape. `id -Gn` is the portable alternative. The Linux guest used for the BusyBox checks does not install GNU `groups`, so the GNU examples for this page ran on macOS with Homebrew coreutils 9.11 only. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU groups
GNU groups is part of GNU coreutils, and this page's checks ran coreutils 9.11. With a user name it prints `user : group group ...`.

## Apple
### About macOS groups
This is the groups in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox groups
BusyBox bundles groups as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `groups` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
