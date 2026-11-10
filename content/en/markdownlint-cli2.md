---
schema_version: 1
id: markdownlint-cli2
title: markdownlint-cli2
summary: Check Markdown files for style and structure problems.
aliases: []
tags: [developer, markdown, lint]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/DavidAnson/markdownlint-cli2#readme
  - https://github.com/DavidAnson/markdownlint/blob/main/doc/Rules.md
---
## Common
### Overview
`markdownlint-cli2` checks Markdown files against a list of rules, each with a code such as `MD012`, and prints a line for every problem with the file, the line number and the rule name. It reads a configuration file to turn rules on or off, and it can fix many problems with `--fix`. It needs Node.js.

### Background
`markdownlint-cli2` is an open-source command-line tool for the markdownlint library, which defines the rules. The tool tested here is markdownlint-cli2 0.24.0 with markdownlint 0.42.0.

### Examples
Check a file that has too many blank lines; the lines that name the rule are counted:

```sh
printf '# Title\n\n\n\ntext\n' > a.md && { markdownlint-cli2 a.md 2>&1 | grep -c MD012; }
```

Show the first problem found, which names the file, the line, the rule and what is wrong:

```sh
printf '# Title\n\n\n\ntext\n' > a.md && { markdownlint-cli2 a.md 2>&1 | grep -m1 MD012 | cut -d' ' -f1,3; }
```

Fix what can be fixed automatically with `--fix`, then check again:

```sh
printf '# Title\n\n\n\ntext\n' > a.md && markdownlint-cli2 --fix a.md > /dev/null 2>&1; markdownlint-cli2 a.md > /dev/null 2>&1 && echo clean
```

### Option: --fix
Fix the problems that can be fixed, in place.

### Option: --config
Use this configuration file.

### Common mistakes
The exit status is 1 when problems are found, so a script with `set -e` stops. A rule you disagree with should be turned off in `.markdownlint-cli2.jsonc` or in the file, with a comment such as `<!-- markdownlint-disable MD013 -->`, and not ignored by habit. The first lines of output name the tool and its version, which changes between versions. `--fix` rewrites the files: use version control. Globs need quotes, so that the tool and not the shell expands them.

### Portability
The rules are the same on every system, and new rules and changed defaults arrive with new versions. The examples were executed on markdownlint-cli2 0.24.0, running on the host's Node.js. The Common examples were executed on markdownlint-cli2 0.24.0 (see the development guide for the exact check).

## Upstream
### About upstream markdownlint-cli2
This is the upstream `markdownlint-cli2` tool, installed with npm. Its version is the one in your project or the global one.
