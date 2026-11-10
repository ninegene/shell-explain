---
schema_version: 1
id: chcon
title: chcon
summary: Change the SELinux security context of files.
aliases: [gchcon]
tags: [files, security, selinux]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/chcon-invocation.html
---
## Common
### Overview
`chcon` changes the SELinux security context of files. SELinux is an extra layer of access control on some Linux systems, such as Fedora and Red Hat: every file and process has a context (user, role, type and level), and policy decides who may do what by context. `chcon` changes the context of a file directly, but the change is lost when the file system is relabeled. The examples are reviewed and not run, because they need an SELinux system.

### Background
`chcon` is part of GNU coreutils; the GNU manual documents it as changing the SELinux security context. macOS and BusyBox do not ship it, so this page has only a GNU view. The examples on this page were reviewed against the manual but not executed, because the check systems do not run SELinux.

### Examples
Give a file the type that a web server may read:

```sh
chcon -t httpd_sys_content_t /srv/www/index.html
```

Copy the context of another file with `--reference`:

```sh
chcon --reference=/srv/www/index.html /srv/www/new.html
```

Change a folder and everything inside it with `-R`:

```sh
chcon -R -t httpd_sys_content_t /srv/www
```

### Option: -t
Set the type part of the context.

### Option: -R
Change folders and everything inside them.

### Option: --reference
Use the context of this file.

### Common mistakes
A change made with `chcon` is temporary: a later `restorecon` or relabel puts back the context that the policy says, so for a lasting change use `semanage fcontext` and then `restorecon`. A wrong type can make a service stop working with a "permission denied" error that is not about file permissions; check the audit log. On a system without SELinux, `chcon` fails.

### Portability
This page covers GNU only, and the examples were reviewed only.

## GNU
### About GNU chcon
GNU chcon is part of GNU coreutils.
