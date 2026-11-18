---
schema_version: 1
id: uname
title: uname
summary: Print the name and version of the operating system and the machine.
aliases: [guname]
tags: [system]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/uname-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=uname&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`uname` prints facts about the system it runs on: the kernel name such as `Linux` or `Darwin`, the kernel release, and the machine type such as `x86_64` or `arm64`. Scripts use it to choose commands that differ between systems. With no option it prints the kernel name.

### Background
The `uname` command appeared in PWB UNIX 1.0, and 4.4BSD was the first Berkeley release to have it, according to the macOS manual. Today the name covers GNU uname, the uname that macOS ships, and the uname applet in BusyBox. They agree on `-s`, `-r`, `-m` and `-a`.

### Examples
Print the kernel name; the plain command gives the same as `-s`:

```sh short="Print the kernel name; the plain command equals -s"
test "$(uname)" = "$(uname -s)" && echo same
```

Print the kernel name:

```sh
uname -s | grep -c .
```

Print the kernel release:

```sh
uname -r | grep -c .
```

Print the machine hardware name, such as `x86_64` or `arm64`:

```sh short="Print the machine hardware name (e.g. arm64)"
uname -m | grep -c .
```

Branch on the system in a script:

```sh
case "$(uname -s)" in Darwin) echo mac;; Linux) echo linux;; *) echo other;; esac | grep -c .
```

### Option: -s
Print the kernel name.

### Option: -r
Print the kernel release.

### Option: -m
Print the machine hardware name.

### Option: -a
Print all the fields.

### Common mistakes
`uname -s` is not the same as the distribution: it says `Linux` for every Linux system. For the distribution, read `/etc/os-release` instead. `-a` prints a line whose fields differ between systems, so do not parse it. The machine name is `arm64` on macOS and `aarch64` on Linux for the same kind of processor.

### Portability
Only `-s`, `-r`, `-m` and `-a` are in the Common examples, and their outputs differ by system, so the examples check only that something is printed. GNU adds `-o` for the operating system name and `-i` and `-p`, which print `unknown` on many systems. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU uname
GNU uname is part of GNU coreutils, and this page's checks ran coreutils 9.11. `uname -o` prints the operating system name, such as `GNU/Linux`.

### Option: -o, --operating-system
Print the operating system name.

## Apple
### About macOS uname
This is the uname in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `uname` also prints the kernel name `Darwin` and has `-o` as a synonym for `-s`.

## BusyBox
### About BusyBox uname
BusyBox bundles uname as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `uname` supports `-s`, `-r`, `-m` and `-a`. `-o`, `-i` and `-p` depend on the build.
