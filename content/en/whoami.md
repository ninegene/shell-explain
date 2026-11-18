---
schema_version: 1
id: whoami
title: whoami
summary: Print the user name you are running as.
aliases: [gwhoami]
tags: [users, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/whoami-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=whoami&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`whoami` prints the name of the user that runs it, the effective user. It is the shortest way to ask "which account am I?", for example after `sudo` or inside a container.

### Background
The macOS manual lists `whoami` as equivalent to `id -un`, without giving a date. Today the name covers GNU whoami, the whoami that macOS ships, and the whoami applet in BusyBox. They agree on printing the user name and taking no options.

### Examples
Print your user name; here only that something was printed is shown:

```sh short="Print your user name"
whoami | grep -c .
```

Check that it matches `id -un`:

```sh
test "$(whoami)" = "$(id -un)" && echo same
```

Use it in a path:

```sh
echo "/home/$(whoami)" | grep -c '^/home/.'
```

### Common mistakes
`whoami` shows the effective user, so `sudo whoami` prints `root`. To see who logged in to the session, use `logname`, or `who am i` where it is available. It takes no options, and extra arguments are an error.

### Portability
`whoami` is not in POSIX; `id -un` is the portable form. Both are used in the Common examples because they exist on all three implementations. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU whoami
GNU whoami is part of GNU coreutils, and this page's checks ran coreutils 9.11.

### GNU notes
GNU `whoami` takes only `--help` and `--version`.

## Apple
### About macOS whoami
This is the whoami in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `whoami` prints the effective user name and takes no options.

## BusyBox
### About BusyBox whoami
BusyBox bundles whoami as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `whoami` prints the effective user name.
