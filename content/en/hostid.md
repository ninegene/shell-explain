---
schema_version: 1
id: hostid
title: hostid
summary: Print a numeric identifier for this host.
aliases: [ghostid]
tags: [system]
implementations: [gnu, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/hostid-invocation.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`hostid` prints a number, in hexadecimal, that identifies the machine. On many systems it is derived from the machine's IP address or set by an administrator, so it is not a reliable unique name. It is an old command, kept for compatibility.

### Background
`hostid` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it. They agree on printing eight hexadecimal digits.

### Examples
Print the identifier and check that it is eight hexadecimal digits:

```sh short="Print the identifier: eight hexadecimal digits"
hostid | grep -c '^[0-9a-f]\{8\}$'
```

### Common mistakes
The value is not guaranteed to be unique or stable: it can change when the network configuration changes, and two machines can print the same one. For a stable machine identity use `/etc/machine-id` on Linux or a cloud instance ID. macOS has no `hostid`.

### Portability
This page has no Apple view: macOS has no `hostid`. The output depends on the machine, so the Common example checks its shape. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU hostid
GNU hostid is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## BusyBox
### About BusyBox hostid
BusyBox bundles hostid as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `hostid` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
