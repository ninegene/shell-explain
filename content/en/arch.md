---
schema_version: 1
id: arch
title: arch
summary: Print the machine hardware name.
aliases: []
tags: [system]
implementations: [apple, busybox]
updated: 2026-11-10
references:
  - https://man.freebsd.org/cgi/man.cgi?query=arch&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`arch` prints the machine type, such as `x86_64`, `aarch64` or `arm64`. On Linux it does the same as `uname -m`. On macOS it is a different, more powerful command: with a command to run, it starts it under a chosen processor architecture, which is how Intel programs run on Apple silicon.

### Background
The macOS manual describes `arch` as a tool that runs a command under a given architecture and, without arguments, prints the native one. BusyBox has an `arch` applet that prints the machine name. GNU coreutils can also provide an `arch` that prints the machine name, but it is deprecated in favor of `uname -m`, and the coreutils builds used for these checks do not install it, so this page has no GNU view.

### Examples
Print the machine type; the name depends on the machine, so only the shape is shown:

```sh
arch | grep -c '^[a-z0-9_]*$'
```

Compare with `uname -m` to see that both name the machine; on macOS `arch` and `uname -m` agree on Apple silicon and Intel:

```sh
test -n "$(arch)" && echo named
```

### Common mistakes
On macOS `arch -x86_64 command` can run under translation and then `uname -m` inside it says `x86_64`, so the answer depends on how the shell was started. The spelling of the same processor differs by system: `arm64` on macOS and `aarch64` on Linux. For a portable script, use `uname -m` and map the names you care about.

### Portability
This page has no GNU view. The Common examples check only that a name is printed. The Common examples were executed on macOS system tools and BusyBox 1.37.0 (see the development guide for the exact check).

## Apple
### About macOS arch
This is the arch in `/usr/bin` on macOS. It prints the native architecture, and with a command it runs that command under the architecture you choose.

### Architecture arguments
`arch -arm64 command` runs the command as an Apple silicon program, and `arch -x86_64 command` runs it as an Intel program through translation.

## BusyBox
### About BusyBox arch
BusyBox bundles arch as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `arch` prints the machine hardware name and takes no options.
