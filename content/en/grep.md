---
schema_version: 1
id: grep
title: grep
summary: Find lines of text that match a pattern.
aliases: [ggrep]
tags: [text, search, logs]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/grep/manual/
  - https://man.freebsd.org/cgi/man.cgi?query=grep&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
Use `grep` to select lines from a file or from another command's output. The pattern comes before the input filename. Start with a quoted literal word; regular expressions add more powerful matching rules.

For the full GNU reference, read the [GNU grep manual](https://www.gnu.org/software/grep/manual/).

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

### Option: -n
Print each matching line with its line number. This helps you locate a result in the original file.

### Option: -i
Ignore letter case when matching. For predictable examples, use an ASCII fixture and the C locale.

### Option: -F
Treat the pattern as fixed text rather than a regular expression.

### Option: -v
Select lines that do not match the pattern.

### Common mistakes
Quote the pattern so the shell does not expand special characters before `grep` sees them. A pattern beginning with a dash can look like an option; use `grep -e '-value' file.txt`.

An exit status of 1 means no lines matched; it does not necessarily mean the command failed to read the file.

### Portability
The examples in Common use options supported by the documented variants. Optional regex engines and recursive-search details differ. Implementation detection does not guarantee every installed version supports every feature; inspect the installed manual for version-specific behavior.

## GNU
### GNU examples
Search a directory recursively and print matching line numbers:

```sh
grep -rn --include='*.log' 'ERROR' logs/
```

`--include` restricts which filenames are searched. The shell must not expand its wildcard before grep receives it.

### Option: --line-number
GNU's long spelling for `-n`. Use `-n` for a shorter portable example.

## Apple
### Apple notes
The macOS system utility is usually `/usr/bin/grep`. GNU grep installed by a package manager may be named `ggrep` or may come first on PATH. Use the detected executable path to distinguish them; the operating system alone is not enough.

## BusyBox
### BusyBox notes
BusyBox applets can vary with build configuration. Prefer the common short options and consult `grep --help` in the actual environment before using optional features.
