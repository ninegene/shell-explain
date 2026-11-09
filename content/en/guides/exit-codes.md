---
schema_version: 1
id: exit-codes
title: How a Command Says Yes or No
subtitle: Exit status, and how to run one command only if another worked
summary: Learn the number every command hands back when it ends, what zero means, and how to chain commands with and without conditions.
series: terminal-shell-commands
order: 7
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, exit-status, shell]
related_commands: [grep]
related_guides: [pipes-and-redirects, processes, input-output-errors]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
  - https://man7.org/linux/man-pages/man3/exit.3.html
---
## Goal
By the end you can read a command's exit status, say what zero and non-zero mean, and use `&&` and `||` to run a command only when the one before it worked or failed.

## The idea in plain words
When a command ends, it hands one small number back to the shell. This is its **exit status**. It is separate from the text the command prints.

Zero means success. Any other number means something did not go as planned. The shell keeps the number of the last command in a special variable called `$?`, so you can look at it with `echo $?`.

Different non-zero numbers can mean different things, and each command's manual says which. Some are the same everywhere. The shell uses 127 when it cannot find the command you typed, and 126 when it finds the name but the file cannot be run. A status above 128 means the command was stopped by a signal.

The shell can use the number to decide what to do next. `&&` means: run the next command only if this one succeeded, with status zero. `||` means: run the next command only if this one did not succeed. A semicolon `;` runs the next command either way.

## A picture
Read the picture from the top down. The status of the first command decides which way the line goes.

```diagram
                      first command
                            |
                    exit status is...
                    /               \
               0 (success)       not 0 (failed)
                  |                   |
     && runs the next one    || runs the next one
     || skips it             && skips it
```

Diagram: a first command ends with an exit status, and the picture splits in two. If the status is zero, meaning success, the operator and-and runs the next command and the operator or-or skips it. If the status is not zero, meaning failure, the operator or-or runs the next command and the operator and-and skips it.

## Another way to see it
Think of a traffic light that the last command leaves behind. Green, which is zero, means go on. Any other colour means stop, and then the pieces that handle trouble get their turn.

| Operator | Runs the next command when... | Example |
| --- | --- | --- |
| `&&` | The last status was 0 | `make && echo built` |
| `\|\|` | The last status was not 0 | `make \|\| echo failed` |
| `;` | Always, whatever the status | `make ; echo done` |

## Walk through it
Here is what happens, step by step, when you run `grep -q cherry fruit.txt && echo found`. The option `-q` makes `grep` print nothing and only report by its exit status.

1. The shell starts `grep`, which searches the file.
2. `grep` finds no line with `cherry`, so it ends with exit status 1. For `grep`, 1 means no line matched.
3. The shell stores 1 in `$?`.
4. The shell reads `&&`, which needs status 0. It sees 1, so it skips `echo found`.
5. Nothing is printed. The line ends with status 1.

If the file had contained `cherry`, `grep` would have ended with 0, and `echo found` would have run.

## Try it
These commands only look at things, except the last one, which writes a small file called `nothing.txt` in the folder you are in. Run it where that is fine. Each one works on its own.

Run a command that succeeds and read its status:

```sh
true; echo $?
```

You should see this:

```text
0
```

Run a command that fails and read its status:

```sh
false; echo $?
```

You should see this:

```text
1
```

Run a command only if the one before it succeeded:

```sh
printf 'apple\nbanana\n' | grep -q apple && echo found
```

You should see this:

```text
found
```

Run a command only if the one before it failed:

```sh
printf 'apple\nbanana\n' | grep -q cherry || echo "cherry not found"
```

You should see this:

```text
cherry not found
```

See the status the shell gives for a command that does not exist. The error message goes into a file, so only the number shows:

```sh
nosuchcommand 2> nothing.txt; echo $?
```

You should see this:

```text
127
```

## Common confusions
- **Zero means failure, like in some programming languages.** In the shell, zero means success. A non-zero number means something went wrong.
- **The exit status is what the command prints.** It is not printed. It is a separate number that the shell keeps, and you read it with `$?`.
- **Non-zero always means disaster.** Not always. `grep` ends with 1 when it finds no match, which is a normal answer and not a fault.
- **`&&` and `;` do the same.** `;` always runs the next command. `&&` runs it only after success.
- **`$?` remembers everything.** It holds only the status of the most recent command, so read it straight away.

## Check yourself
1. What does an exit status of zero mean? Answer: the command succeeded.
2. What does `||` do? Answer: it runs the next command only if the one before did not succeed.
3. Where do you read the status of the last command? Answer: in `$?`, for example with `echo $?`.
4. What status does the shell give when it cannot find a command? Answer: 127.

## Where next
You now know how commands read input, write output, report errors, and say whether they worked. The next guide, `explain learn path-and-environment`, explains how the shell finds a command by name, which is where the status 127 comes from.

For a command that uses exit status in a useful way, see `explain grep`.
