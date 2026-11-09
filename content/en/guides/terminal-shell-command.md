---
schema_version: 1
id: terminal-shell-command
title: "Terminal, Shell, Command: Three Things People Mix Up"
subtitle: What each one is, and how to tell which one a name refers to
summary: Tell the terminal, the shell and a command apart, and see how they work together every time you type a line.
series: terminal-shell-commands
order: 1
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, terminal, shell]
related_commands: [cd, ls]
related_guides: [big-picture, processes, anatomy-of-a-command]
references:
  - https://man7.org/linux/man-pages/man7/pty.7.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
  - https://www.gnu.org/software/bash/manual/
---
## Goal
By the end you can say what a terminal, a shell and a command each are, and tell which one a name like bash refers to. People use the three words as if they meant the same thing, so this guide takes them apart.

## The idea in plain words
Three different things work together every time you type a line.

The **terminal** is the window, or the screen, where text appears and where your key presses go in. It does not understand commands. It only shows text and passes your typing on.

The **shell** is a program that runs inside the terminal. It reads the line you type, works out what you want, and starts the right program. The standard describes the shell as a command language interpreter. Bash and zsh are shells.

A **command** is what you ask for: the line you type. Usually its first word is the name of a program, such as `ls`. A few commands, such as `cd`, are handled by the shell itself.

A terminal today is usually a program in a window. To the shell, it looks like the old kind of terminal, a screen with a keyboard. The shell does not need to know the difference.

## A picture
Read the picture from the top down. The terminal is the outer window, and the shell runs inside it.

```diagram
  +------------------ TERMINAL (the window) ------------------+
  |                                                           |
  |   you type here  ---->  SHELL  (bash, zsh, ...)           |
  |                           |                               |
  |                           |  starts                       |
  |                           v                               |
  |                  the program for your COMMAND             |
  |                  (ls, grep, awk, ...)                     |
  |                           |                               |
  |   results appear here <---+                               |
  |                                                           |
  +-----------------------------------------------------------+
```

Diagram: a large box labelled terminal, the window. Inside it, your typing goes to the shell. The shell starts the program for your command, such as ls, grep or awk. The results come back from that program and appear in the terminal window.

## Another way to see it
Think of a phone call to a travel agent. The phone is the terminal: it carries your voice and plays back the answer, but it does not know about travel. The agent is the shell: the one who understands your request and arranges it. What you ask for, such as book me a train, is the command.

| Word | What it is | Example |
| --- | --- | --- |
| Terminal | The window for typing and reading | A terminal window |
| Shell | The program that reads and runs your line | bash, zsh |
| Command | What you type to ask for something | `ls -l` |
| Prompt | The shell's sign that it is ready | `$` or `%` |

The prompt is the text the shell prints before the place where you type. It looks different on different computers.

## Walk through it
Here is what happens, step by step, from opening a terminal to seeing a result.

1. You open a terminal window. The terminal program starts a shell inside it.
2. The shell prints a prompt to say that it is ready.
3. You type `ls` and press Enter. The terminal passes your key presses to the shell.
4. The shell reads the line. It sees one word, so the command is `ls`.
5. The shell finds the program called `ls` and starts it as a new process.
6. The program writes the names of your files. The terminal shows them in the window.
7. The program ends. The shell prints a new prompt.

Closing the terminal window usually ends the shell, and any program started from it.

## Try it
These commands only look at things, except the last one, which writes a small file called `nothing.txt` in the folder you are in. Run it where that is fine. Each one works on its own.

Ask the shell its own name. In a terminal you will see something like `bash` or `-zsh`, which names the shell you are talking to:

```sh
echo $0
```

Show the kind of terminal the shell thinks it is talking to. This often looks like `xterm-256color`, and it is empty if the shell has not been told:

```sh
echo $TERM
```

Type a word that is not a command. The shell, and not the terminal, is the one that looks the name up and reports that it was not found. The message goes into a file, and the number 127 is the shell's status for a command it could not find:

```sh
nosuchcommand 2> nothing.txt; echo $?
```

You should see this:

```text
127
```

## Common confusions
- **The terminal runs my commands.** It does not. The terminal shows text and passes your typing on. The shell reads the line, and a program does the work.
- **Bash is a terminal.** Bash is a shell. It runs inside a terminal window.
- **Terminal, shell and command line are three names for one thing.** They are different. The command line is the line you type, the shell reads it, and the terminal is where you type it.
- **The shell is only for typing one line at a time.** A shell can also run a file full of commands, which is called a script. The same shell language is used.
- **If I close the terminal, nothing happens.** Usually the shell ends with it, and so do the programs it started.

## Check yourself
1. Which of the three reads the line you typed and starts the program? Answer: the shell.
2. Is bash a terminal, a shell or a command? Answer: a shell.
3. What does the terminal do? Answer: it shows text and passes your key presses to the shell.
4. When you type a word that is not a command, which part says it was not found? Answer: the shell.

## Where next
The next guide, `explain learn anatomy-of-a-command`, takes a command line apart word by word: the name, the options and the arguments.

To see commands that the shell handles itself, look at `explain cd`. For a program that is started for you, look at `explain ls`.
