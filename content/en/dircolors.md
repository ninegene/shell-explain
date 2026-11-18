---
schema_version: 1
id: dircolors
title: dircolors
summary: Set up the colors that ls uses for file types.
aliases: [gdircolors]
tags: [terminal, files, configuration]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/dircolors-invocation.html
---
## Common
### Overview
`dircolors` prints shell commands that set the `LS_COLORS` variable, which GNU `ls --color` uses to choose colors by file type or name. You usually run it from a shell startup file as `eval "$(dircolors -b)"`. With `-p` it prints the built-in table so you can copy and change it.

### Background
`dircolors` is part of GNU coreutils; the GNU manual documents the `LS_COLORS` variable it sets. macOS does not ship it and its `ls` uses the `LSCOLORS` variable in a different format, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Print the commands for a Bourne-style shell with `-b`; the first line sets the variable:

```sh short="Print shell commands that set LS_COLORS (-b)"
dircolors -b | grep -c '^LS_COLORS='
```

Print them for a C shell with `-c`:

```sh
dircolors -c | grep -c '^setenv LS_COLORS'
```

Print the built-in table; the lines that start with `DIR` set the directory color:

```sh short="Print the built-in table (DIR sets directory color)"
dircolors -p | grep -c '^DIR'
```

### Option: -b
Print commands for a Bourne-style shell such as bash or zsh.

### Option: -c
Print commands for a C shell.

### Option: -p
Print the built-in configuration table.

### Common mistakes
`dircolors` only prints text; the colors change when you run its output, for example with `eval`. macOS `ls` ignores `LS_COLORS`: use `ls -G` and `LSCOLORS` there, or GNU `gls --color`. A bad entry in your own table is ignored in silence, so test a change by running `ls --color`.

### Portability
This page covers GNU only. The two variables, `LS_COLORS` for GNU and `LSCOLORS` for BSD, have different formats and are not interchangeable. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU dircolors
GNU dircolors is part of GNU coreutils, and this page's checks ran coreutils 9.11.
