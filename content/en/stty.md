---
schema_version: 1
id: stty
title: stty
summary: Show and change the settings of the terminal.
aliases: [gstty]
tags: [terminal, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/stty-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=stty&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnustty.1.html
---
## Common
### Overview
`stty` shows and changes the settings of the terminal your input comes from: whether typed characters are shown (echo), which key interrupts a program, the size of the window and more. Scripts use it to read a password without showing it, or to put the terminal in raw mode. It needs a real terminal, so most of its examples cannot run in a test.

### Background
Today the name covers GNU stty, the stty that macOS ships, and the stty applet in BusyBox. The macOS manual describes it as setting options for a terminal device interface. They agree on `-a`, on turning settings on by name and off with a leading dash, and on `size`. The check ran on systems without a terminal, so the examples on this page were reviewed against the manuals but not executed.

### Examples
Show all settings of the current terminal:

```sh
stty -a
```

Print the size of the window as rows and columns:

```sh
stty size
```

Turn off echo so typed text is hidden, then turn it back on; do this only in a script that always restores it:

```sh short="Hide typed text, then turn echo back on (restore it!)"
stty -echo; read -r secret; stty echo
```

Return the terminal to reasonable defaults after something left it in a strange state:

```sh short="Reset the terminal to sane defaults"
stty sane
```

### Common mistakes
If a program stops before it restores the terminal (for example after `stty -echo`), the terminal looks dead: type `stty sane` and press Ctrl+J, or `reset`, to recover. `stty` acts on its standard input, so `stty < file` fails with "not a terminal". A script that changes settings should save them with `stty -g` and restore them in a `trap`. `size` is a common extension and is not in POSIX.

### Portability
The examples were reviewed only, because the checks ran without a terminal. `-a` prints a different layout on each system, so do not parse it. The option names follow the same pattern, but a few differ, such as `-f` on macOS and `-F` on GNU to name the device.

## GNU
### About GNU stty
GNU stty is part of GNU coreutils. `-F DEVICE` acts on a named terminal instead of standard input.

### Option: -F, --file
Use this terminal device instead of standard input.

## Apple
### About macOS stty
This is the stty in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `stty` uses `-f DEVICE` where GNU uses `-F`.

## BusyBox
### About BusyBox stty
BusyBox bundles stty as an applet, and anything beyond the common options depends on how it was built.

### BusyBox notes
BusyBox `stty` supports `-a`, `size`, `sane` and named settings. Other options depend on the build.
