---
schema_version: 1
id: ssh-add
title: ssh-add
summary: Add private keys to the SSH agent, or list the keys it holds.
aliases: [sshadd]
tags: [ssh, security, keys]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://man.openbsd.org/ssh-add
---
## Common
### Overview
`ssh-add` gives your private keys to a running `ssh-agent`, so that you type the passphrase of a key once instead of each time you connect. With `-l` it lists the keys the agent holds, with `-L` it prints their public keys and with `-D` it removes them all. It needs a running agent: without one it reports that it cannot connect. This page also answers for `sshadd`, a common misspelling; the real command has a hyphen.

### Background
`ssh-add` is part of OpenSSH, whose manual is on man.openbsd.org. The tool tested here is the OpenSSH that macOS ships. The examples that start their own agent use a short socket path in `/tmp`, because the path of a long temporary folder is too long for a Unix socket.

### Examples
Without an agent, `ssh-add -l` cannot connect and the status is 2:

```sh
env -u SSH_AUTH_SOCK ssh-add -l > /dev/null 2>&1; echo "status $?"
```

Start a private agent, add a throw-away key, and count the keys in the list; the agent stops with the command, and the socket file is removed:

```sh
ssh-keygen -q -t ed25519 -N '' -f k && ssh-agent -a "/tmp/sa$$" sh -c 'ssh-add k 2> /dev/null && ssh-add -l | wc -l | tr -d " "'; rm -f "/tmp/sa$$"
```

Print the public key of the added key; the first field is its type:

```sh
ssh-keygen -q -t ed25519 -N '' -f k && ssh-agent -a "/tmp/sa$$" sh -c 'ssh-add k 2> /dev/null; ssh-add -L | cut -d" " -f1'; rm -f "/tmp/sa$$"
```

Remove all keys again with `-D`; the list is then empty and the status is 1:

```sh
ssh-keygen -q -t ed25519 -N '' -f k && ssh-agent -a "/tmp/sa$$" sh -c 'ssh-add k 2> /dev/null; ssh-add -D 2> /dev/null; ssh-add -l > /dev/null 2>&1; echo "status $?"'; rm -f "/tmp/sa$$"
```

### Option: -l
List the fingerprints of the keys in the agent.

### Option: -L
Print the public keys held by the agent.

### Option: -D
Remove all keys from the agent.

### Option: -d
Remove one key.

### Common mistakes
An agent holds your key in memory and anyone who can use the agent's socket can use your key, so do not forward an agent (`ssh -A`) to a host you do not trust. `ssh-add` with no arguments adds the default keys, such as `~/.ssh/id_ed25519`, and asks for a passphrase on a terminal. A key with the wrong permissions is refused. On macOS, `--apple-use-keychain` stores the passphrase in the keychain, and that option does not exist in other versions.

### Portability
The options above are the same in OpenSSH on every system, and the status values are 0 (keys listed), 1 (no keys) and 2 (no agent). The examples were executed on the OpenSSH of the macOS host. The Common examples were executed on OpenSSH (the macOS system build) (see the development guide for the exact check).

## Upstream
### About upstream ssh-add
This is the upstream OpenSSH `ssh-add`. It is one program everywhere; the version differs.
