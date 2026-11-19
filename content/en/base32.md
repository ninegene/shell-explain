---
schema_version: 1
id: base32
title: base32
summary: Encode data as base32 text, or decode it back.
aliases: [gbase32]
tags: [text, encoding]
implementations: [gnu]
updated: 2026-11-18
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/base32-invocation.html
  - https://manpages.ubuntu.com/manpages/resolute/man1/gnubase32.1.html
---
## Common
### Overview
`base32` encodes data using only capital letters `A` to `Z` and the digits `2` to `7`, plus `=` for padding. It is longer than base64 but safe in places that ignore case or cannot use symbols. With `-d` it decodes. Like base64, it is an encoding, not encryption.

### Background
Base32 is the encoding described in RFC 4648. `base32` is part of GNU coreutils; macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Encode some text:

```sh
printf hello | base32
```

Decode it back with `-d`:

```sh
printf NBSWY3DP | base32 -d
```

Print the encoded text on one line with `-w 0`:

```sh
printf '%0100d' 0 | base32 -w 0 | grep -c .
```

### Option: -d
Decode instead of encode.

### Option: -w, --wrap
Wrap encoded lines after this many characters; `-w 0` never wraps.

### Common mistakes
Base32 output is about 60 percent bigger than the input, more than base64. Padding with `=` is part of the standard, so keep it when you pass the text on. Do not confuse it with other base32 alphabets, such as the Crockford variant, which `base32` does not produce.

### Portability
This page covers GNU only. `base64` is the encoding to use when a tool must work everywhere. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU base32
GNU base32 is part of GNU coreutils, and this page's checks ran coreutils 9.11.
