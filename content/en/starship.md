---
schema_version: 1
id: starship
title: starship
summary: Show a fast, informative prompt in any shell.
aliases: []
tags: [shell, prompt, productivity]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://starship.rs/config/
---
## Common
### Overview
`starship` draws your shell prompt. It shows what matters at the moment: the folder, the git branch and its state, the version of the language of the project, how long the last command took and whether it failed. You turn it on by putting `eval "$(starship init bash)"` (or the line for your shell) in the shell's startup file, and you change it in `~/.config/starship.toml`.

### Background
`starship` is an open-source tool whose manual is at starship.rs. It is one program for every shell. The tool tested here is starship 1.27.0.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
starship --version
```

Print the setup code for Bash; reviewed only, because the code changes with each release:

```sh
starship init bash
```

### Subcommands
`starship init SHELL` prints the code for your shell. `starship prompt` prints the prompt text once. `starship config` changes a setting. `starship timings` shows how long each module took.

### Common mistakes
The prompt only changes after the line in the startup file runs, so start a new shell after you add it. Many modules look at the folder, so a slow disk or a large git repository can slow the prompt: use `starship timings` to see which module is slow. The configuration is a TOML file, and a mistake in it shows a warning and the default for that part. Fonts: some symbols need a patched font (a Nerd Font) to show correctly.

### Portability
The setup code differs per shell (`bash`, `zsh`, `fish` and others) but the program and the configuration are the same everywhere. The examples are reviewed only; the page was written against starship 1.27.0.

## Upstream
### About upstream starship
This is the upstream `starship` tool. It is one program everywhere; what differs is the version.
