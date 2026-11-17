---
schema_version: 1
id: uvx
title: uvx
summary: Run a Python tool once, in a temporary environment.
aliases: []
tags: [developer, python, packages]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://docs.astral.sh/uv/guides/tools/
---
## Common
### Overview
`uvx` runs a command that is provided by a Python package, in an environment of its own, without installing it into your project. It is the Python partner of `npx`, and an alias for `uv tool run`. It downloads the package the first time and keeps it in a cache.

### Background
`uvx` comes with `uv`, and its guide is in the uv documentation. The tool tested here is the uvx of the uv on the host.

### Examples
Print the installed version, which is the uv version; reviewed only, because the output changes with each release:

```sh
uvx --version
```

Run a tool from a package; the first run downloads it, so this is reviewed only:

```sh
uvx ruff --version
```

Run a tool from a package whose name differs from the command, with `--from`; reviewed only:

```sh
uvx --from httpie http --version
```

### Option: --from
Install the command from this package.

### Common mistakes
A name that is not installed is downloaded from the public package index and run, so check the name before you press Enter. The tool lives in a cache, not in your project, so it is not part of your dependencies; use `uv add` for that. Each run may use a newer version unless you pin one with `name@version`.

### Portability
The behavior is the same on every system. The examples that use the network are reviewed only: they were not run by the tests, and the ones that call a public service were checked by hand on 2026-11-10. The version example is reviewed only too, because its output changes with each release. The Common examples are reviewed only.

## Upstream
### About upstream uvx
This is the upstream `uvx` tool, installed with `uv`.
