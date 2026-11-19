---
schema_version: 1
id: python3
title: python3
summary: Run Python programs and scripts.
aliases: [python]
tags: [developer, python, runtime]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://docs.python.org/3/using/cmdline.html
  - https://docs.python.org/3/
  - https://manpages.ubuntu.com/manpages/resolute/man1/python3.1.html
---
## Common
### Overview
`python3` runs Python. Give it a file to run a script, `-c` to run a line of code, or `-m` to run a module such as `json.tool` or `http.server`. Without arguments it starts an interactive prompt. On many systems the command `python` is missing or points to an old version, so `python3` is the name to use. This page also answers for the name `python`, but `python` is missing on many systems and points to Python 2 on some old ones, so it is not a safe synonym: check `python --version` before you use it in a script.

### Background
`python3` is the command of the Python interpreter. Its options are in the command-line page of the Python documentation. The tool tested here is the Python 3 on the host.

### Examples
Run a line of code with `-c`:

```sh
python3 -c 'print(1+2)'
```

Format JSON with the standard `json.tool` module, in compact form:

```sh
echo '{"a":1}' | python3 -m json.tool --compact
```

Print the installed version; reviewed only, because the output changes with each release:

```sh
python3 --version
```

Read the major version from inside Python:

```sh
python3 -c 'import sys; print(sys.version_info[0])'
```

Serve the current folder over HTTP; this listens on a port and is reviewed only:

```sh
python3 -m http.server 8000
```

### Option: -c
Run this code.

### Option: -m
Run this module as a script, such as `-m json.tool`.

### Option: --version
Print the version.

### Common mistakes
Many systems have several Pythons, and `python3` is the one first on `PATH`; check `which -a python3` and `python3 --version`. Install packages in a virtual environment (`python3 -m venv .venv`) and not with `sudo pip`. Code in `-c` needs shell quotes. The indentation of a script is part of its syntax, so mixing tabs and spaces is an error. `python3 -m http.server` shares the whole folder, so use it only where that is safe.

### Portability
Python 3 behaves the same on macOS and Linux for the options above; features come with newer versions. The examples that run code were executed on the Python 3 of the host; the version and server examples are reviewed only. The Common examples were executed on Python 3 (the version on the host) (see the development guide for the exact check).

## Upstream
### About upstream python3
This is the upstream Python interpreter, started with `python3`. It is one program everywhere; what differs is the version.
