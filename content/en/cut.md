---
schema_version: 1
id: cut
title: cut
summary: Select columns or character ranges from each line.
aliases: [gcut]
tags: [text, files, columns]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/cut-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=cut&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`cut` prints chosen parts of every input line: fields split by a delimiter, or character positions. It is the simple way to pull a column out of delimited text.

### Examples
Print the first comma-separated field:

```sh
cut -d, -f1 data.csv
```

Print the second field:

```sh
cut -d, -f2 data.csv
```

Print fields 1 and 2:

```sh
cut -d, -f1,2 data.csv
```

Print field 2 through the end of the line:

```sh
cut -d, -f2- data.csv
```

Print the first three characters of each line:

```sh
cut -c1-3 notes.txt
```

Split text from a pipeline on spaces; quote the space delimiter:

```sh
echo 'a b c' | cut -d' ' -f2
```

### Option: -d
Set the field delimiter to this single character. The default is a tab.

### Option: -f
Select fields by number: `1`, `1,3`, `2-`, or `-2`. Lines without the delimiter are printed whole unless you suppress them.

### Option: -c
Select character positions, for example `1-3`. In some locales a character is more than one byte.

### Common mistakes
`cut` has no notion of quoted fields, so a comma inside quotes in CSV still splits the line. It cannot reorder fields: `-f2,1` prints them in file order. Runs of spaces are separate delimiters, so use `awk` for whitespace-aligned columns.

### Portability
`-d`, `-f`, and `-c` with the ranges above are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### GNU examples
Print every field except the first. `--complement` is a GNU extension:

```sh
cut --complement -d, -f1 data.csv
```

## Apple
### Apple notes
macOS `cut` is the BSD implementation. It has no long options such as `--complement` or `--delimiter`; use the short spellings.

## BusyBox
### BusyBox notes
BusyBox `cut` supports `-d`, `-f`, and `-c`; other options depend on the build.
