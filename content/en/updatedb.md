---
schema_version: 1
id: updatedb
title: updatedb
summary: Build or refresh the index that locate searches.
aliases: [gupdatedb]
tags: [files, search, index]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/findutils/manual/html_node/Invoking-updatedb.html
---
## Common
### Overview
`updatedb` walks the file system and writes a database of all the file names it finds, which `locate` then searches. Normally it runs once a day as a scheduled job and needs administrator rights to index the whole system. For a small tree, a private database with `--output` and `--localpaths` needs no rights, which is how the examples use it.

### Background
`updatedb` is part of GNU findutils, whose manual describes it. macOS does not ship it (it has `/usr/libexec/locate.updatedb` instead, which is a different program), and the BusyBox 1.37.0 build used for this page has no `updatedb` applet, so this page has only a GNU view. The tool tested here is the version 4.10.0.

### Examples
Index one folder into a private database file; no output means it worked, and the file exists:

```sh
mkdir -p tree && printf a > tree/a.txt && updatedb --localpaths="$PWD/tree" --output=db && test -s db && echo built
```

Search the private database with `locate -d`; the answer is a full path, so the folder is cut off:

```sh
mkdir -p tree && printf a > tree/a.txt && updatedb --localpaths="$PWD/tree" --output=db && locate -d db a.txt | sed "s#$PWD/##"
```

Leave a folder out of the index with `--prunepaths`; the file in the skipped folder is not found:

```sh
mkdir -p tree/skip && printf a > tree/a.txt && printf s > tree/skip/s.txt && updatedb --localpaths="$PWD/tree" --prunepaths="$PWD/tree/skip" --output=db && locate -d db -c txt
```

### Option: --localpaths
The folders to index, separated by spaces.

### Option: --output
Write the database to this file instead of the system one.

### Option: --prunepaths
Folders to leave out of the index.

### Common mistakes
The system database is rebuilt by a scheduled job, often at night, so new files appear the next day unless you run `updatedb` yourself with the right to read everything. The database lists names of files you cannot read, so on a shared machine the system one is built by a user with limited rights on purpose; do not make one as an administrator for other people to search. `--output` replaces the file without asking. A very large tree takes minutes and a lot of disk reads.

### Portability
This page covers GNU only. On macOS the index of file names is Spotlight (`mdfind`), and the old BSD index is built by `/usr/libexec/locate.updatedb`. The examples were executed on GNU findutils 4.10.0 on the host. The Common examples were executed on GNU findutils 4.10.0 on the host (see the development guide for the exact check).

## GNU
### About GNU updatedb
GNU updatedb is part of GNU findutils, and this page's checks ran version 4.10.0.
