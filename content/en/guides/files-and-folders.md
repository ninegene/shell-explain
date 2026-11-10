---
schema_version: 1
id: files-and-folders
title: "The Map Your Shell Walks: Files and Folders"
subtitle: Paths, the current folder, and how to move around
summary: Learn how files and folders are arranged in a tree, how to read a path, and how absolute and relative paths differ.
series: terminal-shell-commands
order: 3
level: beginner
reading_minutes: 6
updated: 2026-11-08
tags: [basics, files, paths]
related_commands: [cd, ls, mkdir]
related_guides: [terminal-shell-command, anatomy-of-a-command, processes]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap03.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
---
## Goal
By the end you can read a path, tell an absolute path from a relative one, and say where the shell is working right now. You can also move around with `cd`.

## The idea in plain words
A **file** holds data. A **directory**, which most people call a folder, holds files and other directories. Directories inside directories make a tree, and the top of the tree is the **root directory**, written `/`.

A **path** is a way to reach a file or folder, written as names separated by `/`. A path can start from the top, or from where you are.

An **absolute path** starts with `/`, so it begins at the root. It names the same place from wherever you stand.

A **relative path** does not start with `/`. It begins at the **current directory**, the folder the shell is working in right now. It names different places depending on where you stand.

Every folder has two special names. A single dot `.` means this folder, and two dots `..` mean the parent folder, the one above. In the shell, `~` stands for your home directory, your own personal folder.

Three commands go with this. `pwd` prints the current directory, `ls` lists what is in a folder, and `cd` changes the current directory.

## A picture
Read the picture from the top down. This is an example tree, not your real disk. Each branch is a folder or a file inside the one above it.

```diagram
   /                          the root directory
   +-- docs/
   |     +-- notes/
   |     |     +-- today.txt
   |     +-- todo.txt
   +-- photos/
         +-- trip.jpg

   if you are in /docs:
     absolute path to a file    /docs/notes/today.txt
     relative path to it        notes/today.txt
     relative path up and over  ../photos/trip.jpg
```

Diagram: a tree. The root directory has two folders, docs and photos. Docs holds a folder called notes, which holds today dot txt, and also a file todo dot txt. Photos holds trip dot jpg. Below the tree are three paths for someone standing in the docs folder: the absolute path slash docs slash notes slash today dot txt, the relative path notes slash today dot txt, and the relative path dot dot slash photos slash trip dot jpg, which goes up one level and over.

## Another way to see it
Think of addresses. A full postal address, with the country and the city, is an absolute path. It works from anywhere. Telling a friend to go two doors down is a relative path. It only makes sense from where they stand now.

| Path | Kind | Where it leads |
| --- | --- | --- |
| `/docs/todo.txt` | Absolute | From the root, whatever the current folder |
| `notes/today.txt` | Relative | Down into notes, from the current folder |
| `./todo.txt` | Relative | A file in the current folder |
| `../photos` | Relative | Up one folder, then into photos |

## Walk through it
Here is what happens, step by step, when you run these three lines in a folder called `practice`.

1. You run `mkdir -p docs/notes`. The shell starts `mkdir`, which makes the folder `docs` and, inside it, `notes`. The path `docs/notes` is relative, so it starts from the current folder.
2. You run `cd docs/notes`. The `cd` command is handled by the shell itself, because only the shell can change its own current folder. The shell is now working in `practice/docs/notes`.
3. You run `cd ..`. The two dots mean the parent folder, so the shell moves up to `practice/docs`.
4. You run `pwd`. It prints the absolute path of the current folder, starting with `/` and ending with `practice/docs`.
5. You run `ls`. With no argument, `ls` lists the current folder, so it shows `notes`.

A relative path is always read from the current folder. Moving the current folder changes what a relative path means.

## Try it
These commands make a folder called `practice` and a small file in it, in the folder you are in. Run them where that is fine. Each one works on its own.

Print the absolute path of the folder you are in. It starts with `/`, and the rest differs on every computer:

```sh
pwd
```

Make folders, move into them and back, and print the name of the folder you are in each time. The command `basename` prints only the last part of a path:

```sh
mkdir -p practice/docs
cd practice/docs
basename "$(pwd)"
cd ..
basename "$(pwd)"
```

You should see this:

```text
docs
practice
```

Reach one file by a relative path and by an absolute path. Both name the same file, so both print the same text. The variable `$PWD` holds the absolute path of the current folder:

```sh
mkdir -p practice/docs
echo hi > practice/docs/note.txt
cat practice/docs/note.txt
cat "$PWD/practice/docs/note.txt"
```

You should see this:

```text
hi
hi
```

Show the two special names that every folder has. The option `-a` makes `ls` include names that start with a dot. In an empty folder you see only the dot and the double dot:

```sh
mkdir -p practice
ls -a practice
```

You should see this:

```text
.
..
```

## Common confusions
- **Folder and directory are different.** They mean the same thing. This guide says directory when it matches the manuals, and folder otherwise.
- **A path that does not start with a slash is wrong.** It is a relative path. It begins at the current folder.
- **A relative path works from anywhere.** It does not. It depends on the current folder, so the same relative path can name two different files.
- **`cd` changes the file or folder.** It only changes where the shell is working. Nothing on the disk is touched.
- **`ls` tells me where I am.** It lists what is in a folder. Use `pwd` to see where you are.
- **`..` means up two levels.** It means up one level, to the parent. Use `../..` for two.

## Check yourself
1. What does a path that starts with `/` mean? Answer: it is an absolute path, so it starts at the root.
2. What do `.` and `..` mean? Answer: `.` is the current folder, and `..` is its parent.
3. Which command prints the current folder? Answer: `pwd`.
4. Why does a relative path give different results in different folders? Answer: it starts from the current folder, so moving the current folder changes where it points.

## Where next
You now know the commands, the options and the paths. The next guide, `explain learn text-and-binary`, looks at what is inside a file. Later guides cover how the shell finds a program by name and who may use a file.

To practise, look at the command pages for `cd`, `ls` and `mkdir`.
