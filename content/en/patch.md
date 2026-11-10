---
schema_version: 1
id: patch
title: patch
summary: Apply a diff to files, to update them or to undo a change.
aliases: [gpatch]
tags: [files, text, compare]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/diffutils/manual/html_node/Invoking-patch.html
  - https://man.freebsd.org/cgi/man.cgi?query=patch&sektion=1
---
## Common
### Overview
`patch` reads a difference made by `diff` and applies it to the files it names, so one file or a whole source tree can be updated without sending the new copy. It is how many projects exchange changes by email or attach them to a bug report. With `-R` it undoes a patch.

### Background
Today the name covers GNU patch and the patch that macOS ships (Apple's build of a BSD patch). The BusyBox 1.37.0 build used for this page has no `patch` applet. They agree on applying a unified diff, on `-p`, on `-R` and on `--dry-run`.

### Examples
Make a patch with `diff -u`, then apply it to the old file; `patch` names the file it changed, and the result is the new content:

```sh
printf 'a\nb\n' > old.txt && printf 'a\nc\n' > new.txt && { diff -u old.txt new.txt > change.patch; patch old.txt change.patch > /dev/null && cat old.txt; }
```

Undo a patch with `-R`; the file returns to the old content:

```sh
printf 'a\nb\n' > old.txt && printf 'a\nc\n' > new.txt && { diff -u old.txt new.txt > change.patch; cp new.txt target.txt; patch -R target.txt change.patch > /dev/null && cat target.txt; }
```

Check that a patch applies without changing anything, with `--dry-run`:

```sh
printf 'a\nb\n' > old.txt && printf 'a\nc\n' > new.txt && { diff -u old.txt new.txt > change.patch; patch --dry-run old.txt change.patch > /dev/null && cat old.txt; }
```

Strip the first part of the path in the patch\'s file names with `-p1`; the names `a/f` and `b/f` become `f`:

```sh
printf 'x\n' > f && printf 'y\n' > g && { diff -u --label a/f --label b/f f g > d.patch; patch -p1 < d.patch > /dev/null && cat f; }
```

### Option: -p
Strip this many leading directory names from the file names in the patch, as in `-p1`.

### Option: -R
Treat the patch as reversed: undo it.

### Option: --dry-run
Report what would happen, without changing any file.

### Common mistakes
A patch applies to the version it was made from; if the file has changed, `patch` reports a "hunk failed" and writes the parts it could not apply to a `.rej` file. Use `--dry-run` first, and keep a copy or use version control, because `patch` changes files in place. Without `-p`, a path in the patch is used as written, which often points at a file that does not exist on your machine; `-p1` is the usual choice for patches made with `git diff`. The messages `patch` prints differ between implementations.

### Portability
This page has no BusyBox view: BusyBox has no `patch` applet. `--dry-run` is accepted by both GNU and macOS patch. The Common examples hide the messages, which differ ("checking file" on GNU, "patching file" on macOS). The Common examples were executed on GNU patch 2.8 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU patch
GNU patch is a separate package from GNU diffutils, and this page's checks ran version 2.8. It can also apply context and normal diffs, and `--backup` keeps the original.

## Apple
### About macOS patch
This is the patch in the system folders of macOS. It is a BSD patch, version 2.0-12u11-Apple on the release tested.
