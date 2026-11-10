---
schema_version: 1
id: git
title: git
summary: Track changes to files in a project, and share them.
aliases: []
tags: [git, version-control, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://git-scm.com/docs
  - https://git-scm.com/book/en/v2
---
## Common
### Overview
`git` records the history of a folder: every commit is a saved snapshot with a message, and you can go back, compare, branch off and merge. A folder with a history is a repository, kept in a hidden `.git` folder. Git works on your own computer without a network, and shares commits with other repositories (`clone`, `fetch`, `push`) when you ask. The examples create a repository in a temporary folder.

### Background
`git` is an open-source tool whose documentation is at git-scm.com. It has many subcommands; this page covers the everyday ones and their output. The tool tested here is git 2.56.0.

### Examples
Create a repository named `main`, and see that nothing is changed yet; the `-q` flags keep it quiet:

```sh
git init -q -b main && git status --short | wc -l | tr -d ' '
```

Stage a new file with `git add`; the status shows `A` for added:

```sh
git init -q -b main && printf hi > a.txt && git add a.txt && git status --short
```

Commit it with a message, and read the history; the name and address are given with `-c` so no setup is needed:

```sh
git init -q -b main && printf hi > a.txt && git add a.txt && git -c user.name=Demo -c user.email=demo@example.com commit -q -m first && git log --format=%s
```

See what changed since the last commit, as a summary:

```sh
git init -q -b main && printf 'a\n' > f && git add f && git -c user.name=D -c user.email=d@e.c commit -q -m one && printf 'b\n' >> f && git diff --stat | tail -n 1
```

Make a branch and switch to it, then ask which branch you are on:

```sh
git init -q -b main && git -c user.name=D -c user.email=d@e.c commit -q --allow-empty -m one && git switch -q -c feature && git branch --show-current
```

List several commit messages on one line, newest first:

```sh
git init -q -b main && git -c user.name=D -c user.email=d@e.c commit -q --allow-empty -m one && git -c user.name=D -c user.email=d@e.c commit -q --allow-empty -m two && git log --format=%s | paste -sd, -
```

Mark a commit with a tag, and list the tags:

```sh
git init -q -b main && git -c user.name=D -c user.email=d@e.c commit -q --allow-empty -m one && git tag v1 && git tag
```

Ask whether you are inside a repository, which is useful in scripts:

```sh
git init -q -b main && git rev-parse --is-inside-work-tree
```

### Option: -C
Run as if started in this folder, as in `git -C path status`.

### Option: -c
Set a configuration value for this one command, as in `-c user.name=Demo`.

### Common mistakes
A commit needs a name and an email; if none is set, `git commit` stops and tells you to set them with `git config`. `git add` only stages: nothing is saved in history until `git commit`. Commands that rewrite history, such as `git reset --hard`, `git clean -fd` and `git push --force`, can destroy work and cannot always be undone; read what they will do first. Never commit passwords or keys: they stay in the history even after the file is removed. Old tutorials use `master` as the first branch name; the name is a setting, so say it with `git init -b main`.

### Portability
`git` has one upstream, and the subcommands in this page behave the same on every system, but newer subcommands such as `git switch` and `git restore` need a recent version (git 2.23 or later). The examples were executed on git 2.56.0, in a temporary folder with a clean home folder, so your own settings do not change them. The Common examples were executed on git 2.56.0 (see the development guide for the exact check).

## Upstream
### About upstream git
This is the upstream `git` command. It is one program everywhere; what differs between systems is its version and the settings in `~/.gitconfig`.
