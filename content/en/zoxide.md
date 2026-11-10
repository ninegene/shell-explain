---
schema_version: 1
id: zoxide
title: zoxide
summary: Jump to folders you use often with a few letters.
aliases: []
tags: [shell, navigation, productivity]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/ajeetdsouza/zoxide#readme
---
## Common
### Overview
`zoxide` remembers the folders you visit and lets you jump to one with a short part of its name. After you set it up in your shell with `zoxide init`, the command `z proj` takes you to the folder that matches `proj` best, judged by how often and how recently you went there. It is a smarter `cd`.

### Background
`zoxide` is an open-source tool whose manual is in its repository. It keeps a small database in your data folder. The tool tested here is zoxide 0.10.0.

### Examples
Print the version:

```sh
zoxide --version
```

Tell it about a folder with `zoxide add`, then ask for the best match for a word; the answer is a full path, so the start is cut off:

```sh
mkdir p && zoxide add "$PWD/p" && zoxide query p | sed "s#$PWD/##"
```

Print the shell code that sets it up; the first line is a comment:

```sh
zoxide init bash | head -n 1 | cut -c1-30
```

### Subcommands
`zoxide init SHELL` prints the code that defines `z` and `zi`. `zoxide add PATH` records a folder. `zoxide query WORDS` prints the best match. `zoxide remove PATH` forgets a folder.

### Common mistakes
The `z` command does not exist until you put `eval "$(zoxide init bash)"` (or the line for your shell) in your shell's startup file. It learns from your visits, so a folder you have never visited is not found: use `cd` once. A word can match several folders; `zi` lets you choose with `fzf`. The database is a file, and removing a deleted folder from it is done with `zoxide remove` or with `zoxide query` that skips missing ones.

### Portability
The commands are the same on every system. The examples use a throw-away data folder, because the tests run with a clean home folder. They were executed on zoxide 0.10.0. The Common examples were executed on zoxide 0.10.0 (see the development guide for the exact check).

## Upstream
### About upstream zoxide
This is the upstream `zoxide` tool. It is one program everywhere; what differs is the version.
