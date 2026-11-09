---
schema_version: 1
id: path-and-environment
title: 'Why "Command Not Found" Happens'
subtitle: How the shell finds programs, and what environment variables are
summary: Learn how the shell uses PATH to find a program by name, what environment variables are, and why command not found does not always mean the program is missing.
series: terminal-shell-commands
order: 8
level: beginner
reading_minutes: 6
updated: 2026-11-08
tags: [basics, path, environment]
related_commands: [ls, cd]
related_guides: [terminal-shell-command, processes, exit-codes]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap08.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
  - https://man7.org/linux/man-pages/man7/environ.7.html
---
## Goal
By the end you can explain how the shell finds a program from its name, what an environment variable is, and what to check when you see command not found.

## The idea in plain words
When you type `ls`, the shell has to find the program called `ls` on the disk. It does not search the whole disk. It looks in a short list of folders.

That list is held in a variable called **PATH**. A **variable** is a name that holds a piece of text. PATH holds folder names separated by colons, such as `/usr/local/bin:/usr/bin:/bin`. The shell looks in the first folder for a program with that name, then the second, and so on, and runs the first one it finds. If it finds none, it says command not found.

Some commands never need PATH. A command that contains a slash, such as `./script`, names its place directly, so the shell does not search. And built-in commands such as `cd` are part of the shell.

A variable that every program carries along is an **environment variable**. When a program starts another, the new program gets a copy of the parent's environment. That is how PATH, HOME and others reach every program you run. A variable you set in the shell stays private to the shell until you mark it for export with `export`.

## A picture
Read the picture from the top down. The shell checks each folder in order and stops at the first match. The folder names are an example.

```diagram
   you type:   grep

   PATH = /usr/local/bin : /usr/bin : /bin
              |               |         |
              |               |         +-- never checked
              |               +------------ grep is here: run it, stop
              +---------------------------- no grep here, look further

   if no folder has it:   command not found
```

Diagram: the command grep and a PATH made of three folders: slash usr slash local slash bin, slash usr slash bin, and slash bin. The shell looks in the first folder and finds no grep. It looks in the second folder and finds grep, so it runs that one and stops. The third folder is never checked. If no folder had grep, the shell would report command not found.

## Another way to see it
Think of looking for a book in a few rooms of a house, always in the same order. You check the study, then the lounge, then the hall. You take the first copy you find. If the book is on a shelf in the garage, you will not find it, even though the house has it, because the garage is not on your list.

| Variable | What it holds | Example |
| --- | --- | --- |
| `PATH` | Folders searched for programs | `/usr/bin:/bin` |
| `HOME` | Your home directory | Your own folder |
| `TERM` | The kind of terminal in use | `xterm-256color` |

You can read any variable with a dollar sign. `echo $HOME` prints your home folder.

## Walk through it
Here is what happens, step by step, when you type `grep` and the shell cannot find it.

1. The shell reads the line. The command name is `grep`, which contains no slash.
2. The shell is not running a built-in with that name, so it reads the variable PATH.
3. It looks for a file called `grep` in the first folder of PATH, then the second, and so on.
4. It finds no file called `grep` that it may run in any of them.
5. The shell gives up. It writes command not found to its standard error and sets the exit status to 127.

The program may well exist somewhere else on your disk. It just is not in a folder on the list.

## Try it
These commands only look at things, except the last one, which writes a small file called `nothing.txt` in the folder you are in. Run it where that is fine. Each one works on its own.

Print your PATH. You will see folders separated by colons. They differ between computers:

```sh
echo $PATH
```

Ask the shell how it would run a name. For a built-in such as `cd`, which only the shell can do, it answers with the name itself:

```sh
command -v cd
```

You should see this:

```text
cd
```

Ask about a program. The shell prints the full path of the file it found by searching PATH. The path differs between systems:

```sh
command -v ls
```

Ask about a name that does not exist. The `command -v` check prints nothing and fails, so the `||` part runs:

```sh
command -v nosuchcommand || echo "not found"
```

You should see this:

```text
not found
```

See the difference between a variable that is only the shell's and one that is exported. The child shell started by `sh -c` is a separate program, so it only sees exported variables:

```sh
GUIDE_COLOR=blue
sh -c 'echo "child sees [$GUIDE_COLOR]"'
export GUIDE_COLOR
sh -c 'echo "child sees [$GUIDE_COLOR]"'
```

You should see this:

```text
child sees []
child sees [blue]
```

Give a variable to one command only, by putting it before the command. The next command does not see it:

```sh
GUIDE_SIZE=large sh -c 'echo "size is [$GUIDE_SIZE]"'
sh -c 'echo "now it is [$GUIDE_SIZE]"'
```

You should see this:

```text
size is [large]
now it is []
```

Make the search fail on purpose. A shell whose PATH is a folder that does not exist cannot find `ls`, so it reports status 127. The error message goes into a file:

```sh
sh -c 'PATH=/nonexistent; ls' 2> nothing.txt; echo $?
```

You should see this:

```text
127
```

## Common confusions
- **Command not found means the program is not installed.** It means the shell did not find it in the folders on PATH. The program may be installed in a folder that is not on the list.
- **PATH holds programs.** It holds folder names. The programs are files inside those folders.
- **The order of PATH does not matter.** It does. The shell takes the first match, so an earlier folder wins when two have the same name.
- **A variable I set is seen by every program.** Only variables that you export are passed on, and only to the programs started after that. Each program gets its own copy, so a change inside it does not come back.
- **A name with a slash is searched for in PATH.** It is not. `./script` names its place directly, so the shell does not search.

## Check yourself
1. What does PATH hold? Answer: a list of folders, separated by colons, in which the shell looks for programs.
2. What does the shell do when it finds the name in the first folder? Answer: it runs that one and stops looking.
3. What exit status does the shell set when it cannot find a command? Answer: 127.
4. Why might `grep` be installed but still give command not found? Answer: its folder is not on PATH.

## Where next
The last guide in this series, `explain learn permissions`, looks at who may read, write and run a file, which also explains the other common error, permission denied.

To move around the folders that PATH is made of, look at `explain ls` and `explain cd`.
