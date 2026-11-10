---
schema_version: 1
id: security
title: security
summary: Work with the macOS keychains, certificates and security settings.
aliases: []
tags: [macos, security, keychain]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://ss64.com/mac/security.html
  - https://developer.apple.com/documentation/security
---
## Common
### Overview
`security` is a command that only exists on macOS. It reads and changes the keychains where macOS stores passwords, keys and certificates, and the trust settings for certificates. A script can use it to fetch a stored password (`find-generic-password`) or add one (`add-generic-password`) without keeping a secret in a file. Because it handles secrets, all the examples on this page are reviewed only and none was run against a real keychain.

### Background
`security` is part of macOS, and the manual is the built-in help (`security help`) and Apple's documentation of the Security framework. This page covers the command on macOS only; other systems have other tools, such as `secret-tool` or `pass`. Nothing here was run, because the examples touch your login keychain and can show a permission prompt.

### Examples
List the keychains that are searched, one per line; this only reads:

```sh
security list-keychains
```

Show a stored password for an item named `example-service`; prints the secret, so use it only in a place where that is safe:

```sh
security find-generic-password -s example-service -w
```

Store a password for an item; the secret goes on the command line, which other users can see, so prefer typing it when asked:

```sh
security add-generic-password -a alice -s example-service -w
```

Delete an item:

```sh
security delete-generic-password -s example-service
```

Show the help for one command:

```sh
security help find-generic-password
```

### Subcommands
`security list-keychains` shows the keychains in use. `security find-generic-password` reads an item. `security add-generic-password` stores one. `security delete-generic-password` removes one. `security help COMMAND` shows the usage of a command.

### Common mistakes
The `-w` option prints the password, and the `-w` option of `add-generic-password` takes it from the command line, which shows up in the list of processes and in your shell history: leave the value off to be asked for it. The first use of an item by a program can show a window asking to allow access. A script run over SSH cannot answer such a window, and the keychain may be locked. Do not store secrets in a shell variable or log; a leaked build log with a secret is a common source of leaks.

### Portability
The command exists only on macOS, and its subcommands vary with the version of macOS. The examples are reviewed only.

## Upstream
### About upstream security
This is the macOS `security` tool. It has no GNU or BSD versions to choose between, and it exists only on macOS.
