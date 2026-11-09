---
schema_version: 1
id: tar
title: tar
summary: Create, list and extract archive files.
aliases: [gtar]
tags: [files, compression, archive]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/tar/manual/tar.html
  - https://man.freebsd.org/cgi/man.cgi?query=bsdtar&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`tar` bundles files and directories into one archive file, lists an archive, or extracts it. One mode letter chooses the job: `-c` creates, `-t` lists and `-x` extracts. `-f` names the archive, and `-z` compresses with gzip.

### Examples
Create an archive from a directory and list its contents:

```sh
tar -cf logs.tar logs && tar -tf logs.tar
```

Archive several named files:

```sh
tar -cf backup.tar app.log notes.txt && tar -tf backup.tar
```

Create a gzip-compressed archive and list it. The `-z` flag works in both modes:

```sh
tar -czf logs.tgz logs && tar -tzf logs.tgz
```

Extract an archive into another directory; `-C` changes directory first and the target must exist:

```sh
mkdir restore && tar -cf logs.tar logs && tar -xf logs.tar -C restore && ls restore/logs
```

Print one member to standard output without extracting anything:

```sh
tar -cf logs.tar logs && tar -xOf logs.tar logs/server.log
```

Leave files out by pattern. Put the option before the names it applies to:

```sh
tar --exclude='*.log' -cf logs.tar logs && tar -tf logs.tar
```

### Option: -c
Create a new archive from the files and directories named.

### Option: -t
List the members of an archive.

### Option: -x
Extract members. With names after the archive name, extract only those.

### Option: -f
The next argument is the archive file name. Without it, `tar` uses a device or standard input and output, depending on the build.

### Option: -z
Filter the archive through gzip.

### Common mistakes
Extracting an untrusted archive can write outside the target directory; list it first. A leading slash or `..` in member names is stripped or refused differently between implementations. Always put the archive name straight after `-f`, as in `-cf name.tar`, never `-fc name.tar`.

### Portability
The Common examples use `-c`, `-t`, `-x`, `-f`, `-z`, `-C`, `-O` and `--exclude`. macOS ships `bsdtar` from libarchive, which also reads many other formats. GNU tar adds `-J`, `--transform` and many other long options, and BusyBox builds enable only a subset. The Common examples were executed on GNU tar, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### GNU examples
Rename members while archiving with `--transform`, a sed-style expression. This option is a GNU extension:

```sh
tar -cf logs.tar --transform 's,^logs,archive,' logs && tar -tf logs.tar
```

## Apple
### Apple notes
macOS `tar` is `bsdtar` from libarchive. It lists the same member names for the common examples but supports different long options and formats than GNU tar.

## BusyBox
### BusyBox notes
BusyBox `tar` supports the options in the common examples; compression flags and long options depend on the build.
