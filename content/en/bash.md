---
schema_version: 1
id: bash
title: bash
summary: Start a Bash shell, run a script, or run one line of Bash.
aliases: []
tags: [shell, scripts, bash]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/bash/manual/bash.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/bash.1.html
---
## Common
### Overview
This page is about the `bash` program: the command that starts the shell, runs a script or runs one line with `-c`. The commands Bash has built in, such as `cd` and `echo`, have their own pages. Bash is the usual shell on Linux and was the login shell on older macOS. A script that begins with `#!/bin/bash` is run by this program.

### Background
`bash` is the GNU shell, and its manual describes the options of the program and the language. The tool tested here is the Bash of the host, 5.2. macOS ships Bash 3.2, which is much older and lacks many features, so a script that must run on both should avoid newer Bash features.

### Examples
Run one line of Bash with `-c`; the result of arithmetic is printed:

```sh
bash -c 'echo $((1+2))'
```

Give a default for an unset variable with `${name:-value}`:

```sh
bash -c 'echo ${y:-no}'
```

The first argument after the code becomes `$0`:

```sh
bash -c 'echo "$0"' name
```

Run commands read from standard input:

```sh
echo 'echo piped' | bash
```

Check the syntax of a script without running it, with `-n`:

```sh
printf 'echo hi\n' > s.sh && bash -n s.sh && echo syntax-ok
```

Use a feature of Bash that `sh` lacks, the `[[ ]]` test with a pattern:

```sh
bash -c 'if [[ abc == a* ]]; then echo match; fi'
```

### Option: -c
Run this string as a command line.

### Option: -n
Read the script and check its syntax, without running it.

### Option: -x
Print each command before it runs, which helps with debugging.

### Option: -e
Stop at the first command that fails.

### Common mistakes
A script run with `sh` is not run by Bash, even if you wrote Bash features in it: use `bash script.sh` or a `#!/bin/bash` line. On macOS the `bash` in `/bin` is version 3.2, and many Bash 4 and 5 features (such as associative arrays) fail there. Code in `-c` needs shell quotes. `bash -n` checks only syntax; use `shellcheck` for more. A Bash script is not portable to every shell, so say in the first line which shell it needs.

### Portability
Only options that work in Bash 3.2 and 5.x are in the examples. The Common behavior of the shell language is described in the manual, and features differ by version. The examples were executed on Bash 5.2.37. The Common examples were executed on Bash 5.2.37 (see the development guide for the exact check).

## Upstream
### About upstream bash
This is the upstream GNU Bash program. The version differs between systems, most of all between macOS (3.2) and Linux (5.x).
