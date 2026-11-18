---
schema_version: 1
id: wget
title: wget
summary: Download files from the web, without a browser.
aliases: [gwget]
tags: [network, web, download]
implementations: [gnu, busybox]
updated: 2026-11-16
references:
  - https://www.gnu.org/software/wget/manual/wget.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`wget` downloads a file from a web address and saves it in the current folder, with the name the address ends in. It works without a terminal, so it is used in scripts and over SSH, and it can continue a broken download. Compared with `curl`, `wget` is simpler for saving files and can fetch a whole site; `curl` is better for talking to APIs.

### Background
GNU Wget is a project of the GNU system, documented in the Wget manual. BusyBox has a smaller `wget` applet with fewer options, and macOS does not ship `wget` (use `curl` there, or install it). The tool tested here is GNU Wget 1.25.0. The examples that download are reviewed only, because the checks never call another computer; each was checked once by hand against the live service on 2026-11-10.

### Examples
Print the text from a short answer to the screen: `-q` hides the progress report, and `-O -` writes to standard output:

```sh short="Print a short answer to the screen (-q, -O -)"
wget -q -O - https://api.github.com/zen
```

Save a download under a name of your choice with `-O`:

```sh
wget -q -O out.txt https://raw.githubusercontent.com/octocat/Hello-World/master/README && cat out.txt
```

Find out whether an address answers, without saving anything, with `--spider`:

```sh
wget -q --spider https://api.github.com && echo up
```

Continue a download that was cut off, with `-c`:

```sh
wget -c https://example.org/big.iso
```

Use the exit status: 8 means that the server answered with an error such as 404:

```sh
wget -q -O /dev/null https://api.github.com/nonexistent-path; echo "status $?"
```

### Option: -q
Quiet: no progress report.

### Option: -O
Write the download to this file, or `-` for standard output.

### Option: -c
Continue a partly downloaded file.

### Option: --spider
Check that the address exists, without downloading it.

### Common mistakes
Without `-O`, an existing file of the same name is not replaced: `wget` saves the new one as `name.1`, which is easy to miss in a script. A download that stops halfway leaves a part of the file; `-c` continues it. Put the address in quotes when it contains `&` or `?`. `wget` follows redirects by default, so the file you get may come from another site. `wget -r` can download a whole site and put load on it: use it with care and with `--wait`. Do not put passwords in the address or on the command line.

### Portability
Only `-q`, `-O` and `-c` are in the Common examples, and BusyBox 1.37.0 supports a smaller set than GNU Wget. The Common examples are reviewed only, and so are the help, version and usage examples in the sections below, because their text changes with each release.

## GNU
### About GNU wget
GNU Wget is the full-featured tool, tested here as version 1.25.0.

### GNU examples
Read the help text; reviewed only, because it changes with each release:

```sh
wget --help
```

Print the installed version; reviewed only, because the output changes with each release:

```sh
wget --version
```

### Option: -r
Download recursively, following links.

### Option: --wait
Wait this many seconds between downloads.

## BusyBox
### About BusyBox wget
BusyBox bundles wget as an applet, and anything beyond the common options depends on how it was built. Its download examples are reviewed only, like the common ones, and the example below is reviewed only too.

### BusyBox examples
Run it without arguments to see its usage; reviewed only, because the text depends on the version and build:

```sh
wget
```

### BusyBox notes
BusyBox `wget` supports `-q`, `-O` and `-c`, and downloads over HTTP, HTTPS (when built with it) and FTP. It has no `--version`, `--spider` or recursive mode.
