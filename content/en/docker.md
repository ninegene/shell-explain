---
schema_version: 1
id: docker
title: docker
summary: Build and run programs in containers.
aliases: []
tags: [containers, developer, devops]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://docs.docker.com/reference/cli/docker/
---
## Common
### Overview
`docker` builds and runs containers: lightweight, isolated environments that hold a program with everything it needs. An image is the packaged recipe; a container is a running copy of it. The `docker` command is a client: it talks to a background service (the daemon), which on macOS runs inside a virtual machine that a tool such as Docker Desktop, Colima or Lima provides. If no daemon is running, most commands fail with a message that they cannot connect.

### Background
`docker` is the command of Docker, and its reference is on docs.docker.com. The tool tested here is Docker 29.9.2, with no daemon running, so the examples that need one (`docker run`, `docker ps`, `docker build`) are reviewed only. No example on this page calls a cloud service: the executed ones read the tool's own help and settings in a clean home folder with no credentials, and the ones that would call a service are reviewed only.

### Examples
Print the version of the client; this needs no daemon:

```sh
docker --version | grep -c '^Docker version'
```

Read the usage of a command without running it, here `docker run`:

```sh
docker run --help | head -n 1
```

List the contexts, which are the daemons the client can talk to; the default one is named `default`:

```sh
docker context ls --format '{{.Name}}'
```

Run a short-lived container and remove it afterwards with `--rm`; reviewed only because it needs a daemon and downloads an image:

```sh
docker run --rm alpine echo hello
```

List running containers; reviewed only:

```sh
docker ps
```

Build an image from the `Dockerfile` in this folder and name it with `-t`; reviewed only:

```sh
docker build -t demo .
```

Open a shell inside a running container; reviewed only:

```sh
docker exec -it demo sh
```

### Subcommands
`docker run` starts a container from an image. `docker ps` lists containers. `docker build` makes an image from a `Dockerfile`. `docker exec` runs a command in a running container. `docker images` lists images. `docker context` chooses which daemon to use.

### Common mistakes
A container is removed only when you remove it: stopped containers and unused images pile up and use disk space, so add `--rm` to short runs and clean up with `docker system prune`, after reading what it will delete. A port or folder you publish with `-p` or `-v` opens a path from the container to your machine, so mount only what is needed. Do not put secrets in an image or a `Dockerfile`, because they stay in the layers. A container that is stopped loses its changes unless they were written to a volume. Commands that run as root inside the container are not the same as being root on your machine, but a mounted folder can be changed by it.

### Portability
The client commands above are the same everywhere; what differs is the daemon behind them (Docker Desktop, Colima, Lima or the Docker Engine on Linux). The first three examples were executed on Docker 29.9.2 with no daemon. The Common examples were executed on Docker 29.9.2 (client only) (see the development guide for the exact check).

## Upstream
### About upstream docker
This is the upstream Docker command-line client. The version differs between systems, and it needs a daemon for most of its work.
