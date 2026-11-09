---
schema_version: 1
id: input-output-errors
title: Three Streams Every Program Has
subtitle: Standard input, standard output and standard error, and where each one goes
summary: Learn the three channels that every program starts with, where each one goes by default, and why errors travel on their own channel instead of mixing with results.
series: terminal-shell-commands
order: 5
level: beginner
reading_minutes: 8
updated: 2026-11-08
tags: [basics, streams, input-output]
related_commands: [cat, grep, sort]
related_guides: [what-is-an-os, processes, pipes-and-redirects]
references:
  - https://man7.org/linux/man-pages/man3/stdin.3.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
---
## Goal
By the end you can name the three streams every program has, say where each one goes by default, and explain why a program keeps its error messages apart from its results.

## The idea in plain words
A **stream** is a flow of data, one piece after another, like water through a tap. When a program starts, it already has three streams open.

**Standard input**, also written **stdin**, is where the program reads from. By default it comes from your keyboard.

**Standard output**, or **stdout**, is where the program writes its normal results. By default it goes to your screen.

**Standard error**, or **stderr**, is where the program writes error messages and warnings. By default it also goes to your screen.

The two outputs are separate on purpose. Results and complaints can be handled differently, for example by saving the results in a file while the complaints still reach you.

The streams are only attached to the keyboard and screen by default. The program that starts another one, such as your shell, can point them somewhere else: at a file, or at another program. The program does not need to know.

### Why is there a separate error stream?
Imagine a program with only one output. Everything it said would come out in one flow: the results you asked for, and also its complaints, warnings and progress notes.

That works while the flow goes to your screen. You read it all and sort it out in your head. But the shell can send the flow somewhere else, into a file or into another program. Then the mixture causes trouble.

- **The results get dirty.** Save the output in a file, and the file holds error messages between the real lines. Another program reading that file would treat a complaint as data.
- **The errors can vanish.** If the flow goes into a file or a pipe, you do not see it. Something could go wrong and nobody would notice.
- **You cannot handle them differently.** You may want the results in one file and the problems in another. You may want to throw the problems away, or to look only at the problems.

A separate error stream solves all three. The results flow wherever you send them, and they stay clean. The errors keep their own route. It leads to your screen unless you choose otherwise, so you still see them. In short, programs and files read the results, and people read the errors.

The manual page for these streams describes the third one as the stream for diagnostic and error messages. That is the reason it exists.

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

Now see why the separate route matters. Here `cat` reads two files, but one of them does not exist, and its output goes into `sort`:

```diagram
   cat fruit.txt nosuchfile                 sort
  +------------------------+            +-----------+
  |                        | stdout (1) |           |
  |  reads the two files   | ---------> |  sorts    | --> apple, pear
  |                        | (the pipe) |           |
  +------------------------+            +-----------+
             |
             |  stderr (2): does not go into the pipe
             +------------------------------> the error message
                                              (on your screen)
```

Diagram: two boxes joined by a pipe. The first, cat, reads fruit.txt and a file that does not exist. Its standard output, number 1, goes through the pipe into the second box, sort, which writes apple and pear to the screen. Its standard error, number 2, takes a separate route that bypasses the pipe and ends at the screen as an error message. The sorted results contain no error text.

## Another way to see it
Think of a shop counter. You hand your order across, which is the input. You receive your purchases, which is the output. If something is wrong, the clerk calls out a problem, which is the error. The purchases and the problem are different things, even though both reach you.

Or think of a factory belt. The finished products ride the belt to the packing box. A worker who spots a fault does not put a note inside a product box. The worker rings a bell, so the manager hears about it and the boxes stay clean. The belt is standard output, and the bell is standard error.

| Stream | Number | What it carries | By default |
| --- | --- | --- | --- |
| Standard input | 0 | What the program reads | Your keyboard |
| Standard output | 1 | Normal results | Your screen |
| Standard error | 2 | Error messages | Your screen |

The numbers matter later. They are how a command line says which stream it means.

Here is the case for the separate error stream in one table. It compares what happens with one shared output against what happens with a separate error stream.

| When you are... | One shared output | Separate stderr |
| --- | --- | --- |
| Saving results to a file | Errors land inside it | The file stays clean |
| Piping into another program | It reads errors as data | Errors skip the pipe |
| Hiding output from the screen | You miss the errors | You still see errors |
| Keeping a log of problems | Hard to pick them out | Save stream 2 alone |

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

Now follow a command that has both results and an error. The file `fruit.txt` holds the lines `pear` and `apple`, and the second file does not exist. You run `cat fruit.txt nosuchfile | sort`.

1. The shell starts `cat` and `sort` and joins them with a pipe. The standard output of `cat` goes into the pipe. Its standard error stays attached to your terminal.
2. `cat` reads `fruit.txt` and writes `pear` and `apple` to its standard output. They enter the pipe.
3. `cat` tries `nosuchfile` and fails. It writes a message to standard error. The message does not enter the pipe. It goes straight to your screen.
4. `sort` reads `pear` and `apple` from the pipe, puts them in order, and writes `apple` and `pear` to its standard output, which is your screen.
5. You see the sorted results and the error message, each by its own route. The results hold no error text, so whatever reads them next gets only data.

## Try it
These commands write small files in the folder you are in, so run them where that is fine. Each one works on its own.

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

See the separate error stream protect the results. Here `cat` fails on the second file. Its error goes to `errors.txt`, so only the two sorted lines come out of the pipe:

```sh
printf 'pear\napple\n' > fruit.txt
cat fruit.txt nosuchfile 2> errors.txt | sort
```

You should see this:

```text
apple
pear
```

Now see what would happen without it. Count the lines that come out of the pipe, first with the error kept apart, then with it merged into the output by `2>&1`. When the error is merged, its message counts as a third line of results:

```sh
printf 'pear\napple\n' > fruit.txt
cat fruit.txt nosuchfile 2> errors.txt | awk 'END {print NR}'
cat fruit.txt nosuchfile 2>&1 | awk 'END {print NR}'
```

You should see this:

```text
2
3
```

The first count is the two real lines. The second count includes the error message, which is exactly the mess that a separate error stream prevents.

## Common confusions
- **Input always means the keyboard.** That is only the default. A program's input can come from a file or from another program.
- **Output and error are the same.** They are two streams. Both reach the screen by default, which hides the difference, but you can send one somewhere else and keep the other.
- **Error output means the program failed.** Not always. Programs also use standard error for warnings and progress messages, and a program can write to it and still succeed.
- **A program could just print errors to standard output.** Then the errors would be saved in your files and fed into the next program as if they were results, and they would disappear when you hide the output. The separate stream avoids all three.
- **If I send output to a file, I lose the error messages.** You do not. Only standard output moves. Standard error still goes to your screen unless you redirect it as well.
- **A program only reads the file I name.** Many programs read standard input when you name no file, as `grep` did above.
- **Output is only what I see on the screen.** The screen is the default place. The output can also go to a file or into another program.

## Check yourself
1. Which stream do a program's normal results go to? Answer: standard output.
2. Which stream do error messages go to? Answer: standard error.
3. Where does standard input come from by default? Answer: your keyboard.
4. Who decides where the three streams are attached? Answer: the program that starts the other one, such as the shell.
5. Name two problems that would appear if errors and results shared one stream. Answer: for example, error messages end up inside saved files, and the next program in a pipeline reads them as data. Errors could also go unseen when the output is hidden.
6. In `cat a.txt b.txt | sort`, does an error message from `cat` go through the pipe? Answer: no. The pipe carries only standard output, and the error goes to your screen.

## Where next
The next guide, `explain learn pipes-and-redirects`, shows how to point these streams at files and at other programs.

To see real commands that read standard input, look at `explain cat`, `explain grep` and `explain sort`.
