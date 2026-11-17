---
schema_version: 1
id: uv
title: uv
summary: Manage Python projects, packages and virtual environments quickly.
aliases: []
tags: [developer, python, packages]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://docs.astral.sh/uv/
---
## Common
### Overview
`uv` is a single tool for most Python project work: it creates projects, makes virtual environments, installs packages, runs scripts and manages Python versions. It is much faster than `pip` and `venv` and replaces them in many workflows. `uvx` is its short form for running a tool once.

### Background
`uv` is an open-source tool whose documentation is at docs.astral.sh/uv. It is written in Rust and installed as one program. The tool tested here is the uv on the host.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
uv --version
```

Start a new project in a folder; `-q` keeps it quiet, and the files it makes are listed:

```sh
uv init -q demo && ls demo
```

Create a virtual environment and check that its folder exists:

```sh
uv venv -q && test -d .venv && echo venv
```

Add a dependency to a project; this downloads, so it is reviewed only:

```sh
uv init -q demo && cd demo && uv add requests
```

Run a command inside the project environment; reviewed only because it may download Python or packages:

```sh
uv run python -c 'print(1)'
```

### Subcommands
`uv init` creates a project. `uv venv` creates a virtual environment. `uv add` adds a dependency to the project. `uv run` runs a command in the project environment.

### Common mistakes
`uv` stores files in a cache and may download a Python version if the project asks for one, so the first run can be slow and needs the network. A project has a `pyproject.toml` and a lock file; commit both. `uv pip` is a drop-in for `pip`, and it is not the same as `uv add`, which edits the project. Mixing system `pip` and `uv` in one environment can confuse both.

### Portability
The commands are the same on every system; the tool is new and changes between versions, so check `uv --version`. The examples that use the network are reviewed only: they were not run by the tests, and the ones that call a public service were checked by hand on 2026-11-10. The examples that start a project and create a virtual environment were executed on the uv of the host; the version example is reviewed only. The Common examples were executed on uv (the version on the host) (see the development guide for the exact check).

## Upstream
### About upstream uv
This is the upstream `uv` tool. It is one program everywhere; what differs is the version.
