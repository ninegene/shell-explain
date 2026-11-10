---
schema_version: 1
id: ssh-agent
title: ssh-agent
summary: Hold your SSH keys in memory so you type a passphrase only once.
aliases: []
tags: [ssh, security, keys]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://man.openbsd.org/ssh-agent
---
## Common
### Overview
`ssh-agent` is a small background program that keeps your decrypted private keys in memory and hands them to `ssh` when it needs to sign something. You start it once for a session, and then add keys with `ssh-add`. It tells other programs where to find it through the variable `SSH_AUTH_SOCK`. Most desktop systems start one for you at login.

### Background
`ssh-agent` is part of OpenSSH, whose manual is on man.openbsd.org. The tool tested here is the OpenSSH that macOS ships. The examples run the agent for the length of one command, so no agent stays behind.

### Examples
Run a command with a private agent; the variable that points to it is set, and the agent stops when the command ends:

```sh
ssh-agent -a "/tmp/sa$$" sh -c 'echo "${SSH_AUTH_SOCK:+set}"'; rm -f "/tmp/sa$$"
```

A new agent holds no keys, so `ssh-add -l` says so with status 1:

```sh
ssh-agent -a "/tmp/sa$$" sh -c 'ssh-add -l > /dev/null 2>&1; echo "status $?"'; rm -f "/tmp/sa$$"
```

Start an agent for the whole session and set the variables in the current shell; reviewed only because it leaves a program running:

```sh
eval "$(ssh-agent -s)"
```

Stop the agent that the variables point to; reviewed only:

```sh
ssh-agent -k
```

### Option: -a
The path of the socket to use.

### Option: -s
Print commands for a Bourne-style shell.

### Option: -k
Stop the running agent.

### Common mistakes
Starting `ssh-agent` bare prints the commands that set the variables and does not set them: use `eval "$(ssh-agent -s)"`, or run a command under it. An agent that is started in a script and not stopped stays in memory, so use `ssh-agent -k` or the "run a command" form. Anyone who can reach the socket can use the keys, so the socket file is private to you and an agent should not be forwarded to untrusted hosts. The paths of sockets can be too long on systems with long temporary folder names; `-a` lets you choose a short one.

### Portability
The options above are the same in OpenSSH on every system. The first two examples were executed on the OpenSSH of the macOS host; the other two are reviewed only. The Common examples were executed on OpenSSH (the macOS system build) (see the development guide for the exact check).

## Upstream
### About upstream ssh-agent
This is the upstream OpenSSH `ssh-agent`. It is one program everywhere; the version differs.
