---
schema_version: 1
id: fzf
title: fzf
summary: Pick a line from a list by typing a few letters.
aliases: []
tags: [search, interactive, productivity]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://github.com/junegunn/fzf#readme
---
## Common
### Overview
`fzf` is a fuzzy finder. You give it a list of lines, such as file names or command history, and it shows an interactive list where you type a few letters, in any order, and the list shrinks to the lines that match. You press Enter to choose. It prints the chosen line, which you can use in a command. With `-f` it filters without the interface, so it also works in scripts.

### Background
`fzf` is an open-source tool whose manual is in its repository. Many shell setups bind it to Ctrl+R (history) and Ctrl+T (files). The tool tested here is fzf 0.74.5. The interactive example needs a terminal and is reviewed only.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
fzf --version
```

Filter a list without the interface, with `-f`; the letters `an` match lines that contain `a` then `n`:

```sh short="Filter a list without the interface (-f)"
printf 'apple\nbanana\ncherry\n' | fzf -f an
```

Match at the start of a line with a `^`:

```sh
printf 'apple\nbanana\ncherry\n' | fzf -f '^ch'
```

Match exactly, with a leading quote:

```sh
printf 'apple\nbanana\ncherry\n' | fzf -f "'err"
```

Choose a file interactively and open it in an editor; reviewed only because it needs a terminal:

```sh
vim "$(find . -type f | fzf)"
```

### Option: -f
Filter: print the matches for this query and exit, without the interface.

### Option: --exact
Exact matching instead of fuzzy.

### Option: -m
Allow choosing several lines.

### Common mistakes
Without `-f`, `fzf` needs a terminal, and the chosen line goes to standard output, so use it with `$(...)`. An empty choice (pressing Esc) returns nothing and status 130, so check before using the result. `fzf` matches letters in order but not next to each other, so a short query can match many lines: add a leading quote for an exact match or `^` for the start. The default list of files comes from `find`, and `FZF_DEFAULT_COMMAND` changes it, for example to `fd`.

### Portability
The matching rules and the options above are the same on every system, and new options arrive with new versions. The examples that filter were executed on fzf 0.74.5 (see the development guide for the exact check); the version and interactive examples are reviewed only.

## Upstream
### About upstream fzf
This is the upstream `fzf` tool. It is one program everywhere; what differs is the version.
