---
schema_version: 1
id: diff
title: diff
summary: Show the differences between two files or two folders.
aliases: [gdiff]
tags: [files, text, compare]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/diffutils/manual/html_node/Invoking-diff.html
  - https://man.freebsd.org/cgi/man.cgi?query=diff&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`diff` compares two files line by line and prints the lines that differ, in a form that `patch` can apply. It exits with status 0 when the files are the same, 1 when they differ and 2 on an error, so it also answers "are these the same?". The unified format (`-u`) is the one used by version control and code review. With `-r` it compares two folders, file by file.

### Background
Today the name covers GNU diff, the diff that macOS ships (based on FreeBSD diff), and the diff applet in BusyBox. They agree on `-u`, `-q`, `-r`, `-i`, `-w` and `-B`, and on the exit statuses. The default format and the side-by-side `-y` are not in every implementation.

### Examples
Compare two files in the unified format; the first two lines name the files and carry times, so `tail` starts at the hunk:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && { diff -u f1 f2 | tail -n +3; }
```

Only say whether the files differ, with `-q`; the status is 1 when they do:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && diff -q f1 f2; echo "status $?"
```

Use the exit status: 0 means the same:

```sh
printf 'a\n' > f1 && cp f1 f2 && diff f1 f2 && echo same
```

Use the exit status the other way: 1 means different, so `||` runs:

```sh
printf 'a\n' > f1 && printf 'b\n' > f2 && { diff -q f1 f2 > /dev/null || echo differ; }
```

Ignore differences in the amount of white space with `-w`:

```sh
printf 'x  y\n' > w1 && printf 'x y\n' > w2 && diff -w w1 w2 && echo same-ignoring-space
```

Compare two folders with `-r`; identical folders give no output and status 0:

```sh
mkdir d1 d2 && printf 'x\n' > d1/f && printf 'x\n' > d2/f && diff -r d1 d2 && echo same
```

### Option: -u
Print the differences in the unified format, with context lines and `+`/`-` marks.

### Option: -q
Only report whether the files differ.

### Option: -r
Compare folders recursively.

### Option: -i
Ignore case differences.

### Option: -w
Ignore all white space.

### Option: -B
Ignore changes that only add or remove blank lines.

### Common mistakes
A status of 1 is not an error: `diff` uses it for "different", so a script with `set -e` stops on files that differ, and `diff f1 f2 || true` is needed. Without `-u`, the output is in the old normal format, which `patch` can also read but which has no context. The first two lines of `-u` output contain file times, so they differ between runs: strip them when you compare outputs. `diff -r` also lists files that exist in only one folder. Binary files are reported as "differ" without a line diff.

### Portability
Only the unified format and the options above are in the Common examples. The default format, `-y` (side by side), `-c` and `-e` are shown in the GNU and Apple sections; BusyBox 1.37.0 supports a smaller set. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU diff
GNU diff is part of GNU diffutils, and this page's checks ran version 3.12. It has the usual three output formats and many more options.

### GNU examples
Print in the default normal format, which shows the line numbers and the lines:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && diff f1 f2; test $? -eq 1 && echo differ
```

### Option: -y, --side-by-side
Show the two files in two columns.

## Apple
### About macOS diff
This is the diff in the system folders of macOS. It is Apple's version, based on FreeBSD diff, and accepts the options above.

### Apple examples
Print in the default normal format, which shows the line numbers and the lines:

```sh
printf 'a\nb\n' > f1 && printf 'a\nc\n' > f2 && diff f1 f2; test $? -eq 1 && echo differ
```

### Option: -y
Show the two files in two columns.

## BusyBox
### About BusyBox diff
BusyBox bundles diff as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `diff` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
