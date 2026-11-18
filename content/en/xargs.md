---
schema_version: 1
id: xargs
title: xargs
summary: Build and run commands from standard input.
aliases: [gxargs]
tags: [text, files, pipelines]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/findutils/manual/html_node/xargs-options.html
  - https://man.freebsd.org/cgi/man.cgi?query=xargs&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`xargs` reads words from standard input and runs a command with them as arguments. It starts as many runs as it needs to stay under the system's argument-length limit. By default words are split on blanks and newlines, so names with spaces need `-0`.

### Background
The `xargs` utility appeared in PWB Unix, the Programmer's Workbench, according to the FreeBSD and macOS manuals. It builds command lines from its input. Today the name covers GNU xargs from the findutils package, the xargs that macOS ships, and the xargs applet in BusyBox. They agree on `-n`, `-I`, `-0`, and `-r`, and differ on parallel runs with `-P` and the `-d` delimiter.

### Examples
Pass every input line to one command call:

```sh
printf 'a\nb\n' | xargs echo
```

Run the command with at most two arguments each time:

```sh short="Run with at most two arguments each time (-n 2)"
printf 'a\nb\nc\n' | xargs -n 2 echo
```

Run the command once for each input item by putting it in a `{}` placeholder:

```sh short="Run once per item, with a {} placeholder"
printf 'one\ntwo\n' | xargs -I {} echo item-{}
```

Read null-separated items so a space stays inside one name; the null separator is what `find -print0` produces:

```sh short="Read null-separated items, as from find -print0"
printf 'a b\0c\0' | xargs -0 -n 1 echo
```

Run `wc` on the files named on standard input:

```sh
printf 'app.log\nnotes.txt\n' | xargs wc -l
```

Do nothing when the input is empty, instead of running the command once with no arguments:

```sh short="Do nothing when the input is empty (-r)"
printf '' | xargs -r echo ran; echo finished
```

### Option: -n
Use at most this many arguments per command run.

### Option: -I
Replace the given string, such as `{}`, with each input item and run the command once per item. Input is split on newlines only.

### Option: -0
Split input on null bytes instead of blanks and newlines. Use it with `find -print0`.

### Option: -r
Skip the command when the input is empty. GNU runs the command once with no arguments without `-r`; BSD does not run it.

### Common mistakes
Without `-0`, a name with a space becomes two arguments. Without `-r`, GNU `xargs` runs the command once with no arguments for empty input, so `find ... | xargs rm` can run a bare `rm`. `xargs` runs the program directly, not through a shell, so pipes and redirections in the command do not work unless you run `sh -c`.

### Portability
The Common examples use `-n`, `-I`, `-0` and `-r`. Parallel runs with `-P` and the line-based `-d` delimiter are not portable. The Common examples were executed on GNU findutils 4.10.0, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU xargs
GNU xargs comes from the GNU findutils package, which also provides `find`. This page's checks ran findutils 4.10.0. It adds extensions such as `-d` to choose the input delimiter.

### GNU examples
Split input on commas with `-d`. This option is a GNU extension:

```sh
printf 'a,b,c' | xargs -d , -n 1 echo
```

## Apple
### About macOS xargs
This is the xargs in `/usr/bin` on macOS. It does not run the command for empty input and accepts `-r` silently. It has no `-d`, so convert the delimiter with `tr` first.

### Apple notes
macOS `xargs` is the BSD implementation. It does not run the command for empty input and it accepts `-r` silently. It has no `-d`; convert the delimiter with `tr` first.

## BusyBox
### About BusyBox xargs
BusyBox bundles xargs as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `xargs` supported the options in the common examples in the tested build (1.37.0); other options depend on the build.
