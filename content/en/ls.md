---
schema_version: 1
id: ls
title: ls
summary: List files and directories.
aliases: [gls]
tags: [files, directories]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/ls-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=ls&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`ls` lists the names in a directory, or the files you name. Names are sorted, and when the output goes to a pipe or file each name is on its own line.

### Examples
List the current directory:

```sh
ls
```

List the contents of one directory:

```sh
ls logs
```

List only the files that match a pattern; the shell expands `*.log` before `ls` runs:

```sh
ls *.log
```

Show a directory itself instead of its contents. The trailing slash in the pattern matches only directories:

```sh
ls -d */
```

Include entries that start with a dot, including `.` and `..`:

```sh
ls -a logs
```

List directories and everything below them. Naming two directories keeps the output the same on every implementation, with a header line for each:

```sh
ls -R logs projects
```

Reverse the sort order:

```sh
ls -r
```

Show the permission string of a file; `-l` is the long format and `cut` keeps its first ten characters:

```sh
ls -l notes.txt | cut -c1-10
```

### Option: -a
Include hidden entries, those whose names begin with a dot.

### Option: -d
List directories themselves rather than their contents.

### Option: -l
Use the long format: permissions, link count, owner, group, size, modification time and name. The columns are padded differently between implementations, so scripts should not parse them.

### Option: -R
List subdirectories recursively.

### Option: -r
Reverse the sort order.

### Common mistakes
Do not parse `ls` output in scripts. Names can contain spaces or newlines, and the long format varies by implementation; use a glob or `find` instead. `ls` with no operands lists the current directory, but an unmatched pattern such as `ls *.log` passes the literal text and fails if nothing matches.

### Portability
The Common examples use only `-a`, `-d`, `-l`, `-R` and `-r`, which all three implementations accept. With a single directory, BSD `ls -R` leaves out the header line that GNU prints. Colour and sorting options such as `--color` and `-t` differ between GNU and BSD. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### GNU examples
List directories before files; this option is a GNU extension:

```sh
ls --group-directories-first
```

## Apple
### Apple notes
macOS `ls` is the BSD implementation. Colour is enabled with `-G` rather than GNU's `--color`, and it has no `--group-directories-first`.

## BusyBox
### BusyBox notes
BusyBox `ls` supports the options in the common examples; GNU long options and colour depend on the build.
