---
schema_version: 1
id: ptx
title: ptx
summary: Make a permuted index, listing each word with the text around it.
aliases: [gptx]
tags: [text, index]
implementations: [gnu]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/ptx-invocation.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuptx.1.html
---
## Common
### Overview
`ptx` builds a permuted index, also called a keyword-in-context index. For each word of the input, it prints a line with that word lined up in a column and the words that came before and after it. It is a way to look up where a word is used. It is an old tool, mostly of interest for text indexing.

### Background
`ptx` is part of GNU coreutils; the GNU manual documents it as producing a permuted index of file contents. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Make an index of one short line, with an output width of 40 characters; the keyword column is in the middle, and the spacing is compared with blanks collapsed:

```sh short="Make an index of one short line, 40 characters wide"
printf 'the quick fox\n' | ptx -w 40
```

Count the lines: there is one for each word in the input:

```sh
printf 'the quick fox\n' | ptx -w 40 | grep -c .
```

### Option: -w
Set the width of the output in characters, as in `-w 40`.

### Common mistakes
The output is made of columns padded with spaces, so it is meant to be read on a screen, not parsed. Common words such as `the` get an entry too, unless you give a list of words to ignore with `-i`. The width defaults to 72 and has to be wide enough for the context, or the context is cut. Other tools, such as `grep -n`, are simpler when you only need the places where a word is used.

### Portability
This page covers GNU only. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU ptx
GNU ptx is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-i FILE` takes a list of words to leave out, and `-o FILE` a list of words to include only.
