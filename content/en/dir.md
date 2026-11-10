---
schema_version: 1
id: dir
title: dir
summary: List files in columns, like ls -C with escaped names.
aliases: [gdir]
tags: [files, directories]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/dir-invocation.html
---
## Common
### Overview
`dir` lists the files of a folder in columns, the way `ls -C -b` does. It exists for people who expect the name from other systems. It takes the same options as `ls`, so everything on the `ls` page applies, except that its defaults are different: columns even when output is not a terminal, and special characters shown as escapes.

### Background
`dir` is part of GNU coreutils; the GNU manual describes it as a version of `ls` with different defaults. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
List a folder in columns:

```sh
mkdir d && touch d/a d/b && dir d
```

Show a name with a space as an escaped name, which stays on one line:

```sh
mkdir d && touch 'd/a b' && dir d
```

### Common mistakes
`dir` output is meant for people, not for scripts, as with `ls`; use a glob or `find` in scripts. The escaped names (`a\ b`) are not the same as the raw names. On other systems `dir` may be a different program, or not exist at all.

### Portability
This page covers GNU only. Use `ls` for anything that must run on macOS or BusyBox. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU dir
GNU dir is part of GNU coreutils, and this page's checks ran coreutils 9.11.
