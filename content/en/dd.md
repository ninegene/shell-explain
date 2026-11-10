---
schema_version: 1
id: dd
title: dd
summary: Copy bytes from an input to an output, with control over block sizes and conversion.
aliases: [gdd]
tags: [files, binary, disks]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/dd-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=dd&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`dd` copies data from an input to an output in blocks you choose, and can skip bytes, stop after a number of blocks, and convert case or other things on the way. It is used to make exact copies of devices and to cut pieces out of a file. Its syntax is unusual: arguments are written as `name=value`, with no dashes. It prints a summary of what it copied on standard error.

### Background
A `dd` command appeared in Version 5 AT&T UNIX, according to the macOS manual, which also says it is expected to be a superset of POSIX.2. Today the name covers GNU dd, the dd that macOS ships, and the dd applet in BusyBox. They agree on `if=`, `of=`, `bs=`, `count=` and `skip=`.

### Examples
Copy the first three bytes from standard input; `2>/dev/null` hides the summary that `dd` prints on standard error:

```sh
printf abcdef | dd bs=3 count=1 2>/dev/null
```

Skip two bytes and copy the next three, with a block size of one byte:

```sh
printf abcdef | dd bs=1 skip=2 count=3 2>/dev/null
```

Copy one file to another with `if=` (input file) and `of=` (output file):

```sh
dd if=notes.txt of=out.txt 2>/dev/null && head -n 1 out.txt
```

Copy with a larger block size, which is faster for big files:

```sh
dd if=notes.txt of=out.txt bs=4096 2>/dev/null && head -n 1 out.txt
```

### Operands
`if=FILE` reads from a file instead of standard input. `of=FILE` writes to a file instead of standard output, replacing it if it exists. `bs=N` sets the block size in bytes for reading and writing. `count=N` copies only that many blocks. `skip=N` skips that many blocks of the input first.

### Common mistakes
`dd` writes where you tell it, and `of=` can be a disk: mixing up `if=` and `of=` with a device can destroy data with no warning and no undo, so check both names twice and never run an `of=/dev/...` command you do not understand. The summary on standard error is not part of the copy, but it shows up on the screen unless you hide it. `dd` has no `-` options, so `dd --help` is the way to see what GNU supports.

### Portability
Only `if=`, `of=`, `bs=`, `count=` and `skip=` are in the Common examples. The conversions and the `status=` operand differ between implementations; the examples that use them are in the GNU and Apple sections. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU dd
GNU dd is part of GNU coreutils, and this page's checks ran coreutils 9.11. `status=none` hides the summary, and `status=progress` shows progress while it copies.

### GNU examples
Change the case of the copied text with `conv=ucase`, and hide the summary with `status=none`:

```sh
printf abc | dd conv=ucase status=none
```

### Operands: status and conv
`status=none` prints no summary, `status=noxfer` leaves out the rate line and `status=progress` shows a live count. `conv=ucase`, `conv=lcase` and `conv=sync` convert the data as it is copied.

## Apple
### About macOS dd
This is the dd in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Change the case of the copied text with `conv=ucase`, and hide the summary with `status=none`:

```sh
printf abc | dd conv=ucase status=none
```

### Operands: conv
`conv=ucase`, `conv=lcase` and `conv=sync` convert the data as it is copied.

### Apple notes
macOS `dd` accepts `status=none` and `status=noxfer` on the release tested.

## BusyBox
### About BusyBox dd
BusyBox bundles dd as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `dd` supports `if=`, `of=`, `bs=`, `count=` and `skip=`. The conversions and `status=` depend on the build.
