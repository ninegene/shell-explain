---
schema_version: 1
id: free
title: free
summary: Show how much memory is free and used.
aliases: []
tags: [system, memory]
implementations: [procps, busybox]
updated: 2026-11-10
references:
  - https://gitlab.com/procps-ng/procps
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`free` prints the amount of memory in the machine: how much is in total, used, free and shared, and how much is held by the system as cache (buff/cache) and could be given back to programs (available). The `available` column is the one that answers "is the machine short of memory?"; a small `free` is normal, because Linux fills spare memory with cache.

### Background
`free` is a Linux command and comes from procps-ng, with a version in BusyBox. macOS does not ship it: there, `vm_stat` and the `top` heading show memory. The two versions agree on the `Mem:` row and a `total used free` heading.

### Examples
Print the memory table and check that it has a `Mem:` row:

```sh
free | grep -c '^Mem'
```

Print the heading line and check for the `total` column:

```sh
free | head -n 1 | grep -c total
```

Show the sizes in megabytes with `-m`:

```sh
free -m | grep -c '^Mem'
```

### Option: -m
Show sizes in megabytes.

### Option: -k
Show sizes in kilobytes, the default.

### Common mistakes
A low `free` number is not a problem by itself: Linux uses spare memory as cache, so look at `available`. The units differ: older versions print kilobytes, and `-h` (human readable) is in procps but not in every BusyBox build. Containers may show the memory of the host, not the limit of the container. Swap is shown on its own `Swap:` row.

### Portability
This page has no Apple view: macOS has no `free`. The numbers differ on every machine, so the Common examples check only the shape. The Common examples were executed on procps-ng and BusyBox 1.37.0 (see the development guide for the exact check).

## BusyBox
### About BusyBox free
BusyBox bundles free as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `free` supports the options used in the Common examples. Others depend on the build.

## Procps
### About procps-ng free
This is free from procps-ng, the package that supplies the process tools of most Linux distributions, and this page's checks ran the version in Alpine 3.24. `-h` prints sizes with units such as `Gi`, and `-s N` repeats every N seconds.
