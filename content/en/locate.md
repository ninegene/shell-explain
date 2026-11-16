---
schema_version: 1
id: locate
title: locate
summary: Find files by name in a prebuilt index, very quickly.
aliases: [glocate]
tags: [files, search, index]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/findutils/manual/html_node/Invoking-locate.html
  - https://man.freebsd.org/cgi/man.cgi?query=locate&sektion=1
---
## Common
### Overview
`locate` looks for file names in a database of all the paths on the system, which makes it much faster than `find`. The price is that the database is built ahead of time (by `updatedb`, usually once a day), so it does not know about files made since, and it may show files that are already deleted. Use `find` when you need the current state.

### Background
Today the name covers GNU locate (part of GNU findutils) and the locate that macOS ships, which is a different program with a different database. The BusyBox 1.37.0 build used for this page has no `locate` applet. The GNU examples build a small private database with `updatedb`, so they need no system database; the macOS examples are reviewed only, because macOS builds its database only when you enable a system service.

### Examples
Find a file by part of its name in the system database; reviewed only, because the answer depends on the machine and on whether the database exists:

```sh
locate notes.txt
```

### Common mistakes
A fresh machine or a file created a moment ago may not be in the database: run `updatedb` (with administrator rights for the system one) or use `find`. The database can show files that you are not allowed to open, or leave out files in folders it was told to skip. The pattern matches anywhere in the whole path, so `locate log` finds many things; add `-b` (GNU) to match only the file name. On macOS the system database is off by default.

### Portability
The two programs share only the idea and the `-i` and `-c` options, so the examples are in the sections below. The common example is reviewed only. The GNU examples were executed on GNU findutils 4.10.0 on the host.

## GNU
### About GNU locate
GNU locate is part of GNU findutils, and this page's checks ran version 4.10.0. It reads a database made by GNU `updatedb`.

### GNU examples
Build a database of a small folder tree, then search it; the answer is a full path, so the folder is cut off to keep the example short:

```sh
mkdir -p tree/sub && printf a > tree/a.txt && printf b > tree/sub/b.txt && updatedb --localpaths="$PWD/tree" --output=db && locate -d db a.txt | sed "s#$PWD/##"
```

Count the matches instead of listing them, with `-c`:

```sh
mkdir -p tree/sub && printf a > tree/a.txt && printf b > tree/sub/b.txt && updatedb --localpaths="$PWD/tree" --output=db && locate -d db -c txt
```

Ignore case with `-i`:

```sh
mkdir -p tree/sub && printf a > tree/a.txt && printf b > tree/sub/b.txt && updatedb --localpaths="$PWD/tree" --output=db && locate -d db -i B.TXT | sed "s#$PWD/##"
```

Stop after a number of matches with `-l`:

```sh
mkdir -p tree/sub && printf a > tree/a.txt && printf b > tree/sub/b.txt && updatedb --localpaths="$PWD/tree" --output=db && locate -d db -l 1 txt | wc -l | tr -d ' '
```

### Option: -d, --database
Use this database instead of the system one.

### Option: -i, --ignore-case
Ignore case.

### Option: -c, --count
Print the number of matches.

### Option: -b, --basename
Match only the file name part.

### Option: -l, --limit
Stop after this many matches.

## Apple
### About macOS locate
This is the locate in `/usr/bin` on macOS, a BSD program with its own database format, which the program `/usr/libexec/locate.updatedb` builds. The usage line of the version tested is `locate [-0Scims] [-l limit] [-d database] pattern ...`.

### Apple notes
On macOS the database is built by a system service that is off by default, so a plain `locate` may say that the database does not exist. The examples for macOS are reviewed only. Use `mdfind -name` (Spotlight) or `find` when `locate` is not set up.
