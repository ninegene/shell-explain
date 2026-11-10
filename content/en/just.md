---
schema_version: 1
id: just
title: just
summary: Run the commands of a project from a file named justfile.
aliases: []
tags: [developer, build, scripts]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://just.systems/man/en/
---
## Common
### Overview
`just` runs named recipes written in a file called `justfile`, much like `make` but without its rules about files and dates. A recipe is a name followed by an indented list of commands. `just` is used to keep a project's everyday commands (build, test, release) in one readable place that the whole team shares.

### Background
`just` is an open-source command runner whose manual is at just.systems. It is a separate program from `make`, with a simpler syntax, and is installed on request. The tool tested here is the just on the host.

### Examples
Run a recipe by its name; the command line is printed first, then its output:

```sh
printf 'hi:\n\techo hello\n' > justfile && just hi 2>&1
```

List the recipes in the file:

```sh
printf 'hi:\n\techo hello\n' > justfile && just --list
```

Show what a recipe would run without running it, with `--dry-run`:

```sh
printf 'hi:\n\techo hello\n' > justfile && just --dry-run hi 2>&1
```

Pass an argument to a recipe by naming a parameter:

```sh
printf 'greet name:\n\techo hello {{name}}\n' > justfile && just greet world 2>&1
```

### Option: --list
List the recipes.

### Option: --dry-run
Print the commands without running them.

### Option: --justfile
Use this file instead of the default `justfile`.

### Common mistakes
The lines of a recipe must be indented, with the same tab or spaces on every line, and `just` complains about a mix. Each line of a recipe runs in its own shell, so `cd` on one line does not apply to the next unless the recipe is written as one script. `just` prints each command before it runs it, to standard error; add `@` in front of a line to hide it. Run it in the folder with the `justfile`, or in a folder below it: `just` finds the file in the parent folders.

### Portability
The syntax is the same on every system, and new features come with new versions. Recipes run with `sh` by default, so what they run depends on the system. The examples were executed on the just of the host. The Common examples were executed on just (the version on the host) (see the development guide for the exact check).

## Upstream
### About upstream just
This is the upstream `just` tool. It is one program everywhere; what differs is the version.
