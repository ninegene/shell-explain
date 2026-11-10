---
schema_version: 1
id: expr
title: expr
summary: "Evaluate an expression: arithmetic, comparisons and text matching."
aliases: [gexpr]
tags: [shell, scripts, math]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/expr-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=expr&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`expr` evaluates the words you give it as an expression and prints the result. Each operator and number is a separate argument, so there must be spaces between them. It was the way to do arithmetic in old shell scripts, before `$(( ))` existed, and it is still seen in them. Its exit status is 1 when the result is 0 or empty.

### Background
An `expr` utility first appeared in the Programmer's Workbench (PWB/UNIX), according to the macOS manual, which also says it conforms to POSIX.1-2008. Today the name covers GNU expr, the expr that macOS ships, and the expr applet in BusyBox. They agree on arithmetic, comparisons and the `:` pattern match.

### Examples
Add two numbers; put spaces around the operator:

```sh
expr 2 + 3
```

Divide with whole numbers, which drops the remainder:

```sh
expr 7 / 2
```

Multiply; the star must be quoted or escaped so the shell does not expand it:

```sh
expr 3 \* 4
```

Compare two numbers; the result is `1` for true and `0` for false:

```sh
expr 3 \< 5
```

Match text with a regular expression; the result is the number of matching characters:

```sh
expr abc : 'a.*'
```

Use the exit status; it is 1 when the result is 0, so `||` runs:

```sh
expr 1 - 1 || echo zero
```

### Common mistakes
Every token needs its own argument: `expr 2+3` prints `2+3`, not 5. Characters the shell treats specially, such as `*`, `<`, `>`, `(` and `|`, must be quoted. A result of 0 or an empty string gives exit status 1, which can stop a script that uses `set -e`. In a modern shell script, `$((2 + 3))` is simpler and needs no process.

### Portability
Only arithmetic, comparison and `:` matching are in the Common examples. GNU adds `length`, `substr`, `index` and `match` keywords, which BSD `expr` may not have. In a shell script, `$(( ))` is POSIX and works the same everywhere. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU expr
GNU expr is part of GNU coreutils, and this page's checks ran coreutils 9.11. It has the keywords `length`, `substr`, `index` and `match`.

### GNU examples
Count the characters in a word with `length`:

```sh
expr length hello
```

## Apple
### About macOS expr
This is the expr in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `expr` follows POSIX arithmetic and `:` matching. Prefer `$(( ))` or `${#var}` for length.

## BusyBox
### About BusyBox expr
BusyBox bundles expr as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `expr` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
