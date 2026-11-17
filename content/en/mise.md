---
schema_version: 1
id: mise
title: mise
summary: Install and switch between versions of developer tools, and set project environments.
aliases: []
tags: [developer, tools, environment]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://mise.jdx.dev/
---
## Common
### Overview
`mise` installs the tools a project needs (Node.js, Python, Go and many more) at the version the project asks for, and puts them on your `PATH` when you are in the project folder. It also sets environment variables per project, and runs tasks. A file named `mise.toml` in the project says which tools and values the project needs, so a new person gets the same setup.

### Background
`mise` is an open-source tool whose manual is at mise.jdx.dev. It plays the role of version managers such as `nvm` and `pyenv`, and of `direnv`, in one program. The tool tested here is mise 2026.10.19. Its commands that install tools download them, so those are reviewed only.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
mise --version
```

Set an environment variable for the project; this writes `mise.toml`:

```sh
mise set FOO=bar > /dev/null 2>&1 && cat mise.toml
```

Run a command with the project environment applied:

```sh
mise set FOO=bar > /dev/null 2>&1 && mise exec -- sh -c 'echo $FOO' 2> /dev/null
```

Install a tool at a version and use it in this folder; this downloads, so it is reviewed only:

```sh
mise use node@lts
```

List the tools that are installed; reviewed only because the list depends on your machine:

```sh
mise ls
```

### Subcommands
`mise use` picks a tool and its version for this folder. `mise exec` runs a command with the tools and variables applied. `mise set` writes an environment variable into `mise.toml`. `mise ls` lists installed tools.

### Common mistakes
A project file asks for a tool; `mise` may download it the first time, and a script that runs `mise exec` can therefore need the network. Trust the `mise.toml` of a project before you let it run tasks or set variables, because it can run code; `mise trust` marks a file as trusted. The tool is changing quickly, and the output text changes between versions. Running `mise` from a folder without a config file just uses the global one.

### Portability
The commands and the file format are the same on every system. Because versions move fast, check `mise --version` against the manual. The examples that use the network are reviewed only: they were not run by the tests, and the ones that call a public service were checked by hand on 2026-11-10. The two examples that set a variable were executed on mise 2026.10.19; the version example is reviewed only. The Common examples were executed on mise 2026.10.19 (see the development guide for the exact check).

## Upstream
### About upstream mise
This is the upstream `mise` tool. It is one program everywhere; what differs is the version.
