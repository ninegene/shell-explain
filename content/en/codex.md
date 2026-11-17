---
schema_version: 1
id: codex
title: codex
summary: Work with an AI coding assistant from the terminal.
aliases: []
tags: [ai, developer, assistant]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://github.com/openai/codex#readme
---
## Common
### Overview
`codex` starts the Codex command-line assistant, an AI tool that reads and changes the code in a project folder at your request, and runs commands with your approval. Run `codex` in a project for an interactive session, or give it a task as an argument. The text of your task and the files it reads are sent to a service, so the examples that start a task are reviewed only.

### Background
`codex` is the command of the Codex CLI, with its documentation in its repository. The tool tested here is codex-cli 0.164.0. The version and help examples are reviewed only too, because their text changes with each release.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
codex --version
```

Read the help text; reviewed only, because it changes with each release:

```sh
codex --help
```

Give it a task and run it without an interactive session; reviewed only because it sends text to a service and can change files:

```sh
codex exec "Add a unit test for the parser"
```

Start an interactive session in the current project; reviewed only:

```sh
codex
```

### Subcommands
`codex exec` runs one task without the interactive screen. `codex` alone starts an interactive session. `codex --help` lists the other commands and options.

### Common mistakes
Everything `codex` reads can leave your machine, so keep secrets out of the folder it works in. It can change files and run commands: read what it proposes before you approve, work in a version-controlled folder, and use its sandbox and approval settings to limit what it may do. Do not run it with approvals switched off in a folder that has anything you cannot lose. The answers can be wrong; test the changes. A session uses your account and may cost money or use a quota.

### Portability
The tool changes quickly; check `codex --help` for your version. The Common examples are reviewed only.

## Upstream
### About upstream codex
This is the upstream `codex` command. It is one program everywhere; the version changes often.
