---
schema_version: 1
id: sleep
title: sleep
summary: Pause for a given time.
aliases: [gsleep]
tags: [scripts, time]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/sleep-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=sleep&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`sleep` does nothing for the number of seconds you give, and then exits. Scripts use it to wait between tries, to space out requests or to give a service time to start. It prints nothing.

### Background
A `sleep` command appeared in Version 4 AT&T UNIX, according to the macOS manual, which also says support for non-integer intervals and for units other than seconds is an extension. Today the name covers GNU sleep, the sleep that macOS ships, and the sleep applet in BusyBox. They agree on a number of seconds.

### Examples
Wait one second, then print a message:

```sh
sleep 1 && echo done
```

Wait zero seconds, which returns at once:

```sh
sleep 0 && echo ok
```

Wait between steps in a loop:

```sh
for i in 1 2; do echo "step $i"; sleep 1; done
```

Run something after a delay, in the background, and wait for it:

```sh
(sleep 1; echo later) & wait
```

### Common mistakes
`sleep 1` waits at least a second, not exactly one, and a script that must match a clock should compare times instead. A fractional value such as `0.5` or a suffix such as `1m` works on many systems but not all, so keep to whole seconds in portable scripts. Waiting in a loop for a service is fragile; check for the service itself, for example by trying it, rather than assuming a fixed delay is enough.

### Portability
Only whole seconds are in the Common examples. GNU accepts `s`, `m`, `h` and `d` suffixes and fractions, the macOS release tested accepts them too, and BusyBox accepts them when it is built with the fancy sleep option. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU sleep
GNU sleep is part of GNU coreutils, and this page's checks ran coreutils 9.11. It accepts a suffix for minutes, hours and days, and several values are added together.

### GNU examples
Wait a tenth of a second:

```sh
sleep 0.1 && echo done
```

## Apple
### About macOS sleep
This is the sleep in the system folders of macOS. It follows the BSD implementation.

### Apple notes
On the macOS release tested, `sleep` accepts fractions such as `0.1` and a unit suffix such as `m`. Older releases accept whole seconds only, so keep to whole seconds when you must support them.

## BusyBox
### About BusyBox sleep
BusyBox bundles sleep as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `sleep` accepts whole seconds; fractions and suffixes depend on the build.
