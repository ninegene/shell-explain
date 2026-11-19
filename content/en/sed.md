---
schema_version: 1
id: sed
title: sed
summary: Transform text one line at a time using editing rules.
aliases: [gsed]
tags: [text, replace, files]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/sed/manual/
  - https://man.freebsd.org/cgi/man.cgi?query=sed&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/sed.1.html
---
## Common
### Overview
Use `sed` to print transformed text without changing the original file. A substitution rule has the form `s/pattern/replacement/`. Preview the result before choosing an in-place edit.

### Background
The name is short for stream editor: `sed` edits text as it flows from a file or pipe instead of opening it in an editor. A `sed` command, written by L. E. McMahon, appeared in Version 7 Unix, according to the FreeBSD and macOS manuals. Today the name covers separate programs: GNU sed, the sed that macOS ships, and the compact sed inside BusyBox. They agree on the substitution and printing examples below and differ on in-place editing, extensions, and regular-expression details.

### Examples
Replace the first occurrence of old on each line:

```sh
sed 's/old/new/' notes.txt
```

Replace every occurrence on each line:

```sh
sed 's/old/new/g' notes.txt
```

Print only lines 2 through 4:

```sh
sed -n '2,4p' notes.txt
```

Delete the lines that match a pattern from the output:

```sh
sed '/second/d' notes.txt
```

Print only the lines that match a pattern:

```sh
sed -n '/third/p' notes.txt
```

Skip the first line, for example a header:

```sh
sed '1d' notes.txt
```

Print only the last line:

```sh
sed -n '$p' notes.txt
```

Apply two rules in order:

```sh
sed -e 's/old/new/' -e 's/second/2nd/' notes.txt
```

Use another delimiter when the pattern contains slashes:

```sh
echo '/usr/local/bin' | sed 's#/usr/local#/opt#'
```

Use an extended regular expression to replace runs of digits:

```sh
echo 'a1b22' | sed -E 's/[0-9]+/N/g'
```

### Option: -n
Suppress automatic printing. Pair it with a rule such as `p` to explicitly print selected lines.

### Option: -e
Add an editing rule. Repeat `-e` to apply multiple rules in order.

### Option: -E
Use extended regular expressions, so `+`, `?`, `|`, and `()` work without backslashes. GNU, macOS, and BusyBox sed all accept `-E`.

### Common mistakes
The substitution flag `g` changes every match on a line, not the number of files edited. Replacement text treats `&` specially: it expands to the matched text. Preview rules containing punctuation carefully.

### Portability
In-place editing is not uniform across implementations. The examples below create backup files rather than discarding the original. Version-sensitive extensions are outside the common examples. The Common examples were executed on GNU sed 4.10, macOS sed, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU sed
GNU sed is the GNU Project's sed. This page's checks ran GNU sed 4.10 on macOS (Homebrew) and 4.9 on Alpine Linux. It accepts `-i` with or without a backup suffix.

### GNU examples
Edit a disposable copy in place and retain the previous contents as notes.txt.bak:

```sh
sed -i.bak 's/old/new/g' notes.txt
```

### Option: -i, --in-place
Write changes back to the file. GNU accepts `-i` without a backup suffix, but use a backup while learning. An existing backup path can be overwritten; practice in a temporary directory.

## Apple
### About macOS sed
This is the sed in `/usr/bin` on macOS. Its manual carries the same Version 7 history note as the FreeBSD manual, but this page does not claim that it matches any particular FreeBSD version. The difference that matters most in practice is that `-i` always takes a backup-suffix argument, and an empty one means no backup.

### Apple examples
Retain a backup while editing a disposable file:

```sh
sed -i '.bak' 's/old/new/g' notes.txt
```

### Option: -i
Apple sed expects an argument for the backup suffix. The empty argument in `-i ''` means no backup; it is not the same syntax as GNU's usual suffix-free form. Use a backup until you have verified the edit.

## BusyBox
### About BusyBox sed
BusyBox bundles a compact sed as an applet. Which extensions it has depends on how BusyBox was built; this page's examples ran on BusyBox 1.37.0.

### BusyBox examples
Keep a backup while editing a disposable file:

```sh
sed -i.bak 's/old/new/g' notes.txt
```

### Option: -i
Edit in place with an optional attached backup suffix. BusyBox features depend on the build; check its local help before relying on an extension.
