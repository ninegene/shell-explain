---
schema_version: 1
id: lsof
title: lsof
summary: List the open files of processes, including network connections.
aliases: []
tags: [process, system, network]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://man7.org/linux/man-pages/man8/lsof.8.html
---
## Common
### Overview
`lsof` (list open files) shows which files each process has open. On Unix almost everything is a file: regular files, folders, pipes and network sockets, so `lsof` answers questions such as "which program is using this port?" and "who has this file open?". Without options it lists everything, which is long; most uses give it `-p` for a process, `-i` for network connections or a file name.

### Background
`lsof` is an open-source tool of the lsof project. The macOS build is lsof 4.91 with changes by Apple, and most Linux distributions offer the upstream version as a package. The tool tested here is the system lsof of macOS.

### Examples
List the open files of your own shell, and show the heading; the first column is `COMMAND`:

```sh
lsof -p $$ | head -n 1 | cut -d' ' -f1
```

Print only the process number with `-t`, which is useful as input for `kill`:

```sh
lsof -t -p $$ | grep -c "^$$\$"
```

Ask for the current folder of a process in a form that is easy to read, with `-F n` (the name field):

```sh
lsof -a -p $$ -d cwd -Fn | tail -n 1 | cut -c1-2
```

Find the process that listens on a port; this needs the right to see other users' processes for a full answer, and is reviewed only:

```sh
lsof -i :8080
```

Find the processes that have a file open:

```sh
lsof /var/log/system.log
```

### Option: -p
Only this process number.

### Option: -i
Only network connections, as in `-i :80` or `-i TCP`.

### Option: -t
Print only process numbers.

### Option: -a
Combine the selection options with "and".

### Option: -F
Print in a field format for programs.

### Option: -n
Do not look up host names.

### Option: -P
Do not look up port names.

### Common mistakes
Without administrator rights you see only your own processes, so a port held by another user can look free. `lsof -i` can take a while because it looks up host names; add `-n -P` to skip the lookups. The output is for people: for scripts use `-F` (field output) or `-t`. A file that was deleted but is still open shows with a link count of 0 and still takes disk space. Each option you add narrows the list by "or" unless you add `-a`, which makes them "and".

### Portability
The options above work in the macOS and Linux versions of lsof, and the columns and the file types differ a little. The first three examples were executed on lsof 4.91 (the macOS build); the other two depend on what is running and are reviewed only. The Common examples were executed on lsof 4.91 (macOS build) (see the development guide for the exact check).

## Upstream
### About upstream lsof
This is the upstream `lsof` tool. The macOS build is Apple's version of lsof; on Linux it comes from the lsof package.
