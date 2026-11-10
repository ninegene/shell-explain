---
schema_version: 1
id: same-command-many-versions
title: One Command, Many Versions
subtitle: Why sed -i works one way on a Mac and another way on Linux
summary: See why one command name can stand for several different programs, what a standard does and does not promise, and how to write commands that work everywhere.
series: why-so-many-versions
order: 3
level: beginner
reading_minutes: 7
updated: 2026-11-09
tags: [history, portability, commands]
related_commands: [sed, head, date, stat]
related_guides: [family-tree, whos-who, anatomy-of-a-command]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sed.html
  - https://busybox.net/about.html
  - https://man7.org/linux/man-pages/man7/standards.7.html
  - https://www.gnu.org/gnu/gnu-history.html
---
## Goal
By the end you can explain why a command that works on one computer fails on another, tell which kind of difference you are seeing, and write the common core that works on all of them.

## The idea in plain words
A command name such as `sed` is a name for a job, not for one particular program. Behind the name, different systems keep different programs. Each is an **implementation** of the same command.

The implementations differ for four reasons, and the earlier guides gave you the history behind them.

1. **They have separate origins.** The BSD tools descend from the ones Berkeley shipped. The GNU Project wrote its own, aiming at compatibility with Unix. BusyBox wrote tiny versions of its own.
2. **They have different goals.** BusyBox is built for small and embedded systems, and by its own description its tools generally have fewer options than their GNU cousins.
3. **The standard covers only a core.** POSIX describes what every system should do. A tool can add much more, and each family added different extras.
4. **The system's builder chooses.** A computer ships with whichever set of tools its maker picked, and you can install others beside them.

The common core is the safe part. An extra option is where the trouble starts, and a command that uses one works only where that option exists.

## A picture
Read the picture from left to right. One command name leads to different programs, depending on the system, and the same option can mean something different in each.

```diagram
   you type:  sed -i ...
        |
        +--> macOS sed     -i needs an argument:  -i '.bak'  or  -i ''
        |
        +--> GNU sed       -i has an optional suffix:  -i.bak  or  -i
        |
        +--> BusyBox sed   -i has an optional suffix:  -i.bak  or  -i
```

Diagram: the command sed with the option dash i leads to three different programs. On macOS, sed needs an argument after dash i, for example dash i then quote dot bak, or dash i then empty quotes for no backup. GNU sed takes an optional suffix attached to the option, as in dash i dot bak, or dash i alone. BusyBox sed does the same: dash i dot bak, or dash i alone.

## Another way to see it
Think of plug sockets in different countries. The idea is the same everywhere: you plug in a device. But the shape of the plug differs, so a plug that works in one country needs an adapter in another. A command's core is the wall, and an extra option is the shape of the plug.

The table shows the same job done on different systems. Each entry here was run on macOS's own tools, on GNU coreutils 9.11 and on BusyBox 1.37.0, as part of the pages in this project.

| Job | GNU and BusyBox | macOS |
| --- | --- | --- |
| Edit, keep a backup | `sed -i.bak 's/a/b/' f` | `sed -i '.bak' 's/a/b/' f` |
| Date from seconds | `date -d @86400` | `date -r 86400` |
| Print a file's size | `stat -c '%s' f` | `stat -f '%z' f` |
| Drop last two lines | `head -n -2 f` (GNU) | none; use sed or tail |

The standard does not describe `-i` at all. Its description of `sed` lists only the options `-E`, `-e`, `-f` and `-n`. In-place editing is an extension that each family added in its own way.

## Walk through it
Here is what happens, step by step, when you type the same line on two computers: `sed -i 's/old/new/g' notes.txt`.

1. On a computer with GNU `sed`, the shell starts GNU `sed`. It sees `-i` and treats it as in-place editing with no backup. The text `'s/old/new/g'` is the script, and `notes.txt` is the file. The BusyBox `sed` tested in this project did the same.
2. On a Mac, the shell starts the macOS `sed`. It sees `-i` and expects the next word to be the backup suffix. It takes `'s/old/new/g'` as that suffix, and then `notes.txt` is taken as the script. That is not a valid script, so `sed` stops with an error and leaves the file as it was.
3. The same line printed no complaint in step 1 and an error in step 2. Neither program is broken. Each follows its own rule for `-i`.
4. To write a line that both accept, use the core only. Send the result to a new file with `>`, and keep the original. A line such as `sed 's/old/new/g' notes.txt > notes.new` works the same everywhere.
5. When you must use an extra option, look it up for the system you are on. This is why this project's pages ask which implementation you use, and why `explain sed` first works out which `sed` you have.

Step 2 is what the macOS `sed` did when this project ran that exact line: it printed an error saying that the characters after the command were extra, and the file was unchanged. The exact error text depends on the version.

## Try it
These commands write small files in the folder you are in, so run them where that is fine. Each one works on its own. Each makes its own small file, and the name `sample.txt` is unusual so that it does not replace a file of yours.

Use only the common core: print the second line. This works the same on every system:

```sh
printf 'old old\nsecond\nthird\n' > sample.txt
sed -n '2p' sample.txt
```

You should see this:

```text
second
```

Make an edited copy instead of editing in place. This works the same on every system, and your original stays as it was:

```sh
printf 'old old\nsecond\nthird\n' > sample.txt
sed 's/old/new/g' sample.txt > sample.new
cat sample.new
```

You should see this:

```text
new new
second
third
```

Ask your `sed` for its version, and keep the first line. The answer depends on your system, so there is no fixed output. A GNU `sed` names itself. The BusyBox one says that it is not GNU `sed`. The macOS one reports an illegal option, because it has no `--version`:

```sh
sed --version | head -n 1
```

## Common confusions
- **The same command name means the same program.** It does not. The name is shared, and the program behind it depends on the system.
- **If a command fails on my Mac, I typed it wrong.** Maybe not. The option may exist only in the GNU version. Check the manual for your system.
- **The standard makes every system behave the same.** It covers a core. POSIX `sed`, for example, has no `-i`, so each family added its own.
- **The GNU version is the real one.** There is no single real one. GNU's has more options, and the standard's core is what all of them share.
- **Every Linux system has the same tools.** It does not. The kernel is Linux, but the tools might be GNU or BusyBox or something else, and BusyBox generally has fewer options.
- **A different version means a worse version.** It means a different set of choices. BusyBox leaves options out to stay small.

## Check yourself
1. Why can one command name mean several programs? Answer: different systems keep different implementations of it, from different origins and with different goals.
2. Does POSIX describe `sed -i`? Answer: no. Its `sed` has only `-E`, `-e`, `-f` and `-n`.
3. What is the safe way to write a command that works everywhere? Answer: use the options in the standard's core, and make a new file instead of an in-place edit.
4. How does `explain` choose which documentation to show? Answer: it looks at the program itself to work out which implementation it is.

## Where next
You now know why `sed` on one machine is not `sed` on another. Later guides in this series will explain how the tools reach your computer, and how to find out which ones you have.

For the exact differences, one command at a time, see `explain sed`, `explain head`, `explain date` and `explain stat`.
