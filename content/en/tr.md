---
schema_version: 1
id: tr
title: tr
summary: Translate, squeeze, or delete characters from standard input.
aliases: [gtr]
tags: [text, pipes, characters]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/tr-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=tr&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`tr` reads standard input and writes it back after replacing, deleting, or squeezing characters. It works on single characters, not words or patterns, and it never takes a file name.

### Examples
Change lowercase letters to uppercase:

```sh
echo 'hello' | tr 'a-z' 'A-Z'
```

Do the same with character classes:

```sh
echo 'hello' | tr '[:lower:]' '[:upper:]'
```

Squeeze repeated spaces into one:

```sh
echo 'a   b    c' | tr -s ' '
```

Delete every letter o:

```sh
echo 'hello world' | tr -d 'o'
```

Keep only digits by deleting everything else (the newline is deleted too):

```sh
echo 'abc123' | tr -cd '0-9'
```

Turn colons into newlines:

```sh
echo 'a:b:c' | tr ':' '\n'
```

Turn the commas of a file into tabs by redirecting it in:

```sh
tr ',' '\t' < data.csv
```

### Option: -d
Delete the characters in the set instead of translating them.

### Option: -s
Squeeze runs of a repeated character from the last set into a single one.

### Option: -c
Use the complement of the first set, meaning every character not listed.

### Common mistakes
`tr` reads standard input only, so use `tr ... < file` or a pipe, not `tr ... file`. Sets are characters, not words: `tr 'cat' 'dog'` maps each letter separately. Keep both sets the same length when translating, because behavior for unequal sets differs between implementations. Quote sets so the shell does not expand brackets or ranges.

### Portability
`-d`, `-s`, and `-c` with ranges and the classes above are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). Multibyte-character behavior depends on the locale and implementation.

## GNU
### GNU notes
GNU `tr` accepts the ranges, classes, and escapes shown here. Confirm anything beyond them, such as equivalence classes, on other implementations.

## Apple
### Apple notes
macOS `tr` is the BSD implementation. Escapes such as `'\n'` and classes such as `[:lower:]` work as shown here; behavior with multibyte text is locale dependent and not covered.

## BusyBox
### BusyBox notes
BusyBox `tr` supports ranges, classes, and `-c`, `-d`, `-s` as shown; other features depend on the build configuration.
