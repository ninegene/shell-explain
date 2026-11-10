---
schema_version: 1
id: ssh-keygen
title: ssh-keygen
summary: Create and inspect SSH keys.
aliases: []
tags: [ssh, security, keys]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://man.openbsd.org/ssh-keygen
---
## Common
### Overview
`ssh-keygen` makes the key pairs that SSH uses to log in without a password: a private key that you keep secret, and a public key (ending in `.pub`) that you put on the servers you want to reach. It can also show the fingerprint of a key, extract the public key from a private one and change a passphrase.

### Background
`ssh-keygen` is part of OpenSSH, whose manual is on man.openbsd.org. The tool tested here is the OpenSSH that macOS ships.

### Examples
Make a key pair of the type `ed25519` with no passphrase, in a file named `key`; `-q` keeps it quiet and the public key is `key.pub`:

```sh
ssh-keygen -q -t ed25519 -N '' -f key && ls key key.pub
```

Show the size of the key from its fingerprint with `-l`; the first field is the number of bits:

```sh
ssh-keygen -q -t ed25519 -N '' -f key && ssh-keygen -lf key.pub | cut -d' ' -f1
```

Print the public key from the private one, with `-y`; it starts with the key type:

```sh
ssh-keygen -q -t ed25519 -N '' -f key && ssh-keygen -y -f key | cut -d' ' -f1
```

Show the type at the end of the fingerprint line:

```sh
ssh-keygen -q -t ed25519 -N '' -f key && ssh-keygen -lf key.pub | awk '{print $NF}'
```

### Option: -t
The type of key to make, such as `ed25519` or `rsa`.

### Option: -f
The file to write or read.

### Option: -N
The new passphrase; an empty string means none.

### Option: -l
Show the fingerprint of a key.

### Option: -y
Print the public key of a private key.

### Option: -q
Quiet.

### Common mistakes
Never share or upload the private key (the file without `.pub`), and give it a passphrase for any key you keep on a laptop; the examples use an empty one only because they are throw-away keys in a temporary folder. If `-f` names an existing file, `ssh-keygen` asks before replacing it, so it can lock up a script that has no terminal. Permissions of the private key must be private (mode 600), or `ssh` refuses to use it. `ed25519` is a good default today; very old servers may need `rsa`.

### Portability
The options above are the same in OpenSSH on every system. Key types differ with the version. The examples were executed on the OpenSSH of the macOS host. The Common examples were executed on OpenSSH (the macOS system build) (see the development guide for the exact check).

## Upstream
### About upstream ssh-keygen
This is the upstream OpenSSH `ssh-keygen`. It is one program everywhere; the version differs.
