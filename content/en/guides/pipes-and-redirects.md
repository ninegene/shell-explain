---
schema_version: 1
id: pipes-and-redirects
title: Snapping Commands Together
subtitle: Redirect a stream into a file, or pipe one command's output into the next
summary: Learn the symbols that send output into files and connect small commands into a pipeline.
series: terminal-shell-commands
order: 6
level: beginner
reading_minutes: 6
updated: 2026-11-08
tags: [basics, pipes, redirection]
related_commands: [cat, grep, sort, head]
related_guides: [input-output-errors, processes, exit-codes]
references:
  - https://man7.org/linux/man-pages/man7/pipe.7.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
---
## Goal
By the end you can send a command's output into a file, read input from a file, and connect several commands into a pipeline. This guide uses the three streams from the previous guide.

## The idea in plain words
Most commands do one small job. The shell lets you snap them together, so the result of one becomes the input of the next.

A **redirect** changes where one of a program's streams goes. You can send standard output into a file, read standard input from a file, or send standard error into a file. The program does not know. The shell sets it up before the program starts.

A **pipe** connects the standard output of one command to the standard input of the next. The pipe is a one-way channel: what the first command writes in at one end comes out at the other end for the second command to read. A line of commands joined with pipes is a **pipeline**.

The symbols are short. Here are the ones to learn first.

| Symbol | What it does | Example |
| --- | --- | --- |
| `>` | Writes output to a file, replacing it | `sort > out.txt` |
| `>>` | Adds output to the end of a file | `echo x >> out.txt` |
| `<` | Reads input from a file | `sort < fruit.txt` |
| `2>` | Writes errors to a file | `cat nofile 2> err.txt` |
| `\|` | Pipes output into the next command | `sort \| head -n 2` |

The `>` symbol replaces the file's old contents without asking, so it needs care. Use `>>` to add to the end instead.

## A picture
Read the picture from left to right. Each command's output is the next command's input, and only the last command's output reaches the screen.

```diagram
   printf             sort               head -n 2
  +--------+  pipe  +--------+  pipe   +----------+
  | writes | -----> | sorts  | ------> | keeps the|  ----->  screen
  | 3 lines|        | them   |         | first two|
  +--------+        +--------+         +----------+
   pear, apple,      apple, fig,        apple, fig
   fig               pear
```

Diagram: three commands in a row joined by pipes. The first, printf, writes three lines: pear, apple and fig. The second, sort, puts them in order: apple, fig, pear. The third, head with the option dash n 2, keeps the first two lines, apple and fig, and those are the only lines that reach the screen.

## Another way to see it
Think of an assembly line. Each worker does one job on what the worker before passed along. A redirect is like putting the finished item into a box on a shelf, which is a file, instead of handing it on.

| Idea | Assembly line | Shell |
| --- | --- | --- |
| Worker | Does one job | A command |
| Hand to the next worker | The conveyor belt | A pipe |
| Put in a box on a shelf | Store the result | A redirect to a file |

## Walk through it
Here is what happens, step by step, when you run `printf 'pear\napple\nfig\n' | sort | head -n 2`.

1. The shell reads the line and sees two pipes, so there are three commands.
2. The shell makes two pipes. Each has a read end and a write end.
3. The shell starts the three commands as three processes, at the same time.
4. It points the standard output of `printf` at the first pipe, and the standard input of `sort` at the other end of it. It does the same for `sort` and `head` with the second pipe.
5. `printf` writes its three lines. `sort` reads them, puts them in order, and writes them on.
6. `head` reads the sorted lines, keeps the first two, and writes them. Its output is not redirected, so it goes to your screen.
7. When the last command ends, the shell shows the prompt again.

Because the commands run together, the second one can start work as soon as the first writes something.

## Try it
These commands write small files in the folder you are in, so run them where that is fine. Each one works on its own.

Send the output of one command into the next:

```sh
printf 'pear\napple\nfig\n' | sort
```

You should see this:

```text
apple
fig
pear
```

Write into a file, then read the file back as input:

```sh
printf 'pear\napple\nfig\n' > fruit.txt
sort < fruit.txt
```

You should see this:

```text
apple
fig
pear
```

Add to a file with `>>`:

```sh
echo one > list.txt
echo two >> list.txt
cat list.txt
```

You should see this:

```text
one
two
```

Chain three commands:

```sh
printf 'pear\napple\nfig\n' | sort | head -n 2
```

You should see this:

```text
apple
fig
```

Send both results and errors to one file. The `2>&1` means: send standard error to the same place as standard output. Put it after the `>` part:

```sh
cat nosuchfile > both.txt 2>&1
test -s both.txt && echo "both.txt has the error message"
```

You should see this:

```text
both.txt has the error message
```

## Common confusions
- **`>` adds to a file.** It replaces the file. Use `>>` to add to the end.
- **A pipe sends a file.** It sends the output text, not a file or a file name. The next command reads that text as its input.
- **A redirect hides errors.** `command > file` only moves standard output. Errors still reach the screen, because standard error is a different stream.
- **A pipe and a redirect are the same.** A redirect connects a stream to a file. A pipe connects it to another program.
- **Any command can go after a pipe.** The command has to read standard input. A command that only works on file names it is given will ignore what the pipe sends.
- **The order of `2>&1` does not matter.** It does. Redirections are applied from left to right, so `> both.txt 2>&1` sends both streams to the file, but `2>&1 > both.txt` sends errors to the screen.

## Check yourself
1. What is the difference between `>` and `>>`? Answer: `>` replaces the file, and `>>` adds to the end.
2. What does a pipe connect? Answer: the standard output of one command to the standard input of the next.
3. How do you save only the error messages of a command in a file? Answer: use `2>` and the file name.
4. In `printf ... | sort | head`, which command's output reaches the screen? Answer: the last one, `head`.

## Where next
The next guide, `explain learn exit-codes`, shows how a command tells the shell whether it worked, which lets you chain commands with `&&` and `||`.

To see more commands that work well in a pipeline, look at `explain sort`, `explain grep`, `explain head` and `explain cat`.
