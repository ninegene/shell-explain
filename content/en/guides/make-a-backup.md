---
schema_version: 1
id: make-a-backup
title: Make a Backup You Can Trust
subtitle: Copy a folder, pack it into one file, and check that the copy is good
summary: Make a backup of a folder by copying it, or by packing it into one archive file with tar, then restore it somewhere else and check that it matches the original.
series: getting-things-done
order: 4
level: beginner
reading_minutes: 7
updated: 2026-11-09
tags: [backup, archive, files]
related_commands: [cp, tar, gzip, mkdir]
related_guides: [files-and-folders, rename-and-move, exit-codes]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html
  - https://www.gnu.org/software/tar/manual/tar.html
  - https://man.freebsd.org/cgi/man.cgi?query=bsdtar&sektion=1
---
## Goal
By the end you can make a copy of a folder in two ways, restore it into a different place, and check that the restored copy matches the original. A backup you have never restored is only a hope.

## The idea in plain words
A **backup** is a second copy of your files, kept so that you can get them back if the first copy is lost, damaged or changed by mistake. The command you use does not matter as much as three habits:

1. **Make the copy.** There are two common routes. You can copy the folder as it is with `cp -R`. Or you can pack it into one file called an **archive** with `tar`, which is easier to move, store and send.
2. **Keep it somewhere else.** A copy on the same disk is lost if the disk is. This guide stays inside one folder so it is safe to try, but a real backup goes to another disk or another computer.
3. **Check that it works.** Restore the copy somewhere new and compare it with the original.

An archive made with `tar` holds the files and their folder names, with no compression by default. Adding `-z` compresses it with gzip, which usually makes the file smaller.

## A picture
Read the picture from left to right. The original folder goes into a backup, and the backup comes back out as a restored copy that you compare with the original.

```diagram
   folder        backup                      restored copy
   keep/   -->   keep.copy/        (route 1: cp -R)
   keep/   -->   keep.tar          (route 2: tar -cf)
   keep/   -->   keep.tgz          (route 2: tar -czf, compressed)

   keep.tar  -->  restore/keep/  -->  diff -r keep restore/keep
```

Diagram: the folder named keep leads to three possible backups. Route one copies the folder with cp dash capital r into keep dot copy. Route two packs it with tar dash c f into one file called keep dot tar. The compressed form of route two uses tar dash c z f and makes keep dot tgz. The last line shows a restore: the tar file is unpacked into restore slash keep, and diff dash r compares it with the original.

## Another way to see it
Think of moving house. Route 1 carries the furniture over piece by piece, and the new place looks like the old one. Route 2 packs everything into boxes: one container, easy to put in a van or a store room. Either way, the real test is unpacking the boxes and checking that nothing is missing.

The table compares the routes.

| Route | Result | Good for |
| --- | --- | --- |
| `cp -R` | A second folder | Quick local copy |
| `tar -cf` | One archive file | Moving and storing |
| `tar -czf` | One smaller file | Saving space |

The standard names `-R` as the option for copying a folder and its contents. It no longer specifies `-r`, though some systems still accept it, so use `-R` when you want a command that works everywhere.

## Walk through it
Here is what happens, step by step, when you make an archive of a folder, restore it, and check it.

1. `tar -cf keep.tar keep` packs the folder `keep` and everything in it into one file, `keep.tar`. The `c` means create, and `-f` says that the next word is the archive's name.
2. `tar -tf keep.tar` lists what is inside without unpacking it. This is the quickest check that the archive holds the files you expect.
3. `mkdir restore` makes an empty folder for the test. Never restore on top of the original.
4. `tar -xf keep.tar -C restore` unpacks the archive into `restore`. The `x` means extract, and `-C restore` changes into that folder first. The folder must already exist.
5. `diff -r keep restore/keep` compares the two folders file by file. With nothing to report, it prints nothing and exits with status 0. Then `&& echo 'backup matches'` runs, and you see the result.

An archive can write files outside the folder you name, if it comes from someone else and holds odd paths. List an archive you did not make before you extract it.

## Try it
These commands write small files into a folder called `trybackup` inside the folder you are in, and then back them up and restore them. The name is unusual so it does not replace a folder of yours. Each block makes its own files, so each one works on its own.

Copy a folder with `cp -R`:

```sh
mkdir -p trybackup/keep/docs
printf 'one\n' > trybackup/keep/a.txt
printf 'two\n' > trybackup/keep/docs/b.txt
cp -R trybackup/keep trybackup/keep.copy
ls trybackup/keep.copy
```

You should see this:

```text
a.txt
docs
```

Pack a folder into one archive, then list the files in it. `grep` keeps the lines that name files, and `sort` makes the order the same everywhere:

```sh
mkdir -p trybackup/keep/docs
printf 'one\n' > trybackup/keep/a.txt
printf 'two\n' > trybackup/keep/docs/b.txt
cd trybackup
tar -cf keep.tar keep
tar -tf keep.tar | grep txt | sort
```

You should see this:

```text
keep/a.txt
keep/docs/b.txt
```

Compress it with `-z`. The name `.tgz` is a habit that tells people it is a compressed archive:

```sh
mkdir -p trybackup/keep
printf 'one\n' > trybackup/keep/a.txt
cd trybackup
tar -czf keep.tgz keep
tar -tzf keep.tgz | grep txt
```

You should see this:

```text
keep/a.txt
```

Restore into a new folder and check that the copy matches. Nothing is printed by `diff` when the folders are the same, so the `echo` tells you the result:

```sh
mkdir -p trybackup/keep/docs
printf 'one\n' > trybackup/keep/a.txt
printf 'two\n' > trybackup/keep/docs/b.txt
cd trybackup
tar -cf keep.tar keep
mkdir -p restore
tar -xf keep.tar -C restore
diff -r keep restore/keep && echo 'backup matches'
```

You should see this:

```text
backup matches
```

## Common confusions
- **A backup is any second copy.** A copy on the same disk is lost with the disk. Keep one somewhere else.
- **If the command finished, the backup is good.** Check it. Restore into a new folder and compare with `diff -r`.
- **`cp` copies folders by itself.** It needs `-R` for a folder, and it stops with an error without it.
- **An archive is compressed.** Not unless you add `-z`. A plain `tar` archive holds the files as they are.
- **Extracting is always safe.** Restore into an empty folder, and list an archive from someone else first.
- **The archive's name decides its type.** The name `.tgz` is only a habit. The `-z` option is what compresses.

## Check yourself
1. What is the difference between `cp -R` and `tar -cf`? Answer: `cp -R` makes a second folder, and `tar -cf` packs the folder into one archive file.
2. How do you see what is in an archive without unpacking it? Answer: `tar -tf name.tar`.
3. How do you check that a restored copy matches the original? Answer: `diff -r original restored`, which prints nothing when they match.
4. Why restore into a new folder? Answer: so the test cannot damage the original.

## Where next
This ends the first group of goal-first guides. You can find a word, build a report, tidy files and make a backup you have checked, and each one uses the ideas from the earlier series.

For the details of each command, see `explain cp`, `explain tar`, `explain gzip` and `explain mkdir`.
