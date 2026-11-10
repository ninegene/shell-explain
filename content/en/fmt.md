---
schema_version: 1
id: fmt
title: fmt
summary: Reflow paragraphs of text to a chosen line width.
aliases: [gfmt]
tags: [text, formatting]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/fmt-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=fmt&sektion=1
---
## Common
### Overview
`fmt` reads text and joins and breaks its lines so each is about the same length, which is useful for tidying a paragraph. Unlike `fold`, it breaks only at spaces and it also joins short lines into longer ones. The default width is about 75 characters.

### Background
The `fmt` command appeared in 3BSD, and the version on macOS is a complete rewrite that appeared in FreeBSD 4.4, according to the macOS manual. Today the name covers GNU fmt and the fmt that macOS ships. The BusyBox 1.37.0 build used for this page has no `fmt` applet. They agree on `-w` for the width.

### Examples
Reflow text to at most 15 characters per line:

```sh
echo 'the quick brown fox jumps' | fmt -w 15
```

Join short lines into one paragraph:

```sh
printf 'one\ntwo\nthree\n' | fmt
```

Keep separate paragraphs apart, since an empty line ends a paragraph:

```sh
printf 'one\ntwo\n\nthree\nfour\n' | fmt
```

### Option: -w
Set the maximum line width, as in `-w 15`.

### Common mistakes
`fmt` treats a line that starts with extra spaces as the start of a new paragraph, so indented lists may not reflow the way you expect. It does not understand Markdown or code, so do not run it on source files. The goal width is a little shorter than the maximum, so lines may be slightly shorter than `-w` suggests.

### Portability
This page has no BusyBox view: BusyBox has no `fmt` applet. Only `-w` is in the Common examples; Both have `-s` (split only), but `-p` means different things: GNU takes a prefix, and macOS allows indented paragraphs. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check). The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU fmt
GNU fmt is part of GNU coreutils, and this page's checks ran coreutils 9.11. `-p PREFIX` reformats only lines that start with the prefix, which is useful for comment blocks.

### Option: -p, --prefix
Reformat only lines that begin with this prefix, and keep the prefix.

## Apple
### About macOS fmt
This is the fmt in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `fmt` accepts `-w`, and also a width as a plain number such as `-15`.
