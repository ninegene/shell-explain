---
schema_version: 1
id: find-a-word
title: "Needle in a Haystack: Find a Word in Many Files"
subtitle: Three routes to every line that contains a word
summary: Find the lines that contain a word in one file or in a whole folder of files, and choose between grep on its own, grep with a folder, and find with grep.
series: getting-things-done
order: 1
level: beginner
reading_minutes: 6
updated: 2026-11-09
tags: [search, text, files]
related_commands: [grep, find, sort]
related_guides: [pipes-and-redirects, exit-codes, files-and-folders]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/grep.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/find.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html
---
## Goal
By the end you can show every line that contains a word, in one file or in a whole folder of files, and you can say which of three routes suits the job.

## The idea in plain words
You have a word, and somewhere in a pile of files it appears. You want the lines that contain it, and the name of the file each line came from. A computer can read every line much faster than you can, so you ask a command to do the reading.

The command for this job is `grep`. It reads text, keeps the lines that match a pattern, and prints them. There are three ways to point it at your files:

1. **Name the file.** `grep needle notes.txt` searches one file.
2. **Name a folder and add `-r`.** `grep -r needle folder` searches every file under that folder, and the folders inside it.
3. **Let `find` pick the files first.** `find` lists files that match a rule, such as a name ending in `.txt`, and hands those names to `grep`.

All three print the same kind of answer. The difference is which files get read.

## A picture
Read the picture from top to bottom. Each route starts from the same goal and differs in how it chooses the files.

```diagram
   goal: every line with "needle"

   route 1   grep needle FILE            one file you already know
   route 2   grep -r needle FOLDER       every file under a folder
   route 3   find FOLDER -name '*.txt' -exec grep needle {} +
                                         only the files find picks
```

Diagram: three routes to the same goal, which is every line containing the word needle. Route one is grep, the word, then one file name, and it searches the one file you already know. Route two is grep dash r, the word, then a folder, and it searches every file under that folder. Route three is find with a folder and a name rule, followed by an exec option that runs grep on the files find picks, so it searches only those files.

## Another way to see it
Think of a library. Route 1 is opening one book you know and running your finger down the page. Route 2 is walking every shelf and checking every book. Route 3 is asking the librarian for only the paperbacks, and checking those.

The table compares the three routes. It lists what each is for and what to watch for.

| Route | Use it when | Watch out |
| --- | --- | --- |
| 1: grep FILE | You know the file | Only that file is read |
| 2: grep -r | Any file under a folder | POSIX grep lists no -r |
| 3: find + grep | Only some files | The longest to type |

The standard's list of `grep` options has no `-r`, so a script that must run everywhere cannot count on it. The macOS, GNU and BusyBox `grep` programs tested in this project all accepted it.

## Walk through it
Here is what happens, step by step, when you type route 3: `find wordhunt -name '*.txt' -exec grep -n needle {} +`.

1. The shell hands `find` the folder `wordhunt`, the rule `-name '*.txt'` and the action after `-exec`. The quotes stop the shell from expanding `*.txt` before `find` sees it.
2. `find` walks the folder and every folder inside it, and keeps the files whose names end in `.txt`.
3. The `{} +` ending tells `find` to gather the names it kept and run `grep` once with all of them as arguments. The standard describes this form as running the utility once per set of pathnames.
4. `grep` reads each file and keeps the lines that contain `needle`. With more than one file, it puts the file name and, because of `-n`, the line number in front of each line.
5. The order depends on the system, because `find` does not promise one. Adding `| sort` at the end makes the output the same everywhere.

If nothing matches, `grep` prints nothing and exits with status 1. That is not an error. The standard says 0 means a line was selected, 1 means none was, and anything above 1 means an error. The guide on exit codes explains how to use that in a script.

## Try it
These commands write a few small files into a folder called `wordhunt` inside the folder you are in, and then search them. The name is unusual so it does not replace a folder of yours. Each block makes its own files, so each one works on its own.

Search one file for a word. Only the matching line is printed:

```sh
mkdir -p wordhunt/docs
printf 'alpha\nbeta needle\n' > wordhunt/a.txt
grep needle wordhunt/a.txt
```

You should see this:

```text
beta needle
```

Search a whole folder with `-r`, and show line numbers with `-n`. The capital N in `Needle caps` does not match, because a search is case sensitive. Sorting the output makes the order the same everywhere:

```sh
mkdir -p wordhunt/docs
printf 'alpha\nbeta needle\n' > wordhunt/a.txt
printf 'needle here\nNeedle caps\n' > wordhunt/docs/b.txt
grep -rn needle wordhunt | sort
```

You should see this:

```text
wordhunt/a.txt:2:beta needle
wordhunt/docs/b.txt:1:needle here
```

Ignore case with `-i`, and list only the names of the files that match with `-l`:

```sh
mkdir -p wordhunt/docs
printf 'alpha\nbeta needle\n' > wordhunt/a.txt
printf 'needle here\nNeedle caps\n' > wordhunt/docs/b.txt
printf 'nothing\n' > wordhunt/docs/c.log
grep -rli needle wordhunt | sort
```

You should see this:

```text
wordhunt/a.txt
wordhunt/docs/b.txt
```

Let `find` pick only the `.txt` files, then search those. The `.log` file is never read:

```sh
mkdir -p wordhunt/docs
printf 'alpha\nbeta needle\n' > wordhunt/a.txt
printf 'needle here\nNeedle caps\n' > wordhunt/docs/b.txt
printf 'needle in a log\n' > wordhunt/docs/c.log
find wordhunt -name '*.txt' -exec grep -n needle {} + | sort
```

You should see this:

```text
wordhunt/a.txt:2:beta needle
wordhunt/docs/b.txt:1:needle here
```

Ask a yes-or-no question. With `-q`, `grep` prints nothing and only sets its exit status, so the `||` runs when the word is missing:

```sh
mkdir -p wordhunt
printf 'alpha\nbeta needle\n' > wordhunt/a.txt
grep -q zebra wordhunt/a.txt || echo 'not found'
```

You should see this:

```text
not found
```

## Common confusions
- **No output means something broke.** It usually means nothing matched. Check the exit status: 1 means no match, and above 1 means an error.
- **The word is plain text.** The pattern is a pattern, so a dot matches any character. Add `-F` to search for the exact characters you typed.
- **A capital letter matches a small one.** It does not, unless you add `-i`.
- **`grep -r` reads only text files.** It reads every file under the folder, including binary ones.
- **I can leave out the quotes around the pattern.** Quote it, so the shell does not change it before `grep` sees it. Quote `'*.txt'` for `find` for the same reason.
- **The output order is fixed.** It is not. Add `| sort` when the order matters.

## Check yourself
1. Which command prints the lines of one file that contain a word? Answer: `grep word file`.
2. How do you search every file under a folder? Answer: `grep -r word folder`, or `find` with `-exec grep`.
3. What does exit status 1 from `grep` mean? Answer: no line matched. It is not an error.
4. Why quote `'*.txt'` in the `find` command? Answer: so the shell does not expand it, and `find` receives the pattern itself.

## Where next
You can now find a word wherever it hides. The next guide takes a list of things, counts them and sorts them into a small report.

For the details of each command, see `explain grep`, `explain find` and `explain sort`.
