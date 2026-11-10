---
schema_version: 1
id: prettier
title: prettier
summary: Format JavaScript, JSON, CSS, Markdown and other files in a fixed style.
aliases: []
tags: [developer, format, web]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://prettier.io/docs/en/cli.html
  - https://prettier.io/docs/en/options.html
---
## Common
### Overview
`prettier` rewrites code in a fixed style, so a team stops arguing about layout. It understands JavaScript, TypeScript, JSON, CSS, HTML, Markdown, YAML and more. With a file name it prints the formatted result; with `--check` it only says whether files are already formatted, and with `--write` it changes them.

### Background
`prettier` is an open-source tool written in JavaScript, so it needs Node.js to run. Its options are in the CLI documentation, and many are set once in a `.prettierrc` file. The tool tested here is Prettier 3.9.10.

### Examples
Format a file and print the result; the extra spaces and the missing semicolon are fixed:

```sh
printf 'a  =  1\n' > a.js && prettier a.js
```

Check a file without changing it; the status is 1 when it needs formatting, and the message goes to standard error:

```sh
printf 'a  =  1\n' > a.js && { prettier --check a.js 2> /dev/null || echo needs-formatting; }
```

Format text from standard input, naming the language with `--parser`:

```sh
printf '{"a":1,"b":[1,2]}\n' | prettier --parser json
```

Rewrite a file in place with `--write`, then print it:

```sh
printf 'a  =  1\n' > a.js && prettier --write a.js > /dev/null && cat a.js
```

### Option: --check
Report whether files are formatted, and set the exit status.

### Option: --write
Rewrite the files in place.

### Option: --parser
Choose the language when it cannot be told from a file name.

### Common mistakes
Without `--write`, `prettier` prints the formatted file and leaves the original alone, which is easy to miss. It has few options on purpose, so most disagreements are settled by its defaults. A `.prettierrc` and an `.editorconfig` in a parent folder change the result, which can make your output differ from an example. Files that are listed in `.prettierignore` are skipped. It does not fix mistakes in the code: it only changes layout, and a syntax error stops it.

### Portability
The result can change between major versions, so teams pin a version. It runs wherever Node.js runs. The examples were executed on Prettier 3.9.10. The Common examples were executed on Prettier 3.9.10 (running on Node.js) (see the development guide for the exact check).

## Upstream
### About upstream prettier
This is the upstream `prettier` tool. It is installed with npm, so its version is the one in your project or the global one.
