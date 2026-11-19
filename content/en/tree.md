---
schema_version: 1
id: tree
title: tree
summary: Show a folder and everything inside it as an indented tree.
aliases: []
tags: [files, directories, developer]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://oldmanprogrammer.net/source.php?dir=projects/tree
  - https://manpages.ubuntu.com/manpages/resolute/man1/tree.1.html
---
## Common
### Overview
`tree` prints a folder and the files and folders inside it, as a tree with lines showing which item belongs to which. It ends with a count of directories and files. It is a quick way to see the layout of a project or to paste a layout into a document.

### Background
`tree` is a small open-source command, and its manual page documents the options in this page. It is not installed by default on macOS or on many Linux systems; install it with a package manager. The tool tested here is tree v2.3.2.

### Examples
Show a folder as a tree; `--charset=ascii` draws the lines with plain characters so they copy well:

```sh
mkdir -p src/sub && printf 'x\n' > src/m.rs && printf 'y\n' > src/sub/n.md && tree --charset=ascii src
```

Show only folders, with `-d`:

```sh
mkdir -p src/sub logs && tree -d --charset=ascii .
```

Limit the depth to one level, with `-L 1`:

```sh
mkdir -p src/sub logs && printf 'x\n' > f.txt && tree -L 1 --charset=ascii .
```

Leave out items that match a pattern, with `-I`:

```sh
mkdir -p src logs && printf 'x\n' > f.txt && tree -I 'logs' -L 1 --charset=ascii .
```

### Option: -L
Show at most this many levels.

### Option: -d
Show directories only.

### Option: -a
Include hidden files, whose names start with a dot.

### Option: -I
Leave out names that match this pattern.

### Option: --charset
Choose the character set for the lines, such as `ascii`.

### Common mistakes
Hidden files are not shown unless you add `-a`. On a big folder, the output can be huge: limit it with `-L` or leave out folders such as `node_modules` with `-I`. The default lines use box-drawing characters, which look wrong in some terminals and when saved to a file with the wrong encoding; use `--charset=ascii`. The counts at the end include only what was shown, and in the version tested they count the folder you started from as a directory too.

### Portability
`tree` is not on every system by default. Some versions differ in the exact summary line (`1 directory, 2 files` or `1 directories, 2 files`). The examples were executed on tree v2.3.2. The Common examples were executed on tree v2.3.2 (see the development guide for the exact check).

## Upstream
### About upstream tree
This is the upstream `tree` command. It is one program, but its version and some of its options differ between packages.
