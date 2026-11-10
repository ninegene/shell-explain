---
schema_version: 1
id: family-tree
title: "A Short History: Unix, BSD, GNU, Linux"
subtitle: Who made what, and in what order
summary: Follow the story from the first Unix to the systems on your computer today, and see why the same tools come in several versions.
series: why-so-many-versions
order: 1
level: beginner
reading_minutes: 7
updated: 2026-11-09
tags: [history, unix, linux]
related_commands: [sed, ls]
related_guides: [permissions, whos-who, same-command-many-versions]
references:
  - https://man7.org/linux/man-pages/man7/standards.7.html
  - https://www.gnu.org/gnu/gnu-history.html
  - https://www.gnu.org/gnu/initial-announcement.html
  - https://www.kernel.org/linux.html
  - https://docs.freebsd.org/en/books/handbook/introduction/
  - https://developer.apple.com/library/archive/documentation/Darwin/Conceptual/KernelProgramming/Architecture/Architecture.html
---
## Goal
By the end you can place Unix, BSD, System V, GNU, Linux, FreeBSD and macOS on a rough timeline, and say who made each one. This is the story behind a fact you will meet again and again: the same command can behave differently on different computers.

## The idea in plain words
Almost every command in this series comes from one family of operating systems, the Unix family. It began with a series of numbered releases from AT&T Bell Labs. Version 7, released in 1979, is a turning point, because after it Unix systems split into two main dialects: **BSD** and **System V**.

**BSD** stands for Berkeley Software Distribution. It was the University of California at Berkeley's own version of Unix. **System V** was the version sold by AT&T.

Then two projects set out to build their own Unix-like systems. The **GNU Project**, announced by Richard Stallman in 1983, aimed to write a complete Unix-compatible system and give it away free. **Linux** is a kernel that Linus Torvalds developed from 1991. Put the GNU tools and the Linux kernel together and you get a complete system.

The BSD line carried on too. **FreeBSD** is based on a release from the Berkeley group, and a large part of the technology in the core of **macOS** comes from FreeBSD.

So there is more than one family line, and each one wrote or rewrote its own tools. That is why `sed` on a Mac and `sed` on a Linux machine are cousins, not twins.

## A picture
Read the picture from top to bottom. Lines show which system came from which. The years are the dates the sources give.

```diagram
   Unix from AT&T Bell Labs, Version 1 ... Version 7 (1979)
                         |
        +----------------+----------------+
        |                                 |
   BSD (Berkeley)                    System V (AT&T)
   3BSD, 4BSD   1980                 SVr1   1983
   4.2BSD       1983                 SVr4   1989
   4.4BSD       1993
        |
        +--> FreeBSD   (from 4.4BSD-Lite)
        +--> Darwin and macOS   (much of it from FreeBSD)

   GNU: announced 1983, work began 1984. A Unix-compatible system.
   Linux: a kernel from 1991, free software in 1992, from scratch.
   GNU tools + the Linux kernel  =  a complete GNU/Linux system
```

Diagram: a tree with Unix from AT&T Bell Labs at the top, ending at Version 7 in 1979. It splits into two lines. The left line is BSD from Berkeley, with 3BSD and 4BSD in 1980, 4.2BSD in 1983 and 4.4BSD in 1993. FreeBSD comes from 4.4BSD-Lite, and much of the technology in Darwin and macOS comes from FreeBSD. The right line is System V from AT&T, from SVr1 in 1983 to SVr4 in 1989. Below the tree are two separate efforts. GNU was announced in 1983, with work beginning in 1984, as a Unix-compatible system. Linux is a kernel from 1991, free software in 1992, written from scratch. Combining the GNU tools with the Linux kernel gives a complete GNU/Linux system.

## Another way to see it
Think of languages. Latin split into several languages. They share many words, and a speaker of one can often guess the other, but a sentence that works in one can fail in another. Unix split in a similar way, and later systems were written to behave like it without copying its code.

| Name | When | Who | In one line |
| --- | --- | --- | --- |
| Version 7 Unix | 1979 | AT&T Bell Labs | Unix split in two after it |
| BSD | 1980 to 1993 | UC Berkeley | One Unix dialect |
| System V | 1983 to 1989 | AT&T | The other Unix dialect |
| GNU | 1983 on | Richard Stallman | A free Unix-compatible system |
| Linux | 1991 on | Linus Torvalds | A Unix-like kernel |
| FreeBSD | 1990s on | FreeBSD Project | A system from 4.4BSD-Lite |

## Walk through it
Here is the story in order, as one list of events.

1. In 1979, Version 7 Unix is released by AT&T Bell Labs. After this, Unix systems diverge into two dialects, BSD and System V.
2. Berkeley releases 3BSD and 4BSD in 1980 and 4.1BSD in 1981. In 1983 comes 4.2BSD, the first Berkeley release with a TCP/IP stack and the sockets API.
3. In 1983, AT&T releases System V. On September 27, 1983, Richard Stallman announces the GNU Project: a complete Unix-compatible system, to be given away free.
4. The plan starts slowly. The work begins in January 1984, and the Free Software Foundation is founded in October 1985 to raise money for it.
5. Berkeley releases 4.3BSD in 1986. AT&T releases System V Release 3 in 1986 and Release 4 in 1989.
6. By 1990, the GNU Project has found or written all the major components of its system except the kernel.
7. In 1991, Linus Torvalds develops the Linux kernel. It is made free software in 1992. Combining it with the GNU system gives a complete operating system.
8. In 1993, Berkeley releases 4.4BSD, its last major release. FreeBSD is based on 4.4BSD-Lite, a release from the same Berkeley group. Apple's Darwin, the core of macOS, derives a large portion of its technology from FreeBSD.

The people who wrote GNU chose to make it Unix-compatible. Their reason was that the overall design was already proven and portable, and compatibility makes it easy for Unix users to switch.

## Try it
These commands only look at things. Each one works on its own.

Ask your computer what kind of kernel it runs. On a Mac you will see `Darwin`. On a Linux system you will see `Linux`. Other systems print their own names:

```sh
uname -s
```

Look at the family resemblance. These three top-level folders exist on macOS, which is on the BSD side, and on Linux. This project's tests found them on macOS and in an Alpine Linux guest:

```sh
ls -d /bin /etc /usr
```

You should see this:

```text
/bin
/etc
/usr
```

## Common confusions
- **Linux is Unix.** Linux is described as a clone of Unix, written from scratch by Linus Torvalds with help from many others. It is built to behave like Unix, and it aims towards the POSIX standard.
- **GNU and Linux are the same thing.** They are not. GNU is a project that produced a large set of tools. Linux is a kernel. The GNU Project calls the complete system made from both GNU/Linux.
- **macOS is based on Linux.** It is not. The core of macOS comes from Apple's Darwin, which derives much of its technology from FreeBSD, with Mach and technology that Apple acquired from NeXT. Apple's own description of this is archived documentation from when the system was called OS X, so check Apple's current pages for today's details.
- **Unix is one system.** After Version 7 it split into two main dialects, and many systems came later that are Unix-like but built separately.
- **History does not matter when I only type commands.** It does, because each family wrote its own versions of the tools, and they differ in the details. The next guides explain how.

## Check yourself
1. Which two dialects did Unix split into after Version 7? Answer: BSD and System V.
2. Who announced the GNU Project, and in which year? Answer: Richard Stallman, in 1983.
3. What is Linux, in one word? Answer: a kernel.
4. Which system is FreeBSD based on? Answer: the 4.4BSD-Lite release from Berkeley.

## Where next
The next guide, `explain learn whos-who`, sorts out the names from this story: which ones are kernels, which are complete systems, which are collections of tools, and which are standards.

After that, `explain learn same-command-many-versions` shows the effect of this history on the commands you type.
