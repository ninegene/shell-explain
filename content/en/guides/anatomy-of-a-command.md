---
schema_version: 1
id: anatomy-of-a-command
title: Anatomy of a Command
subtitle: The command name, its options and its arguments, word by word
summary: Take a command line apart into its name, options and arguments, and see why quoting changes how many arguments there are.
series: terminal-shell-commands
order: 2
level: beginner
reading_minutes: 6
updated: 2026-11-08
tags: [basics, commands, options]
related_commands: [ls, sort, head]
related_guides: [terminal-shell-command, input-output-errors, files-and-folders]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap12.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
  - https://man7.org/linux/man-pages/man2/execve.2.html
---
## Goal
By the end you can split a command line into its parts: the command name, the options, and the arguments. You can also say why quotes change the number of arguments.

## The idea in plain words
A command line is a list of words separated by spaces. The shell splits your line into those words before it starts anything.

The first word is the **command name**: which program to run.

Words that start with a dash, such as `-l`, are **options**. An option changes how the command works. People also call options **flags**. Some options take a value of their own, as in `-n 2`, where the `2` belongs to the option `-n`.

The remaining words are the **arguments**, in everyday talk. They say what to work on, such as a file name. The official manuals call these words operands, and they call the word after an option its option-argument. This guide says arguments for everything after the command name.

Several short options can be joined behind one dash. `ls -la` means the same as `ls -l -a`. Some commands also accept long options written with two dashes, such as `--help`. These are not the same on every system, so a long option that works on Linux may not exist on macOS.

The shell can be told to keep words together with quotes. `'one two'` is a single word, so it is one argument.

## A picture
Read the picture from left to right. Each word of the line has a role.

```diagram
     head    -n    2    notes.txt
      |      |     |        |
      |      |     |        +-- argument: the file to read
      |      |     +----------- value of the option -n (how many lines)
      |      +----------------- option: changes how head works
      +------------------------ command name: the program to run
```

Diagram: the command line head dash n 2 notes dot txt, with an arrow from each word to its role. Head is the command name, the program to run. Dash n is an option that changes how head works. The 2 is the value of that option, how many lines to show. Notes dot txt is the argument, the file to read.

## Another way to see it
Think of ordering coffee. The command name is what you ask for, a coffee. Options say how: large, no sugar. The arguments say what to use it on or for, such as for here, or the name on the cup.

| Part | What it is | In `ls -l notes.txt` |
| --- | --- | --- |
| Command name | The program to run | `ls` |
| Option | Changes how it works | `-l` |
| Argument | What to work on | `notes.txt` |

An option is a way to ask for a variation of the same command. An argument is the thing the command acts on.

## Walk through it
Here is what happens, step by step, when you run `grep -i error app.log`.

1. The shell splits the line into four words: `grep`, `-i`, `error` and `app.log`.
2. The shell takes the first word as the command name and finds the program `grep`.
3. The shell starts `grep` and hands it the list of words that came after the name. The kernel passes that list to the new program.
4. `grep` looks at the list. It sees that `-i` starts with a dash, so it treats it as an option and ignores the case of letters.
5. The next word, `error`, is the pattern to look for. The last word, `app.log`, is the file to search.
6. `grep` searches the file and writes the matching lines.

The shell only splits the words. The program decides which of them are options. That is why two commands can treat the same letter differently.

## Try it
These commands only look at things. Each one works on its own.

Give `echo` three arguments. It prints them with a space between:

```sh
echo one two three
```

You should see this:

```text
one two three
```

Now show the arguments one per line. The `printf` command here prints each argument on its own line, so the shell's splitting is visible. There are three arguments:

```sh
printf '%s\n' one two three
```

You should see this:

```text
one
two
three
```

Put quotes around the same words. They become one argument:

```sh
printf '%s\n' 'one two three'
```

You should see this:

```text
one two three
```

Use an option. The option `-r` tells `sort` to reverse the order:

```sh
printf 'b\na\nc\n' | sort -r
```

You should see this:

```text
c
b
a
```

Give an option a value. The option `-n 2` tells `head` to show two lines, and `2` is its value:

```sh
printf '1\n2\n3\n4\n' | head -n 2
```

You should see this:

```text
1
2
```

Join two options behind one dash. This checks that `ls -la` and `ls -l -a` print the same thing:

```sh
test "$(ls -la)" = "$(ls -l -a)" && echo same
```

You should see this:

```text
same
```

## Common confusions
- **Options have to come first.** Often they do, but the rule belongs to each command. Read its manual before assuming.
- **A dash word is always an option.** A file can be named with a leading dash, which confuses commands. Most commands accept a lone `--` to mean that the options are over.
- **The quotes are part of the argument.** They are not. The shell removes them, and the program only sees the words.
- **The shell knows what the options mean.** It does not. It splits the words and passes them on. The program reads the options.
- **Option, flag and switch are different.** They are three names for the same thing. This guide says option.
- **A long option works on every system.** Long options with two dashes are not standard. A command on macOS may not accept the long option that its Linux version does.

## Check yourself
1. In `sort -r fruit.txt`, which word is the command name? Answer: `sort`.
2. In the same line, what is `-r`? Answer: an option.
3. How many arguments are in `printf '%s\n' 'a b'`? Answer: two, `'%s\n'` and `'a b'`, because the quotes keep `a b` together.
4. Who decides what an option means, the shell or the program? Answer: the program.

## Where next
The next guide, `explain learn files-and-folders`, looks at the most common kind of argument: a path that names a file or a folder.

To see how real commands document their options, look at `explain ls`, `explain sort` and `explain head`.
