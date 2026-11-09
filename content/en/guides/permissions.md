---
schema_version: 1
id: permissions
title: Who May Read, Write, Run
subtitle: How to read the letters in ls -l, and what chmod changes
summary: Learn the three kinds of permission, the three groups of people they apply to, and how to read and change them.
series: terminal-shell-commands
order: 9
level: beginner
reading_minutes: 6
updated: 2026-11-08
tags: [basics, permissions, files]
related_commands: [ls, chmod, stat]
related_guides: [files-and-folders, processes, exit-codes]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html
  - https://man7.org/linux/man-pages/man7/inode.7.html
  - https://man7.org/linux/man-pages/man1/chmod.1.html
---
## Goal
By the end you can read a permission string such as `-rw-r--r--`, say who may read, write and run a file, and change the permissions with `chmod`.

## The idea in plain words
Every file has an **owner**, and it belongs to a **group** of users. The system keeps a short list of **permissions** for the file. They answer one question: who may do what?

There are three kinds of permission. **Read** lets you look at the contents. **Write** lets you change them. **Execute** lets you run the file as a program.

And there are three classes of people. The **owner** is the user who owns the file. The **group** is the other users in the file's group. **Others** is everyone else.

Each class has its own read, write and execute settings, so there are nine settings in all. The system checks which class you belong to, and uses that class's settings.

On a folder, the letters mean something a little different. Read lets you list the names inside. Write lets you add or remove entries. Execute lets you go into the folder and use it as part of a path.

An account called the **superuser** is allowed to read and write any file. Even the superuser can only run a file that has at least one execute setting turned on.

## A picture
Read the picture from left to right. The ten characters at the start of an `ls -l` line say what the file is and who may do what.

```diagram
     - r w -   r - -   r - -     note.txt
     | \___/   \___/   \___/
     |   |       |       |
     |   |       |       +-- others: read only
     |   |       +---------- group: read only
     |   +------------------ owner: read and write
     +---------------------- type: - is a file, d is a directory
```

Diagram: the ten characters dash r w dash r dash dash r dash dash, followed by the name note dot txt. The first character, a dash, says the item is a file; a d would mean a directory. The next three characters, r w dash, are the owner's permissions: read and write. The next three, r dash dash, are the group's: read only. The last three, r dash dash, are for others: read only.

## Another way to see it
Think of a house. The owner holds every key. The family, the group, can come in and look around. Visitors, the others, can look through the window but not touch anything. Each person has the rights that fit their role.

| Letter | On a file | On a folder |
| --- | --- | --- |
| `r` | Read the contents | List the names inside |
| `w` | Change the contents | Add or remove entries |
| `x` | Run it as a program | Go in, and use it in a path |

You can also write the nine settings as three numbers. Read is 4, write is 2 and execute is 1. Add them for each class: 644 means rw- for the owner (4 plus 2), and r-- for the group and for others.

## Walk through it
Here is what happens, step by step, when you try to run a script that is not yet marked as a program.

1. You type `./hello.sh`. The shell asks the kernel to run that file.
2. The kernel looks at the file's permissions. It decides which class you are in, and checks the execute setting for that class.
3. No class has execute turned on, so the kernel refuses. The shell prints a permission denied message and sets the exit status to 126.
4. You run `chmod +x hello.sh`. The `chmod` command changes the permissions so that the execute settings are on.
5. You run `./hello.sh` again. The kernel finds the execute setting, so it starts the program, and the script prints its text.

Permissions guard the file, but the owner can change them. Normally only the owner, or the superuser, may run `chmod` on a file.

## Try it
These commands write small files in the folder you are in, so run them where that is fine. Each one works on its own.

Make a file and give it the usual permissions for a text file. The `cut` command keeps the first ten characters of the `ls -l` line, which are the permission string:

```sh
echo hi > note.txt
chmod 644 note.txt
ls -l note.txt | cut -c1-10
```

You should see this:

```text
-rw-r--r--
```

Now take away the group's and the others' access, so that only the owner may read and write:

```sh
echo hi > note.txt
chmod 600 note.txt
ls -l note.txt | cut -c1-10
```

You should see this:

```text
-rw-------
```

Try to run a script before and after turning on its execute permission. The first attempt is refused, and the message goes into a file:

```sh
printf '#!/bin/sh\necho hello from a script\n' > hello.sh
./hello.sh 2> denied.txt || echo "could not run it yet"
chmod +x hello.sh
./hello.sh
```

You should see this:

```text
could not run it yet
hello from a script
```

See that a folder starts with the letter `d`. The option `-d` makes `ls` describe the folder itself and not what is inside it:

```sh
ls -ld . | cut -c1
```

You should see this:

```text
d
```

## Common confusions
- **Execute only matters for programs.** It also matters for folders. Without it, you cannot go into a folder or use it in a path.
- **Others means everybody.** It means everyone who is not the owner and not in the group. The owner and the group have their own settings.
- **Permissions stop the superuser.** The superuser can read and write any file. Only the execute setting still has to be on.
- **`chmod 777` fixes everything.** It makes the file writable by everyone, which is risky. Turn on only what you need.
- **Permission denied means the file is missing.** It means the file is there, but your class does not have the setting that the action needs.

## Check yourself
1. What do the ten characters at the start of an `ls -l` line show? Answer: the type of the item, and then the owner's, the group's and others' read, write and execute settings.
2. What does `-rw-r--r--` allow? Answer: the owner may read and write, and the group and others may only read.
3. What does the number 6 stand for in a mode such as 600? Answer: read plus write, because 4 plus 2 is 6.
4. Why does `./hello.sh` fail before `chmod +x`? Answer: the file has no execute setting, so the system refuses to run it.

## Where next
You have now met the main ideas of using a shell: the terminal and the shell, commands and options, paths, streams, pipes, exit status, text and binary files, PATH and permissions. A later series will look at why the same command differs between systems.

To look at permissions on real files, see `explain ls`, `explain chmod` and `explain stat`.
