---
schema_version: 1
id: shuf
title: shuf
summary: Shuffle lines into a random order, or pick random lines.
aliases: [gshuf]
tags: [text, random]
implementations: [gnu, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/shuf-invocation.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`shuf` reads lines and prints them in a random order. With `-n` it prints only some of them, so it can pick a random line from a list. With `-e` it shuffles the words you give on the command line instead of reading input. Because the order is random, the examples sort the output or count lines.

### Background
`shuf` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it. They agree on shuffling lines, on `-n` and on `-e`.

### Examples
Shuffle the numbers 1 to 5; sorting the result shows that none was lost:

```sh
seq 5 | shuf | sort -n | paste -sd, -
```

Shuffle the words given as arguments with `-e`:

```sh
shuf -e a b c | sort | paste -sd, -
```

Pick three random lines from a file with `-n`:

```sh
seq 6 > n6.txt && shuf -n 3 n6.txt | grep -c .
```

Pick one random line; here the file has one line, so the answer is known:

```sh
printf 'only\n' | shuf -n 1
```

### Common mistakes
`shuf` is not a source of security-grade randomness; use a proper tool for keys. Without `-n`, every line is read into memory. Run twice, it gives different results, so a test that compares output needs `sort` or a count. `-n` larger than the number of lines just prints all of them.

### Portability
This page has no Apple view: macOS has no `shuf`. Common substitutes there are `sort -R` (which groups equal lines together) and `jot -r`. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check). The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU shuf
GNU shuf is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can draw with repeats and from a numeric range.

### GNU examples
Shuffle a numeric range, without a file, with `-i`:

```sh
shuf -i 1-5 | sort -n | paste -sd, -
```

Allow the same item to be picked more than once with `-r`:

```sh
shuf -r -n 5 -e x | paste -sd, -
```

### Option: -i, --input-range
Use the numbers in a range such as `1-5` as the lines.

### Option: -r, --repeat
Pick with replacement, so lines can repeat.

## BusyBox
### About BusyBox shuf
BusyBox bundles shuf as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `shuf` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
