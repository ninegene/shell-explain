---
schema_version: 1
id: tsort
title: tsort
summary: Sort items so that each comes after the items it depends on.
aliases: [gtsort]
tags: [text, sorting, dependencies]
implementations: [gnu, apple]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/tsort-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=tsort&sektion=1
---
## Common
### Overview
`tsort` does a topological sort: it reads pairs of names, each meaning "the first must come before the second", and prints the names in an order that respects every pair. It is useful for dependencies, such as which package to build first. A loop in the pairs, where A needs B and B needs A, has no valid order.

### Background
The `tsort` command appeared in Version 7 AT&T UNIX, according to the macOS manual. Today the name covers GNU tsort and the tsort that macOS ships. The BusyBox 1.37.0 build used for this page has no `tsort` applet. They agree on reading pairs of names and printing an order.

### Examples
Order three names from two pairs:

```sh
printf 'a b\nb c\n' | tsort
```

Give more pairs, including two that start from the same name; any order that respects every pair is valid, so only one is shown here:

```sh
printf 'a b\na c\n' | tsort | head -n 1
```

A name on its own as a pair of the same name keeps it in the output, with no ordering constraint:

```sh
printf 'a a\nb c\n' | tsort | sort | paste -sd, -
```

### Common mistakes
The order is not unique when several names are independent, so two implementations can print different valid orders for the same input. A loop in the pairs is an error on GNU, with a message and a non-zero exit status, and on macOS the check printed a warning and still produced output, so do not rely on one behavior. Names are separated by blanks, so a name with a space needs a different tool.

### Portability
This page has no BusyBox view: BusyBox has no `tsort` applet. The Common examples avoid loops and unspecified orders. The Common examples were executed on GNU coreutils 9.11 and macOS system tools (see the development guide for the exact check).

## GNU
### About GNU tsort
GNU tsort is part of GNU coreutils, and this page's checks ran coreutils 9.11. A loop in the input is reported and the exit status is 1.

## Apple
### About macOS tsort
This is the tsort in the system folders of macOS. It follows the BSD implementation.

### Apple notes
macOS `tsort` treats a loop differently from GNU: in the check it printed an order for a loop and the exit status was 0.
