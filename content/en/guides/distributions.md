---
schema_version: 1
id: distributions
title: "Ubuntu, Debian, Alpine: What Is a Distribution?"
subtitle: Why two Linux computers can have different tools, and how new ones arrive
summary: Learn what a Linux distribution adds to the kernel, what a package manager does, and why the choices of one distribution decide which version of a command you get.
series: why-so-many-versions
order: 4
level: beginner
reading_minutes: 7
updated: 2026-11-09
tags: [linux, distributions, packages]
related_commands: [sed, awk]
related_guides: [whos-who, same-command-many-versions, path-and-environment, which-tool-do-i-have]
references:
  - https://www.kernel.org/linux.html
  - https://www.debian.org/intro/about
  - https://www.debian.org/doc/manuals/debian-faq/pkg-basics.en.html
  - https://ubuntu.com/about
  - https://alpinelinux.org/about/
  - https://man7.org/linux/man-pages/man5/os-release.5.html
  - https://manpages.debian.org/trixie/apt/apt.8.en.html
  - https://docs.brew.sh/Installation
---
## Goal
By the end you can say what a Linux distribution is, name three well-known ones, explain what a package manager does, and see why the choices a distribution makes decide which version of a command you get.

## The idea in plain words
Linux is a kernel. The earlier guides explained that a kernel manages the hardware but is not something you can sit down and use by itself. The kernel.org site says the same thing in its own words: the kernel is only one component of a working system, so a newcomer needs a **distribution**, which it describes as a complete Linux system.

A distribution, often shortened to **distro**, is a kernel plus everything around it, chosen and put together by one group of people. It usually includes:

1. The **kernel**.
2. The **commands and the shell**, such as GNU tools or BusyBox.
3. The **libraries** the programs share.
4. An **installer**, to put it all on a computer.
5. A **package manager**, to add and update software afterwards.
6. A **collection of packages**, which the group tests and keeps up to date.

Different groups make different choices, and that is why the tools differ. Alpine Linux says it is built on musl libc and BusyBox. Debian describes itself as a volunteer project that ships more than 70,000 packages. The kernel can be the same in both, and the commands you type can still behave differently.

## A picture
A distribution is a stack of layers, and one group picks each layer. Read the picture from the bottom up. The kernel is at the bottom, and the package manager sits on top so you can add more.

```diagram
   one distribution = layers chosen by one group

   +--------------------------------------------+
   |  package manager    apt, apk, ...          |  adds more software
   +--------------------------------------------+
   |  commands and shell   GNU tools, BusyBox   |  what you type
   +--------------------------------------------+
   |  libraries            shared code          |  used by programs
   +--------------------------------------------+
   |  Linux kernel                              |  runs the hardware
   +--------------------------------------------+
```

Diagram: four layers drawn as stacked boxes. From the bottom up they are the Linux kernel, which runs the hardware, the libraries, which are shared code used by programs, the commands and shell, such as GNU tools or BusyBox, which are what you type, and the package manager, such as apt or apk, which adds more software. One group chooses all four layers, and together they are one distribution.

## Another way to see it
Think of a car. The engine is the kernel: nothing moves without it, but nobody drives an engine. A distribution is the whole car, with seats, dashboard, wheels and a garage that supplies spare parts. Two makers can fit the same engine into different cars, and each car feels different to drive.

The table lists three distributions and what their own sites say about them. It describes what each says about itself, not a ranking.

| Name | Who makes it | What its own site says |
| --- | --- | --- |
| Debian | Volunteers | Begun 1993; 70,000+ packages |
| Ubuntu | Canonical | Desktop, server, cloud, devices |
| Alpine | Non-commercial | Small, simple, secure; musl, BusyBox |

Ubuntu and Debian are related. Ubuntu's own site says that Mark Shuttleworth gathered a small team of Debian developers who founded Canonical to create Ubuntu. They are separate distributions.

Many other distributions exist. The `os-release` file, covered below, is how a system names itself.

## Walk through it
A **package** is a bundle of the files needed for a set of related commands or features. Debian's FAQ puts it this way: packages generally contain all of the files necessary to implement a set of related commands or features. A **package manager** is the tool that finds a package, fetches it, checks that the other packages it needs are present, and puts the files in place.

Here is what happens, step by step, when someone adds a tool on a Debian-style system.

1. They ask `apt` to refresh its list. The `apt` manual says that `update` downloads package information from all configured sources, and the other commands use that data.
2. They ask `apt install` for the package by name. The manual says it installs one or more packages.
3. Behind `apt`, a lower-level tool called `dpkg` unpacks the `.deb` file. Debian's FAQ says that `dpkg` checks dependencies and stops with an error if a required package is missing, unless the user forces the install.
4. The files land in folders on the system. Some are in folders that are on the `PATH`, so the shell can now find the command by name.
5. The next time the person types the name, the shell finds it, as the `PATH` guide explained.

Other systems use other tools for the same job. Alpine's package manager is `apk`. On a Mac, Homebrew is a package manager that can be added, and its documentation says that on a Mac with an Apple Silicon processor it installs into `/opt/homebrew`. That is how a Mac can end up with GNU tools beside its own: Homebrew can install GNU versions, and whichever folder comes first on the `PATH` wins.

Installing changes the whole system, so this guide does not ask you to do it. The commands below only look.

## Try it
These commands only look at things. They print the system's name, the file where a Linux system identifies itself, and which package managers are available. The answers depend on your computer, so there is no fixed output.

Print the name of the operating system. A Mac says `Darwin` and a Linux computer says `Linux`:

```sh
uname -s
```

Look for the file where a Linux system names its distribution. The `os-release` manual describes it as holding the name and version, in lines like `NAME=` and `ID=`. A Mac has no such file, so this prints a note instead:

```sh
if [ -r /etc/os-release ]; then
  cat /etc/os-release
else
  echo 'no os-release file here'
fi
```

Ask which package managers this computer has. Each name that is missing gets a note:

```sh
for pm in apt dpkg apk dnf brew; do
  command -v "$pm" || echo "$pm: not found"
done
```

## Common confusions
- **Linux is a complete operating system.** Linux is the kernel. A distribution is the complete system built around it.
- **All Linux computers have the same tools.** They share a kernel but not necessarily the commands. One may use GNU tools and another BusyBox.
- **Ubuntu is Debian.** Ubuntu was founded by Debian developers, but it is a separate distribution.
- **A package manager is only for Linux.** Homebrew shows the idea on a Mac: it is a tool that fetches and installs packages.
- **`apt install` works everywhere.** It exists on Debian-style systems. Alpine uses `apk`, and a Mac has neither.
- **A distribution is a program.** It is a bundle of many programs, plus the choices and the people who keep them up to date.

## Check yourself
1. What does a distribution add to the Linux kernel? Answer: the commands and shell, the libraries, an installer, a package manager and a collection of packages.
2. What does a package manager do? Answer: it finds a package, fetches it, checks the other packages it needs and puts its files in place.
3. Why can two Linux computers have a different `sed`? Answer: each distribution picks its own tools, such as GNU or BusyBox.
4. Which file lets a Linux system say what distribution it is? Answer: `/etc/os-release`.

## Where next
You now know that a distribution is a bundle of choices, and that this is one reason the same command differs from machine to machine. The last guide in this series puts it to work: how to find out which tool you really have, in about a minute.

For the common differences one command at a time, see `explain sed` and `explain awk`.
