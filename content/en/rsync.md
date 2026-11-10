---
schema_version: 1
id: rsync
title: rsync
summary: Copy files and folders, sending only what changed.
aliases: []
tags: [files, backup, network]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://rsync.samba.org/documentation.html
  - https://github.com/kristapsdz/openrsync#readme
---
## Common
### Overview
`rsync` copies a folder to another place, local or on another computer over SSH, and transfers only the parts that are different, so repeated copies are fast. It can keep permissions and times (`-a`), show what it would do (`-n`) and delete files that no longer exist in the source (`--delete`). It is the usual tool for backups and for uploading a site.

### Background
`rsync` is the command of the rsync project. Note that macOS ships `openrsync`, a smaller separate program that speaks the same protocol and installs as `/usr/bin/rsync`; the tool tested here is that one, and the options above work in both. Install the rsync project's program from a package manager for all its options.

### Examples
Copy a folder to another folder; `-a` keeps permissions and times and copies folders recursively. The trailing slash on the source copies its contents:

```sh
mkdir src dst && printf x > src/f && rsync -a src/ dst/ && cat dst/f
```

Copy a folder recursively without keeping permissions, with `-r`:

```sh
mkdir src dst && printf x > src/f && rsync -r src/ dst && ls dst
```

Remove files from the destination that are not in the source, with `--delete`:

```sh
mkdir src dst && printf x > src/f && touch dst/extra && rsync -a --delete src/ dst/ && ls dst
```

Try a copy without doing it, with `-n`; the destination is not created:

```sh
mkdir src && printf x > src/f && rsync -an src/ dst/ > /dev/null 2>&1; test -e dst/f || echo not-copied
```

Copy to another computer over SSH; reviewed only because it needs a remote host:

```sh
rsync -a src/ user@example.org:backup/
```

### Option: -a
Archive mode: recursive, and keep permissions, times and links.

### Option: -r
Recurse into folders.

### Option: -n
Dry run: show what would happen.

### Option: --delete
Delete files from the destination that are not in the source.

### Option: -v
Print the files transferred.

### Common mistakes
A trailing slash on the source matters: `src/` copies the contents of `src`, and `src` copies the folder itself into the destination. `--delete` removes files in the destination, so try it with `-n` first; a wrong destination can lose data. `-a` does not copy hard links, extended attributes or access lists unless you add `-H`, `-X` and `-A`. Remote copies need SSH access. `rsync` compares by size and time by default, so a file changed without changing either is missed; use `-c` for a checksum compare.

### Portability
The common options are in both rsync and openrsync, which is what macOS ships. Features such as `--info`, `--progress` and compression options differ, and openrsync supports a smaller set. The local examples were executed on openrsync (the macOS `/usr/bin/rsync`). The Common examples were executed on openrsync (the macOS /usr/bin/rsync) (see the development guide for the exact check).

## Upstream
### About upstream rsync
This is the `rsync` command found on the path. On macOS that is `openrsync`, which supports the main options; elsewhere it is usually the rsync project's program.
