---
schema_version: 1
id: tail
title: tail
summary: Print the last lines or bytes of files, or follow a growing file.
aliases: [gtail]
tags: [text, files, logs]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/tail-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=tail&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`tail` prints the end of each file you name, ten lines by default. It is the usual tool for reading the newest log lines, and with `-f` it keeps watching a file as it grows.

### Examples
Print the last ten lines (all of a short file):

```sh
tail notes.txt
```

Print the last two lines:

```sh
tail -n 2 notes.txt
```

Print everything from line 3 onward:

```sh
tail -n +3 notes.txt
```

Print the last six bytes:

```sh
tail -c 6 notes.txt
```

Print the last line of several files; each is introduced by a header:

```sh
tail -n 1 app.log logs/server.log
```

Watch a log as new lines arrive; press Ctrl+C to stop:

```sh
tail -f app.log
```

### Option: -n
Print this many lines. `-n +3` instead starts at line 3 and prints the rest.

### Option: -c
Print this many bytes from the end.

### Option: -f
Keep the file open and print new data as it is appended. It never exits on its own; stop it with Ctrl+C. It follows the open file, so a log that is rotated to a new file may need `-F` where available.

### Common mistakes
`-n 3` means the last three lines, but `-n +3` means start at line 3, a different thing. `tail -f` blocks, so do not use it where a script must continue. Do not use `tail` on a file that is also being rewritten from the start if you need a consistent snapshot.

### Portability
`-n`, `-c`, and `-f` are in the Common examples. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check). The `-f` example is reviewed but not executed here because it never ends.

## GNU
### GNU notes
GNU `tail` adds long options such as `--follow` and `--lines`. Confirm other builds support an option before relying on it.

## Apple
### Apple notes
macOS `tail` is the BSD implementation. `-F` is supported; some GNU long options such as `--follow` are not.

## BusyBox
### BusyBox notes
BusyBox `tail` supports `-n`, `-c`, and `-f`; `-F` and other features depend on the build configuration.
