---
schema_version: 1
id: sdiff
title: sdiff
summary: Compare two files side by side, or merge them interactively.
aliases: [gsdiff]
tags: [files, text, compare]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/diffutils/manual/html_node/Invoking-sdiff.html
  - https://man.freebsd.org/cgi/man.cgi?query=sdiff&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/sdiff.1.html
---
## Common
### Overview
`sdiff` shows two files in two columns, with a mark between them: nothing for the same line, `|` for a changed line, `<` for a line only on the left and `>` for a line only on the right. With `-o FILE` it asks you, line group by line group, which side to keep and writes the merged file.

### Background
Today the name covers GNU sdiff and the sdiff that macOS ships. The BusyBox 1.37.0 build used for this page has no `sdiff` applet. They agree on the two-column view and on `-w` for the width.

### Examples
Show two files side by side; the `|` marks a changed line, and the exact spacing is compared with blanks collapsed:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && sdiff f1 f2; test $? -eq 1 && echo differ
```

Set the total width of the output with `-w`:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && sdiff -w 20 f1 f2; test $? -eq 1 && echo differ
```

Show only the lines that differ with `-s`, leaving out the common ones:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && sdiff -s f1 f2; test $? -eq 1 && echo differ
```

### Option: -w
Set the width of the output in characters.

### Option: -s
Do not print lines that are the same.

### Option: -o
Merge interactively and write the result to this file.

### Common mistakes
`sdiff` exits with status 1 when the files differ, like `diff`. The columns are padded with spaces and tabs, which differ between systems, so do not parse the output. Lines longer than half the width are cut. The interactive mode with `-o` waits for you to type a choice for each block, so it is not for scripts.

### Portability
This page has no BusyBox view: BusyBox has no `sdiff` applet. The layout of the columns (tabs or spaces) differs between GNU and macOS, so the examples compare output with runs of blanks collapsed. The Common examples were executed on GNU diffutils 3.12 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU sdiff
GNU sdiff is part of GNU diffutils, and this page's checks ran version 3.12.

## Apple
### About macOS sdiff
This is the sdiff in the system folders of macOS.
