---
schema_version: 1
id: vdir
title: vdir
summary: List files in the long format, like ls -l with escaped names.
aliases: [gvdir]
tags: [files, directories]
implementations: [gnu]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/vdir-invocation.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuvdir.1.html
---
## Common
### Overview
`vdir` lists files in the long format, the way `ls -l -b` does: permissions, owner, size, time and name. Like `dir` it exists for compatibility with other systems, and it takes the same options as `ls`.

### Background
`vdir` is part of GNU coreutils; the GNU manual describes it as a version of `ls` with different defaults. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
List a folder in the long format; the first line is the `total`, then one line per file:

```sh short="Long format: a total, then a line per file"
mkdir d && touch d/a d/b && vdir d | grep -c .
```

Show only the names, by taking the last field of each line after the total:

```sh
mkdir d && touch d/a d/b && vdir d | tail -n +2 | awk '{print $NF}'
```

### Common mistakes
The long format changes between systems and with the locale, so do not parse it in a script. `vdir` is the same as `ls -l`, so most people just type `ls -l`. A name with spaces makes the last-field trick above fail.

### Portability
This page covers GNU only. Use `ls -l` for anything that must run on macOS or BusyBox. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU vdir
GNU vdir is part of GNU coreutils, and this page's checks ran coreutils 9.11.
