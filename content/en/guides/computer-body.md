---
schema_version: 1
id: computer-body
title: "The Computer's Body: CPU, Memory, Disk"
subtitle: What each part does, and why your files live in two places
summary: Meet the three parts that do the work of every command, and see why a computer needs both memory and a disk.
series: how-a-computer-works
order: 2
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, hardware]
related_commands: [ls]
related_guides: [big-picture, who-is-boss]
references:
  - https://docs.kernel.org/admin-guide/mm/concepts.html
---
## Goal
By the end you can say what the CPU, the memory and the disk each do, and explain why a computer keeps things in two different places. This guide is about the bottom layer from the big picture: the hardware.

## The idea in plain words
Three parts do almost all the work behind a command.

The **CPU** (central processing unit, also called the processor) carries out instructions, one tiny step at a time, very fast. It is the worker.

The **memory** holds what the computer is using right now: the programs that are running and the data they are working on. It is fast, but it is **volatile**, which means it loses everything when the power goes off. It is the desk you work at.

The **disk** keeps your files for the long term, such as a hard disk or a solid-state drive (SSD). It keeps them when the computer is switched off. It is slower than memory. It is the filing cabinet.

A program waits on the disk until you run it. Then the computer copies it into memory, and the CPU runs it from there.

## A picture
Read the picture from left to right: things start on the disk, move into memory, and the CPU works on them.

```diagram
  DISK                  MEMORY                 CPU
  the filing cabinet    the desk               the worker
  keeps files when      holds what is in       does the steps,
  the power is off      use right now          one after another

  +-----------+         +-----------+
  | a program |  copy   | a program |  read   +-------------+
  | your file | ------> | file data | ------> | run the     |
  | ...       |         |           | <------ | steps       |
  +-----------+         +-----------+  write  +-------------+
   slower, lasting       faster, forgets        very fast
```

Diagram: three parts side by side. On the left is the disk, which keeps files and is slower. In the middle is the memory, which holds what is in use and is faster. On the right is the CPU, which reads from and writes to memory. Arrows show a program and file data being copied from the disk to memory, and the CPU reading and writing the data in memory.

## Another way to see it
Think of cooking. The disk is the pantry, the memory is the counter, and the CPU is you chopping and stirring.

| Part | Its job | Kept when power is off? |
| --- | --- | --- |
| CPU | Does the steps of a program | It is not for keeping things |
| Memory | Holds what is in use now | No |
| Disk | Stores files for later | Yes |

You cannot cook from the pantry shelf. You bring the ingredients to the counter first, and you put the leftovers back when you are done.

## Walk through it
Here is what happens, step by step, when you run a command and save its result.

1. You type a command. The shell asks the operating system to start the program.
2. The operating system copies the program from the disk into memory.
3. The CPU reads the program from memory and carries out its steps.
4. The program reads a file. The data comes from the disk into memory first, so the CPU can use it.
5. The program works on the data, and the results sit in memory.
6. You save, or the program writes to a file. Only now are the results copied from memory back to the disk.
7. You switch the computer off. Everything in memory disappears. Everything on the disk is still there tomorrow.

The kernel documentation describes the same idea: memory is volatile, and data usually gets into memory by reading it from files.

## Try it
These commands only look at things and change nothing.

See how full the disk holding your current folder is. The numbers will differ on every computer:

```sh
df -h .
```

See how much disk space your current folder takes up:

```sh
du -sh .
```

The first command asks about a whole disk, and the second about one folder. Both numbers describe the filing cabinet, not the desk.

## Common confusions
- **Memory and storage are the same.** They are not. In everyday talk the word memory sometimes means the disk. In a terminal, memory means the fast, temporary part, and the disk is where files are kept.
- **Closing a window saves my work.** It does not. Work only reaches the disk when you save it.
- **A faster CPU makes everything faster.** Not always. A CPU that waits for slow disk reads is a worker waiting at the pantry.
- **The disk is slow because it is old.** It is slower by design: it trades speed for keeping things safely when the power is off.

## Check yourself
1. Which part forgets everything when the power goes off? Answer: the memory, because it is volatile.
2. Where does a program wait before you run it? Answer: on the disk.
3. Which part carries out the steps of a program? Answer: the CPU.
4. Why is your unsaved work lost in a power cut? Answer: it was only in memory, and the disk never received it.

## Where next
The next guide, `explain learn who-is-boss`, shows what wakes these three parts up when you press the power button, and who is in charge once they are running.

To look at the disk side from a terminal, see the command page for `ls` with `explain ls`.
