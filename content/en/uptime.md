---
schema_version: 1
id: uptime
title: uptime
summary: Show how long the system has been running and the load.
aliases: [guptime]
tags: [system, time]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/uptime-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=uptime&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`uptime` prints one line: the current time, how long the machine has been running, how many users are logged in, and the load averages for the last 1, 5 and 15 minutes. A load average is roughly the number of programs waiting to run; compare it with the number of processors.

### Background
The `uptime` command appeared in 3.0BSD, according to the macOS manual. Today the name covers GNU uptime, the uptime that macOS ships, and the uptime applet in BusyBox. They agree on printing one line that includes `up`.

### Examples
Print the line and check that it contains `up`; the numbers change all the time:

```sh short="Check that the line contains up"
uptime | grep -c up
```

Show only the load averages by cutting the line at its label; the word `load` appears on all three systems:

```sh short="Show only the load averages, by cutting at the label"
uptime | grep -c load
```

### Common mistakes
The exact layout, the number of users and the wording of the days differ between systems, so do not parse the line in a portable script. Read `/proc/loadavg` on Linux or use `sysctl -n vm.loadavg` on macOS for the numbers alone. A load average above the number of processors means programs are waiting.

### Portability
Only the presence of `up` and `load` is in the Common examples. GNU, macOS and BusyBox lay out the line differently. The Linux guest used for the BusyBox checks does not install GNU `uptime`, so the GNU examples for this page ran on macOS with Homebrew coreutils 9.11 only. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU uptime
GNU uptime is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS uptime
This is the uptime in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox uptime
BusyBox bundles uptime as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `uptime` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
