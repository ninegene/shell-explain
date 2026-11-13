---
schema_version: 1
id: pkill
title: pkill
summary: Send a signal to processes found by name or by their command line.
aliases: []
tags: [process, system]
implementations: [apple, procps, busybox]
updated: 2026-11-12
references:
  - https://gitlab.com/procps-ng/procps
  - https://man.freebsd.org/cgi/man.cgi?query=pkill&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`pkill` is `pgrep` followed by `kill`: it finds the processes that match a pattern and sends them a signal, `TERM` by default, which asks them to stop. It is a quick way to stop a program by name without looking up its number. Because it acts on every match, check what the pattern matches first with `pgrep`.

### Background
Today the name covers pkill from procps-ng, the pkill that macOS ships and the pkill applet in BusyBox. They agree on the pattern, on `-f` and on `-s` or a signal name. The exit status is 0 when at least one process matched and 1 when none did.

### Examples
Stop a background process by a word in its command line. Wait up to five seconds for it to start; the trap cleans it up if the search or signal fails:

```sh
(
	sleep 7$$ &
	p=$!
	trap 'kill "$p" 2>/dev/null || :' 0
	n=0
	until pgrep -f "sleep 7$$" >/dev/null; do
		n=$((n + 1))
		[ "$n" -lt 5 ] || exit 1
		sleep 1
	done
	pkill -f "sleep 7$$" || exit 1
	wait "$p" 2>/dev/null || :
	trap - 0
	echo 'done'
) 2>/dev/null
```

Use the exit status: 1 means that no process matched:

```sh
pkill -f "no-such-process-$$" || echo none
```

Start another process, stop it, and check that it is gone afterwards:

```sh
(
	sleep 7$$ &
	p=$!
	trap 'kill "$p" 2>/dev/null || :' 0
	n=0
	until pgrep -f "sleep 7$$" >/dev/null; do
		n=$((n + 1))
		[ "$n" -lt 5 ] || exit 1
		sleep 1
	done
	pkill -f "sleep 7$$" || exit 1
	wait "$p" 2>/dev/null || :
	trap - 0
	pgrep -f "sleep 7$$" || echo gone
) 2>/dev/null
```

### Option: -f
Match the pattern against the full command line, not only the name.

### Option: -x
Match the name exactly.

### Option: -u
Only processes of this user.

### Option: -9
Send `KILL`, which cannot be caught.

### Common mistakes
`pkill` acts on every process that matches, which can be more than you expect: `pkill python` stops every Python program of yours. Look with `pgrep -l` first. Use a name only a single program has, or `-x`. `pkill -f` also matches the command line of other programs that merely mention the word, such as an editor with the file open. Signals other than `TERM` give the program less chance to clean up. You can only signal your own processes unless you are an administrator.

### Portability
Only `-f` and the exit status are in the Common examples, and they start their own process with a number that is unique, so nothing else is stopped. The Common examples were executed on macOS system tools, procps-ng and BusyBox 1.37.0 (see the development guide for the exact check).

## Apple
### About macOS pkill
This is the pkill in the system folders of macOS. It follows the BSD implementation.

## BusyBox
### About BusyBox pkill
BusyBox bundles pkill as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `pkill` supports the options used in the Common examples. Others depend on the build.

## Procps
### About procps-ng pkill
This is pkill from procps-ng, the package that supplies the process tools of most Linux distributions, and this page's checks ran the version in Alpine 3.24. It also has `-e` to print what it stopped.
