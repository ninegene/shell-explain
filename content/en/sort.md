---
schema_version: 1
id: sort
title: sort
summary: Sort lines of text.
aliases: [gsort]
tags: [text, files, sorting]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/sort-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=sort&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`sort` reads lines from files or standard input and prints them in order. By default it compares text using the current locale, so punctuation, case, and numbers sort as text.

### Examples
Sort lines alphabetically:

```sh
sort fruit.txt
```

Sort in reverse order:

```sh
sort -r fruit.txt
```

Sort numbers by value instead of as text:

```sh
sort -n nums.txt
```

Sort and drop duplicate lines:

```sh
sort -u fruit.txt
```

Sort a comma-separated file by its second column; the header sorts first because it has no number:

```sh
sort -t, -k2,2n data.csv
```

Find the largest number:

```sh
sort -rn nums.txt | head -n 1
```

### Option: -r
Reverse the order.

### Option: -n
Compare as numbers, not text. Without it, `10` sorts before `9`.

### Option: -u
Output only the first of lines that compare equal, which also removes duplicates.

### Option: -t
Set the field separator for `-k` to this single character, as in `-t,`.

### Option: -k
Sort by a field or character range. `-k2,2n` means start and stop at field 2, numerically; without the `,2` the key runs to the end of the line.

### Common mistakes
Without `-n`, numbers sort as text, so `10` comes before `9`. Order depends on the locale: the examples here use the C locale (`LC_ALL=C`) for predictable results, and other locales can sort case and punctuation differently. `sort file > file` empties the file before sorting; write to another file.

### Portability
`-r`, `-n`, `-u`, `-t`, and `-k` are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). Version sort `-V`, human-numeric `-h`, and stability options differ between builds.

## GNU
### GNU notes
GNU sort adds options such as `-h` (human-readable sizes like 2K or 1G) and `-V` (version order). Confirm that a macOS or BusyBox build provides an option before depending on it.

## Apple
### Apple notes
macOS `sort` identifies itself as `2.3-Apple` and is derived from GNU sort, so most common options match; confirm rarely used long options before relying on them.

## BusyBox
### BusyBox notes
BusyBox `sort` supports the common options above. Options such as `-h`, `-V`, and `--parallel` depend on the build configuration.
