---
schema_version: 1
id: factor
title: factor
summary: Print the prime factors of numbers.
aliases: [gfactor]
tags: [numbers, math]
implementations: [gnu, busybox]
updated: 2026-11-12
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/factor-invocation.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`factor` breaks a whole number into its prime factors, and prints the number, a colon and the factors. Without arguments it reads numbers from standard input, one per line. It is mostly a curiosity, but it is quick for checking whether a number is prime.

### Background
`factor` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it. They agree on the output format `N: p1 p2 ...`.

### Examples
Factor a number:

```sh
factor 12
```

Factor several numbers at once:

```sh
factor 12 17
```

Read the numbers from standard input:

```sh
echo 30 | factor
```

Check whether a number is prime: a prime has exactly one factor, itself:

```sh
factor 17 | awk '{print NF - 1}'
```

### Common mistakes
A prime prints as `17: 17`, which is easy to read as a mistake. Very large numbers can take a long time, and the limit differs between implementations. `factor` accepts only non-negative whole numbers.

### Portability
This page has no Apple view: macOS has no `factor`. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU factor
GNU factor is part of GNU coreutils, and this page's checks ran coreutils 9.11. It supports numbers larger than 64 bits when built with a big-number library, and coreutils 9.11 has `-h` for exponents; older versions may not support it.

### GNU examples
Print repeated factors as powers with `-h`:

```sh
factor -h 12
```

### Option: -h, --exponents
Print a repeated factor as `p^e`.

## BusyBox
### About BusyBox factor
BusyBox bundles factor as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `factor` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
