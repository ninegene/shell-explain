---
schema_version: 1
id: lefthook
title: lefthook
summary: Run checks automatically when you commit or push with git.
aliases: []
tags: [git, developer, hooks]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://lefthook.dev/
---
## Common
### Overview
`lefthook` manages git hooks: small programs that git runs at certain moments, such as before a commit. You describe the checks in `lefthook.yml` (a formatter, a linter, a test), and `lefthook install` puts the hooks into the repository, so that everyone who clones it runs the same checks. It can run commands in parallel and only on the files that changed.

### Background
`lefthook` is an open-source tool whose documentation is at lefthook.dev. It is a single program without runtime dependencies. The tool tested here is lefthook 2.2.2.

### Examples
Print the version:

```sh
lefthook version
```

Write a configuration with one command for `pre-commit`, then install the hooks, and check that git now has the hook file:

```sh
git init -q -b main && printf 'pre-commit:\n  commands:\n    hello:\n      run: echo hello\n' > lefthook.yml && lefthook install > /dev/null 2>&1 && test -f .git/hooks/pre-commit && echo installed
```

Run the hook by hand, without committing; the output names the command:

```sh
git init -q -b main && printf 'pre-commit:\n  commands:\n    hello:\n      run: echo hello\n' > lefthook.yml && lefthook install > /dev/null 2>&1 && lefthook run pre-commit 2>&1 | grep -c hello
```

### Subcommands
`lefthook install` writes the hooks into the repository. `lefthook run NAME` runs a hook by hand. `lefthook uninstall` removes the hooks.

### Common mistakes
Run `lefthook install` again after you change the list of hooks, or git keeps running the old ones. A failing command in a hook stops the commit; `git commit --no-verify` skips the hooks, which you should do only with good reason. Hooks are not copied by `git clone`: each person must run `lefthook install`, so many projects do this in a setup step. Commands run from the repository root and need their programs on the `PATH` of whoever commits.

### Portability
The configuration file is the same on every system, and details change with versions. The examples were executed on lefthook 2.2.2. The Common examples were executed on lefthook 2.2.2 (see the development guide for the exact check).

## Upstream
### About upstream lefthook
This is the upstream `lefthook` tool. It is one program everywhere; what differs is the version.
