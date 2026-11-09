---
schema_version: 1
id: big-picture
title: From Keypress to Screen
subtitle: The big picture of what happens when you type a command
summary: Follow one command from your fingers, through the layers of the computer, to the text on your screen.
series: how-a-computer-works
order: 1
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, overview]
related_commands: [ls, cat]
related_guides: []
references:
  - https://www.gnu.org/software/bash/manual/
  - https://docs.kernel.org/
---
## Goal
By the end you can name the layers between your fingers and the text on your screen, and say what each layer does. This guide looks at the whole stack at once. Later guides zoom in on one layer each.

## The idea in plain words
When you type a command, nothing magic happens. Your request passes through a chain of helpers. Each helper does one small job and hands the work to the next one.

A **command** is a short instruction you type, such as `ls`, which lists the files in a folder. The **terminal** is the window where you type. The **shell** is the helper that reads what you typed and starts the right program. A **program** is a set of instructions that does one job. The **operating system** is the software that manages the computer for every program, and its core is called the **kernel**. The **hardware** is the physical machine: the processor (CPU), the memory, and the disk.

## A picture
Read the picture from the top down. Your typing goes down through the layers, and the answer comes back up the same way.

```diagram
 YOU                      type:  ls
   |  keys
   v
 KEYBOARD and SCREEN      the hardware you touch
   |
   v
 TERMINAL                 a window that shows text, passes keys on
   |  text
   v
 SHELL                    reads your line, finds the program, starts it
   |  "run ls"            (bash, zsh, ...)
   v
 PROGRAM (command)        ls, grep, awk ... does one job
   |  asks politely: "list this folder"   (system calls)
   v
 OPERATING SYSTEM         the kernel: shares out CPU, memory, disk
   |
   v
 HARDWARE                 CPU, memory, disk (firmware starts it)
```

Diagram: seven layers stacked from top to bottom. You are at the top. Under you are the keyboard and screen, then the terminal, then the shell, then the program, then the operating system, and the hardware is at the bottom.

## Another way to see it
Think of a restaurant. Each layer has a partner there.

| Layer | In a restaurant |
| --- | --- |
| You | The customer with an order |
| Terminal | The ordering tablet: shows the menu, sends what you tap |
| Shell | The waiter: understands the order, tells the right cook |
| Program | A specialist cook: one dish done well |
| Kernel | The kitchen manager: shares the stoves and fridges fairly |
| Hardware | The stoves, fridges and counters |

An analogy is only a guide. One waiter serves many tables at once, but a shell handles one line of typing at a time.

## Walk through it
Here is what happens, step by step, when you type `ls` and press Enter.

1. You press the keys l, s and Enter. The keyboard sends signals to the computer.
2. The terminal passes your key presses to the shell and shows the letters on the screen.
3. The shell reads the line. It sees one word, `ls`, and looks through a list of folders (called the PATH) for a program with that name.
4. The shell asks the kernel to start that program.
5. The kernel loads the program from the disk into memory and gives it a turn on the CPU.
6. The program asks the kernel which files are in the current folder. The kernel reads the disk and answers.
7. The program writes the names as its output.
8. The terminal draws the names on your screen. The shell shows its prompt again, ready for the next command.

## Try it
These commands are safe: they only look at things and change nothing.

List the files in the folder you are in. Your names will differ from anyone else's:

```sh
ls
```

Make the shell print a word back to you. The shell starts the work itself here, and the terminal shows the result:

```sh
echo hello
```

You should see this:

```text
hello
```

Show the path of the folder you are in:

```sh
pwd
```

## Common confusions
- **The terminal runs my commands.** It does not. The terminal only shows text and passes your key presses on. The shell reads the line, and a program does the work.
- **Terminal, shell and command are the same thing.** They are three different things that work together. The next series takes them apart one by one.
- **The shell is the operating system.** It is not. The shell is one program that runs on top of the operating system.
- **Bash is a terminal.** Bash is a shell. It runs inside a terminal window.

## Check yourself
1. Which layer reads the line you typed and starts the program? Answer: the shell.
2. Which layer decides which program gets the CPU and memory? Answer: the operating system, through its kernel.
3. Where do the names of your files come from when you run `ls`? Answer: the program asks the kernel, and the kernel reads them from the disk.
4. Is the terminal the thing that does the work of your command? Answer: no. It shows text and passes your keys on.

## Where next
The next guides zoom in on each layer: the parts of the computer, what an operating system does, and how the terminal, the shell and a command differ.

To see real commands, open the command pages. For example, run `explain ls` or `explain cat` in your terminal.
