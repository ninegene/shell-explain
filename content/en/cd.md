---
schema_version: 1
id: cd
title: cd
summary: Change the current directory of your interactive shell.
aliases: []
tags: [shell, navigation, directories]
implementations: [bash]
references:
  - https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html
---
## Common
### Overview
`cd` is a shell builtin because changing directory must affect the current shell. A separate executable cannot change its parent shell's directory. This reference documents Bash behavior; another shell can differ.

### Examples
Move into a project directory:

```sh
cd projects
```

Move to the parent directory:

```sh
cd ..
```

Return to your home directory:

```sh
cd
```

### Common mistakes
Quote directory names containing spaces: `cd 'project notes'`. Your current directory changes how relative paths are interpreted. Confirm your location with `pwd` before running a command that modifies files.

### Portability
A standalone explain binary cannot inspect your parent shell's aliases, functions, or exact builtin implementation. Choose Bash explicitly when needed; a detected SHELL environment variable is not proof of the active shell.

## Bash
### Bash examples
Return to the previous directory:

```sh
cd -
```

Bash uses OLDPWD for this shortcut and prints the new directory.

### Option: -P
Resolve symbolic links physically rather than retaining their logical path. `pwd -P` can help compare the physical location.

### Option: -L
Follow logical directory paths, preserving symbolic-link components when processing the path. This is Bash's default behavior.
