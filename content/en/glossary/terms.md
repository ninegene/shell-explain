---
schema_version: 1
id: glossary
title: Words You Will See
summary: Short definitions of the words beginners meet when they start using a terminal.
updated: 2026-11-09
---

## Absolute path
A path that starts with `/`, so it names a file or folder from the top of the tree, wherever you are.

See: files-and-folders

## Argument
A word you type after a command's name, such as a file name or an option. In everyday talk every word after the name is an argument.

See: input-output-errors, anatomy-of-a-command

## ASCII
An old, simple encoding that gives each of 128 characters a number. The letter A is 65. UTF-8 keeps the same numbers for these characters.

See: text-and-binary

## Binary file
A file whose bytes are not meant to be read as characters, such as a photo, a program or a compressed archive.

See: text-and-binary

## Boot
To start a computer: the steps from pressing the power button until the operating system is running.

See: who-is-boss

## Boot loader
A small program that finds the kernel on the disk, loads it into memory and starts it. On Linux, GRUB is a common boot loader.

See: who-is-boss

## BSD
Short for Berkeley Software Distribution: the versions of Unix released by the University of California at Berkeley, from 3BSD in 1980 to 4.4BSD in 1993, and the family of systems derived from them.

See: family-tree, whos-who

## BusyBox
One small program that holds tiny versions of many common Unix tools. It is used on small and embedded systems, and its tools generally have fewer options than the GNU ones.

See: whos-who, same-command-many-versions

## Byte
A number from 0 to 255. A file is a row of bytes.

See: text-and-binary

## Command
A short instruction you type into the shell. It is usually the name of a program, followed by extra words that tell the program what to do.

See: big-picture, terminal-shell-command

## Command line
The line you type into the shell: a command name followed by its options and arguments.

See: anatomy-of-a-command, terminal-shell-command

## CPU
The part of the computer that carries out instructions, one tiny step at a time, very fast. It is also called the processor.

See: big-picture, computer-body

## Current directory
The folder the shell is working in. A relative path starts from it. `pwd` prints it, and `cd` changes it.

See: files-and-folders

## Darwin
The open core of macOS, according to Apple: its kernel and the BSD command-line environment, with technology derived from FreeBSD, Mach and NeXT.

See: whos-who, family-tree

## Directory
A folder: a place that holds files and other directories.

See: files-and-folders

## Disk
Long-term storage that keeps your files when the computer is turned off, such as a hard disk or a solid-state drive (SSD). It is slower than memory.

See: big-picture, computer-body

## Driver
The part of the kernel that knows how to work one kind of hardware, such as a disk or a keyboard.

See: what-is-an-os

## Encoding
The rule that says which number stands for which character.

See: text-and-binary

## Environment variable
A named piece of text that a program carries and hands to the programs it starts, such as PATH or HOME.

See: path-and-environment

## Executable file
A file that has the execute permission turned on, so the system will run it as a program.

See: permissions

## Exit status
The number a command hands back to the shell when it ends. Zero means success, and any other number means something went wrong.

See: exit-codes

## Firmware
Small built-in software stored on the hardware itself. It starts the computer when you switch it on and hands control to the operating system.

See: big-picture, who-is-boss

## Flag
Another word for an option.

See: anatomy-of-a-command

## FreeBSD
A complete Unix-like operating system based on the 4.4BSD-Lite release from the University of California at Berkeley.

See: family-tree, whos-who

## GNU
A project announced by Richard Stallman in 1983, with work starting in January 1984, to write a complete free system compatible with Unix. It also names the large set of tools the project produced.

See: family-tree, whos-who

## Hardware
The physical parts of a computer that you can touch: the processor, the memory, the disk, the keyboard and the screen.

See: big-picture, computer-body

## Home directory
Your personal folder. In the shell, `~` stands for it.

See: files-and-folders

## Implementation
One particular program that does the job a command name stands for. GNU sed, macOS sed and BusyBox sed are three implementations of sed.

See: same-command-many-versions, whos-who

## Kernel
The core of the operating system. It decides which program gets the CPU and memory, and it talks to the hardware on behalf of programs.

See: big-picture, what-is-an-os, who-is-boss

## Linux
A Unix-like kernel written from scratch by Linus Torvalds, first developed in 1991 and made free software in 1992. The GNU Project calls a complete system built from GNU tools and this kernel GNU/Linux.

See: family-tree, whos-who

## macOS
Apple's operating system for the Mac. It is built on Darwin and adds graphical layers on top.

See: whos-who, family-tree

## Memory
Fast, temporary storage for the programs and data the computer is working on right now. What it holds is lost when the power goes off.

See: big-picture, computer-body

## Operating system
The software that manages the computer for all other programs. It shares out the CPU, memory and disk, and it provides files and a way to run programs. Examples include macOS, Windows, and the many systems built around the Linux kernel.

See: big-picture, what-is-an-os

## Option
A word that starts with a dash, such as `-l`, and changes how a command works. Some options take a value of their own.

See: anatomy-of-a-command

## Parent process
The process that started another one. The one it started is its child.

See: processes

## Path
The way to reach a file or folder, written as names separated by `/`.

See: files-and-folders

## PATH variable
The environment variable that lists the folders, separated by colons, where the shell looks for a program by name.

See: path-and-environment

## Permission
A setting that says whether the owner, the group or everyone else may read, write or execute a file.

See: permissions

## PID
Short for process ID: the number the kernel gives each running process so it can be told apart from the others.

See: processes

## Pipe
A one-way channel that connects the standard output of one command to the standard input of the next. In the shell it is written with the `|` symbol.

See: pipes-and-redirects

## Pipeline
Commands joined with pipes, so that the output of each one becomes the input of the next.

See: pipes-and-redirects

## POSIX
A standard that describes what a Unix-like system should provide, including the commands and their options. The first, POSIX.1-1988, was followed by POSIX.1-1990. Richard Stallman coined the name.

See: family-tree, same-command-many-versions

## Process
A program that is running. It has its own PID and its own memory.

See: processes

## Program
A set of instructions the computer can run to do one job, such as listing files. A command usually starts a program.

See: big-picture, processes

## Prompt
The text the shell prints to show that it is ready for your next command.

See: terminal-shell-command

## RAM
Short for random-access memory. It is the computer's memory: fast, temporary storage that is lost when the power goes off.

See: computer-body

## Redirection
Sending one of a program's streams to or from a file instead of the keyboard and screen, for example with `>` or `<`.

See: pipes-and-redirects, input-output-errors

## Relative path
A path that does not start with `/`. It starts from the current directory.

See: files-and-folders

## Root directory
The top folder, written `/`. Every other file and folder is inside it.

See: files-and-folders

## Shell
A program that reads the commands you type and starts the programs they name. Bash and zsh are shells.

See: big-picture, terminal-shell-command

## Software
Instructions, stored as files, that tell a computer what to do. The operating system and every program you run are software.

See: who-is-boss

## stderr
Short for standard error: the stream a program writes its error messages and warnings to. It is kept apart from standard output so that results stay clean when they are saved or piped, and so that errors are still seen. By default it goes to your screen.

See: input-output-errors

## stdin
Short for standard input: the stream a program reads from. By default it comes from your keyboard.

See: input-output-errors

## stdout
Short for standard output: the stream a program writes its normal results to. A pipe or a `>` redirect carries this stream, and not standard error. By default it goes to your screen.

See: input-output-errors

## Stream
A flow of data, one piece after another. Every program starts with three streams: input, output and error.

See: input-output-errors

## Superuser
An account that may read and write any file and skips most permission checks. It is often called root.

See: permissions

## System call
A request a program makes to the kernel, such as asking it to open a file. It is how a program asks the kernel to do something for it.

See: what-is-an-os

## System V
The version of Unix sold by AT&T, first released in 1983. With BSD, it is one of the two main dialects Unix split into after Version 7.

See: family-tree

## Terminal
A window or screen that shows text and passes your key presses to the shell.

See: big-picture, terminal-shell-command

## Text file
A file whose bytes stand for characters, arranged in lines.

See: text-and-binary

## Unix
The operating system from AT&T Bell Labs, whose Version 7 was released in 1979, and the family of systems descended from it.

See: family-tree

## Unix-like
Behaving like Unix without being a direct descendant. Linux, a clone written from scratch, and the GNU system, written to be compatible with Unix, are Unix-like.

See: family-tree, whos-who

## UTF-8
The common way to write text as bytes today. ASCII characters keep their numbers, and other characters take two or more bytes.

See: text-and-binary

## Variable
A name that holds a piece of text, such as PATH.

See: path-and-environment

## Volatile
Describes memory that loses its contents when the power goes off.

See: computer-body
