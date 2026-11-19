---
schema_version: 1
id: ssh-copy-id
title: ssh-copy-id
summary: Install your public key on a server so you can log in with the key.
aliases: [sshcopyid]
tags: [ssh, keys, remote]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://github.com/openssh/openssh-portable/blob/master/contrib/ssh-copy-id.1
  - https://manpages.ubuntu.com/manpages/resolute/man1/ssh-copy-id.1.html
---
## Common
### Overview
`ssh-copy-id` adds your public key to the file `~/.ssh/authorized_keys` of an account on another computer. After that, `ssh` can log in to that account with the key and no password. It logs in once with the password to do this, so the server must accept passwords for the first time. It is a script around `ssh`. This page also answers for `sshcopyid`, a common misspelling; the real command has hyphens.

### Background
`ssh-copy-id` is a script that comes with OpenSSH on most systems, and its manual is included in the portable OpenSSH source. The tool tested here is the one that macOS ships. Its use connects to a host, so its examples are reviewed only.

### Examples
Print the usage text; reviewed only, because the text differs between versions:

```sh
ssh-copy-id -h
```

Install the default public key on a host; reviewed only because it connects:

```sh
ssh-copy-id alice@example.org
```

Install a specific public key with `-i`; reviewed only:

```sh
ssh-copy-id -i ~/.ssh/id_ed25519.pub alice@example.org
```

Only show what would be done, with `-n` (a dry run); reviewed only:

```sh
ssh-copy-id -n alice@example.org
```

### Option: -i
The public key file to install.

### Option: -n
Dry run: show the keys that would be added.

### Option: -p
The port to connect to.

### Option: -f
Force: install without checking whether the key is already there.

### Common mistakes
It copies your public key, so keep the private key where it is: never copy the private key. It needs an existing key pair: make one with `ssh-keygen` first. After a successful run, test the login in a second window before you close the first one and before you turn off passwords on the server, or you can lock yourself out. A key with a wrong comment or an old key can pile up in `authorized_keys` if you run the command many times; it skips keys that are already there.

### Portability
The script and its options are the same on every system that has it; the macOS and the Linux versions print slightly different messages. The examples are reviewed only; the page was written against ssh-copy-id (the macOS system build).

## Upstream
### About upstream ssh-copy-id
This is the upstream OpenSSH `ssh-copy-id` script. It is one program everywhere; the version differs.
