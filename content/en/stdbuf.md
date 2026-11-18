---
schema_version: 1
id: stdbuf
title: stdbuf
summary: Run a command with changed buffering of its input and output.
aliases: [gstdbuf]
tags: [process, pipes, scripts]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/stdbuf-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=stdbuf&sektion=1
---
## Common
### Overview
`stdbuf` runs a command with a different buffering mode for its standard streams. A program writing into a pipe usually collects its output in a buffer and sends it in big pieces, which makes a pipeline look stuck. `stdbuf -oL` asks the program to send each line at once, and `-o0` to send every byte at once. It works by preloading a small library, so it only helps programs that use the C library's standard I/O.

### Background
The `stdbuf` utility first appeared in FreeBSD 8.4, according to the macOS manual, which also says the idea comes from the GNU command of the same name. Today the name covers GNU stdbuf and the stdbuf that macOS ships. The BusyBox 1.37.0 build used for this page has no such applet. They agree on `-i`, `-o` and `-e`.

### Examples
Run a command with unbuffered output:

```sh
stdbuf -o0 echo hi
```

Make a program in a pipeline send each line as soon as it is ready; the text is unchanged:

```sh short="Send each line as soon as it is ready; text is unchanged"
printf 'a\nb\n' | stdbuf -oL cat
```

Change the buffering of the input too, with `-i`:

```sh
echo hi | stdbuf -i0 cat
```

### Option: -i
Set the buffering of standard input: `0` for none, `L` for lines, or a size.

### Option: -o
Set the buffering of standard output.

### Option: -e
Set the buffering of standard error.

### Common mistakes
`stdbuf` has no effect on programs that do their own buffering or do not use the C library (many programs written in Go or Rust, or with a static build), and on commands that start with a shell script that runs another program. It cannot fix a pipeline where the first program never writes a line. The `-o0` form is slower, since each byte is a separate write. On macOS, system programs that are protected by System Integrity Protection may ignore the preloaded library.

### Portability
This page has no BusyBox view: BusyBox has no `stdbuf` applet. Some programs have their own option for line buffering, such as `grep --line-buffered` and `sed -u`; those are more reliable. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU stdbuf
GNU stdbuf is part of GNU coreutils, and this page's checks ran coreutils 9.11.

## Apple
### About macOS stdbuf
This is the stdbuf in the system folders of macOS. It follows the BSD implementation.
