---
schema_version: 1
id: ssh
title: ssh
summary: Log in to another computer, or run a command on it, over an encrypted connection.
aliases: []
tags: [ssh, network, remote]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://man.openbsd.org/ssh
  - https://man.openbsd.org/ssh_config
---
## Common
### Overview
`ssh` (secure shell) connects to another computer and gives you a shell there, or runs one command and returns its output. Everything that crosses the network is encrypted. It is the standard way to administer servers, and many other tools (`scp`, `rsync`, `git` over SSH) use it underneath. You log in with a password or, better, with a key pair made by `ssh-keygen`.

### Background
`ssh` is part of OpenSSH, whose manual is on man.openbsd.org. The settings it reads from `~/.ssh/config` are described in the `ssh_config` manual. The tool tested here is the OpenSSH that macOS ships. Commands that connect to a host are reviewed only, because the checks never contact another computer.

### Examples
Show the version; it starts with `OpenSSH_`:

```sh
ssh -V 2>&1 | cut -c1-8
```

Print the settings `ssh` would use for a host, without connecting, with `-G`; here the host name:

```sh
ssh -G example.org 2> /dev/null | grep '^hostname '
```

Set the user and the port on the command line, and see them in the settings:

```sh
ssh -G -p 2222 -l alice example.org 2> /dev/null | grep -E '^(user|port) '
```

Set one option for this call with `-o`, here asking never to prompt:

```sh
ssh -G -o BatchMode=yes example.org 2> /dev/null | grep -i '^batchmode '
```

Log in to a host; reviewed only because it connects:

```sh
ssh alice@example.org
```

Run one command on a host and print its output; reviewed only:

```sh
ssh alice@example.org uptime
```

Forward a local port to a service on the host; reviewed only:

```sh
ssh -L 8080:localhost:80 alice@example.org
```

Reach a host through a jump host with `-J`; reviewed only:

```sh
ssh -J jump.example.org alice@internal.example.org
```

### Option: -p
The port to connect to, if not 22.

### Option: -l
The user name to log in as.

### Option: -i
The private key file to use.

### Option: -o
Set a configuration option for this call.

### Option: -L
Forward a local port to a remote address.

### Option: -J
Connect through this jump host.

### Option: -G
Print the settings that would be used, and exit.

### Option: -V
Print the version.

### Common mistakes
The first time you connect, `ssh` shows the fingerprint of the host and asks you to trust it: check it through another channel if the host matters, and never answer "yes" to a changed fingerprint without finding out why. Keep the private key private and use a passphrase; `ssh-agent` remembers it for you. `ssh host command` passes the command to the remote shell, so quotes and `$` are interpreted twice: first by your shell, then by theirs. A forwarded port opens a door into a network: bind it to `localhost` unless you mean otherwise. Do not put passwords on the command line.

### Portability
The options above are the same in OpenSSH on every system, and the configuration file is read from `~/.ssh/config`. Newer versions drop old algorithms, so a very old server may need options to connect. The examples that print settings were executed on the OpenSSH of the macOS host; the ones that connect are reviewed only. The Common examples were executed on OpenSSH (the macOS system build) (see the development guide for the exact check).

## Upstream
### About upstream ssh
This is the upstream OpenSSH `ssh`. It is one program everywhere; the version and the default algorithms differ.
