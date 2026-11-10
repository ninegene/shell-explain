---
schema_version: 1
id: rg
title: rg
summary: Search the files under a folder for a pattern, fast, ignoring what .gitignore ignores.
aliases: [ripgrep]
tags: [search, text, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/BurntSushi/ripgrep/blob/master/GUIDE.md
---
## Common
### Overview
`rg` (ripgrep) searches the text in all the files under the current folder, or the places you name, and prints the lines that match. It is built for searching source code: it is fast, it skips files that `.gitignore` lists, hidden files and binary files, and it uses regular expressions by default. It prints the file name and line number when it searches many files.

### Background
`rg` is an open-source tool whose guide is in its repository. It is a separate program from `grep`, with different defaults, and is installed on request, not by default on most systems. The tool tested here is ripgrep 15.2.0.

### Examples
Search one file; the matching line is printed:

```sh
printf 'one\ntwo\nthree\n' > f.txt && rg two f.txt
```

Add the line number with `-n`:

```sh
printf 'one\ntwo\nthree\n' > f.txt && rg -n two f.txt
```

Count the lines that match, with `-c`:

```sh
printf 'one\ntwo\nthree\n' > f.txt && rg -c o f.txt
```

Ignore case with `-i`:

```sh
printf 'one\ntwo\nthree\n' > f.txt && rg -i TWO f.txt
```

Search several places; each match starts with its file name:

```sh
mkdir logs && printf 'two\n' > logs/a.log && printf 'one\ntwo\n' > f.txt && rg two logs f.txt | sort
```

List the files that contain a match, with `-l`:

```sh
mkdir logs && printf 'two\n' > logs/a.log && printf 'one\n' > f.txt && rg -l two .
```

Match whole words only with `-w`; the status is 1 when nothing matches:

```sh
printf 'one\ntwo\n' > f.txt && { rg -w tw f.txt || echo none; }
```

List the files rg would search, with `--files`, and restrict by name with `-g`:

```sh
mkdir src && printf 'x\n' > src/m.rs && printf 'y\n' > f.txt && rg -g '*.rs' --files
```

### Option: -n
Show line numbers.

### Option: -i
Ignore case.

### Option: -w
Match whole words only.

### Option: -l
Print only the names of files with a match.

### Option: -c
Print the number of matching lines in each file.

### Option: -g
Include or exclude files by a glob, as in `-g '*.rs'`.

### Option: -F
Treat the pattern as a fixed string, not a regular expression.

### Option: --files
List the files that would be searched.

### Common mistakes
`rg` skips files that are ignored by `.gitignore`, hidden files and binary files, so a file you expect may not be searched; add `--no-ignore`, `--hidden` or `-uuu` to include them. The pattern is a regular expression, so `.`, `(` and `[` have meanings; use `-F` for plain text. When its output is not a terminal, `rg` leaves out the file name for a single file and the heading layout, which can surprise a script. The order of results from several files can vary between runs, because it searches in parallel; add `--sort path` for a stable order.

### Portability
`rg` is a different program from `grep`, so options such as `-r` (which means replace in `rg`) do not mean what they do in `grep`. The examples were executed on ripgrep 15.2.0. The Common examples were executed on ripgrep 15.2.0 (see the development guide for the exact check).

## Upstream
### About upstream rg
This is the upstream ripgrep tool, installed as `rg`. It is one program on every system; the version is what differs.
