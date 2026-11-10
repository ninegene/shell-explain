---
schema_version: 1
id: whos-who
title: "Who's Who: What Each Name Really Means"
subtitle: Kernel, complete system, set of tools, standard or product
summary: Sort Unix, POSIX, BSD, FreeBSD, Darwin, macOS, GNU, Linux and BusyBox into what each one actually is.
series: why-so-many-versions
order: 2
level: beginner
reading_minutes: 6
updated: 2026-11-09
tags: [history, unix, linux]
related_commands: [sed, cat]
related_guides: [family-tree, same-command-many-versions, what-is-an-os]
references:
  - https://www.kernel.org/linux.html
  - https://www.gnu.org/gnu/gnu-history.html
  - https://docs.freebsd.org/en/books/handbook/introduction/
  - https://developer.apple.com/library/archive/documentation/Darwin/Conceptual/KernelProgramming/Architecture/Architecture.html
  - https://busybox.net/about.html
  - https://man7.org/linux/man-pages/man7/standards.7.html
---
## Goal
By the end you can say what kind of thing each of these names is: a kernel, a complete system, a set of tools, a standard, or a company's product. People use these names as if they were all the same kind of thing, and they are not.

## The idea in plain words
Names in this area sit on different levels, and mixing the levels causes most of the confusion.

A **kernel** is the core that manages the hardware, as the earlier guide on the operating system explained. A **complete system** is a kernel plus all the tools around it, so you can sit down and use it. A **set of tools** is a collection of commands, such as `ls` and `sed`, that can run on top of a kernel. A **standard** is a written description of how things should behave, so that different systems can agree. A **product** is something a company sells or ships.

One name can mean two levels at once. People say Linux for the kernel, and also for a whole system built around it. This guide says which one it means each time.

The point to keep is this: the tools that you type, such as `grep`, are separate from the kernel under them. The kernel and the tools can come from different places.

## A picture
Read the picture from the top down. Each column is one way to build a system. Every column has a kernel at the bottom and tools above it, but the pieces come from different families.

```diagram
   macOS             Linux with GNU tools      a small device
   -----------       -------------------       ----------------
   graphical layers  GNU tools                 BusyBox (one tiny
                     (ls, sed, grep, ...)      program with
   BSD tools         ----------------------    tiny versions of
   (ls, sed, ...)    Linux kernel              many tools)
   -----------                                 ----------------
   Darwin kernel                               Linux kernel
```

Diagram: three columns, each a way to build a system. The first is macOS, with graphical layers on top, then BSD tools such as ls and sed, then the Darwin kernel at the bottom. The second is Linux with GNU tools, with the GNU tools such as ls, sed and grep on top and the Linux kernel at the bottom. The third is a small device, with BusyBox on top, a single tiny program holding small versions of many tools, and the Linux kernel at the bottom.

## Another way to see it
Think of a car. The kernel is the engine. The complete system is the whole car. The tools are the dashboard controls you use every day. A standard is the road code that every car must respect. And a product is a particular brand of car that a company sells.

| Name | What kind of thing it is | Who |
| --- | --- | --- |
| Unix | The original system, and a family | AT&T Bell Labs, then many |
| POSIX | A standard | Ratified by IEEE and ISO/IEC |
| BSD | Releases, and a family of systems | UC Berkeley |
| FreeBSD | A complete system | The FreeBSD Project |
| Darwin | The open core of macOS | Apple |
| macOS | A product that adds to Darwin's core | Apple |
| GNU | A project, and its set of tools | Richard Stallman and others |
| Linux | A kernel | Linus Torvalds and others |
| BusyBox | One program that holds many tools | Maintained by Denys Vlasenko |

The last column names who started each one. For BusyBox it names the maintainer listed on its own site. Many people have worked on each. Apple's description of Darwin and macOS that this guide relies on is archived documentation from when the system was called OS X, so check Apple's current pages for today's details.

## Walk through it
Here is how to work out what a computer is running, one step at a time.

1. Ask for the kernel's name with `uname -s`. A Mac says `Darwin`. A Linux system says `Linux`. This tells you about the bottom layer only.
2. If the answer is `Darwin`, the system is macOS, or Apple's open-source Darwin on its own. Its command line tools come from the BSD side of the family, on top of the Darwin kernel.
3. If the answer is `Linux`, you know the kernel but not the tools. The tools could be the GNU ones, or BusyBox, or something else.
4. To learn which tools you have, ask a tool. Run `sed --version`. In this project's tests, a GNU `sed` (version 4.9) printed a line that began `sed (GNU sed)`. The BusyBox `sed` (version 1.37.0) printed `This is not GNU sed version 4.0`. The `sed` on macOS rejected the option with an error. That is three different answers to one question.
5. Because the answer comes from asking the tool and not from guessing, this project's `explain` command works the same way. It looks at the program itself before it picks the right documentation.

The kernel's name never settles the question of which tools you have. That is why Linux alone does not tell you what `sed` will do.

## Try it
These commands only look at things. Each one works on its own.

Ask for the kernel's name and use it to decide what to print. The `test` command checks the text and sets the exit status that `&&` and `||` use, as in the guide on exit codes. You will see one of the two sentences, depending on your computer. The backslash at the end of a line tells the shell that the command goes on:

```sh
test "$(uname -s)" = Darwin \
  && echo "kernel: Darwin" \
  || echo "kernel: not Darwin"
```

Ask `sed` for its version, and keep only the first line. A GNU `sed` names itself. The BusyBox one says that it is not GNU `sed`. On macOS you will see an error about an illegal option, and a usage line, because that `sed` has no such option. The answer depends on your system, so there is no fixed output:

```sh
sed --version | head -n 1
```

## Common confusions
- **Linux is a complete operating system.** The name Linux belongs to the kernel. A complete system adds tools. The GNU Project calls a complete system made of its tools and that kernel GNU/Linux.
- **macOS is Unix-based, so it is the same as Linux.** They are separate families. macOS includes the BSD command-line environment, and much of the technology in its core derives from FreeBSD. Linux is a separate kernel written from scratch.
- **BusyBox is a version of Linux.** It is not. It is a single small program that holds tiny versions of many common Unix tools, usually with fewer options than the GNU ones. It runs on top of a kernel, often Linux.
- **POSIX is a program.** It is a standard: a written description. Systems aim to follow it, and Linux aims towards POSIX and Single UNIX Specification compliance.
- **GNU is only a license.** GNU is a project that set out to write a complete Unix-compatible system, and it produced many of the tools you use.

## Check yourself
1. Is Linux a kernel or a complete system? Answer: the name Linux belongs to a kernel. A complete system also needs tools.
2. What is Darwin? Answer: Apple's open core of macOS, with a kernel and the BSD command-line environment.
3. Is POSIX a program or a standard? Answer: a standard.
4. What is BusyBox? Answer: one small program that holds tiny versions of many tools.

## Where next
The next guide, `explain learn same-command-many-versions`, puts this to use: it shows why one command name, such as `sed`, can mean different programs with different options.

To see commands from the different families side by side, look at `explain sed` and `explain cat`.
