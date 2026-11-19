---
schema_version: 1
id: shellcheck
title: shellcheck
summary: Find bugs and risky constructs in shell scripts.
aliases: []
tags: [shell, scripts, lint]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://www.shellcheck.net/wiki/
  - https://github.com/koalaman/shellcheck#readme
  - https://manpages.ubuntu.com/manpages/resolute/man1/shellcheck.1.html
---
## Common
### Overview
`shellcheck` reads a shell script and reports mistakes without running it: unquoted variables, broken tests, commands that behave differently in `sh` and `bash`, and many more. Each finding has a code such as `SC2086`, and the website shellcheck.net has a page for every code that explains the problem and the fix. It exits with status 0 when it finds nothing.

### Background
`shellcheck` is an open-source tool, and its wiki documents every check. It reads the `#!` line of a script to decide which shell it is written for. The tool tested here is ShellCheck 0.11.0.

### Examples
Check a script that has an unquoted variable; `-f gcc` prints one line per finding in a form that editors read, and `cut` keeps the message:

```sh
printf '#!/bin/sh\necho $1\n' > t.sh && shellcheck -f gcc t.sh | cut -d: -f4-
```

Check a clean script: nothing is printed and the status is 0:

```sh
printf '#!/bin/sh\necho "$1"\n' > ok.sh && shellcheck ok.sh && echo clean
```

Use the exit status: 1 means that something was found:

```sh
printf '#!/bin/sh\necho $1\n' > t.sh && { shellcheck -f gcc t.sh > /dev/null || echo findings; }
```

Report only errors with `-S error`, ignoring notes and warnings:

```sh
printf '#!/bin/sh\necho $1\n' > t.sh && shellcheck -S error t.sh && echo ok-at-error-level
```

### Option: -f
Choose the output format, such as `gcc`, `json` or `tty`.

### Option: -S
Report only findings at this severity or higher: `error`, `warning`, `info` or `style`.

### Option: -s
Say which shell the script is written for, such as `bash` or `sh`.

### Common mistakes
A script without a `#!` line makes `shellcheck` guess the shell, so add one or use `-s bash`. A finding is not always a bug, but each one is worth reading: the wiki page for the code says when it is safe to ignore it, with `# shellcheck disable=SC2086` above the line. `shellcheck` does not run the script and cannot see values that only exist at run time. It checks shell, and not programs that the script calls.

### Portability
The output format and the codes are the same on every system, but new codes appear in new versions, so a script that is clean in one version can have findings in the next. The examples were executed on ShellCheck 0.11.0. The Common examples were executed on ShellCheck 0.11.0 (see the development guide for the exact check).

## Upstream
### About upstream shellcheck
This is the upstream `shellcheck` tool. It is one program everywhere; what differs is the version.
