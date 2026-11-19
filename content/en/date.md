---
schema_version: 1
id: date
title: date
summary: Print the date and time in a chosen format.
aliases: [gdate]
tags: [time, system]
implementations: [gnu, apple, busybox]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/date-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=date&sektion=1
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnudate.1.html
---
## Common
### Overview
`date` prints the current date and time. A format that starts with `+` chooses the layout, for example `+%Y-%m-%d` for `2026-11-08`. `-u` uses UTC instead of the local time zone. Giving `date` a bare number such as `110812002026` tries to set the system clock, so only ever pass it a `+` format unless you mean to change the time.

The current time changes on every run, so the examples below pipe the result through `grep -c`, which prints `1` when the output has the expected shape.

### Background
A `date` command appeared in Version 1 Unix, according to the FreeBSD and macOS manuals. Today the name covers GNU date, the date that macOS ships, and the date applet in BusyBox. They agree on printing the current time with a `+` format, `-u`, and the conversions the examples use, and differ on how to read or convert a specific time: GNU uses `-d`, while macOS uses `-r` and `-j -f`.

### Examples
Print today's date as year-month-day, such as `2026-11-08`. The `grep` checks the shape:

```sh
date +%Y-%m-%d | grep -c '^[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}$'
```

Print the time of day, such as `14:05:09`:

```sh
date +%H:%M:%S | grep -c '^[0-9][0-9]:[0-9][0-9]:[0-9][0-9]$'
```

Print the time as seconds since 1970-01-01 UTC, handy for timestamps and comparisons:

```sh short="Print seconds since 1970-01-01 UTC"
date +%s | grep -c '^[0-9]*$'
```

Print an ISO 8601 timestamp in UTC, such as `2026-11-08T22:05:09Z`:

```sh short="Print an ISO 8601 timestamp in UTC"
date -u +%Y-%m-%dT%H:%M:%SZ | grep -c 'Z$'
```

Print the name of the time zone in use; with `-u` it is always `UTC`:

```sh short="Print the time zone name; with -u it is UTC"
date -u +%Z
```

### Option: -u
Use Coordinated Universal Time instead of the local time zone.

### Common mistakes
A format without the leading `+` is read as a new time to set, not as a layout. `%Z` prints an abbreviation such as `PST`, which is ambiguous; use `%z` for a numeric offset when a script needs one. Showing or converting a time given on the command line differs between implementations, so the Common examples only print the current time.

### Portability
The Common examples use `-u` and the `%Y %m %d %H %M %S %s %Z` conversions, which all three implementations accept. Reading or converting a specific time (`-d` on GNU, `-r` and `-j -f` on macOS) is different on each, and the examples are in the sections below. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU date
GNU date is part of GNU coreutils, and this page's checks ran coreutils 9.11. It converts a specific time with `-d`, including the `@seconds` form.

### GNU examples
Convert a count of seconds since 1970 into a date with `-d @`. This is a GNU extension:

```sh
date -u -d @86400 +%Y-%m-%d
```

Parse a date and time written in words or digits and print it as seconds since 1970:

```sh short="Parse a date and print it as seconds since 1970"
date -u -d '2026-11-08 12:00' +%s
```

## Apple
### About macOS date
This is the date in `/bin` on macOS. It converts a specific time with `-r` for seconds since 1970, or `-j -f` with an input format.

### Apple examples
Convert seconds since 1970 into a date with `-r`. On GNU, `-r` takes a file name instead:

```sh
date -u -r 86400 +%Y-%m-%d
```

Parse a time with `-j -f`; `-j` stops `date` from trying to set the clock, and `-f` gives the input format:

```sh
date -u -j -f '%Y-%m-%d %H:%M:%S' '2026-11-08 12:00:00' +%s
```

## BusyBox
### About BusyBox date
BusyBox bundles date as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox examples
BusyBox accepts `-d @seconds` like GNU:

```sh
date -u -d @86400 +%Y-%m-%d
```
