---
schema_version: 1
id: pulumi
title: pulumi
summary: Define and deploy cloud infrastructure with ordinary programming languages.
aliases: []
tags: [cloud, infrastructure, devops]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://www.pulumi.com/docs/iac/cli/commands/
---
## Common
### Overview
`pulumi` creates, updates and deletes cloud resources described in a program written in a language you already use (TypeScript, Python, Go and others). The program describes the wanted state, `pulumi preview` shows what would change, and `pulumi up` makes the change. It records what it has made in a state, so it can change or delete it later. A change to a stack can create, replace or destroy real resources.

### Background
`pulumi` is the command of Pulumi, whose commands are described in its CLI reference. The tool tested here is Pulumi v3.272.0. No example on this page is executed: the version and help examples are reviewed only because their text changes with each release, and the others would call a service. Commands that change infrastructure (`pulumi up`, `pulumi destroy`, `pulumi refresh`) are reviewed only, and are never run against real resources by this project.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
pulumi version
```

Read the help text; reviewed only, because it changes with each release:

```sh
pulumi --help
```

Read the help for `pulumi new`, which creates a project from a template; reviewed only, because it changes with each release:

```sh
pulumi new --help
```

Create a project from a template; reviewed only because it asks questions and may call the service:

```sh
pulumi new aws-typescript
```

Show what a deployment would change, without changing anything; reviewed only because it needs a project, a login and credentials:

```sh
pulumi preview
```

Apply the changes; this can create, replace or destroy real resources, so it is reviewed only:

```sh short="Apply the changes"
pulumi up
```

Remove everything a stack created; destructive, reviewed only:

```sh
pulumi destroy
```

### Subcommands
`pulumi new` creates a project from a template. `pulumi preview` shows the planned changes. `pulumi up` applies them. `pulumi destroy` removes the resources of a stack. `pulumi stack` selects and lists stacks.

### Common mistakes
Always run `pulumi preview` and read it before `pulumi up`: the lines marked as replace or delete can destroy data, such as a database. `pulumi destroy` and `pulumi refresh` change real resources too. The state decides what Pulumi thinks exists; if it is lost or out of date, Pulumi can try to create things that already exist. Secrets in the configuration are encrypted only if you set them with `--secret`. Use a separate stack for each environment, and check which stack is selected (`pulumi stack`) before a change.

### Portability
The commands are the same on every system; the languages and providers come as packages with their own versions. The Common examples are reviewed only.

## Upstream
### About upstream pulumi
This is the upstream Pulumi CLI. It is one program everywhere; what differs is the version.
