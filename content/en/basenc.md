---
schema_version: 1
id: basenc
title: basenc
summary: Encode or decode data with several base-N encodings.
aliases: [gbasenc]
tags: [text, encoding]
implementations: [gnu]
updated: 2026-11-10
references:
  - https://www.gnu.org/software/coreutils/manual/html_node/basenc-invocation.html
---
## Common
### Overview
`basenc` is the general form of `base64` and `base32`: you pick the encoding with an option such as `--base64`, `--base32` or `--base16`, and it encodes standard input as text. With `-d` it decodes. It also has URL-safe and other variants that `base64` has no option for.

### Background
`basenc` is part of GNU coreutils and was added after `base64` and `base32`; the GNU manual documents the encodings it supports. macOS does not ship it, and the BusyBox 1.37.0 build used for this page has no such applet, so this page has only a GNU view.

### Examples
Encode text as base64:

```sh
printf hello | basenc --base64
```

Encode it as base16, which is hexadecimal in capital letters:

```sh
printf hello | basenc --base16
```

Encode it in the URL-safe base64 alphabet, which uses `-` and `_` instead of `+` and `/`:

```sh
printf '\377\376' | basenc --base64url
```

Decode with `-d` and the same encoding option:

```sh
printf aGVsbG8= | basenc --base64 -d
```

### Option: --base64
Use standard base64.

### Option: --base64url
Use base64 with the URL-safe alphabet.

### Option: --base32
Use standard base32.

### Option: --base16
Use hexadecimal with capital letters.

### Option: -d
Decode instead of encode.

### Common mistakes
You must give an encoding option; `basenc` with none is an error. Decode with the same encoding that was used to encode, or the result is wrong or an error. URL-safe base64 may leave out `=` padding in other tools, but `basenc` keeps it.

### Portability
This page covers GNU only. Where a script must run on macOS or BusyBox, use `base64` instead. The Common examples were executed on GNU coreutils 9.11 (see the development guide for the exact check).

## GNU
### About GNU basenc
GNU basenc is part of GNU coreutils, and this page's checks ran coreutils 9.11.
