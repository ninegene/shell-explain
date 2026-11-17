---
schema_version: 1
id: colima
title: colima
summary: Run Linux containers on macOS without Docker Desktop.
aliases: []
tags: [containers, macos, devops]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://github.com/abiosoft/colima#readme
---
## Common
### Overview
`colima` starts a small Linux virtual machine on macOS or Linux and sets it up to run containers, so the `docker` command works without Docker Desktop. After `colima start`, `docker` finds the daemon that runs inside the machine. You choose the number of processors and the memory with options, and `colima stop` ends it. Under the hood it uses Lima.

### Background
`colima` is an open-source tool whose documentation is in its repository. It runs the virtual machine through `limactl`. The tool tested here is colima v0.10.3. Starting it downloads an image and runs a virtual machine, so those commands are reviewed only.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
colima version
```

List the machines it manages; the first column of the heading is `PROFILE`:

```sh
colima list 2> /dev/null | head -n 1 | cut -d' ' -f1
```

Start a machine with 4 processors and 8 GB of memory; reviewed only because it downloads and starts a virtual machine:

```sh
colima start --cpu 4 --memory 8
```

Show whether it is running; reviewed only:

```sh
colima status
```

Stop the machine; reviewed only:

```sh
colima stop
```

### Subcommands
`colima start` creates and starts the machine. `colima stop` stops it. `colima status` shows its state. `colima list` lists the machines. `colima delete` removes one.

### Common mistakes
A running machine takes processor time and memory from the host; stop it when you do not need containers. The size is set when the machine is created: change it with `colima stop`, then `colima start` with new options, or `colima delete` to start again. The first start takes minutes, because it downloads an image. Containers and images live inside the machine, so `colima delete` removes them too. Docker contexts matter: `docker context ls` shows which daemon `docker` talks to.

### Portability
The commands are the same on every system where it runs (macOS and Linux). The second example was executed on colima v0.10.3 (see the development guide for the exact check); the others are reviewed only.

## Upstream
### About upstream colima
This is the upstream `colima` tool. It is one program everywhere; what differs is the version.
