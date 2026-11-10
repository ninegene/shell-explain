---
schema_version: 1
id: fd
title: fd
summary: Find files and folders by name, in a simpler and faster way than find.
aliases: []
tags: [files, search, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/sharkdp/fd#readme
---
## Common
### Overview
`fd` finds files and folders whose names match a pattern. Without options it searches the whole tree below the current folder for names that contain the pattern, ignoring case when the pattern is all lower case, and it skips hidden files and what `.gitignore` lists. It is a friendlier alternative to `find`, for the common job of finding a file by name.

### Background
`fd` is an open-source tool whose manual is in its repository. It is a separate program from `find`, with a different syntax, and is installed on request. The tool tested here is fd 10.5.0.

### Examples
Find files whose name contains `txt`:

```sh
mkdir src && printf 'x\n' > f.txt && printf 'y\n' > src/n.md && fd txt
```

Find files with an extension, with `-e`:

```sh
mkdir -p src/sub && printf 'x\n' > src/sub/n.md && fd -e md
```

Find only folders, with `-t d`:

```sh
mkdir -p logs src/sub && fd -t d | sort
```

Search inside one folder by naming it after the pattern; the pattern `.` matches everything:

```sh
mkdir logs && printf 'two\n' > logs/a.log && fd . logs
```

Search a folder for a pattern; the folder comes second:

```sh
mkdir src && printf 'x\n' > src/m.rs && fd m src
```

### Option: -e
Only names with this extension.

### Option: -t
Only this type: `f` for files, `d` for directories, `l` for links.

### Option: -H
Include hidden files.

### Option: -I
Do not respect `.gitignore`.

### Common mistakes
`fd` skips hidden files and the files that `.gitignore` lists, so a file you expect may not appear; add `-H` and `-I` to include them. Its pattern is a regular expression matched against the name, not a glob, so a star does not mean what it does in a shell; use `fd -e txt`, or `fd -g` with a glob. Folders are printed with a trailing slash. The order of results is not guaranteed to be the same between runs, because it searches in parallel.

### Portability
`fd` is a different program from `find`, and its options do not map one to one. The examples were executed on fd 10.5.0. On Debian and Ubuntu the program is installed under the name `fdfind`. The Common examples were executed on fd 10.5.0 (see the development guide for the exact check).

## Upstream
### About upstream fd
This is the upstream `fd` tool. It is one program on every system, though some Linux packages name it `fdfind`.
