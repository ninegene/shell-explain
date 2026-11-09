---
schema_version: 1
id: what-is-an-os
title: What an Operating System Really Does
subtitle: Why programs ask the kernel instead of touching the hardware
summary: See the jobs the kernel does for every program, and why a program must ask it before it can open a file or use memory.
series: how-a-computer-works
order: 4
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, kernel]
related_commands: [cat, ls]
related_guides: [big-picture, who-is-boss, processes]
references:
  - https://man7.org/linux/man-pages/man2/syscalls.2.html
  - https://docs.kernel.org/admin-guide/mm/concepts.html
  - https://man7.org/linux/man-pages/man7/boot.7.html
---
## Goal
By the end you can list the jobs an operating system does for every program, and explain why programs ask the kernel for things instead of reaching for the hardware themselves.

## The idea in plain words
Imagine fifty programs sharing one computer. If each one grabbed the processor, the memory and the disk for itself, they would trample each other's work. Someone has to keep order. That is the job of the operating system, and its core, the kernel.

The kernel does four jobs for every program.

It shares out the **CPU**, so each program gets a turn. It gives each program **memory** and keeps programs from reading each other's memory. It looks after **files** on the disk. And it **talks to the hardware** through **drivers**, the parts of the kernel that know how to work one kind of device.

A program never reaches past the kernel. When it needs something, it makes a **system call**: a request to the kernel, such as open this file. The system call is the fundamental interface between a program and the kernel.

## A picture
Read the picture from the top down. The line in the middle is the border that programs may only cross by asking.

```diagram
   PROGRAMS     ls    grep    your shell    a browser
                 |      |         |             |
   - - - - - - -  system calls: "open", "read", "write"  - - - - -
                 |      |         |             |
                 v      v         v             v
   KERNEL       shares the CPU - gives memory - manages files
                 |
                 |  drivers: one for each kind of device
                 v
   HARDWARE     CPU    memory    disk    keyboard    screen
```

Diagram: three bands stacked from top to bottom. At the top are programs such as ls, grep, a shell and a browser. A dashed line below them is the border where system calls such as open, read and write cross into the kernel. In the middle is the kernel, which shares the CPU, gives out memory and manages files. At the bottom is the hardware, reached only through the kernel's drivers.

## Another way to see it
Think of a bank. You, the customer, never walk into the vault. You ask a teller, and the teller does the work.

| Job | Without an OS | With one |
| --- | --- | --- |
| Running | Each hogs the CPU | Each gets a turn |
| Memory | Programs overwrite each other | Each has its own space |
| Files | Each learns the disk layout | One service handles files |
| Devices | Each needs every driver | Kernel drivers do it once |

## Walk through it
Here is what happens, step by step, when you run `cat hello.txt` to print a small file.

1. The shell asks the kernel to start the program `cat`.
2. `cat` asks the kernel, with a system call, to open the file `hello.txt`.
3. The kernel checks that the file exists and that you may read it, then hands `cat` a number that stands for the open file.
4. `cat` asks the kernel to read from that file.
5. The kernel gets the bytes from the disk through its disk driver and copies them into memory that belongs to `cat`.
6. `cat` asks the kernel to write the text to its output.
7. The kernel passes the text on to your terminal, which shows it on screen.

At no step does `cat` touch the disk or the screen itself.

## Try it
These commands make a small folder and a small file to work on. Run them in a place where a new folder called `practice` is fine. They delete nothing.

Make a folder, make a file in it, and print the file. Each line sends requests to the kernel:

```sh
mkdir practice
cd practice
echo hello > hello.txt
cat hello.txt
```

You should see this:

```text
hello
```

Making the folder, writing the file and reading it back were all system calls. You never had to say where on the disk the file lives. The kernel decided that.

## Common confusions
- **The operating system is what I see on the screen.** The desktop and windows are programs that run on top of it. The operating system is the layer underneath that they all ask for help.
- **Programs talk to the hardware directly.** They do not. They ask the kernel, and the kernel's drivers talk to the hardware.
- **The kernel is a program I can start.** It starts before anything else, as the previous guide showed, and it keeps running underneath everything.
- **Linux is a whole operating system.** Linux is the name of a kernel. The complete systems built around it are covered later in this series.

## Check yourself
1. Name two jobs the kernel does for every program. Answer: for example, sharing out the CPU, and giving each program its own memory.
2. What is a system call? Answer: a request a program makes to the kernel.
3. What is a driver? Answer: the part of the kernel that knows how to work one kind of hardware.
4. When `cat` prints a file, who reads the bytes from the disk? Answer: the kernel does, on behalf of `cat`.

## Where next
The next guide, `explain learn processes`, looks at how the kernel keeps track of every running program.

To try more file commands, see `explain cat` and `explain ls`.
