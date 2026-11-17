---
schema_version: 1
id: vim
title: vim
summary: Edit text files in a terminal, with commands you type.
aliases: []
tags: [editor, text, developer]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://vimhelp.org/
---
## Common
### Overview
`vim` is a text editor that runs in the terminal. It has modes: in normal mode, keys are commands (`dd` deletes a line, `x` a character); in insert mode (press `i`) keys type text; `:` starts a command such as `:w` (write), `:q` (quit) and `:wq` (both). It is on almost every Unix system, which makes it the editor you can count on over an SSH connection. Vim also has a scripted mode that applies commands to a file without showing anything, which is used in the examples.

### Background
`vim` is the open-source editor "Vi IMproved". Its manual is the built-in help (`:help`) and the site vimhelp.org. macOS ships a build of vim, and many people install a newer one. The tool tested here is the vim on the host. Examples that open the editor need a terminal and are reviewed only.

### Examples
Replace text in a file with an ex command, without opening the editor; `-es` is silent mode, `-u NONE` ignores your settings:

```sh
printf 'a\n' > f && vim -es -u NONE -c '%s/a/b/' -c 'wq' f && cat f
```

Delete the lines that match a pattern, with `g/pattern/d`:

```sh
printf 'one\ntwo\n' > g && vim -es -u NONE -c 'g/two/d' -c 'wq' g && cat g
```

Type text at the end of a line with a normal-mode command:

```sh
printf 'x\n' > h && vim -es -u NONE -c 'normal Ahello' -c 'wq' h && cat h
```

Print the version and the features it was built with; reviewed only, because the output differs between builds:

```sh
vim --version
```

Open a file to edit it; the editor takes over the terminal, so this is reviewed only:

```sh
vim notes.txt
```

### Option: -es
Silent ex mode: run commands without showing the editor.

### Option: -u
Use this startup file instead of your settings; `NONE` uses none.

### Option: -c
Run this command after the first file is loaded.

### Option: -R
Open read-only.

### Common mistakes
The biggest surprise is that you cannot quit by typing: press `Esc` first, then `:q!` to quit without saving or `:wq` to save and quit. In a script, always give `-u NONE` or `-N` so that your settings do not change the result, and `-es` so that nothing is shown. A swap file named `.name.swp` stays after a crash and makes the next open ask what to do. Vim can change the line endings and the encoding of a file; look at `:set fileformat?` when this matters. `vi` is often vim in compatible mode, which behaves differently.

### Portability
The editing commands are the same everywhere, and features (such as the `normal` command in scripted mode) come with versions: macOS ships an older vim than a package manager. The scripted examples were executed on Vim (the version on the host) (see the development guide for the exact check); the version and editor examples are reviewed only.

## Upstream
### About upstream vim
This is the upstream `vim` editor. It is one program everywhere; what differs is the version and how it was built (features such as clipboard and scripting).
