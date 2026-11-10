---
schema_version: 1
id: hostname
title: hostname
summary: Show or set the name of this computer on the network.
aliases: []
tags: [system, network]
implementations: [apple, busybox]
updated: 2026-11-10
references:
  - https://man.freebsd.org/cgi/man.cgi?query=hostname&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`hostname` prints the name of the computer. Only an administrator can change it, and doing so is outside this page. The name is used on networks and in the shell prompt, and it is not the same as a user name.

### Background
The `hostname` command appeared in 4.2BSD, according to the macOS manual. This page covers the hostname that macOS ships and the hostname applet in BusyBox. GNU coreutils does not include `hostname`: on many Linux systems it comes from a separate package, such as net-tools or inetutils, whose options differ, so this page has no GNU view.

### Examples
Print the name and check that there is one:

```sh
test -n "$(hostname)" && echo set
```

Print the short form, without the domain part, with `-s`:

```sh
hostname -s | grep -c .
```

Use it in a file name:

```sh
echo "backup-$(hostname -s).tar" | grep -c '^backup-.'
```

### Option: -s
Print the name with any domain part removed.

### Common mistakes
A computer can have several names: the name from the network configuration, the name on a local network and a name from a service. `hostname` shows one of them, so it may differ from what other computers call yours. The output of `hostname -f` or `-i` depends on the system and its network, so avoid them in a script that must be portable.

### Portability
Only plain `hostname` and `-s` are in the Common examples, and they are available on macOS and BusyBox. GNU `hostname` programs from other packages differ in their options. The Common examples were executed on macOS system tools and BusyBox 1.37.0 (see the development guide for the exact check).

## Apple
### About macOS hostname
This is the hostname in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `hostname` has `-s` for the short name and `-f` for the fully qualified name.

## BusyBox
### About BusyBox hostname
BusyBox bundles hostname as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `hostname` supports `-s`. `-f` and `-d` depend on the build.
