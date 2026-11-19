---
schema_version: 1
id: jq
title: jq
summary: Read and transform JSON on the command line.
aliases: []
tags: [json, text, data]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://jqlang.org/manual/
  - https://jqlang.org/
  - https://manpages.ubuntu.com/manpages/resolute/man1/jq.1.html
---
## Common
### Overview
`jq` reads JSON and prints parts of it, or a new JSON made from it. You give it a small program, called a filter, as the first argument. The filter `.` means "the whole input", `.name` means the field `name`, `.[0]` the first item of a list, and a `|` sends the result of one step to the next. It is the standard tool for working with the answers of web APIs.

### Background
`jq` is an open-source command with its own manual at jqlang.org, which describes the filter language used in this page. The tool tested here is jq 1.8.2.

### Examples
Read one field of an object:

```sh
echo '{"a":1,"b":[1,2,3]}' | jq .a
```

Print a string without its quotes, with `-r` (raw output):

```sh
echo '{"n":"x"}' | jq -r .n
```

Count the items in a list:

```sh
echo '{"b":[1,2,3]}' | jq '.b | length'
```

Print each item of a list on its own line, in compact form with `-c`:

```sh
echo '{"b":[1,2,3]}' | jq -c '.b[]' | paste -sd, -
```

Keep only the items that match a condition, then take one field:

```sh
echo '[{"n":"a","v":1},{"n":"b","v":2}]' | jq -r '.[] | select(.v > 1) | .n'
```

Build new JSON from a variable with `--arg`, using `-n` for no input:

```sh short="Build new JSON from a variable (--arg, -n)"
jq -n --arg x hi -c '{greeting:$x}'
```

Sort a list:

```sh
echo '[3,1,2]' | jq -c sort
```

List the keys of an object:

```sh
echo '{"a":1}' | jq -c keys
```

Use it on the answer of an API:

```sh
curl -s https://api.github.com/users/octocat | jq -r '.login, .type'
```

### Option: -r
Raw output: print strings without quotes.

### Option: -c
Compact output: one line per result.

### Option: -n
Use no input; run the filter once with `null`.

### Option: -s
Slurp: read all inputs into one list.

### Option: -e
Set the exit status from the last result: 1 for `false` or `null`.

### Option: --arg
Set a variable to a string, as in `--arg name value`.

### Common mistakes
A filter needs quotes in the shell, because `.`, `|`, `[` and `$` mean things to it. `jq .name` prints `"value"` with quotes; add `-r` when you want the text only. A missing field gives `null`, not an error, so check with `-e` or `// "default"`. `jq` reads all of the input as JSON, so a stray log line breaks it. Use `--arg` rather than putting a shell variable inside the filter, to avoid injection. Numbers larger than about 17 digits can lose precision.

### Portability
`jq` has one upstream, so the filter language is the same everywhere, but features differ by version, so check `jq --version` before relying on a newer function. The examples that use the network were checked by hand against the live service on 2026-11-10 and are not run by the tests, because the answers come from another computer and can change. The local examples were executed on jq 1.8.2. The Common examples were executed on jq 1.8.2 (see the development guide for the exact check).

## Upstream
### About upstream jq
This is the upstream `jq` tool. The command and the filter language are the same on every system; what differs is the version.
