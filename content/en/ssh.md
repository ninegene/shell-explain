---
schema_version: 1
id: ssh
title: ssh
summary: Log in to another computer, or run a command on it, over an encrypted connection.
aliases: []
tags: [ssh, network, remote]
implementations: [upstream]
updated: 2026-11-18
references:
  - https://man.openbsd.org/ssh
  - https://man.openbsd.org/ssh_config
  - https://manpages.ubuntu.com/manpages/resolute/man1/ssh.1.html
---
## Common
### Overview
`ssh` (secure shell) connects to another computer and gives you a shell there, or runs one command and returns its output. Everything that crosses the network is encrypted. It is the standard way to administer servers, and many other tools (`scp`, `rsync`, `git` over SSH) use it underneath. You log in with a password or, better, with a key pair made by `ssh-keygen`.

### Background
`ssh` is part of OpenSSH, whose manual is on man.openbsd.org. The settings it reads from `~/.ssh/config` are described in the `ssh_config` manual. The tool tested here is the OpenSSH that macOS ships. Commands that connect to a host are reviewed only, because the checks never contact another computer.

### Examples
Log in to a host; the first time, check the fingerprint it shows:

```sh short="Log in; the first time, check the host fingerprint"
ssh user@host
```

Use a different port with `-p` (`scp` uses a capital `-P`):

```sh short="Use a different port (-p; scp uses a capital -P)"
ssh -p 2222 user@host
```

Use a specific key with `-i`:

```sh short="Use a specific key (-i)"
ssh -i ~/.ssh/id_ed25519 user@host
```

Run one command on the host and come back:

```sh short="Run one command and come back"
ssh user@host uptime
```

Quote a longer remote command so your own shell does not split it:

```sh short="Quote a remote command so your shell leaves it alone"
ssh user@host 'cd app && git pull'
```

Forward a local port with `-L`, so `http://localhost:8080` reaches port 80 on the host:

```sh short="Local forward (-L): localhost:8080 reaches host port 80"
ssh -L 8080:localhost:80 user@host
```

Forward a port back with `-R`, so port 9000 on the host reaches your port 3000:

```sh short="Remote forward (-R): host port 9000 to your port 3000"
ssh -R 9000:localhost:3000 user@host
```

Make a SOCKS proxy on `localhost:1080` with `-D` that sends traffic through the host:

```sh short="SOCKS proxy (-D) on localhost:1080 through the host"
ssh -D 1080 user@host
```

Open only a tunnel, with no shell (`-N`), in the background (`-f`):

```sh short="Tunnel only (-N, no shell), in the background (-f)"
ssh -N -f -L 5432:db:5432 user@host
```

Go through a jump host with `-J`:

```sh short="Go through a jump host (-J)"
ssh -J jump@bastion user@host
```

Debug a connection with `-v`; `-vv` and `-vvv` say more:

```sh short="Debug a connection (-v; -vv and -vvv say more)"
ssh -v user@host
```

Test a GitHub key; it should say that you have authenticated:

```sh short="Test a GitHub key; it should say you authenticated"
ssh -T git@github.com
```

Keep an idle session from being dropped; `-o` sets any config option:

```sh short="Keep an idle session alive (-o sets any config option)"
ssh -o ServerAliveInterval=30 user@host
```

Print the settings `ssh` would use for a host, without connecting, with `-G`; here the host name:

```sh
ssh -G example.org 2> /dev/null | grep '^hostname '
```

Set the user and the port on the command line, and see them in the settings:

```sh
ssh -G -p 2222 -l alice example.org 2> /dev/null | grep -E '^(user|port) '
```

Set one option for this call with `-o`, here asking never to prompt:

```sh
ssh -G -o BatchMode=yes example.org 2> /dev/null | grep -i '^batchmode '
```

Show the installed version:

```sh
ssh -V
```

### Option: -p
The port to connect to, if not 22.

### Option: -l
The user name to log in as.

### Option: -i
The private key file to use.

### Option: -o
Set a configuration option for this call.

### Option: -L
Forward a local port to a remote address.

### Option: -J
Connect through this jump host.

### Option: -G
Print the settings that would be used, and exit.

### Option: -V
Print the version.

### Common mistakes
The first time you connect, `ssh` shows the fingerprint of the host and asks you to trust it: check it through another channel if the host matters, and never answer "yes" to a changed fingerprint without finding out why. Keep the private key private and use a passphrase; `ssh-agent` remembers it for you. `ssh host command` passes the command to the remote shell, so quotes and `$` are interpreted twice: first by your shell, then by theirs. A forwarded port opens a door into a network: bind it to `localhost` unless you mean otherwise. Do not put passwords on the command line.

### Portability
The options above are the same in OpenSSH on every system, and the configuration file is read from `~/.ssh/config`. Newer versions drop old algorithms, so a very old server may need options to connect. The examples that print settings with `-G` were executed on OpenSSH (the macOS system build) (see the development guide for the exact check); the version example and the ones that connect are reviewed only.

## Upstream
### About upstream ssh
This is the upstream OpenSSH `ssh`. It is one program everywhere; the version and the default algorithms differ.
