---
schema_version: 1
id: rename-and-move
title: "Tidy Up: Rename and Move Files Without Losing Any"
subtitle: One command for two jobs, and the habits that keep your files safe
summary: Rename one file, move files into a folder, rename many files with a loop, and avoid overwriting a file by accident, using mv and a shell loop.
series: getting-things-done
order: 3
level: beginner
reading_minutes: 7
updated: 2026-11-09
tags: [files, rename, safety]
related_commands: [mv, mkdir, cp, rm, ls]
related_guides: [files-and-folders, find-a-word, exit-codes]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/mv.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/cp.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html
---
## Goal
By the end you can rename a file, move files into a folder, rename many files at once, and avoid the one mistake that destroys data: moving a file on top of another file.

## The idea in plain words
A file's name and its folder are part of where it lives. Renaming a file and moving it to another folder look like different jobs, but the system treats them the same way: the file gets a new path. One command, `mv`, does both.

The command takes a source and a destination:

1. **Destination is a new name.** `mv draft.txt final.txt` renames the file.
2. **Destination is a folder that exists.** `mv draft.txt archive/` moves the file into that folder, under the same name.
3. **Both at once.** `mv draft.txt archive/final.txt` moves it and renames it.

The standard says that when the last argument names an existing folder, each source is moved into it under its own name. Otherwise the last argument is the new name.

The danger is the case where the destination file already exists. By default, `mv` replaces it without asking, and the old contents are gone. The habit that prevents this is to check first, or to ask `mv` to prompt with `-i`.

## A picture
Read the picture from top to bottom. The same command takes different actions, depending on what the last word is.

```diagram
   mv  SOURCE  DESTINATION

   DESTINATION is...             result
   -----------------------       ------------------------------
   a name that is free      -->  the file is renamed
   a folder that exists     -->  the file goes inside the folder
   a file that exists       -->  that file is replaced, silently
```

Diagram: the mv command takes a source and a destination, and the result depends on what the destination is. If the destination is a name that is free, the file is renamed. If it is a folder that exists, the file goes inside the folder. If it is a file that exists, that file is replaced, silently.

## Another way to see it
Think of a coat check. Handing over a ticket with a new name on it renames the coat. Handing it to a different cloakroom moves it. But if you hand it to a hook that already holds a coat, the old coat falls on the floor, and nobody tells you.

The table lists the jobs and the safe form of each one.

| Job | Command | Safer habit |
| --- | --- | --- |
| Rename | `mv old new` | Check that new is free |
| Move into a folder | `mv file dir/` | End with a slash |
| Rename many | a `for` loop | Print the names first |
| Copy, keep original | `cp file copy` | Use when unsure |

Ending the destination with a slash, as in `dir/`, makes your meaning clear: you mean a folder. If the folder does not exist, `mv` reports an error instead of renaming your file to that name.

## Walk through it
Here is what happens, step by step, when you rename every `.txt` file in a folder to `.md`: `for f in tidy/*.txt; do mv "$f" "${f%.txt}.md"; done`.

1. The shell expands `tidy/*.txt` into the list of matching file names before the loop starts. Files that do not end in `.txt` are not in the list.
2. The loop takes the names one at a time and calls each one `f`.
3. `${f%.txt}` is the standard shell way to remove `.txt` from the end of the value. The result, with `.md` added, is the new name. For `tidy/one.txt`, the new name is `tidy/one.md`.
4. `mv` renames the file. The quotes around `"$f"` and the new name keep a name with a space in it in one piece.
5. The loop repeats for each name, and then ends.

The shell has already done the matching, so the loop runs once per file it found. If nothing matches, the standard says the shell leaves `tidy/*.txt` as typed, and `mv` then complains that the file does not exist. The zsh shell reports an error instead. Check with `ls` first.

To stay safe when the new name might exist, test for it. `[ -e name ]` is true when something is called that, so `if [ -e new ]; then echo exists; else mv old new; fi` moves only when the name is free.

## Try it
These commands write small files into a folder called `tidy` inside the folder you are in, and then move and rename them. The name is unusual so it does not replace a folder of yours. Each block makes its own files, so each one works on its own.

Rename a file. `ls` shows the folder afterwards:

```sh
mkdir -p tidy
printf 'v1\n' > tidy/draft.txt
mv tidy/draft.txt tidy/final.txt
ls tidy
```

You should see this:

```text
final.txt
```

Move two files into a folder. The trailing slash says the destination is a folder:

```sh
mkdir -p tidy/old
printf 'a\n' > tidy/a.log
printf 'b\n' > tidy/b.log
mv tidy/a.log tidy/b.log tidy/old/
ls tidy/old
```

You should see this:

```text
a.log
b.log
```

Rename many files with a loop. Every `.txt` file gets a `.md` name:

```sh
mkdir -p tidy
printf 'one\n' > tidy/one.txt
printf 'two\n' > tidy/two.txt
for f in tidy/*.txt; do mv "$f" "${f%.txt}.md"; done
ls tidy
```

You should see this:

```text
one.md
two.md
```

See the danger. The destination already holds the word `keep`, and `mv` replaces it without a word. The old contents are gone:

```sh
mkdir -p tidy
printf 'keep\n' > tidy/target.txt
printf 'new\n' > tidy/source.txt
mv tidy/source.txt tidy/target.txt
cat tidy/target.txt
```

You should see this:

```text
new
```

Avoid the danger by checking first. This time the name is taken, so nothing moves and the original survives:

```sh
mkdir -p tidy
printf 'keep\n' > tidy/target.txt
printf 'new\n' > tidy/source.txt
if [ -e tidy/target.txt ]; then
  echo 'exists; not moving'
else
  mv tidy/source.txt tidy/target.txt
fi
cat tidy/target.txt
```

You should see this:

```text
exists; not moving
keep
```

## Common confusions
- **Renaming and moving are different commands.** They are one command. A new name renames, and an existing folder moves.
- **`mv` asks before it overwrites.** By default it does not. Add `-i` if you want it to ask.
- **`mv` makes a copy.** It does not. It gives the file a new path, so the old path is gone. Use `cp` to keep the original.
- **The loop sees files that do not match.** The shell chooses the files before the loop starts, so only matching names are in the list.
- **A name with a space is fine without quotes.** It splits into two words. Quote `"$f"`.
- **The trailing slash is decoration.** It tells `mv` that the destination is a folder, and gives an error instead of a surprise when the folder is missing.

## Check yourself
1. What does `mv a.txt b.txt` do when `b.txt` already exists? Answer: it replaces `b.txt` with the contents of `a.txt`, without asking.
2. How do you move a file into a folder and keep its name? Answer: `mv file folder/`.
3. How do you stop `mv` from overwriting by accident? Answer: test with `[ -e name ]` first, or add `-i` so it asks.
4. Why quote `"$f"` in the loop? Answer: so a name with a space stays one word.

## Where next
You can now rename and move files with care. The last guide in this series makes a copy of a folder that you can rely on if something goes wrong.

For the details of each command, see `explain mv`, `explain mkdir`, `explain cp` and `explain ls`.
