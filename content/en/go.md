---
schema_version: 1
id: go
title: go
summary: Build, run, test and manage Go programs.
aliases: []
tags: [developer, go, build]
implementations: [upstream]
updated: 2026-11-10
references:
  - https://pkg.go.dev/cmd/go
  - https://go.dev/doc/
---
## Common
### Overview
`go` is the one command for Go work: it compiles programs (`go build`), runs them (`go run`), tests them (`go test`), checks them (`go vet`) and manages the modules that a project depends on (`go mod`). A Go project is a folder with a `go.mod` file that names the module and its dependencies.

### Background
`go` is the command of the Go toolchain, and its manual is on pkg.go.dev. The examples work in a temporary folder with no network. The tool tested here is Go 1.27.3.

### Examples
Print the version and check its form:

```sh
go version | grep -c '^go version go'
```

Run a small program from one file with `go run`:

```sh
printf 'package main\nimport "fmt"\nfunc main(){fmt.Println("hi")}\n' > main.go && go run main.go
```

Start a module with `go mod init`; the first line of `go.mod` names it:

```sh
go mod init demo > /dev/null 2>&1 && head -n 1 go.mod
```

Build a program to a file and run the file:

```sh
printf 'package main\nimport "fmt"\nfunc main(){fmt.Println("hi")}\n' > main.go && go mod init demo > /dev/null 2>&1 && go build -o hello . && ./hello
```

Check a package for suspicious code with `go vet`; the status is 0 when it finds nothing:

```sh
printf 'package main\nfunc main(){}\n' > main.go && go mod init demo > /dev/null 2>&1 && go vet ./... && echo vet-ok
```

### Subcommands
`go build` compiles the package. `go run` compiles and runs a program. `go test` runs the tests. `go vet` reports suspicious code. `go mod` manages modules, as in `go mod init` and `go mod tidy`.

### Common mistakes
The first line of a program must be `package main` and it needs a `main` function, or `go run` says it cannot run a package. Run `go` commands in the folder that holds `go.mod`, or tell it with `-C`. Downloading modules needs the network, and the files go to a cache shared by all your projects. `go run` builds to a temporary place each time, so use `go build` for a program you run often. A newer `go` may download a different toolchain, as the `toolchain` line of `go.mod` asks.

### Portability
The commands are the same on every system; new language features come with new versions. The examples were executed on Go 1.27.3. The Common examples were executed on Go 1.27.3 (see the development guide for the exact check).

## Upstream
### About upstream go
This is the upstream `go` tool. It is one program everywhere; what differs is the version.
