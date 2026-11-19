---
schema_version: 1
id: rmt
title: rmt
summary: Let tar and other programs use a tape drive on another computer.
aliases: []
tags: [archive, tape, network]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/tar/manual/html_node/Remote-Tape-Server.html
  - https://manpages.ubuntu.com/manpages/resolute/man8/rmt-tar.8.html
---
## Common
### Overview
`rmt` is the remote magnetic tape server. A program such as `tar` or `dump` that is told to write to a tape on another host (`user@host:/dev/tape`) starts `rmt` there, through `rsh` or `ssh`, and talks to it in a small command language to open, read, write and rewind the device. People do not usually type `rmt` themselves.

### Background
`rmt` is a program of the tape tools. GNU tar documents the protocol in its manual, as the remote tape server. The tool is not installed on macOS or on the host used for these checks (`command -v rmt` finds nothing), so no example on this page could be run. All examples are reviewed against the manual, and not executed.

### Examples
Write an archive to a tape drive on another host through SSH; `tar` starts `rmt` on the remote side. Reviewed only because it needs a remote host with a tape drive:

```sh
tar --rsh-command=/usr/bin/ssh -cf backup@tapehost.example.org:/dev/st0 project/
```

List the contents of an archive on a remote tape; reviewed only:

```sh
tar --rsh-command=/usr/bin/ssh -tf backup@tapehost.example.org:/dev/st0
```

Name the `rmt` program to use on the remote side with `--rmt-command`; reviewed only:

```sh
tar --rsh-command=/usr/bin/ssh --rmt-command=/usr/sbin/rmt -cf backup@tapehost.example.org:/dev/st0 project/
```

### Common mistakes
Tape and remote tape drives are rare today, and most backups use files and `rsync`. `rmt` runs on the remote host with the rights of the account you log in to, so give that account only the rights it needs. The location of the program differs between systems (`/usr/sbin/rmt`, `/etc/rmt`, `/usr/libexec/rmt`), which is why `tar` has an option for it. The older `rsh` is not encrypted, so use `ssh`.

### Portability
The protocol is the same on every system that has the program; the path differs. The examples are reviewed only.

## Upstream
### About upstream rmt
This is the upstream remote tape server. It is not part of macOS or of the host used for these checks, so the page is documentation only.
