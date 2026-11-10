---
schema_version: 1
id: bat
title: bat
summary: Show files with syntax highlighting, line numbers and paging.
aliases: []
tags: [files, text, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/sharkdp/bat#readme
---
## Common
### Overview
`bat` prints files like `cat`, but at a terminal it adds syntax highlighting, line numbers and a frame, and pages long files. When its output goes into a pipe or a file it prints the plain text only, so it can replace `cat` without changing scripts. It also marks lines changed since the last git commit.

### Background
`bat` is an open-source tool whose manual is in its repository. It is a separate program from `cat` and is installed on request. The tool tested here is bat 0.26.1.

### Examples
Print a file without any decoration, with `-p` (plain):

```sh
printf 'one\ntwo\n' > f.txt && bat -p f.txt
```

Show line numbers; the option `--decorations=always` keeps them even when the output is a pipe:

```sh
printf 'one\ntwo\n' > f.txt && bat --decorations=always --style=numbers f.txt
```

Say which language to highlight with `-l`, for text that has no file name:

```sh
echo '{"a":1}' | bat -l json --style=plain
```

Show only a line range with `-r`:

```sh
printf 'one\ntwo\nthree\n' > f.txt && bat -p -r 2:2 f.txt
```

### Option: -p
Plain: no line numbers, no frame.

### Option: -l
Set the language used for highlighting.

### Option: --style
Choose which decorations to show, such as `numbers`, `grid` or `plain`.

### Option: -r
Print only this line range, as in `-r 10:20`.

### Common mistakes
When output goes to a pipe `bat` behaves like `cat`, which is good for scripts, but it can look as if the options do nothing; add `--decorations=always` to force them. On Debian and Ubuntu the command is named `batcat` because `bat` was already taken. Its pager is `less`; press `q` to quit. Do not use `bat` for binary files you want to see unchanged: it shows them with escapes.

### Portability
`bat` is a different program from `cat`, so options such as `-A` and `-s` do not mean what they do in `cat`. The examples were executed on bat 0.26.1. The Common examples were executed on bat 0.26.1 (see the development guide for the exact check).

## Upstream
### About upstream bat
This is the upstream `bat` tool. It is one program on every system, though Debian and Ubuntu name it `batcat`.
