---
schema_version: 1
id: base64
title: base64
summary: Encode data as text, or decode it back.
aliases: [gbase64]
tags: [text, encoding]
implementations: [gnu, apple, busybox]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/base64-invocation.html
  - https://man.freebsd.org/cgi/man.cgi?query=base64&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
`base64` turns any data into plain letters, digits, `+`, `/` and `=`, so binary content can travel through systems that only handle text, such as email or a JSON field. With `-d` it turns that text back into the original data. It is an encoding, not encryption: anyone can decode it.

### Background
Base64 is the encoding described in RFC 4648. Today the name covers GNU base64, the base64 that macOS ships, and the base64 applet in BusyBox. They agree on encoding standard input and on `-d` for decoding.

### Examples
Encode some text:

```sh
printf hello | base64
```

Decode it back with `-d`:

```sh
printf aGVsbG8= | base64 -d
```

Encode a file, and decode the result into a new file. The file is read from standard input, because macOS `base64` has no plain file argument:

```sh
printf hello > in.txt && base64 < in.txt > in.b64 && base64 -d < in.b64 > out.txt && cat out.txt
```

Check that decoding undoes encoding for text with a newline:

```sh
printf 'a\nb\n' | base64 | base64 -d
```

### Option: -d
Decode instead of encode. Older macOS releases spell it `-D`.

### Common mistakes
Base64 is not secret: use real encryption for secrets. `echo text | base64` encodes the newline that `echo` adds, so the result differs from `printf text | base64`. Output line wrapping differs between implementations (see below), which matters when a program expects one line. The standard alphabet uses `+` and `/`; the URL-safe form uses `-` and `_` and needs extra steps.

### Portability
Only encoding, `-d` and short inputs are in the Common examples, which fit on one line on every implementation. GNU and BusyBox wrap long output at 76 characters and macOS does not. The Common examples were executed on GNU coreutils 9.11, macOS system tools, and BusyBox 1.37.0 (see the development guide for the exact check).

## GNU
### About GNU base64
GNU base64 is part of GNU coreutils, and this page's checks ran coreutils 9.11. It wraps encoded lines at 76 characters unless told otherwise.

### GNU examples
Print the encoded text on one line with `-w 0`; this matters for long input, which GNU would otherwise wrap:

```sh
printf '%0100d' 0 | base64 -w 0 | grep -c .
```

### Option: -w, --wrap
Wrap encoded lines after this many characters; `-w 0` never wraps.

## Apple
### About macOS base64
This is the base64 that macOS ships. It follows the BSD implementation.

### Apple examples
Wrap the output at 76 characters with `-b`, as GNU does by default; the default on macOS is one line:

```sh
printf '%0100d' 0 | base64 -b 76 | grep -c .
```

### Option: -b, --break
Break encoded lines after this many characters; the default is not to break them.

## BusyBox
### About BusyBox base64
BusyBox bundles base64 as an applet, and anything beyond the common options depends on how it was built. This page's examples ran on BusyBox 1.37.0.

### BusyBox notes
BusyBox `base64` wraps encoded lines at 76 characters by default and supports `-d`. `-w` depends on the build.
