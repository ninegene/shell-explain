---
schema_version: 1
id: actionlint
title: actionlint
summary: Check GitHub Actions workflow files for mistakes.
aliases: []
tags: [developer, github, lint]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://github.com/rhysd/actionlint#readme
---
## Common
### Overview
`actionlint` reads GitHub Actions workflow files (the YAML files in `.github/workflows`) and reports mistakes before you push them: wrong keys, missing required fields, bad expressions in `${{ }}`, wrong step inputs and more. It can also run `shellcheck` on the scripts in `run:` steps. It exits with status 0 when nothing is found.

### Background
`actionlint` is an open-source tool whose manual is in its repository. It does not talk to GitHub; it checks the file alone. The tool tested here is actionlint 1.7.12.

### Examples
Check a correct workflow: nothing is printed and the status is 0:

```sh
printf 'on: push\njobs:\n  a:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n' > wf.yml && actionlint wf.yml && echo clean
```

Check a workflow whose job has no `runs-on`; the first line names the file, the line and the problem:

```sh
printf 'on: push\njobs:\n  a:\n    steps:\n      - run: echo hi\n' > bad.yml && { actionlint bad.yml | head -n 1; }
```

Check that the tool reports its version:

```sh
actionlint -version | head -n 1 | grep -c '^[0-9]'
```

### Option: -version
Print the version.

### Option: -format
Choose the output format with a template.

### Common mistakes
Run `actionlint` from the repository, with no file names, to check every file in `.github/workflows`. It cannot know the secrets or variables of your repository, so it may not catch a wrong name there. If `shellcheck` is installed, `actionlint` runs it on `run:` steps, and the results appear as extra findings. A new action version can have inputs that an old `actionlint` does not know, so update both.

### Portability
The output format and the checks are the same on every system, and new checks arrive with new versions. The examples were executed on actionlint 1.7.12. The Common examples were executed on actionlint 1.7.12 (see the development guide for the exact check).

## Upstream
### About upstream actionlint
This is the upstream `actionlint` tool. It is one program everywhere; what differs is the version.
