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

### Option: -name
Match the filename against a pattern. Quote a wildcard pattern so find, rather than your shell, expands it.

### Option: -type
Restrict results by type. Use `f` for regular files and `d` for directories.

### Common mistakes
Always include a starting directory for portable commands. `find . -name '*.log'` differs from leaving the wildcard unquoted. Filenames can contain spaces or newlines, so do not split find output on whitespace to pass it to another tool.

### Portability
Actions that delete files or run other programs need special care and are not used here. Features such as traversal limits, time units, and regex matching differ between implementations and versions.

## GNU
### GNU examples
Look in the current directory and one level below it:

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
