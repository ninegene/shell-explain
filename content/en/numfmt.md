---
schema_version: 1
id: numfmt
title: numfmt
summary: Convert numbers to and from human-readable sizes such as 1.5K or 2M.
aliases: [gnumfmt]
tags: [numbers, text, formatting]
implementations: [gnu]
updated: 2026-11-12
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/numfmt-invocation.html
---
## Common
### Overview
`numfmt` rewrites numbers. It turns a plain number into a readable one with a unit (1500 into 1.5k or 1.5K) with `--to`, or reads such a number and prints the plain value with `--from`. It reads numbers from its arguments or, one per line, from standard input, and with `--field` it changes one column of a line and leaves the rest alone.

### Background
`numfmt` is part of GNU coreutils; the GNU manual documents it as the tool for converting numbers to and from human-readable strings. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Print a number with an SI suffix, where a thousand is shown with a `k` or `K` suffix:

```sh short="Print a number with an SI suffix"
numfmt --to=si 1500
```

Print it with an IEC suffix, where a kilo is 1024:

```sh
numfmt --to=iec 1048576
```

Read a number with a suffix and print it in full:

```sh
numfmt --from=si 2K
```

Convert numbers read from standard input:

```sh
echo 1234567 | numfmt --to=si
```

Add a unit after the suffix with `--suffix`:

```sh
numfmt --to=iec --suffix=B 2048
```

Change only one column of each line; the rest is kept:

```sh
printf '1000 apples\n' | numfmt --field 1 --to=si
```

Choose the number of decimals with a printf-style `--format`:

```sh
numfmt --format %.2f 3.14159
```

### Option: --to
Print numbers in this scale: `si` (powers of 1000), `iec` (powers of 1024) or `iec-i` (with an `i`).

### Option: --from
Read numbers with this scale, so `2K` becomes 2000 or 2048.

### Option: --field
Change only this column of each line.

### Option: --suffix
Add this text after the converted number.

### Option: --format
Print with a printf-style format such as `%.2f`.

### Common mistakes
`si` and `iec` mean different things: 1K is 1000 in `si` and 1024 in `iec`, so a number can differ by a few percent. The suffix can be lower case `k` or upper case `K` depending on the version; do not compare the text with another tool without checking. Output is rounded, and `--round` chooses the direction; the default rounds away from zero, so `3.14159` with two decimals prints `3.15`. A column that is not a number is an error unless you use `--invalid=ignore`.

### Portability
This page covers GNU only. On macOS and BusyBox, `awk` or `printf` can do simple conversions. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU numfmt
GNU numfmt is part of GNU coreutils, and this page's checks ran coreutils 9.11.
