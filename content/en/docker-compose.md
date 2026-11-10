---
schema_version: 1
id: docker-compose
title: docker-compose
summary: Run several containers as one application, described in a file.
aliases: [compose]
tags: [containers, developer, devops]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://docs.docker.com/compose/
---
## Common
### Overview
`docker-compose` starts an application made of several containers (a web server, a database, a cache) from one file, `compose.yaml`. The file lists the services, their images and how they connect, and one command starts, stops or rebuilds all of them. It can also check and print the file, which needs no running daemon, and the examples use that. This page also answers for the short name `compose`, as in `docker compose`.

### Background
`docker-compose` is Docker Compose, documented on docs.docker.com. The tool tested here is Docker Compose v5.6.1, run as the stand-alone command `docker-compose` (the same program is also available as `docker compose`). Commands that start containers need a daemon and are reviewed only. No example on this page calls a cloud service: the executed ones read the tool's own help and settings in a clean home folder with no credentials, and the ones that would call a service are reviewed only.

### Examples
Print the version:

```sh
docker-compose --version | cut -c1-22
```

List the services in a compose file; here one named `web`:

```sh
printf 'services:\n  web:\n    image: nginx\n' > compose.yaml && docker-compose config --services
```

Check that a compose file is valid, with `config -q`; the status is 0 and nothing is printed except what you add:

```sh
printf 'services:\n  web:\n    image: nginx\n' > compose.yaml && docker-compose config -q && echo valid
```

Find the mistake in a file with a misspelled key; the message names the key:

```sh
printf 'services:\n  web:\n    imag: nginx\n' > bad.yaml && docker-compose -f bad.yaml config -q 2>&1 | grep -o 'services.web.imag'
```

Start all services in the background; reviewed only because it needs a daemon:

```sh
docker-compose up -d
```

Stop and remove them; reviewed only:

```sh
docker-compose down
```

### Subcommands
`docker-compose up` creates and starts the services. `docker-compose down` stops and removes them. `docker-compose config` checks and prints the file. `docker-compose ps` lists the services. `docker-compose logs` shows their output.

### Common mistakes
The file name is `compose.yaml` (older projects use `docker-compose.yml`), and the command looks for it in the current folder. Services start in the background with `-d`; without it, the terminal shows their output until you press Ctrl+C. `docker-compose down` removes the containers but keeps named volumes, so data survives; `down -v` deletes the volumes too. Indentation is part of the YAML: a tab or a wrong indent is an error that `config -q` finds before you start. Do not put passwords in the file: use an `.env` file that is not committed.

### Portability
The file format is the same on every system, and keys change with versions. The examples that read the file were executed on Docker Compose v5.6.1. The Common examples were executed on Docker Compose v5.6.1 (see the development guide for the exact check).

## Upstream
### About upstream docker-compose
This is the upstream Docker Compose. It is one program everywhere; what differs is the version.
