---
schema_version: 1
id: gitleaks
title: gitleaks
summary: Scan code and git history for passwords, keys and other secrets.
aliases: []
tags: [security, git, developer]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/gitleaks/gitleaks#readme
---
## Common
### Overview
`gitleaks` looks through files, or the whole history of a git repository, for text that looks like a secret: cloud keys, tokens, private keys and passwords. It prints each finding with the file and line. It is used before a commit or in a continuous-integration run, so that a secret never reaches a public repository.

### Background
`gitleaks` is an open-source tool whose manual is in its repository. It works with a list of rules that match known secret formats, which you can extend. The tool tested here is gitleaks 8.30.1.

### Examples
Print the version:

```sh
gitleaks version
```

Scan a folder that is not a git repository, with `--no-git`; here it holds a made-up token in the format of a GitHub token (it is not a real one, and it is written in two pieces so that scanners do not flag this page), and the status is 1 when something is found:

```sh
printf 'token = %s%s\n' ghp_ R8kT3vQ9mZ2xW5nB7yL4cH6jD1fG0sAe9Pw3 > secret.txt && gitleaks detect --no-git --no-banner -s . > /dev/null 2>&1; echo "status $?"
```

Scan a folder with nothing in it that looks like a secret; the status is 0:

```sh
printf 'hello\n' > note.txt && gitleaks detect --no-git --no-banner -s . > /dev/null 2>&1; echo "status $?"
```

### Option: --no-git
Scan files as they are, not the git history.

### Option: --redact
Hide the secrets in the output.

### Option: -s
The folder or repository to scan.

### Common mistakes
Well-known example keys, such as the key AWS publishes in its documentation, are ignored by default, so a test with them finds nothing. A finding means that a secret may have been exposed: remove it from the code, and then change the secret itself, because it stays in the git history and in every clone. Scanning only the working tree misses secrets in old commits; scan the history (without `--no-git`) too. False findings happen: mark them with an allow list in `.gitleaks.toml` or a `gitleaksignore` file, and not by turning the tool off. The output can contain the secret in plain text; add `--redact` before sharing it.

### Portability
The subcommands and rules are the same on every system, and the command line changed between major versions (older versions use `gitleaks detect`, newer ones also have `gitleaks git` and `gitleaks dir`). The examples were executed on gitleaks 8.30.1. The Common examples were executed on gitleaks 8.30.1 (see the development guide for the exact check).

## Upstream
### About upstream gitleaks
This is the upstream `gitleaks` tool. It is one program everywhere; what differs is the version and its rules.
