---
schema_version: 1
id: sd
title: sd
summary: Find and replace text with regular expressions, in a simpler way than sed.
aliases: []
tags: [text, replace, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/chmln/sd#readme
---
## Common
### Overview
`sd` replaces text. You give it a pattern to find and the text to put in its place, and it reads standard input or changes files in place. The pattern is a regular expression, and groups are written back as `$1`, `$2`. It is meant as an easier replacement for the common uses of `sed`, with no need for slashes or escaping.

### Background
`sd` is an open-source tool whose manual is in its repository. It is a separate program from `sed` and is installed on request. The tool tested here is sd v1.1.0.

### Examples
Replace a word read from standard input:

```sh
echo 'hello world' | sd world there
```

Swap two words with groups; the groups are written `$1` and `$2`:

```sh short="Swap two words with groups ($1, $2)"
echo 'hello world' | sd '(\w+) (\w+)' '$2 $1'
```

Treat the pattern as plain text with `-F`, so the dot is a dot:

```sh short="Treat the pattern as plain text (-F)"
printf 'a.b\n' | sd -F '.' '-'
```

Change a file in place; there is no backup:

```sh
printf 'aaa\n' > s.txt && sd a b s.txt && cat s.txt
```

### Option: -F
Treat the pattern as a fixed string, not a regular expression.

### Option: -p
Preview: print the changes without writing them to the file.

### Common mistakes
With a file argument, `sd` changes the file in place and keeps no backup, so try with `-p` first or use version control. The pattern is a regular expression: a dot matches any character unless you use `-F`. A dollar sign in the replacement starts a group name, so write `$$` for a literal one. `sd` does not use the `s/old/new/` form, so a command copied from a `sed` answer needs rewriting.

### Portability
`sd` is a different program from `sed`: its regular expressions follow the Rust `regex` syntax, which differs from the basic regular expressions of POSIX `sed` in details such as escaping. The examples were executed on sd v1.1.0. The Common examples were executed on sd v1.1.0 (see the development guide for the exact check).

## Upstream
### About upstream sd
This is the upstream `sd` tool. It is one program on every system.
