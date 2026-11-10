---
schema_version: 1
id: top
title: top
summary: Show the busiest processes and the load, updating as you watch.
aliases: []
tags: [process, system, monitoring]
implementations: [apple, procps, busybox]
updated: 2026-11-10
references:
  - https://gitlab.com/procps-ng/procps
  - https://man.freebsd.org/cgi/man.cgi?query=top&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`top` shows a live table of the processes that use the most processor time or memory, with a summary of load and memory above it. It refreshes every few seconds until you press `q`. It is meant for a person at a terminal. For scripts, each version has a batch form that prints once and exits, which is what the examples use.

### Background
Today the name covers top from procps-ng, the top that macOS ships and the top applet in BusyBox. They differ widely in their options: macOS uses `-l` for a number of samples, and procps and BusyBox use `-b -n`. The first example on this page needs a terminal and was reviewed only.

### Examples
Run it at a terminal; keys such as `q` (quit) and `k` (kill) are shown in the top lines:

```sh
top
```

### Common mistakes
`top` takes a few seconds to sample, so a script that calls it is slow. The first sample can show wrong values for the percentage of processor, because it is an average since the process began. Press `q` to quit; if the terminal is left in a strange state, run `reset`. Do not parse its screen output: `ps -o` gives columns that scripts can read.

### Portability
The example that needs a terminal was reviewed only. The batch forms differ: macOS `top -l 1`, procps and BusyBox `top -b -n 1`. Their first line differs too, as the sections below show. The batch examples in the sections below were executed on macOS system tools, procps-ng and BusyBox 1.37.0.

## Apple
### About macOS top
This is the top in the system folders of macOS. It follows the BSD implementation.

### Apple examples
Print one sample and stop with `-l 1`; the number of processes follows `-n`, and `0` prints only the summary. The first line starts with `Processes`:

```sh
top -l 1 -n 0 | head -n 1 | grep -c Processes
```

### Apple notes
macOS `top` has `-l` for a number of samples, `-o` to sort and `-n` for the number of processes to list.

## BusyBox
### About BusyBox top
BusyBox bundles top as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `top` supports `-b` for batch mode and `-n` for the number of samples.

## Procps
### About procps-ng top
This is top from procps-ng, the package that supplies the process tools of most Linux distributions, and this page's checks ran the version in Alpine 3.24.

### Procps examples
Print one sample in batch mode with `-b -n 1`; the first line starts with `top`:

```sh
top -b -n 1 | head -n 1 | grep -c top
```
