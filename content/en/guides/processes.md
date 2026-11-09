---
schema_version: 1
id: processes
title: "Programs That Are Awake: Processes"
subtitle: What a process is, how one starts another, and how to see them
summary: Learn the difference between a program and a process, how the shell starts a process, and how to list the ones running now.
series: how-a-computer-works
order: 5
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, processes]
related_commands: []
related_guides: [big-picture, who-is-boss, what-is-an-os, terminal-shell-command]
references:
  - https://man7.org/linux/man-pages/man7/credentials.7.html
  - https://man7.org/linux/man-pages/man2/fork.2.html
  - https://man7.org/linux/man-pages/man2/execve.2.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/ps.html
---
## Goal
By the end you can explain what a process is, tell it apart from a program, and say how the shell starts one. You can also list the processes on your computer.

## The idea in plain words
A **program** is a file on the disk: a set of instructions that is asleep. A **process** is a program that is awake: it has been started, it is using memory and the CPU, and it has its own number.

You can have one program and many processes. If you open three terminal windows, three shells run, all from the same program file.

The kernel gives every process a **process ID**, or **PID**. A process also knows its **parent process**: the process that started it. The one that was started is its child.

New processes are made by copying. A process asks the kernel to duplicate it, which makes a child that starts out as an exact copy in its own separate memory. Then the child can ask to be replaced by a different program.

## A picture
Read the picture from the top down. The shell starts `ls`, and the arrow shows which process is the parent.

```diagram
   shell                 PID 4000     the parent
     |
     |  1. asks the kernel to copy the shell
     v
   shell (copy)          PID 4001     the child, a copy of the shell
     |
     |  2. asks to be replaced by the program ls
     v
   ls                    PID 4001     same process, now running ls
     |
     |  3. finishes and ends; the parent carries on
     v
   shell                 PID 4000     shows the prompt again
```

Diagram: a shell with process number 4000 asks the kernel to copy it. The copy, number 4001, is first still a shell. It then asks to be replaced by the program ls, keeping the number 4001. When ls finishes, the shell with number 4000 shows the prompt again. The numbers are examples.

## Another way to see it
Think of a recipe book and a kitchen.

| Word | Everyday picture | In the computer |
| --- | --- | --- |
| Program | A recipe in the book | A file on the disk |
| Process | A dish being cooked right now | A program that is running |
| PID | The ticket number on the order | The number the kernel gives a process |
| Parent | The chef who started the dish | The process that started this one |

One recipe can be cooking in three pots at once. Each pot is its own process, with its own ticket number.

## Walk through it
Here is what happens, step by step, when you type `ls` in a shell.

1. The shell reads your line and sees that you want the program `ls`.
2. The shell asks the kernel to duplicate it. The kernel makes a child process and gives it a new PID.
3. The parent shell and the child run in separate memory, and at first the two have the same contents.
4. The child asks the kernel to replace the program it is running with `ls`. It keeps its PID, but now it runs `ls`.
5. `ls` does its work and prints its output.
6. `ls` ends. The parent shell notices, and shows its prompt for the next command.

The very first process, with PID 1, is started by the kernel as the last step of starting up. Every other process descends from it.

## Try it
These commands only look at things and change nothing.

List the processes that belong to your terminal. You will see at least the shell and the `ps` program itself:

```sh
ps
```

Show the PID, the parent's PID and the name of each of those processes. Find a row for `ps`, and see that its parent PID is the PID of your shell:

```sh
ps -o pid,ppid,comm
```

Print the PID of your own shell. The number will differ every time:

```sh
echo $$
```

Start a process that lasts one second and then ends:

```sh
sleep 1
```

## Common confusions
- **Program and process are the same.** A program is the file. A process is that file, running, with its own PID and memory.
- **A PID belongs to a program for ever.** It does not. A PID names one running process, and a new run gets a new number.
- **Closing the terminal window does nothing to the processes.** Usually the shell and the programs started from it end with it.
- **A child process shares memory with its parent.** It starts as a copy, but they run in separate memory, so a change in one does not change the other.

## Check yourself
1. What is the difference between a program and a process? Answer: a program is a file on the disk, and a process is a program that is running.
2. What does the PID tell you? Answer: which running process it is, as a number the kernel gives out.
3. How does the shell start `ls`? Answer: it asks the kernel for a copy of itself, and the copy asks to be replaced by `ls`.
4. Which process has the number 1? Answer: the first process the kernel starts.

## Where next
You now have the whole bottom half of the picture: the hardware, how it starts, what the kernel does, and what a process is. The next series takes apart the terminal, the shell and the commands you type into them. Start with `explain learn terminal-shell-command`, which tells the terminal, the shell and a command apart.

Run `explain glossary process` to see a short definition.
