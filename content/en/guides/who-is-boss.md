---
schema_version: 1
id: who-is-boss
title: "Firmware, Kernel, Programs: Who Is the Boss?"
subtitle: Hardware, firmware, the operating system and software, and who starts whom
summary: Follow a computer from the power button to a running system, and learn which layer starts which.
series: how-a-computer-works
order: 3
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, boot]
related_commands: []
related_guides: [big-picture, computer-body, what-is-an-os]
references:
  - https://man7.org/linux/man-pages/man7/boot.7.html
  - https://docs.kernel.org/
---
## Goal
By the end you can tell hardware, firmware, the operating system and software apart, and explain how each one starts the next. The details of starting a computer differ between systems, so this guide follows the general shape.

## The idea in plain words
Four kinds of thing sit between you and the electricity.

**Hardware** is the physical machine: the processor, the memory, the disk. It does nothing useful alone.

**Firmware** is small built-in software stored on the hardware itself. When you press the power button, it is the first thing that runs. It checks the machine and looks for something to start.

The **operating system** is the software that manages the computer. Its core is the **kernel**. Once it is running, it is the boss: it decides which program gets the CPU and the memory.

**Software** is any set of instructions, stored as files. The operating system is software, and so is every program you run.

Nobody starts alone. Each layer starts the next one, and then the kernel takes charge.

## A picture
Read the picture from the top down. Each step starts the next one.

```diagram
  POWER ON
     |
     v
  FIRMWARE          built into the hardware, runs first
     |              checks the machine, looks for something to start
     v
  BOOT LOADER       a small program found on the disk
     |              finds the kernel, loads it, starts it
     v
  KERNEL            sets up the hardware, then starts the first program
     |              from now on it manages everything
     v
  FIRST PROGRAM     starts the others: logins, services, your shell
     |
     v
  YOUR PROGRAMS     the terminal, the shell, ls, a browser ...
```

Diagram: a chain of six steps from top to bottom. Power on starts the firmware, the firmware starts the boot loader, the boot loader starts the kernel, the kernel starts the first program, and the first program starts your programs.

## Another way to see it
Think of opening a shop in the morning.

| Layer | At the shop |
| --- | --- |
| Hardware | The building, the shelves and the tills |
| Firmware | The alarm panel: runs first and unlocks the door |
| Boot loader | The key holder who lets the manager in |
| Kernel | The manager, who then runs the shop all day |
| Programs | The staff, each doing one job |

The alarm panel does not run the shop. It only gets the day started, and then the manager takes over.

## Walk through it
Here is what happens, step by step, when you switch on a computer. The details vary between systems.

1. You press the power button, and electricity reaches the hardware.
2. Control passes to a program stored in read-only memory on the hardware. This is the firmware, and people often call it the BIOS.
3. The firmware runs a basic self-test. Then it looks for a device that holds a boot loader.
4. The firmware loads the boot loader and hands control to it. On Linux, a common boot loader is called GRUB.
5. The boot loader finds the kernel on the disk, loads it into memory and starts it.
6. The kernel sets up the hardware and the memory, then makes the first program. That program is given the process number 1.
7. The first program starts everything else, including the program that lets you log in and open a shell.

Step 6 uses a word you have not met yet. A process is a running program. The guide on processes explains it.

## Try it
These commands only look at things and change nothing.

Print the name of the operating system's kernel:

```sh
uname -s
```

Print the kernel's version:

```sh
uname -r
```

The first command shows a name such as `Linux` or `Darwin`. The second shows a version number. They report on the kernel, the boss layer, and not on the programs you run on top of it.

## Common confusions
- **The BIOS is the operating system.** It is not. It is firmware: a small program that wakes the machine and hands over to the next step.
- **The kernel is the whole operating system.** It is the core. An operating system also includes the programs that come with it.
- **Software means programs I install.** The operating system is software too, and so is the boot loader.
- **Firmware and software are unrelated.** Firmware is software. It is just stored on the hardware itself, so it is there before anything else runs.

## Check yourself
1. What runs first when you press the power button? Answer: the firmware.
2. Which layer finds the kernel and loads it? Answer: the boot loader.
3. Once the system is running, which layer decides who gets the CPU? Answer: the kernel.
4. Is the operating system software? Answer: yes. All instructions stored as files are software.

## Where next
The next guide, `explain learn what-is-an-os`, looks at what the kernel does all day after it takes charge.

Run `explain glossary kernel` or `explain glossary firmware` to see short definitions of these words.
