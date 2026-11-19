---
schema_version: 1
id: pr
title: pr
summary: Format text for printing, in pages and columns.
aliases: [gpr]
tags: [text, formatting]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/pr-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=pr&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnupr.1.html
---
## Common
### Overview
`pr` prepares text for a printer: it splits it into pages, adds a header with the date and file name to each, and can arrange the text in columns or side by side. In daily use it is mostly handy for its columns and line numbers, with `-t` to drop the page header and footer.

### Background
A `pr` command appeared in Version 1 AT&T UNIX, according to the macOS manual, which also says it is POSIX.1-2001 compatible. Today the name covers GNU pr and the pr that macOS ships. The BusyBox 1.37.0 build used for this page has no `pr` applet. They agree on `-t`, columns and `-m`.

### Examples
Print without the page header and footer, with `-t`:

```sh
printf 'a\nb\n' | pr -t
```

Arrange a list in two columns, down first:

```sh
seq 4 | pr -t -2
```

Put two files side by side with `-m`:

```sh
printf 'a\nb\n' > f1 && printf '1\n2\n' > f2 && pr -t -m f1 f2
```

### Option: -t
Omit the header and footer of each page.

### Option: -m
Print all files side by side, one per column.

### Common mistakes
Without `-t`, `pr` adds a header with the current date and time to every page, so its output changes each time you run it. The column layout uses tabs and padding that differ between systems, so do not parse it. The default page length is 66 lines, and each page is padded with blank lines to that length.

### Portability
This page has no BusyBox view: BusyBox has no `pr` applet. Only `-t`, a column count and `-m` are in the Common examples, compared with runs of blanks collapsed. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU pr
GNU pr is part of GNU coreutils, and this page's checks ran coreutils 9.11. It also numbers lines with `-n`.

## Apple
### About macOS pr
This is the pr in the system folders of macOS. It follows the BSD implementation.
