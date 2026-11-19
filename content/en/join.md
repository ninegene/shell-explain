---
schema_version: 1
id: join
title: join
summary: Combine lines of two sorted files that share a key.
aliases: [gjoin]
tags: [text, files, columns]
implementations: [gnu, apple]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/join-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=join&sektion=1
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnujoin.1.html
---
## Common
### Overview
`join` combines two files line by line where they share a key, like a lookup between two small tables. By default the key is the first word of each line, and both files must be sorted by that key. The output has the key, the rest of the first file's line, and the rest of the second file's line.

### Background
The macOS manual says `join` conforms to IEEE Std 1003.1-2001 (POSIX.1) and gives no date for its first appearance. Today the name covers GNU join and the join that macOS ships. The BusyBox 1.37.0 build used for this page's checks has no `join` applet. The two agree on the key options `-1`, `-2` and `-t`, and on `-a` for unmatched lines.

### Examples
Match two files on their first word:

```sh
printf '1 alice\n2 bob\n3 carol\n' > names.txt && printf '1 admin\n3 user\n' > roles.txt && join names.txt roles.txt
```

Also keep lines of the first file that have no match, with `-a 1`:

```sh
printf '1 alice\n2 bob\n3 carol\n' > names.txt && printf '1 admin\n3 user\n' > roles.txt && join -a 1 names.txt roles.txt
```

Show only the lines of the first file that have no match, with `-v 1`:

```sh
printf '1 alice\n2 bob\n3 carol\n' > names.txt && printf '1 admin\n3 user\n' > roles.txt && join -v 1 names.txt roles.txt
```

Join comma-separated files with `-t`:

```sh
printf '1,alice\n2,bob\n' > names.csv && printf '1,admin\n2,user\n' > roles.csv && join -t, names.csv roles.csv
```

Match on a different column of the second file with `-2`; here the key is the second field:

```sh
printf '1 alice\n2 bob\n' > names.txt && printf 'admin 1\nuser 2\n' > roles2.txt && join -2 2 names.txt roles2.txt
```

### Option: -a
Also print the unmatched lines of file 1 or file 2, as in `-a 1`.

### Option: -v
Print only the unmatched lines of file 1 or file 2, as in `-v 1`.

### Option: -t
Use this character as the field separator instead of blanks.

### Option: -1
Use this field number of file 1 as the key, as in `-1 2`.

### Option: -2
Use this field number of file 2 as the key, as in `-2 2`.

### Common mistakes
Unsorted input gives missing or wrong matches without an error; sort both files on the key first, with the same settings. When you change the key with `-1` or `-2`, sort on that field. Field separators are runs of blanks unless you give `-t`, so a name with a space splits into two fields.

### Portability
Only the options listed above are in the Common examples, and BusyBox has no `join`, so this page has no BusyBox view. GNU adds `--header` and `--check-order`, which macOS does not have. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU join
GNU join is part of GNU coreutils, and this page's checks ran coreutils 9.11. It can treat the first line of each file as a header.

### GNU examples
Keep the header lines at the top of the output with `--header`:

```sh
printf 'id name\n1 alice\n' > names.txt && printf 'id role\n1 admin\n' > roles.txt && join --header names.txt roles.txt
```

### Option: --header
Treat the first line of each input file as field names, and print them first.

## Apple
### About macOS join
This is the join in `/usr/bin` on macOS. It follows the BSD implementation.

### Apple notes
macOS `join` has the key options shown above. It has no `--header` or `--check-order`.
