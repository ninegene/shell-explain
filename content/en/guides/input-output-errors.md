---
schema_version: 1
id: input-output-errors
title: Three Streams Every Program Has
subtitle: Standard input, standard output and standard error, and where each one goes
summary: Learn the three channels that every program starts with, where each one goes by default, and why errors travel separately from results.
series: terminal-shell-commands
order: 5
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, streams, input-output]
related_commands: [cat, grep, sort]
related_guides: [what-is-an-os, processes, pipes-and-redirects]
references:
  - https://man7.org/linux/man-pages/man3/stdin.3.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
---
## Goal
By the end you can name the three streams every program has, say where each one goes by default, and explain why a program writes its results and its error messages to different places.

## The idea in plain words
A **stream** is a flow of data, one piece after another, like water through a tap. When a program starts, it already has three streams open.

**Standard input**, also written **stdin**, is where the program reads from. By default it comes from your keyboard.

**Standard output**, or **stdout**, is where the program writes its normal results. By default it goes to your screen.

**Standard error**, or **stderr**, is where the program writes error messages and warnings. By default it also goes to your screen.

The two outputs are separate on purpose. Results and complaints can be handled differently, for example by saving the results in a file while the complaints still reach you.

The streams are only attached to the keyboard and screen by default. The program that starts another one, such as your shell, can point them somewhere else: at a file, or at another program. The program does not need to know.

## A picture
Read the picture from left to right. Data goes into the program on the left, and comes out on the right.

```diagram
                         +-----------+
   keyboard  --------->  |           |  --------->  screen
              stdin (0)  |  program  |  stdout (1)  (normal results)
                         |           |
                         |           |  --------->  screen
                         +-----------+  stderr (2)  (error messages)
```

Diagram: a box in the middle labelled program. An arrow comes in from the keyboard on the left and is labelled standard input, number 0. Two arrows leave on the right. The upper one is standard output, number 1, for normal results, and the lower one is standard error, number 2, for error messages. Both go to the screen.

## Another way to see it
Think of a shop counter. You hand your order across, which is the input. You receive your purchases, which is the output. If something is wrong, the clerk calls out a problem, which is the error. The purchases and the problem are different things, even though both reach you.

| Stream | Number | What it carries | By default |
| --- | --- | --- | --- |
| Standard input | 0 | What the program reads | Your keyboard |
| Standard output | 1 | Normal results | Your screen |
| Standard error | 2 | Error messages | Your screen |

The numbers matter later. They are how a command line says which stream it means.

## Walk through it
Here is what happens, step by step, when you run `grep apple` and give it no file name.

1. You type `grep apple` and press Enter. The shell starts the program `grep`, and its three streams are attached to your terminal.
2. `grep` looks for a file name among its arguments. An argument is a word you type after the command name. It finds only the word `apple`, which is the pattern to look for.
3. With no file to read, `grep` reads its standard input instead. It waits for you.
4. You type a line, such as `a red apple`, and press Enter. The terminal sends the line to the standard input of `grep`.
5. The line contains `apple`, so `grep` writes it to its standard output. The terminal draws it on your screen.
6. You type a line without `apple`. `grep` reads it and writes nothing.
7. You press Ctrl-D to say that there is no more input. `grep` ends.

If the same program had a problem, such as a file that does not exist, it would write the message to standard error. That also reaches your screen, but along its own path.

## Try it
These commands only look at things. The second one writes a small file called `errors.txt` in the folder you are in, so run it where that is fine.

Give `grep` its input through a pipe instead of the keyboard. The pipe hands what `printf` writes to `grep` as its standard input. The next guide explains pipes:

```sh
printf 'apple\nbanana\n' | grep apple
```

You should see this:

```text
apple
```

Send the error message of a failing command to a file. The error does not appear on the screen, because standard error was pointed at the file:

```sh
cat nosuchfile 2> errors.txt || echo "failed, see errors.txt"
```

You should see this:

```text
failed, see errors.txt
```

Now read the message from the file with `cat errors.txt`. The exact words depend on your system.

## Common confusions
- **Input always means the keyboard.** That is only the default. A program's input can come from a file or from another program.
- **Output and error are the same.** They are two streams. Both reach the screen by default, which hides the difference, but you can send one somewhere else and keep the other.
- **Error output means the program failed.** Not always. Programs also use standard error for warnings and progress messages, and a program can write to it and still succeed.
- **A program only reads the file I name.** Many programs read standard input when you name no file, as `grep` did above.
- **Output is only what I see on the screen.** The screen is the default place. The output can also go to a file or into another program.

## Check yourself
1. Which stream do a program's normal results go to? Answer: standard output.
2. Which stream do error messages go to? Answer: standard error.
3. Where does standard input come from by default? Answer: your keyboard.
4. Who decides where the three streams are attached? Answer: the program that starts the other one, such as the shell.

## Where next
The next guide, `explain learn pipes-and-redirects`, shows how to point these streams at files and at other programs.

To see real commands that read standard input, look at `explain cat`, `explain grep` and `explain sort`.
