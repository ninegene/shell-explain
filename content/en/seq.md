---
schema_version: 1
id: seq
title: seq
summary: Print a sequence of numbers.
aliases: [gseq]
tags: [text, numbers, scripts]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/seq-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=seq&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuseq.1.html
---
## Common
### Overview
`seq` prints numbers, one per line, from a first value to a last value. It is the quick way to make test input, to count in a shell loop, or to number things. With one argument it counts from 1; with two it counts from the first to the second; with three the middle one is the step.

### Background
The `seq` command first appeared in Version 8 AT&T UNIX, according to the macOS manual, and the macOS one is based on the Plan 9 and GNU commands of the same name. Today the name covers GNU seq, the seq that macOS ships, and the seq applet in BusyBox. They agree on counting up and down with a step, and on `-w`; the separator option differs in how it ends the output.

### Examples
Count from 1 to 3:

```sh
seq 3
```

Count between two numbers:

```sh
seq 2 4
```

Count with a step of 2:

```sh
seq 1 2 7
```

Count down by giving a negative step:

```sh
seq 5 -1 1 | paste -sd, -
```

Pad the numbers with zeros to the same width with `-w`:

```sh
seq -w 8 10
```

Use it to repeat something three times:

```sh
for i in $(seq 3); do echo "run $i"; done
```

### Option: -w
Pad numbers with leading zeros so they all have the same width.

### Common mistakes
`seq 5 1` prints nothing, because counting up from 5 never reaches 1; give a negative step, as in `seq 5 -1 1`. `seq` is not in POSIX; on a system without it, `i=1; while [ "$i" -le 3 ]; do ...; i=$((i+1)); done` does the job, and `{1..3}` works in Bash but not in every shell.

### Portability
Only counting with a step and `-w` are in the Common examples. The separator option `-s` ends its output differently on macOS and GNU (see the sections below). The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU seq
GNU seq is part of GNU coreutils, and this page's checks ran coreutils 9.11. It has `-s` to choose the separator and `-f` for a printf-style format.

### GNU examples
Put commas between the numbers; the last number is followed by a newline:

```sh
seq -s, 1 4
```

### Option: -s, --separator
Print this text between numbers instead of a newline. The output ends with a newline.

### Option: -f, --format
Format each number with a printf-style format such as `%03g`.

## Apple
### About macOS seq
This is the seq in `/usr/bin` on macOS. It follows the BSD implementation.

### Apple notes
With `-s`, macOS `seq` prints the separator after every number, including the last, and no newline, so `seq -s, 1 4` ends in a comma. Use `seq 4 | paste -sd, -` for a portable comma-separated list.

## BusyBox
### About BusyBox seq
BusyBox bundles seq as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `seq` supports counting with a step and `-w`; `-s` and `-f` depend on the build.
