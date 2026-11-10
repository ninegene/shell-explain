---
schema_version: 1
id: chroot
title: chroot
summary: Run a command with a different folder as the root of the file system.
aliases: [gchroot]
tags: [system, security, containers]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/chroot-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=chroot&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`chroot` runs a command with a chosen folder as its root, so that `/` means that folder and nothing outside it can be named. The folder must hold everything the command needs, such as the program and its libraries. It is used to build or repair a system from outside and as a building block of containers. It needs administrator rights, so its examples are reviewed and not run.

### Background
Today the name covers GNU chroot, the chroot that macOS ships (in `/usr/sbin`), and the chroot applet in BusyBox. The macOS manual describes it as changing the root directory for a command. They agree on `chroot NEWROOT [COMMAND]`. The examples on this page were reviewed against the manuals but not executed, because they need administrator rights and a prepared folder.

### Examples
Start a shell with `/srv/root` as its root; it needs `/srv/root/bin/sh` and the libraries that shell uses:

```sh
chroot /srv/root /bin/sh
```

Run one command in the new root and return:

```sh
chroot /srv/root /bin/ls /
```

Run as another user inside the new root:

```sh
chroot --userspec=nobody:nogroup /srv/root /bin/id
```

### Common mistakes
`chroot` is not a security boundary: a process with administrator rights can escape it, so do not use it to lock up untrusted code. A missing program or library inside the new root gives a confusing "no such file" error for a file that exists outside. Mount `/proc` and `/dev` into the new root if the command needs them. Running it without enough rights fails with "operation not permitted".

### Portability
The examples were reviewed only. The option for the user and group differs: GNU `--userspec=USER:GROUP`, and the macOS tool has `-u` and `-g`.

## GNU
### About GNU chroot
GNU chroot is part of GNU coreutils. `--userspec=USER:GROUP` runs the command as that user.

## Apple
### About macOS chroot
This is the chroot in the system folders of macOS. It follows the BSD implementation. It is in `/usr/sbin`, and has `-u USER` and `-g GROUP` for the user and groups.

## BusyBox
### About BusyBox chroot
BusyBox bundles chroot as an applet, and anything beyond the common options depends on how it was built.

### BusyBox notes
BusyBox `chroot` supports `NEWROOT [COMMAND]`. Other options depend on the build.
