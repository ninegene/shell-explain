---
schema_version: 1
id: awk
title: awk
summary: Process text line by line with patterns and actions.
aliases: [gawk, mawk]
tags: [text, columns, logs]
implementations: [gnu, apple, busybox, mawk]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/gawk/manual/
  - https://man.freebsd.org/cgi/man.cgi?query=awk&sektion=1
  - https://invisible-island.net/mawk/manpage/mawk.html
  - https://busybox.net/downloads/BusyBox.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gawk.1.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/mawk.1.html
---
## Common
### Overview
`awk` reads its input one line at a time, splits each line into fields, and runs the program on it. The program is a list of `pattern { action }` pairs: when a line matches the pattern, the action runs. Fields are `$1`, `$2` and so on, and `$0` is the whole line. By default fields are separated by runs of blanks; `-F` chooses another separator.

Put the program in single quotes so the shell leaves `$1` alone. `BEGIN { ... }` runs once before the first line and `END { ... }` once after the last. `NR` is the number of the current line, and `NF` is the number of fields in it.

### Background
According to the gawk manual, the name awk comes from the initials of its designers, and awk was written in 1977 at Bell Laboratories. Today the name covers separate programs: gawk, the GNU Project's awk; mawk, whose manual lists Mike Brennan and Thomas E. Dickey as authors; the awk that macOS ships; and the awk applet in BusyBox. They agree on the examples below and differ on extensions such as `gensub`, editing a file in place with `-i inplace`, and CSV mode, so a script that works on one machine can fail on another. Which program a plain `awk` command starts depends on the system.

### Examples
Print the first field of every line. Fields are separated by blanks by default:

```sh short="Print the first field of every line (blank-separated)"
awk '{print $1}' notes.txt
```

Split on commas with `-F,` and print the first column, header line included:

```sh short="Split on commas with -F, and print the first column"
awk -F, '{print $1}' data.csv
```

Print the names whose quantity is above a limit. `NR > 1` skips the header line, and `-v` passes the shell value into the program as the variable `min`:

```sh
awk -F, -v min=5 'NR > 1 && $2 > min {print $1}' data.csv
```

Add up a column. The action runs on every data line, and `END` prints the total after the last one:

```sh
awk -F, 'NR > 1 {sum += $2} END {print sum}' data.csv
```

Print the second field of lines that match a pattern. The pattern goes between slashes and is case sensitive:

```sh short="Print field 2 of lines matching a /pattern/ (case sensitive)"
awk '/ERROR/ {print $2}' app.log
```

Print how many fields each line has:

```sh
awk '{print NF}' notes.txt
```

Format output with `printf`. `%-8s` pads the name to eight characters, `%3d` right-aligns the number, and `printf` adds no newline of its own, so write `\n`:

```sh
awk -F, 'NR > 1 {printf "%-8s %3d\n", $1, $2}' data.csv
```

Count how often each value appears using an array. The order of the `for ... in` loop is not specified, so the lines can come out in any order:

```sh
awk '{count[$1]++} END {for (w in count) print w, count[w]}' fruit.txt
```

### Option: -F
Set the field separator. A single character such as `-F,` splits on that character; a longer value is treated as a regular expression. Quote it if the shell would treat it specially, as in `-F'\t'` for tabs.

### Option: -v
Assign a variable before the program starts, as in `-v min=5`. Use it to pass shell values into the program instead of splicing them into the quoted program text.

### Option: -f
Read the program from a file instead of the command line, as in `awk -f program.awk data.csv`.

### Common mistakes
Fields are numbered from 1, and `$0` is the whole line. A program in double quotes lets the shell expand `$1` before `awk` runs; use single quotes. Comparisons such as `$2 > 5` are numeric when both sides look like numbers and text comparisons otherwise, so a header line can behave differently from the data lines.

A plain separator such as `-F,` does not understand quoted CSV fields that contain commas. The `for ... in` loop visits array keys in no promised order, so sort the output when order matters. An unknown option may be ignored with only a warning on one awk and rejected on another, so read standard error.

### Portability
The Common examples use only features that macOS awk, gawk, mawk, and BusyBox awk all provide: `-F`, `-v`, patterns, `NR`, `NF`, arrays, and `printf`. They were executed on macOS system awk (awk version 20200816), gawk 5.4.0 and 5.3.2, mawk 1.3.4, and BusyBox 1.37.0 (see the development guide for the exact check).

Extensions differ. `gensub` exists in gawk and BusyBox but not in mawk or macOS awk. `strftime` and `systime` exist in gawk, mawk, and BusyBox but not in macOS awk. `-i inplace` and sorted array loops through `PROCINFO` work only in gawk, and `--csv` worked on gawk 5.3.2 and 5.4.0 but not on the other awks. The flag that prints the version also differs, so identify an awk with `explain` instead of parsing it.

## GNU
### About gawk
On this page the GNU view means gawk, the GNU Project's awk. This page's checks ran gawk 5.4.0 on macOS (Homebrew) and 5.3.2 on Alpine Linux. It adds extensions that the other awks lack, including `gensub`, ordered array loops through `PROCINFO`, `-i inplace` for editing a file in place, and `--csv` for quoted CSV input.

### GNU examples
Loop over an array in sorted key order by setting the `sorted_in` entry of `PROCINFO`. Only gawk honors it; other awks ignore it and the order is unspecified:

```sh
awk 'BEGIN {PROCINFO["sorted_in"] = "@ind_str_asc"} {count[$1]++} END {for (w in count) print w, count[w]}' fruit.txt
```

Replace text with `gensub`, which returns the changed text instead of editing `$0`. gawk and BusyBox have it; mawk and macOS awk do not:

```sh
echo 'a1b22' | awk '{print gensub(/[0-9]+/, "N", "g")}'
```

### Option: -i
Load an extension. `-i inplace` rewrites each input file with the program's output, like `sed -i`, so try it on a copy first. Other awks do not have it.

### Option: --csv
Parse input as CSV, including quoted fields that contain commas. It worked on gawk 5.3.2 and 5.4.0; the other awks rejected it or ignored it with a warning, and older gawk releases were not checked.

## Apple
### About macOS awk
This is the awk in `/usr/bin` on macOS. It reports itself as awk version 20200816, and its manual describes it as a pattern-directed scanning and processing language. It has no `gensub`, `strftime`, or `systime`, and it does not understand `-i` or `--csv`: it prints a warning that the option is ignored and carries on, so a command written for gawk can then misread its arguments. Its manual also says that only eight-bit character sets are handled correctly.

## BusyBox
### About BusyBox awk
BusyBox bundles awk as an applet, and options beyond the common ones depend on how it was built. BusyBox 1.37.0 has `gensub`, `strftime`, and `systime`, but it rejects `-i` and `--csv` and has no usable `--version` flag, so it prints its usage text instead.

### BusyBox examples
Format a timestamp with `strftime`. The third argument asks for UTC. gawk and mawk 1.3.4 accept the same call; macOS awk does not have the function:

```sh
awk 'BEGIN {print strftime("%Y-%m-%d", 86400, 1)}'
```

## Mawk
### About mawk
mawk is a separate awk whose manual lists Mike Brennan and Thomas E. Dickey as authors. Version 1.3.4 has `strftime` and `systime` but no `gensub`, and it does not accept `-i` or `--csv`. It prints its version with `--version`, although older releases may not.
