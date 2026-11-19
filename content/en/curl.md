---
schema_version: 1
id: curl
title: curl
summary: Transfer data to or from a server with a URL.
aliases: []
tags: [network, web, http]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://curl.se/docs/manpage.html
  - https://curl.se/docs/
  - https://manpages.ubuntu.com/manpages/resolute/man1/curl.1.html
---
## Common
### Overview
`curl` sends a request to a URL and prints the answer. It speaks HTTP and HTTPS, which is what web pages and web APIs use, and many other protocols, including `file://` for local files. Use it to read an API, download a file, or send data to a service from the terminal or from a script.

### Background
`curl` is an open-source project that provides both the command and the library behind it, libcurl. The command is installed on macOS and on most Linux systems, and the options in this page are in the upstream manual. The tool tested here is curl 8.7.1, the build that macOS ships.

### Examples
Read a local file through a `file://` URL; this works without a network and shows the shape of a call:

```sh
printf 'hello\n' > f.txt && curl -s "file://$PWD/f.txt"
```

Save the answer to a file with `-o`:

```sh
printf 'hello\n' > f.txt && curl -s -o out.txt "file://$PWD/f.txt" && cat out.txt
```

Fetch a short text from the GitHub API; `-s` hides the progress meter:

```sh short="Fetch a short text; -s hides the progress meter"
curl -s https://api.github.com/zen
```

Ask a JSON API for a repository and read one field with `jq`:

```sh
curl -s https://api.github.com/repos/octocat/Hello-World | jq -r .full_name
```

Show only the status line of the answer, with `-I` to ask for the headers only:

```sh short="Show only the status line, using -I"
curl -sI https://api.github.com | head -n 1
```

Print just the HTTP status code, and throw the body away:

```sh
curl -s -o /dev/null -w '%{http_code}\n' https://api.github.com/repos/octocat/Hello-World
```

Say what you accept with a header, as the GitHub API asks:

```sh
curl -s -H 'Accept: application/vnd.github+json' https://api.github.com/repos/octocat/Hello-World | jq -r .name
```

Send data with a POST; here GitHub turns a Markdown text into HTML:

```sh
curl -s -X POST -d '{"text":"**hi**"}' https://api.github.com/markdown
```

Follow redirects with `-L` and download a file by its name, here a short README:

```sh
curl -sL https://raw.githubusercontent.com/octocat/Hello-World/master/README
```

### Option: -s
Silent: hide the progress meter and error messages.

### Option: -S
With `-s`, still show an error message if the call fails.

### Option: -o
Write the answer to this file instead of the screen.

### Option: -O
Write the answer to a file named like the end of the URL.

### Option: -L
Follow redirects.

### Option: -I
Ask for the headers only.

### Option: -H
Add a request header, as in `-H 'Accept: application/json'`.

### Option: -X
Choose the request method, such as `POST`.

### Option: -d
Send this data in the body of a request, which makes it a POST.

### Option: -w
Print a report after the transfer, for example `%{http_code}`.

### Option: -f
Fail with an error code on a server error instead of printing the error page.

### Common mistakes
By default `curl` prints the body of an error page and exits with status 0, so a script should add `-f` (or `--fail`) and check the exit status. Without `-L`, a redirect is not followed and you get the redirect message instead of the page. Put the URL in quotes when it contains `&`, `?` or `*`, which the shell would otherwise act on. Do not put secrets in the URL or on the command line: other users can see a command line, and a URL may be logged. Use a header file or a config file for tokens. Respect a service's rate limits: the GitHub API allows only a small number of calls without signing in.

### Portability
`curl` is not in POSIX, but the options above are stable across versions of the tool. The system `curl` on macOS is a different build from the one on Linux, and its list of supported protocols and its TLS library differ; check `curl --version`. The examples that use the network were checked by hand against the live service on 2026-11-10 and are not run by the tests, because the answers come from another computer and can change. The file examples were executed on curl 8.7.1 (see the development guide for the exact check). The Common examples were executed on curl 8.7.1 (the macOS system build) (see the development guide for the exact check).

## Upstream
### About upstream curl
This is the upstream `curl` tool. It is one program on every system, so there are no GNU or BSD versions to choose between; what differs is the version and the libraries it was built with.
