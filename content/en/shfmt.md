---
schema_version: 1
id: shfmt
title: shfmt
summary: Format shell scripts in one consistent style.
aliases: []
tags: [shell, scripts, format]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://github.com/mvdan/sh#shfmt
  - https://manpages.ubuntu.com/manpages/resolute/man1/shfmt.1.html
---
## Common
### Overview
`shfmt` reads a shell script and prints it again in a standard layout: consistent indentation, spaces around operators and one command per line. It parses the script, so it understands the syntax and does not just move text. With `-w` it rewrites the file, with `-l` it lists files that are not formatted, and with `-d` it shows the difference.

### Background
`shfmt` is part of the open-source `sh` project, which also contains a shell parser. It formats the languages bash, POSIX sh and mksh. The tool tested here is shfmt v3.14.1.

### Examples
Format a script from standard input; the layout of `if` is tidied:

```sh short="Format a script from standard input"
printf 'if true;then echo hi;fi\n' | shfmt
```

Indent with two spaces instead of tabs, with `-i 2`:

```sh
printf 'if true;then\necho hi\nfi\n' | shfmt -i 2
```

List files that are not yet formatted, with `-l`; the status is 1 when it lists one, in the version tested:

```sh
printf 'if true;then echo hi;fi\n' > f.sh && shfmt -l f.sh; echo "status $?"
```

Show what would change, as a diff, with `-d`; the first lines show the file names:

```sh
printf 'if true;then echo hi;fi\n' > f.sh && shfmt -d f.sh | head -n 3 | cut -c1-6
```

### Option: -i
Indent with this many spaces; `0` means tabs.

### Option: -l
List files whose formatting differs.

### Option: -d
Print a diff of what would change.

### Option: -w
Write the result to the file in place.

### Option: -ln
Choose the language: `bash`, `posix`, `mksh` or `bats`.

### Common mistakes
Without `-w`, `shfmt` prints the result and leaves the file alone; with `-w` it rewrites it in place, so use version control. The default is to indent with tabs: use `-i N` for spaces. It guesses the shell dialect from the `#!` line and the file name; give `-ln bash` or `-ln posix` when it cannot. A script with a syntax error is not formatted, and `shfmt` says where the error is.

### Portability
The options are the same on every system, and the default style can change a little between versions. The examples were executed on shfmt v3.14.1. The Common examples were executed on shfmt v3.14.1 (see the development guide for the exact check).

## Upstream
### About upstream shfmt
This is the upstream `shfmt` tool. It is one program everywhere; what differs is the version.
