---
schema_version: 1
id: diff3
title: diff3
summary: Compare three files and show where they differ.
aliases: [gdiff3]
tags: [files, text, compare]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/diffutils/manual/html_node/Invoking-diff3.html
  - https://man.freebsd.org/cgi/man.cgi?query=diff3&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/diff3.1.html
---
## Common
### Overview
`diff3` compares three files at once: your version, the original they both came from, and someone else's version. It reports which lines changed in which file, which is what a version control tool needs to merge two sets of changes. With `-m` it produces the merged file, with conflict markers where both sides changed the same lines.

### Background
Today the name covers GNU diff3 and the diff3 that macOS ships. The BusyBox 1.37.0 build used for this page has no `diff3` applet. They agree on comparing three files and on `-m` for a merge.

### Examples
Compare three files; the lines marked `====` are blocks where the files differ:

```sh
printf 'a\nb\n' > mine && printf 'a\nb\n' > base && printf 'a\nb\nd\n' > yours && diff3 mine base yours
```

Merge changes that do not overlap: you changed line 1 and they added a line at the end:

```sh
printf 'A\nb\n' > mine && printf 'a\nb\n' > base && printf 'a\nb\nd\n' > yours && diff3 -m mine base yours
```

### Common mistakes
The order of the arguments is important: `diff3 MINE OLDER YOURS`, with the common ancestor in the middle. When both sides change the same lines, `-m` writes conflict markers into the output and exits with status 1, and you must fix them by hand. `diff3` merges lines, not meaning, so a clean merge can still be wrong; run your tests. Most people use version control, which does this for them.

### Portability
This page has no BusyBox view: BusyBox has no `diff3` applet. The Common examples avoid the exact layout of the report and a conflict. The Common examples were executed on GNU diffutils 3.12 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU diff3
GNU diff3 is part of GNU diffutils, and this page's checks ran version 3.12.

## Apple
### About macOS diff3
This is the diff3 in the system folders of macOS.
