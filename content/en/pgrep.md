---
schema_version: 1
id: pgrep
title: pgrep
summary: Find processes by name or by their command line and print their numbers.
aliases: []
tags: [process, system]
implementations: [apple, procps, busybox]
updated: 2026-11-12
references:
  - https://gitlab.com/procps-ng/procps
  - https://man.freebsd.org/cgi/man.cgi?query=pgrep&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`pgrep` looks through the running processes and prints the numbers (PIDs) of those that match a pattern. It saves the clumsy `ps | grep` and never lists itself. By default it matches the name of the program; with `-f` it matches the whole command line. Its exit status is 0 when it finds a match and 1 when it does not.

### Background
Today the name covers pgrep from procps-ng, the pgrep that macOS ships and the pgrep applet in BusyBox. They agree on printing process numbers, on `-f` and on the exit statuses. `pkill` is the same search followed by a signal.

### Examples
Find a background process by a word in its command line with `-f`; the number matches `$!`. The bounded wait lets it start, and the trap cleans it up on failure:

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
	pgrep -f "sleep 7$$" | grep -c "^$p\$"
	lookup_status=$?
	kill "$p" 2>/dev/null || :
	wait "$p" 2>/dev/null || :
	trap - 0
	exit "$lookup_status"
)
```

Use the exit status: 1 means that nothing matched, and nothing is printed:

```sh
pgrep -f "no-such-process-$$" || echo none
```

Count the matches by counting the lines of output:

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
	pgrep -f "sleep 7$$" | grep -c .
	lookup_status=$?
	kill "$p" 2>/dev/null || :
	wait "$p" 2>/dev/null || :
	trap - 0
	exit "$lookup_status"
)
```

### Option: -f
Match the pattern against the full command line, not only the name.

### Option: -x
Match the name exactly, not as part of a longer name.

### Option: -l
Print the name next to each number.

### Common mistakes
The pattern is a regular expression, so a name with a dot or a bracket matches more than you meant. Without `-f`, only the first part of a long command name is compared, so a long name may not match. `pgrep -f name` also matches the shell that you started it from if the name is in its own command line; the examples use a number that appears only in the process they start. Use `-x` to avoid matching `sleep` inside `sleeping`. Check the numbers before you pass them to `kill`.

### Portability
Only `-f` and the exit status are in the Common examples. `-l`, `-a` and `-c` print different things on each system. The Common examples were executed on macOS system tools, procps-ng and BusyBox 1.37.0 (see the development guide for the exact check).

## Apple
### About macOS pgrep
This is the pgrep in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Print the name next to the number with `-l`:

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
	pgrep -l -f "sleep 7$$" | cut -d' ' -f2
	lookup_status=$?
	kill "$p" 2>/dev/null || :
	wait "$p" 2>/dev/null || :
	trap - 0
	exit "$lookup_status"
)
```

### Apple notes
macOS `pgrep` has `-l` to show the name and `-lf` for the full command line.

## BusyBox
### About BusyBox pgrep
BusyBox bundles pgrep as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `pgrep` supports the options used in the Common examples. Others depend on the build.

## Procps
### About procps-ng pgrep
This is pgrep from procps-ng, the package that supplies the process tools of most Linux distributions, and this page's checks ran the version in Alpine 3.24. `-c` prints only the count of matches, and `-a` prints the full command line.

### Procps examples
Print only how many processes match, with `-c`:

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
	pgrep -c -f "sleep 7$$"
	lookup_status=$?
	kill "$p" 2>/dev/null || :
	wait "$p" 2>/dev/null || :
	trap - 0
	exit "$lookup_status"
)
```

### Option: -c, --count
Print the number of matches instead of their numbers.
