---
schema_version: 1
id: aws
title: aws
summary: Work with Amazon Web Services from the command line.
aliases: []
tags: [cloud, aws, devops]
implementations: [upstream]
updated: 2026-11-16
references:
  - https://docs.aws.amazon.com/cli/latest/reference/
---
## Common
### Overview
`aws` is the AWS command-line interface. Every command has the form `aws SERVICE ACTION`, for example `aws s3 ls`, and each one sends a signed request to AWS using the credentials in your configuration or environment. The results are JSON by default, which `jq` or the `--query` option can pick apart. Because every call can read, create, change or delete real resources and cost money, treat it with care.

### Background
`aws` is the AWS CLI version 2, documented in the AWS CLI reference. The tool tested here is aws-cli 2.38. No example on this page calls a cloud service: the executed ones read the tool's own settings in a clean home folder with no credentials, and the ones that would call a service are reviewed only. The examples that call AWS (such as `aws sts get-caller-identity`) are reviewed only.

### Examples
Print the installed version; reviewed only, because the output changes with each release:

```sh
aws --version
```

Show where the credentials would come from; in a clean home folder every value is `<not set>`. The columns are compared with blanks collapsed:

```sh
aws configure list 2> /dev/null | head -n 1
```

Ask for a setting that is not set; nothing is printed, and the status is 1:

```sh
aws configure get region; echo "status $?"
```

Read the help of a service; reviewed only, because the text changes with each release:

```sh
aws s3 help
```

Show who the credentials belong to; reviewed only because it calls AWS and needs credentials:

```sh
aws sts get-caller-identity
```

List the buckets of the account; reviewed only:

```sh
aws s3 ls
```

Pick one field of the JSON answer with `--query`; reviewed only:

```sh
aws ec2 describe-instances --query 'Reservations[].Instances[].InstanceId' --output text
```

### Subcommands
`aws configure` sets up credentials and a default region. `aws sts get-caller-identity` shows the current account. `aws s3` works with buckets and objects. `aws ec2` works with virtual machines. `aws help` lists the services.

### Common mistakes
Know which account and region a command will use before you run one that changes something: `aws sts get-caller-identity` shows the account, and `AWS_PROFILE` and `AWS_REGION` change it without a visible sign. Prefer short-lived credentials from a login (`aws sso login`) to long-lived keys, and never put keys in scripts, shell history or a repository. Commands such as `aws s3 rm --recursive` and `aws ec2 terminate-instances` delete real things and have no undo: try with `--dryrun` where it exists. Results can be long, and a service may page them. Calls cost money, and some are slow.

### Portability
The commands are the same on every system. Version 1 and version 2 of the CLI differ in output defaults and login, so check `aws --version`. The second and third examples were executed on aws-cli 2.38 in a clean home folder (see the development guide for the exact check); the others are reviewed only.

## Upstream
### About upstream aws
This is the upstream AWS CLI v2. It is one program everywhere; what differs is the version.
