---
schema_version: 1
id: touch
title: touch
summary: Create empty files or change file timestamps.
aliases: [gtouch]
tags: [files, time]
implementations: [gnu, apple, busybox]
updated: 2026-11-09
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/touch-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=touch&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`touch` sets the access and modification times of the files you name. If a file does not exist, it creates an empty one, which is why people use it to make a new file. It never changes what is inside an existing file.

### Background
A `touch` utility appeared in Version 7 AT&T UNIX, according to the macOS manual. Today the name covers GNU touch, the touch that macOS ships, and the touch applet in BusyBox. They agree on creating files and on `-c`, `-t` and `-r`; they differ in the extra ways to say a time, such as GNU's `-d`.

### Examples
Create an empty file:

```sh
touch new.txt && ls new.txt
```

Create several files in one command:

```sh
touch a.txt b.txt && ls a.txt b.txt
```

Update the time of an existing file; its contents stay the same:

```sh
touch notes.txt && head -n 1 notes.txt
```

Update a file only if it exists, and create nothing; `-c` means "do not create":

```sh
touch -c missing.txt && ls missing.txt 2> /dev/null || echo not created
```

Give a file an old modification time with `-t`, written as year, month, day, hour and minute, then confirm the other file is newer:

```sh
touch -t 200101020304 old.txt && test notes.txt -nt old.txt && echo notes is newer
```

Copy the time of one file onto another with `-r`:

```sh
touch -r notes.txt copy.txt && test ! copy.txt -nt notes.txt && test ! notes.txt -nt copy.txt && echo same time
```

### Option: -c
Do not create a file that does not exist. `touch` still exits with status 0.

### Option: -t
Set the time to `[[CC]YY]MMDDhhmm[.SS]` instead of the current time, as in `-t 200101020304`. Only the modification and access times are set.

### Option: -r
Use the times of another file instead of the current time, as in `-r notes.txt`.

### Common mistakes
`touch` on an existing file looks like it does nothing but it changes the modification time, which can make `make` rebuild things, make backups copy the file again, or reorder `ls -t`. A `touch` of a path whose directory does not exist fails; create the directory with `mkdir -p` first.

### Portability
Only creation, `-c`, `-t` and `-r` are in the Common examples. Ways to name a time with words, such as GNU's `-d 'yesterday'`, are not portable; macOS uses `-d` with a stricter format. Test comparisons such as `-nt` are used here only to show the result; they are not in POSIX `test`, though common shells accept them. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU touch
GNU touch is part of GNU coreutils, and this page's checks ran coreutils 9.11. Its `-d` option understands many phrases for a time.

### GNU examples
Set the modification time from a date written in words or numbers; this is GNU-only, since other versions read `-d` differently:

```sh
touch -d '2001-01-02 03:04' old.txt && test notes.txt -nt old.txt && echo notes is newer
```

### Option: -d
Use the time described by a string, such as `2001-01-02 03:04` or `yesterday`, instead of the current time. GNU accepts a wide range of phrases; the other implementations do not.

## Apple
### About macOS touch
This is the touch in `/usr/bin` on macOS. It follows the BSD implementation, and its `-d` accepts only an ISO 8601 style time.

### Apple notes
macOS `touch -d` takes a stricter format such as `2001-01-02T03:04:00` and does not understand phrases like `yesterday`. Use `-t` for a time that must work on every system.

## BusyBox
### About BusyBox touch
BusyBox bundles touch as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `touch` supports `-c`, `-t` and `-r`. `-d` and some of its formats depend on the build.
