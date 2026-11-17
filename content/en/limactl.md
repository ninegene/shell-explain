---
schema_version: 1
id: limactl
title: limactl
summary: Create and manage Linux virtual machines with Lima.
aliases: [lima]
tags: [virtual-machines, macos, devops]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://lima-vm.io/docs/
---
## Common
### Overview
`limactl` is the command of Lima, which runs Linux virtual machines on macOS and Linux with file sharing and port forwarding set up. You create a machine from a template, start it, open a shell in it with `limactl shell` and stop it. Tools such as Colima use Lima underneath. This page also answers for `lima`, a separate small command that Lima installs and that is an alias for `limactl shell default`: it runs a command, or opens a shell, in the default machine.

### Background
`limactl` is an open-source tool whose documentation is at lima-vm.io. The tool tested here is limactl 2.2.1. Creating and starting machines downloads an image and runs a virtual machine, so those commands are reviewed only.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
limactl --version
```

Read the help text; reviewed only, because it changes with each release:

```sh
limactl help
```

Create and start a machine from a template; reviewed only because it downloads an image:

```sh
limactl start template:default
```

Run a shell inside the machine; reviewed only:

```sh
limactl shell default
```

List the machines and their state; reviewed only because it depends on what exists:

```sh
limactl list
```

Stop and delete a machine; reviewed only:

```sh
limactl delete --force default
```

### Subcommands
`limactl create` makes a machine from a template. `limactl start` starts one. `limactl shell` opens a shell in it. `limactl list` lists machines. `limactl stop` stops one. `limactl delete` removes one.

### Common mistakes
A machine uses disk space (often several GB) and memory while it runs: stop it with `limactl stop`. The first `start` downloads an operating system image and takes minutes. A mounted host folder is visible to the machine, so mount only what the work needs. `limactl list` with no machines prints a warning to standard error and no table. `limactl delete` removes the machine and its disk, including data inside it.

### Portability
The commands are the same on every system where Lima runs; templates and defaults change with versions. The Common examples are reviewed only.

## Upstream
### About upstream limactl
This is the upstream `limactl` tool. It is one program everywhere; what differs is the version.
