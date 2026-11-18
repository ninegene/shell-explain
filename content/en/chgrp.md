---
schema_version: 1
id: chgrp
title: chgrp
summary: Change the group of files.
aliases: [gchgrp]
tags: [files, permissions]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/chgrp-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=chgrp&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`chgrp` changes the group that a file belongs to. A group lets several users share files: the group permission bits of the file apply to its members. Normal users can only change a file to a group they belong to; the examples use your own primary group so they are safe to run.

### Background
The macOS manual says `chgrp` is expected to be POSIX.2 compatible, and gives no date. Today the name covers GNU chgrp, the chgrp that macOS ships, and the chgrp applet in BusyBox. They agree on a group and on `-R`.

### Examples
Set the group to your own primary group:

```sh
chgrp "$(id -g)" notes.txt && echo ok
```

Apply it to a folder and everything inside it with `-R`:

```sh short="Apply it to a folder and everything inside (-R)"
chgrp -R "$(id -g)" projects logs && echo ok
```

Check the result; `ls -l` shows the group in the fourth column, and `id -gn` is the name:

```sh
chgrp "$(id -g)" notes.txt && test "$(ls -l notes.txt | awk '{print $4}')" = "$(id -gn)" && echo yes
```

### Option: -R
Change the group of folders and everything inside them.

### Common mistakes
A normal user can choose only groups they belong to; any other group is refused. `chgrp -R` can touch a very large number of files, so check the path first. Changing the group does not change the permission bits: the group may still lack the access you want, which `chmod g+rw` gives. Use `chown :group` to set the group with `chown`.

### Portability
Only a numeric group and `-R` are in the Common examples. Group names and numbers work, but names can differ between machines. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU chgrp
GNU chgrp is part of GNU coreutils, and this page's checks ran coreutils 9.11. `--reference=FILE` copies the group of another file.

## Apple
### About macOS chgrp
This is the chgrp in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox chgrp
BusyBox bundles chgrp as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `chgrp` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
