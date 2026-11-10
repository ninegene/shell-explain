---
schema_version: 1
id: count-and-sort
title: "Count, Sort, Summarize: Build a Small Report"
subtitle: Turn a list of words into a ranked count, two different ways
summary: Count lines and words, sort them, and build a small report of how often each item appears, using a pipeline of sort and uniq or a single awk program.
series: getting-things-done
order: 2
level: beginner
reading_minutes: 7
updated: 2026-11-09
tags: [report, sorting, counting]
related_commands: [sort, uniq, wc, awk, head]
related_guides: [pipes-and-redirects, find-a-word, input-output-errors]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/sort.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/uniq.html
  - https://pubs.opengroup.org/onlinepubs/9799919799/utilities/wc.html
  - https://www.gnu.org/software/gawk/manual/
---
## Goal
By the end you can count the lines in a file, and build a small report that shows how often each item appears, most common first. You will know two routes to the same report and what each one costs.

## The idea in plain words
Many questions are really "how many of each?" Which page was visited most? Which error showed up most often? Which fruit is in the list most times? The answer is a short report: each item, with a count next to it, biggest first.

A computer does this in three small steps, and each step is a job a small command already does:

1. **Put the same items next to each other.** `sort` puts the lines in order, so equal lines end up side by side.
2. **Collapse the repeats and count them.** `uniq -c` writes each line once, with how many times it occurred. It only compares neighbouring lines, which is why step 1 comes first.
3. **Rank the result.** A second `sort -rn` puts the biggest count first.

This is a **pipeline**, as the earlier guide described. Each command reads what the one before it printed.

A second route does steps 1 and 2 in one go. `awk` keeps a running count for each item in an **array**, a table that maps a name to a number, and prints the table at the end.

## A picture
Read the picture from left to right. The list goes in at the left, and each command changes it a little before the next one reads it.

```diagram
   route 1:  pear      sort       uniq -c      sort -rn
             apple     ----->     ------->     ------->   3 pear
             pear      apple      2 apple      3 pear     2 apple
             fig       apple      1 fig        2 apple    1 fig
             ...       fig        3 pear       1 fig
                       pear ...
   route 2:  awk counts as it reads, then prints one line per item
```

Diagram: route one is a pipeline of three commands. The list of words goes into sort, which puts equal words next to each other. The result goes into uniq with dash c, which writes each word once with its count. A second sort with dash r n puts the biggest count first, giving pear 3, apple 2 and fig 1. Route two is one awk program that counts while it reads and then prints one line per word.

## Another way to see it
Think of counting votes. Route 1 is the paper method: sort the ballots into piles, count each pile, then rank the piles. Route 2 is one person with a tally sheet, adding a mark for each ballot as it comes.

The table compares the two routes.

| Route | Good at | Watch out |
| --- | --- | --- |
| 1: sort, uniq, sort | Short and well known | Padding of the count varies |
| 2: awk | One program, one pass | Output order is unspecified |

With route 2, the order of a `for ... in` loop is not specified, so add a `sort` to put the lines in a known order. With route 1, the standard gives the form of a `uniq -c` line as the count, a space and the line, but systems add leading spaces to line the numbers up, and the amount differs between them. Printing the item first and the count second with `awk '{print $2, $1}'` removes the padding.

## Walk through it
Here is what happens, step by step, when you run `sort fruit.txt | uniq -c | sort -rn` on a file that holds pear, apple, pear, fig, apple and pear, one per line.

1. `sort` reads the six lines and writes them back in order: apple, apple, fig, pear, pear, pear.
2. `uniq -c` compares each line with the one before it. It writes `apple` once with a count of 2, `fig` with 1 and `pear` with 3.
3. The second `sort -rn` sorts by number, which is what `-n` asks for, and `-r` reverses the order. The biggest count comes first.
4. The report appears on your screen, because the last command writes to standard output.

The order of a plain `sort` depends on the language settings of the system, because the standard says it uses the collating sequence of the current locale. Words made only of lower-case letters come out the same everywhere, and that is what the examples use.

To see only the top two, add `| head -n 2`. To count all the lines instead of each item, `wc -l` counts the newline characters in its input.

## Try it
These commands write a small file into a folder called `tally` inside the folder you are in. The name is unusual so it does not replace a folder of yours. Each block makes its own file, so each one works on its own.

Count the lines in the file. `tr -d ' '` removes the padding some systems put before the number:

```sh
mkdir -p tally
printf 'pear\napple\npear\nfig\napple\npear\n' > tally/fruit.txt
wc -l < tally/fruit.txt | tr -d ' '
```

You should see this:

```text
6
```

Show each item once, in order:

```sh
mkdir -p tally
printf 'pear\napple\npear\nfig\napple\npear\n' > tally/fruit.txt
sort tally/fruit.txt | uniq
```

You should see this:

```text
apple
fig
pear
```

Build the report with route 1. The last `awk` step prints the item first and the count second, so the spacing is the same everywhere:

```sh
mkdir -p tally
printf 'pear\napple\npear\nfig\napple\npear\n' > tally/fruit.txt
sort tally/fruit.txt | uniq -c | sort -rn | awk '{print $2, $1}'
```

You should see this:

```text
pear 3
apple 2
fig 1
```

Build the same report with route 2. The array holds a count for each item, and the final `sort` ranks them by the second column, biggest first:

```sh
mkdir -p tally
printf 'pear\napple\npear\nfig\napple\npear\n' > tally/fruit.txt
awk '{c[$1]++} END {for (w in c) print w, c[w]}' tally/fruit.txt \
  | sort -k2,2nr
```

You should see this:

```text
pear 3
apple 2
fig 1
```

See why `uniq` needs sorted input. Without `sort`, no two equal lines are neighbours, so `uniq` removes nothing:

```sh
mkdir -p tally
printf 'pear\napple\npear\nfig\napple\npear\n' > tally/fruit.txt
uniq tally/fruit.txt
```

You should see this:

```text
pear
apple
pear
fig
apple
pear
```

## Common confusions
- **`uniq` removes every repeat.** It removes only repeats that sit next to each other. Run `sort` first.
- **`sort -n` and plain `sort` do the same thing.** Plain `sort` compares as text, so 10 comes before 9. `-n` compares as numbers.
- **The count column looks the same everywhere.** The spacing before the number differs between systems. Use `awk` to reformat it when the exact text matters.
- **`wc -l` counts the lines I can see.** It counts newline characters. A last line without a newline at the end is not counted.
- **The awk loop prints in a fixed order.** It does not. Pipe it to `sort`.
- **A report needs a script.** These are one-line pipelines, and you can retype them when you need them.

## Check yourself
1. Why does `uniq` need `sort` before it? Answer: it compares only neighbouring lines, so equal lines must be next to each other.
2. What does `sort -rn` do? Answer: it sorts by number, biggest first.
3. How do you get only the top three of a ranked report? Answer: add `| head -n 3`.
4. Why add `sort` after the `awk` loop? Answer: the order of the loop is not specified, so the lines could come out in any order.

## Where next
You can now turn a list into a ranked count. The next guide moves and renames files, which is where one careless command can lose something.

For the details of each command, see `explain sort`, `explain uniq`, `explain wc` and `explain awk`.
