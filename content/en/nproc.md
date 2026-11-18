---
schema_version: 1
id: nproc
title: nproc
summary: Print the number of processing units available.
aliases: [gnproc]
tags: [system]
implementations: [gnu, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/nproc-invocation.html
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`nproc` prints how many processors are available to the current process. Scripts use it to decide how many jobs to run in parallel, as in `make -j"$(nproc)"`. The answer can be lower than the number of processors in the machine, for example inside a container with a limit.

### Background
`nproc` is part of GNU coreutils and is also an applet in BusyBox. macOS does not ship it; there, `sysctl -n hw.ncpu` gives the processor count. They agree on printing the number.

### Examples
Print the number of processors; the exact number depends on the machine, so the example only checks that a number is printed:

```sh short="Print the processor count (number shape checked)"
nproc | grep -c '^[0-9][0-9]*$'
```

Use it to choose a number of jobs, for example for `make -j`:

```sh short="Pick a job count, e.g. for make -j"
jobs=$(nproc); test "$jobs" -ge 1 && echo ok
```

### Common mistakes
`nproc` honors limits such as the `OMP_NUM_THREADS` variable in GNU, so it may print less than the machine has; `nproc --all` ignores them. Using every processor for a build can make the machine unresponsive, so leave one or two free. On macOS, `nproc` is not available and a script needs `sysctl -n hw.ncpu` there.

### Portability
This page has no Apple view: macOS has no `nproc`. A portable script can try `nproc`, then `getconf _NPROCESSORS_ONLN`, then `sysctl -n hw.ncpu`. The Common examples were executed on GNU coreutils 9.11 and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU nproc
GNU nproc is part of GNU coreutils, and this page's checks ran coreutils 9.11. `--all` prints the number of installed processors, and `--ignore=N` leaves some out.

### GNU examples
Print the installed count, which is at least as large as the usable count:

```sh
test "$(nproc --all)" -ge "$(nproc)" && echo ok
```

### Option: --all
Print the number of installed processors, not only those usable by the process.

### Option: --ignore
Leave out this many processors from the count.

## BusyBox
### About BusyBox nproc
BusyBox bundles nproc as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `nproc` supports the options used in the Common examples. Others, and GNU long options, depend on the build.
