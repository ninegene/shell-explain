---
schema_version: 1
id: gawkbug
title: gawkbug
summary: Report a bug in gawk, the GNU awk, by mail.
aliases: []
tags: [awk, bugs, gnu]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://www.gnu.org/software/gawk/manual/html_node/Bugs.html
---
## Common
### Overview
`gawkbug` is a small script that helps you send a bug report about `gawk`. It opens an editor with a form that already has the version of gawk and information about your system, you fill in the description, and it sends the mail. It does not change anything on your machine. It exists so that reports contain what the maintainers need.

### Background
`gawkbug` comes with GNU awk (gawk), and the gawk manual explains how to report a bug. The tool tested here is the one with gawk 5.4.1. The examples that open the editor and send the mail are reviewed only: the checks never send anything. The version example is reviewed only too, because its text changes with each gawk release.

### Examples
Print the version of the script; reviewed only, because the output changes with each gawk release:

```sh
gawkbug --version
```

Start a report; it opens your editor with the form and sends the mail when you save and quit; reviewed only:

```sh short="Start a report; it mails it when you save and quit"
gawkbug
```

Mail the report to a different address, given as the last argument; reviewed only:

```sh
gawkbug bug-gawk@example.org
```

### Option: --help
Print the usage text.

### Option: --version
Print the version.

### Common mistakes
Before you report a bug, check that you have the newest version of gawk and read the manual section on bugs: many surprises are documented behavior. Make the smallest program that shows the problem, and say what you expected and what happened. The script needs a working mail command (`sendmail`) to send, and may only save the text when there is none. The form contains your system information; look at it before you send.

### Portability
The script is the same on every system where gawk is installed; it needs an editor (`EDITOR`) and a mail program. The Common examples are reviewed only.

## Upstream
### About upstream gawkbug
This is the upstream `gawkbug` script that comes with gawk. It is a shell script, and its version is the gawk version.
