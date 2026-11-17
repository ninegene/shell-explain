---
schema_version: 1
id: node
title: node
summary: Run JavaScript outside a browser.
aliases: []
tags: [developer, javascript, runtime]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://nodejs.org/docs/latest/api/cli.html
  - https://nodejs.org/docs/latest/api/
---
## Common
### Overview
`node` runs JavaScript on your computer. Give it a file name to run a script, `-e` to run a line of code, or `-p` to run a line and print the result. Without arguments it starts an interactive prompt. Node.js is the base for many command-line tools, such as `npm` and `prettier`.

### Background
`node` is the command of Node.js, an open-source JavaScript runtime. Its command-line options are documented in the CLI page of the manual. The tool tested here is the long-term-support release available on the host.

### Examples
Print the value of an expression with `-p`:

```sh
node -p '2+3'
```

Run a line of code with `-e`:

```sh
node -e 'console.log(1+1)'
```

Run a program read from standard input:

```sh
echo 'console.log(7)' | node
```

Pass arguments to the script; they come after the code in `process.argv`:

```sh
node -e 'console.log(process.argv.slice(1).join(","))' a b
```

Print the installed version; reviewed only, because the output changes with each release:

```sh
node --version
```

### Option: -e
Run this code.

### Option: -p
Run this code and print the result.

### Option: -v
Print the version.

### Common mistakes
The code in `-e` needs shell quotes, because `(`, `)` and `$` mean things to the shell. A script that works in one Node.js version can fail in another: check `node --version`, and use the version your project asks for. `node` does not install packages; `npm` does. A script file is run from top to bottom, and the process ends when nothing more is waiting, so a forgotten open connection keeps it alive.

### Portability
Node.js runs on macOS, Linux and Windows with the same options. Newer language features need newer versions. The first four examples were executed on the long-term-support release on the host; the version example is reviewed only. The Common examples were executed on Node.js (the long-term-support release on the host) (see the development guide for the exact check).

## Upstream
### About upstream node
This is the upstream `node` tool. It is one program everywhere; what differs is the version.
