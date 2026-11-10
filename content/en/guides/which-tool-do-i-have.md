---
schema_version: 1
id: which-tool-do-i-have
title: Which Tool Do I Have?
subtitle: Four questions that tell you which sed, awk or ls is answering
summary: Find out which program a command name starts and which implementation it is, using the shell, the version flag, the file itself and the system's own name.
series: why-so-many-versions
order: 5
level: beginner
reading_minutes: 6
updated: 2026-11-09
tags: [portability, detection, commands]
related_commands: [sed, awk, ls]
related_guides: [same-command-many-versions, distributions, path-and-environment]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/command.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/type.html
  - https://busybox.net/about.html
  - https://man7.org/linux/man-pages/man5/os-release.5.html
  - https://www.gnu.org/software/gawk/manual/
---
## Goal
By the end you can find out, in about a minute, which program a command name starts on your computer and which implementation it is: GNU, BusyBox, macOS or another.

## The idea in plain words
When you type `sed`, the shell looks through the folders in `PATH` and starts the first program it finds with that name. So there are two separate questions:

1. **Which file runs?** The shell answers that.
2. **What kind of program is that file?** The program, or the file itself, answers that.

Four tools cover both questions.

- **Ask the shell** where the name leads, with `command -v` or `type`.
- **Ask the program** for its version, with `--version`.
- **Look at the file**, to see whether it is a link to something else.
- **Ask the system** what it is, with `uname` and, on Linux, the `os-release` file.

No single check is perfect. A program may refuse `--version`, and a name can lead through several folders. Using two checks together is much safer than trusting one.

## A picture
Read the picture from top to bottom. Each step narrows the answer, from the name you typed to the kind of program behind it.

```diagram
   you type:  sed

   1. the shell    command -v sed       -->  /usr/bin/sed
         |
   2. the file     ls -l  /usr/bin/sed  -->  a file, or a link?
         |
   3. the program  sed --version        -->  GNU? BusyBox? an error?
         |
   4. the system   uname -s             -->  Darwin? Linux?
```

Diagram: four numbered steps from top to bottom. Step one, the shell, runs command dash v sed and prints the path of the program that will run. Step two, the file, runs ls dash l on that path to see whether it is a plain file or a link to another program. Step three, the program, runs sed dash dash version and shows whether it is GNU, BusyBox or an error. Step four, the system, runs uname dash s and shows whether the computer is Darwin or Linux.

## Another way to see it
Think of a person at a door saying "I'm from the council." You still want to know which council, so you ask for a badge. The name `sed` is the sentence at the door. The path is the building they walked out of, the version output is the badge, and the file is the stamp on the badge.

The table shows what the first line of `--version` looked like in this project's checks. The GNU `sed` was 4.10 on a Mac and BusyBox was 1.37.0 in a Linux guest. The exact text depends on the version.

| Family | `sed --version` | `awk --version` |
| --- | --- | --- |
| GNU | sed (GNU sed) 4.10 | GNU Awk 5.3.2 |
| BusyBox | This is not GNU sed version 4.0 | unrecognized option |
| macOS | illegal option error | awk version 20200816 |

The same check can give different answers in different places. A Mac's own `sed` has no `--version`, so the check prints an error. That error is itself the clue: a program that does not know the flag is not GNU.

## Walk through it
Here is the whole procedure, step by step, for a `sed` that behaves oddly.

1. Run `command -v sed`. The POSIX standard defines `command -v` to write the path or command that the shell would use. You get one path, such as `/usr/bin/sed`.
2. Run `type sed`. It tells you how the shell would interpret the name: as an alias, a built-in, or a file. In bash and zsh, `type -a sed` lists every match, so you can see when two different programs share the name and the first one wins.
3. Run `sed --version 2>&1 | head -n 1`. The `2>&1` sends error messages into the same pipe, and `head -n 1` keeps the first line. If it names GNU, you have the GNU one. If it says it is not GNU, it is probably BusyBox. If it complains about the option, it is probably the BSD kind that macOS uses.
4. Run `ls -l` on the path from step 1. On a BusyBox system, many commands are links to one program called `busybox`, so a link to it means you have BusyBox's version. A plain file means a program on its own.
5. If that is still unclear, run `uname -s`. It prints `Darwin` on a Mac and `Linux` on a Linux computer. On Linux, the `os-release` file also names the distribution.
6. For the final answer, ask `explain`. It looks at the program itself to choose the right documentation, and `explain sed --details` shows what it found.

If two answers disagree, trust the one closest to the file. A command that works in your interactive shell can still resolve differently inside a script, because an alias or a different `PATH` may apply there.

## Try it
These commands change nothing. Each one works on its own, and the answers depend on your computer, so there is no fixed output.

Find which file runs when you type `sed`:

```sh
command -v sed
```

Ask the shell how it reads the name `sed`. An alias, a built-in and a file each get a different answer:

```sh
type sed
```

Ask `sed` for its version and keep the first line. A GNU `sed` names itself, a BusyBox one says that it is not GNU, and a macOS one reports an illegal option:

```sh
sed --version 2>&1 | head -n 1
```

Look at the file itself. A long listing shows whether it is a plain file or a link:

```sh
ls -l "$(command -v sed)"
```

Print where a link points, or a note when the file is not a link. On a BusyBox system the answer names `busybox`:

```sh
readlink "$(command -v sed)" || echo 'not a symlink'
```

Do the same for `awk`. BusyBox has no usable version flag, so it reports an unrecognized option, which is also a clue:

```sh
awk --version 2>&1 | head -n 1
```

## Common confusions
- **`which` tells me everything.** `which` is not part of the standard. `command -v` is, and `type` shows more than where a file lives.
- **The version number tells me the family.** A number alone does not. Read the name that comes with it, such as GNU or BusyBox.
- **No `--version` means the tool is broken.** Not every implementation has the flag. macOS's `sed` and BusyBox's `awk` do not, and that missing flag is part of the answer.
- **My shell and my scripts see the same command.** Not always. An alias, a function or a different `PATH` can change what a name means.
- **I can tell by the operating system.** Not by itself. A Mac can have GNU tools in front of its own, and a Linux system can use BusyBox.
- **One check is enough.** Two agreeing checks are far stronger than one.

## Check yourself
1. Which command shows the path that a name leads to? Answer: `command -v`, or `type`.
2. What does a link to `busybox` tell you? Answer: that the command is BusyBox's version of it.
3. Why is a failed `--version` still useful? Answer: it shows the program does not know the GNU flag, which is a clue about its family.
4. Which tool picks the right documentation for the program you have? Answer: `explain`, which looks at the program itself.

## Where next
You can now find out which tool answers when you type a command, and why its options may differ from a page you read online. This ends the series on why there are so many versions.

For the exact differences, one command at a time, see `explain sed`, `explain awk` and `explain ls`.
