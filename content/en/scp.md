---
schema_version: 1
id: scp
title: scp
summary: Copy files to or from another computer over SSH.
aliases: []
tags: [ssh, network, files]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://man.openbsd.org/scp
---
## Common
### Overview
`scp` copies files between computers using SSH, so the transfer is encrypted and uses the same keys. A remote place is written `user@host:path`. You can copy from your computer to a host, from a host to your computer, or between two hosts. For repeated copies of big folders, `rsync` is faster, because it sends only what changed.

### Background
`scp` is part of OpenSSH, whose manual is on man.openbsd.org. In current versions `scp` uses the SFTP protocol by default, and the older protocol is still available. The tool tested here is the OpenSSH that macOS ships. Copies to a host are reviewed only, since the checks never contact another computer.

### Examples
Copy a file to another name on your own computer; this needs no host:

```sh
printf data > a.txt && scp a.txt b.txt && cat b.txt
```

Copy a file to a host; reviewed only:

```sh
scp a.txt alice@example.org:backup/
```

Copy a file from a host to the current folder; reviewed only:

```sh
scp alice@example.org:backup/a.txt .
```

Copy a whole folder with `-r`; reviewed only:

```sh
scp -r project alice@example.org:backup/
```

Use another port with `-P` (a capital letter, unlike `ssh -p`); reviewed only:

```sh short="Use another port (-P, a capital, unlike ssh -p)"
scp -P 2222 a.txt alice@example.org:backup/
```

### Option: -r
Copy folders recursively.

### Option: -P
The port to connect to (capital P).

### Option: -p
Keep times and modes.

### Option: -i
The private key file to use.

### Option: -q
Quiet.

### Common mistakes
The port option is `-P` for `scp` and `-p` for `ssh`, and `-p` in `scp` means "keep times and modes", which is a common mistake. A remote path with a `*` or spaces is expanded by the remote shell in the old protocol and not in SFTP mode, so quote it carefully. `scp` overwrites files without asking. It cannot resume a broken transfer: use `rsync` for that. A colon in a local file name makes `scp` think that the part before it is a host; write `./a:b` for such a name.

### Portability
The options above are the same in OpenSSH on every system. The default protocol changed between versions, which changes how paths and wildcards on the remote side work. The local example was executed on the OpenSSH of the macOS host; the others are reviewed only. The Common examples were executed on OpenSSH (the macOS system build) (see the development guide for the exact check).

## Upstream
### About upstream scp
This is the upstream OpenSSH `scp`. It is one program everywhere; the version and the protocol it uses differ.
