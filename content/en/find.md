---
schema_version: 1
id: find
title: find
summary: Locate files and directories using explicit matching rules.
aliases: [gfind]
tags: [files, search, directories]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/findutils/manual/
  - https://man.freebsd.org/cgi/man.cgi?query=find&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
Give `find` a starting directory followed by tests. A dot means the current directory. Combining tests lets you select a file by both its type and its name.

### Examples
Find regular files whose names end in .log:

```sh
find . -type f -name '*.log'
```

Find directories named cache:

```sh
find . -type d -name 'cache'
```

Match a name regardless of letter case:

```sh
find . -type f -iname '*.LOG'
```

Exclude files that match a name:

```sh
find . -type f ! -name '*.log'
```

Match either of two names; the parentheses group the alternatives and must be escaped from the shell:

```sh
find . -type f \( -name '*.log' -o -name '*.txt' \)
```

Find files modified within the last 7 days:

```sh
find . -type f -mtime -7
```

List every directory below the starting point:

```sh
find . -type d
```

Run a read-only command on each match, passing many files per invocation:

```sh
find . -type f -name '*.log' -exec grep -l 'ERROR' {} +
```

### Option: -name
Match the filename against a pattern. Quote a wildcard pattern so find, rather than your shell, expands it.

### Option: -type
Restrict results by type. Use `f` for regular files and `d` for directories.

### Option: -iname
Like `-name`, but ignores letter case.

### Option: -o
Match when either side is true. `-o` binds more loosely than the implied "and" between tests, so group alternatives with escaped parentheses.

### Option: -mtime
Match by modification time in whole days. `-mtime -7` means less than 7 days ago, `-mtime +7` means more than 7 days ago.

### Option: -exec
Run a command for each match. Ending the command with `{} +` passes many files at once; ending with `\;` runs it once per file. Try a read-only command first, because `-exec` can run anything, including deletes.

### Common mistakes
Always include a starting directory for portable commands. `find . -name '*.log'` differs from leaving the wildcard unquoted. Filenames can contain spaces or newlines, so do not split find output on whitespace to pass it to another tool.

### Portability
Actions that delete files need special care and are not used here; the one `-exec` example runs a read-only command. The Common examples were executed on GNU findutils 4.10.0, macOS find, and BusyBox 1.37.0 (see the development guide for the exact check). Features such as traversal limits, time units, and regex matching differ between implementations and versions.

## GNU
### GNU examples
Find matching files up to two levels below the starting directory:

```sh
find . -maxdepth 2 -type f -name '*.log'
```

### Option: -maxdepth
Limit traversal depth; the starting directory is depth 0. Put this traversal option before matching tests. This is an extension, not a universal POSIX option.

## Apple
### Apple notes
Use an explicit starting path with macOS find. The common examples avoid GNU-only assumptions. Check `man find` on the installed macOS release before adopting an extended expression.

## BusyBox
### BusyBox notes
The available tests/actions depend on BusyBox's build configuration. The common name/type examples provide a small baseline; extended predicates need confirmation from the installed help.
