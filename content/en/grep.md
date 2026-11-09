---
schema_version: 1
id: grep
title: grep
summary: Find lines of text that match a pattern.
aliases: [ggrep]
tags: [text, search, logs]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/grep/manual/
  - https://man.freebsd.org/cgi/man.cgi?query=grep&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
Use `grep` to select lines from a file or from another command's output. The pattern comes before the input filename. Start with a quoted literal word; regular expressions add more powerful matching rules.

For the full GNU reference, read the [GNU grep manual](https://www.gnu.org/software/grep/manual/).

### Background
The name comes from the `ed` editor command `g/re/p`: globally search for a regular expression and print the matching lines. `grep` first appeared in Version 6 Unix, according to the FreeBSD manual. Today the name covers several separate programs that share that interface but not their source code: GNU grep, the BSD-derived grep that macOS ships, and a small reimplementation inside BusyBox. They agree on the common options below and differ on extensions such as the Perl-style `-P` patterns, regular-expression engines, and recursive-search details. That is why the sections that follow are split by implementation.

### Examples
Find lines containing ERROR in a log:

```sh
grep 'ERROR' app.log
```

Show matching line numbers and ignore letter case:

```sh
grep -n -i 'error' app.log
```

Match a literal string rather than treating punctuation as a regular expression:

```sh
grep -F 'example.org' app.log
```

Show the lines that do not match:

```sh
grep -v 'ERROR' app.log
```

Count matching lines instead of printing them:

```sh
grep -c 'ERROR' app.log
```

List the names of files that contain a match:

```sh
grep -l 'ERROR' app.log logs/server.log
```

Match whole words only, so `error` does not match inside a longer word:

```sh
grep -w 'error' app.log
```

Match either of two patterns with an extended regular expression:

```sh
grep -E 'ERROR|INFO' app.log
```

Combine several patterns without regular-expression syntax by repeating `-e`:

```sh
grep -e 'INFO' -e 'example' app.log
```

Show one line of context before and after each match:

```sh
grep -C 1 'ERROR' app.log
```

Search every file under a directory:

```sh
grep -r 'ERROR' logs
```

### Option: -n
Print each matching line with its line number. This helps you locate a result in the original file.

### Option: -i
Ignore letter case when matching. For predictable examples, use an ASCII fixture and the C locale.

### Option: -F
Treat the pattern as fixed text rather than a regular expression.

### Option: -v
Select lines that do not match the pattern.

### Option: -c
Print only the number of matching lines for each file, instead of the lines. With several files, each count is prefixed by the filename.

### Option: -l
Print only the names of files that contain at least one match, once each.

### Option: -w
Require the match to be a whole word: it must not be preceded or followed by a letter, digit, or underscore.

### Option: -E
Use extended regular expressions, so `|`, `+`, `?`, and `()` work without backslashes. Quote the pattern.

### Option: -e
Give a pattern explicitly. Repeat it to match any of several patterns, or use it for a pattern that starts with a dash.

### Option: -C
Print the given number of lines of context before and after each match. The value is the next argument, as in `-C 1`.

### Option: -r
Search directories recursively. Name the directory explicitly so the command is portable: without one, GNU grep and macOS grep search the current directory but print paths differently (`a.txt` versus `./a.txt`), and other builds are untested here.

### Common mistakes
Quote the pattern so the shell does not expand special characters before `grep` sees them. A pattern beginning with a dash can look like an option; use `grep -e '-value' file.txt`.

An exit status of 1 means no lines matched; it does not necessarily mean the command failed to read the file.

### Portability
The examples in Common use options that GNU grep, macOS grep, and BusyBox grep all accept; they were executed on GNU grep 3.12, macOS grep (BSD) 2.6.0-FreeBSD, and BusyBox 1.37.0 (see the development guide for the exact check). Optional regex engines and recursive-search details differ. Implementation detection does not guarantee every installed version supports every feature; inspect the installed manual for version-specific behavior.

## GNU
### About GNU grep
GNU grep is the GNU Project's grep. Homebrew installs it on macOS as `ggrep`, and as `grep` in its `gnubin` directory, which acts as GNU grep when placed first on PATH (checked with Homebrew grep 3.12). It accepts the `-P` Perl-style patterns, which macOS grep rejects. Its `--version` output names GNU grep.

### GNU examples
Search a directory recursively and print matching line numbers:

```sh
grep -rn --include='*.log' 'ERROR' logs/
```

`--include` restricts which filenames are searched. The shell must not expand its wildcard before grep receives it.

### Option: --line-number
GNU's long spelling for `-n`. Use `-n` for a shorter portable example.

## Apple
### About macOS grep
The `grep` in `/usr/bin` on macOS is Apple's build of the BSD-derived grep. It describes itself as BSD grep while accepting many GNU options, so its banner can look GNU-compatible even though it is not GNU grep. For example, macOS grep 2.6.0-FreeBSD rejects `-P` as an invalid option. A package-manager GNU grep earlier on PATH changes which one you run.

### Apple notes
The macOS system utility is usually `/usr/bin/grep`. GNU grep installed by a package manager may be named `ggrep` or may come first on PATH. Use the detected executable path to distinguish them; the operating system alone is not enough.

## BusyBox
### About BusyBox grep
BusyBox bundles a compact `grep` as one of its applets, so it is what you get on systems built around BusyBox, such as the Alpine Linux guest this project's Linux check uses. Which options exist depends on how BusyBox was configured when it was built, so check `grep --help` in the actual environment.

### BusyBox notes
BusyBox applets can vary with build configuration. Prefer the common short options and consult `grep --help` in the actual environment before using optional features.
