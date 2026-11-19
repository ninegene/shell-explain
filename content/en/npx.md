---
schema_version: 1
id: npx
title: npx
summary: Run a command from an npm package without installing it for good.
aliases: []
tags: [developer, javascript, packages]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://docs.npmjs.com/cli/commands/npx
  - https://manpages.ubuntu.com/manpages/resolute/man1/npx.1.html
---
## Common
### Overview
`npx` runs a command that is provided by an npm package. If the package is already in your project it uses that one; if not, it downloads it for this one run. It is the quick way to try a tool such as a code formatter or a project generator without a global install.

### Background
`npx` comes with `npm`, and its manual is on docs.npmjs.com. The tool tested here is the npx that comes with the long-term-support Node.js on the host.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
npx --version
```

Run a tool from a package; the first run downloads it, so this is reviewed only:

```sh short="Run a tool from a package"
npx prettier --version
```

Say which package gives the command with `--package`, when the names differ; reviewed only:

```sh short="Name the package for a command (--package)"
npx --package typescript tsc --version
```

### Option: --no-install
Do not download anything; fail if the command is not available.

### Option: --package
Name the package that provides the command.

### Common mistakes
A name that is not installed makes `npx` download it from the public registry and run its code, so a typing mistake can run an unwanted package; check the name before you press Enter. Use `--no-install` to stop it from downloading. The first run is slow, because it downloads. A tool that is in your project is found there first, which can be a different version from the one you expected.

### Portability
The behavior of downloading on a missing name changed between npm versions, so check `npx --version`. The examples that use the network are reviewed only: they were not run by the tests, and the ones that call a public service were checked by hand on 2026-11-10. The version example is reviewed only too, because its output changes with each release. The Common examples are reviewed only.

## Upstream
### About upstream npx
This is the upstream `npx` tool, installed with npm.
