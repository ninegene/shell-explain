---
schema_version: 1
id: runcon
title: runcon
summary: Run a command with a given SELinux security context.
aliases: [gruncon]
tags: [process, security, selinux]
implementations: [gnu]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/runcon-invocation.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnuruncon.1.html
---
## Common
### Overview
`runcon` starts a command in a security context of your choosing, which is an SELinux concept: every process runs in a context that decides what it may touch. It is a tool for testing policy or for starting a service in a specific context. It needs an SELinux system and the right to change context, so its examples are reviewed and not run.

### Background
`runcon` is part of GNU coreutils; the GNU manual documents it as running a command in a specified SELinux security context. macOS and BusyBox do not ship it, so this page has only a GNU view. The examples on this page were reviewed against the manual but not executed, because the check systems do not run SELinux.

### Examples
Run a command in a given type:

```sh
runcon -t httpd_t /usr/sbin/httpd
```

Run it in a full context, written user:role:type:level:

```sh
runcon system_u:system_r:httpd_t:s0 /usr/sbin/httpd
```

Print the context of the current process by running `id` with `-Z`:

```sh
runcon -t user_t id -Z
```

### Option: -t
Set the type part of the context.

### Option: -r
Set the role part of the context.

### Option: -u
Set the user part of the context.

### Common mistakes
The policy decides whether a transition is allowed, so a command can fail with "permission denied" even for an administrator. Changing context does not change the user: file permissions still apply. On a system without SELinux, `runcon` fails.

### Portability
This page covers GNU only, and the examples were reviewed only.

## GNU
### About GNU runcon
GNU runcon is part of GNU coreutils.
