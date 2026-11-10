---
schema_version: 1
id: claude
title: claude
summary: Work with Claude, an AI assistant, from the terminal.
aliases: []
tags: [ai, developer, assistant]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://docs.claude.com/en/docs/claude-code/overview
---
## Common
### Overview
`claude` starts Claude Code, an AI coding assistant that works in your terminal. You describe what you want in plain language, and it reads files, proposes changes, runs commands (with your permission) and explains what it did. Run `claude` in a project folder for an interactive session, or `claude -p "question"` to ask one question and print the answer.

### Background
`claude` is the command of Claude Code, documented on docs.claude.com. The tool tested here is the version 2.1.310. The examples that ask a question or start a session use your account and send text to a service, so they are reviewed only. The executed examples read the tool's own version and help.

### Examples
Print the version; the output ends with the product name:

```sh
claude --version | cut -d' ' -f2-
```

Print the first line of the help text:

```sh
claude --help | head -n 1
```

Ask one question and print the answer, without an interactive session, with `-p`; reviewed only because it sends text to a service as you:

```sh
claude -p "Explain what this script does" < script.sh
```

Start an interactive session in the current project; reviewed only:

```sh
claude
```

Continue the most recent session in this folder; reviewed only:

```sh
claude --continue
```

### Option: -p
Print mode: answer one prompt and exit.

### Option: --continue
Continue the latest session in this folder.

### Option: --version
Print the version.

### Common mistakes
Text you give `claude`, including files it reads, is sent to a service, so keep secrets out of what it can read. It can change files and run commands: read what it proposes before you approve, and use version control so that you can undo changes. Permission modes decide what it may do without asking; start with the strictest one and relax it when you trust the task. A session costs money or uses a quota. The output is not always right, so test its changes.

### Portability
The options above work the same on every system, and the tool changes quickly, so check `claude --help` for your version. The first two examples were executed on version 2.1.310. The Common examples were executed on Claude Code 2.1.310 (see the development guide for the exact check).

## Upstream
### About upstream claude
This is the upstream `claude` command of Claude Code. It is one program everywhere; the version changes often.
