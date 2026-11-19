---
schema_version: 1
id: tee
title: tee
summary: Copy what comes in to a file and to the screen at the same time.
aliases: [gtee]
tags: [files, pipes, streams]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/tee-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=tee&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnutee.1.html
---
## Common
### Overview
`tee` reads its standard input and writes it twice: to standard output and to each file you name. It is the pipe fitting that lets you save what flows through a pipeline without stopping the flow, like a T-shaped pipe, which is where the name comes from.

### Background
The macOS manual says `tee` is expected to be compatible with IEEE Std 1003.2 (POSIX.2), and it does not give a date for its first appearance. Today the name covers GNU tee, the tee that macOS ships, and the tee applet in BusyBox. They agree on writing to files and on `-a`; the handling of interrupts and of write errors differs.

### Examples
Show the text and save a copy of it in a file; `cat` shows what was saved:

```sh short="Show the text and save a copy; cat shows the copy"
echo hello | tee out.txt && cat out.txt
```

Save without showing the text, by sending standard output to the null device:

```sh
echo hello | tee out.txt > /dev/null && cat out.txt
```

Add to a file instead of replacing it with `-a`:

```sh
echo one | tee log.txt > /dev/null && echo two | tee -a log.txt > /dev/null && cat log.txt
```

Save to several files at once:

```sh
echo hi | tee a.txt b.txt > /dev/null && cat a.txt b.txt
```

Keep a copy of the middle of a pipeline. The file keeps the original order, while the pipeline sorts:

```sh
printf 'b\na\n' | tee unsorted.txt | sort && cat unsorted.txt
```

### Option: -a
Append to the files instead of replacing them.

### Common mistakes
Without `-a`, `tee` empties the file first, so `tee log.txt` loses what was in `log.txt`. `sudo echo text > file` fails because the shell, not `sudo`, opens the file; `echo text | sudo tee file > /dev/null` works because `tee` runs as the privileged user. `tee` copies only standard output of the earlier command; to include error messages, join them first with `2>&1`.

### Portability
Only the file names and `-a` are in the Common examples. `-i` to ignore interrupts is in POSIX, and GNU adds `-p` and `--output-error`, which are not on every system. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU tee
GNU tee is part of GNU coreutils, and this page's checks ran coreutils 9.11. It adds options for how write errors are handled, for use in long pipelines.

### GNU examples
Keep going and report the error when a write to a file fails, with `-p` (`--output-error`), instead of quietly stopping on a closed pipe:

```sh
echo hello | tee -p out.txt > /dev/null && cat out.txt
```

### Option: -p, --output-error
Diagnose errors writing to non-pipes and keep going; `--output-error=MODE` chooses the behavior in more detail.

## Apple
### About macOS tee
This is the tee in `/usr/bin` on macOS. It follows the BSD implementation, with `-a` and `-i`.

### Apple notes
macOS `tee` has `-a` and `-i`; it does not have GNU's `-p` or `--output-error`.

## BusyBox
### About BusyBox tee
BusyBox bundles tee as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `tee` supports `-a` and `-i`. GNU long options depend on the build.
