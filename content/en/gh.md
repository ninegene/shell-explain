---
schema_version: 1
id: gh
title: gh
summary: Work with GitHub issues, pull requests and repositories from the terminal.
aliases: []
tags: [github, git, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://cli.github.com/manual/
---
## Common
### Overview
`gh` is the GitHub command-line tool. It creates and reviews pull requests, opens and comments on issues, manages repositories and releases, and can call any GitHub API endpoint with `gh api`. It uses your GitHub login (set up with `gh auth login`), so every call acts as you on GitHub. For plain git work such as commit and push, you still use `git`.

### Background
`gh` is the official GitHub CLI, with its manual at cli.github.com. The tool tested here is gh 2.103.2. The examples that run here read only its own help and report that no login exists, in a clean home folder; the ones that call GitHub are reviewed only, because they use your login.

### Examples
Print the version; the output begins with `gh version`:

```sh
gh --version | head -n 1 | cut -c1-10
```

Print the first line of the help text:

```sh
gh help | head -n 1
```

Ask for the login state; with no login in a clean home folder, the status is 1:

```sh
gh auth status > /dev/null 2>&1; echo "status $?"
```

Read the help of one command group:

```sh
gh pr --help | head -n 1
```

List the open pull requests of the current repository; reviewed only because it calls GitHub as you:

```sh
gh pr list
```

Create a pull request from the current branch, with a title; reviewed only:

```sh
gh pr create --title "Fix the typo" --body "Details here."
```

Call an API endpoint and pick one field with `--jq`; reviewed only:

```sh
gh api repos/octocat/Hello-World --jq .full_name
```

### Subcommands
`gh auth` logs in and shows the login state. `gh pr` works with pull requests. `gh issue` works with issues. `gh repo` works with repositories. `gh api` calls an API endpoint.

### Common mistakes
Everything `gh` does is done as you: `gh pr merge`, `gh repo delete` and `gh issue close` change the real repository, and some cannot be undone. Check the repository it will use (`gh repo view`) before a command that changes anything, because it picks the one of the current folder. A login token is stored in your keychain or a config file: do not share it, and prefer a token with only the scopes needed. In a script, set `GH_TOKEN` to a token with limited rights, not your own login.

### Portability
The commands are the same on every system, and new ones arrive with new versions. The first four examples were executed on gh 2.103.2 in a clean home folder. The Common examples were executed on gh 2.103.2 (see the development guide for the exact check).

## Upstream
### About upstream gh
This is the upstream GitHub CLI. It is one program everywhere; what differs is the version.
