---
schema_version: 1
id: zed
title: zed
summary: Open files and folders in the Zed editor from the terminal.
aliases: []
tags: [editor, developer]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://zed.dev/docs/
---
## Common
### Overview
`zed` is the command-line launcher of the Zed code editor. With a file or folder name it opens it in the editor, which is a window and not a terminal program, so the command returns at once. It lets you open a project from the terminal with `zed .`.

### Background
`zed` is the command that comes with the Zed editor, and its documentation is on zed.dev. The tool tested here is Zed 1.26.2. Because the command opens a window, its examples are reviewed only.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
zed --version
```

Open the current folder as a project; reviewed only because it opens a window:

```sh
zed .
```

Open a file at a line and column; reviewed only:

```sh
zed notes.txt:10:5
```

Wait until the file is closed, for use as an editor of `git`; reviewed only:

```sh
zed --wait notes.txt
```

### Option: --wait
Do not return until the files are closed.

### Option: --new
Open in a new window.

### Option: --version
Print the version.

### Common mistakes
The command is installed with the application, and it may not be on your `PATH` until you choose the install option in the editor. It opens the window and returns, so a script that needs to wait uses `--wait`. When you work over SSH, `zed` opens an editor on the machine you are sitting at only if a connection is set up.

### Portability
The options above are the same on macOS and Linux, and features change quickly with versions. The examples are reviewed only; the page was written against Zed 1.26.2.

## Upstream
### About upstream zed
This is the upstream `zed` launcher. It comes with the editor, and its version is the editor's.
