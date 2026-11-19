---
schema_version: 1
id: pathchk
title: pathchk
summary: Check whether a path name is valid and portable.
aliases: [gpathchk]
tags: [files, paths, scripts]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/pathchk-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=pathchk&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnupathchk.1.html
---
## Common
### Overview
`pathchk` checks a path name for problems: too long, containing a directory that cannot be searched, or, with `-p`, using characters outside the portable set. It is for scripts that create files on other systems and want to know whether the name will work there. It prints nothing and exits 0 when the name is fine.

### Background
A `pathchk` utility appeared in FreeBSD 5.0, according to the macOS manual, which also says it conforms to POSIX.1-2001. Today the name covers GNU pathchk and the pathchk that macOS ships. The BusyBox 1.37.0 build used for this page has no `pathchk` applet. They agree on `-p` and `-P`.

### Examples
Check that a name uses only portable characters; it succeeds silently:

```sh short="Check that a name uses only portable characters"
pathchk -p good/name && echo ok
```

A name with a space is not in the portable character set, so the check fails:

```sh
pathchk -p 'a b' 2> /dev/null || echo not portable
```

Check an absolute path:

```sh
pathchk -p /usr/bin/x && echo ok
```

### Option: -p
Check against the POSIX portable character set and length limits.

### Option: -P
Also reject empty names and names that start with a hyphen.

### Common mistakes
Without `-p`, `pathchk` checks only against the limits of the current system, so a name can pass here and fail elsewhere. The portable set is letters, digits, dot, underscore and hyphen, with a limit of 14 characters per part in the strictest form, and some systems may allow less. `pathchk` does not check that the file exists.

### Portability
This page has no BusyBox view: BusyBox has no `pathchk` applet. Only `-p` is in the Common examples. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU pathchk
GNU pathchk is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS pathchk
This is the pathchk in the system folders of macOS. It follows the BSD implementation.
