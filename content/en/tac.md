---
schema_version: 1
id: tac
title: tac
summary: Print the lines of a file in reverse order.
aliases: [gtac]
tags: [text, files]
implementations: [gnu, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/tac-invocation.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`tac` is `cat` spelled backwards, and it prints a file from the last line to the first. It is handy for reading a log newest-first, or for processing lines in reverse. macOS does not ship `tac`; there, `tail -r` does the same.

### Background
`tac` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it, and uses `tail -r` for the same job. The name is `cat` reversed, and the two agree on reversing lines.

### Examples
Print a file with its lines reversed:

```sh
tac notes.txt
```

Show the last line first and keep only two lines:

```sh
tac notes.txt | head -n 2
```

Reverse the output of another command:

```sh
printf 'a\nb\nc\n' | tac
```

Number the lines from the end, by reversing before `nl`:

```sh
tac notes.txt | nl | head -n 1
```

### Common mistakes
`tac` reverses lines, not the characters inside them; use `rev` for that. It needs the whole input before it can print, so it does not work on endless input such as `tail -f`. A file whose last line has no newline ends up with its first and second lines run together unless you add one.

### Portability
This page has no Apple view: macOS has no `tac`, and `tail -r file` is its equivalent, which GNU `tail` does not have. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU tac
GNU tac is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-s` changes the separator that ends a record, and `-r` treats it as a regular expression.

### Option: -s, --separator
Use this string instead of a newline as the end of each record.

## BusyBox
### About BusyBox tac
BusyBox bundles tac as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `tac` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
