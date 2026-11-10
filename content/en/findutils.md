---
schema_version: 1
id: findutils
title: findutils
summary: The GNU package that provides find, xargs, locate and updatedb.
aliases: []
tags: [package, files, search]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/findutils/manual/html_mono/find.html
---
## Common
### Overview
GNU findutils is a package, not a single command. It supplies four programs that search for files and run commands on the results: `find` (search the file tree by name, type, size or time), `xargs` (build command lines from a list of names), `locate` (search a prebuilt index of names) and `updatedb` (build that index). Each program has its own page on this site; this page explains how the package relates to what you have installed.

### Background
The GNU manual documents the four programs together as the findutils manual. Linux distributions install the package by default or on request. macOS ships its own BSD `find` and `xargs`, and a different `locate`, so the GNU programs are an extra install there; Homebrew installs them with a `g` prefix (`gfind`, `gxargs`, `glocate`, `gupdatedb`) so they do not replace the system ones. BusyBox has small `find` and `xargs` applets and no `locate` or `updatedb`. The version checked here is GNU findutils 4.10.0.

### Examples
Ask a program of the package for its version; the first line names the package. Reviewed only, because this page has no program of its own to run:

```sh
xargs --version | head -n 1
```

The index tools report the same package; reviewed only:

```sh
locate --version | head -n 1
```

### Common mistakes
Having the word `find` on your system does not mean you have GNU findutils: on macOS it is the BSD one, which has a different set of options (no `-printf`, for example), and on a minimal Linux image it can be the BusyBox applet. Check with `find --version`, which GNU supports and BSD does not. A `find` on your `PATH` can also be another program, such as `bfs`. Scripts that use GNU-only options should say so, or call `gfind` on macOS. Use the pages for each program for their options.

### Portability
This is a package overview. The programs are `find` (see the `find` page), `xargs`, `locate` and `updatedb`, each with its own page. The examples are reviewed only, because there is no executable named `findutils` to run.

## GNU
### About GNU findutils
GNU findutils is a GNU project package, and this page was checked against version 4.10.0.
