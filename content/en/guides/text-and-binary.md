---
schema_version: 1
id: text-and-binary
title: "Words or Wiring: Text Files and Binary Files"
subtitle: What is inside a file, and why some files can be read and others cannot
summary: Learn that every file is a row of bytes, what makes a file text, and why printing a binary file gives you gibberish.
series: terminal-shell-commands
order: 4
level: beginner
reading_minutes: 6
updated: 2026-11-08
tags: [basics, files, text]
related_commands: [cat, head, wc]
related_guides: [files-and-folders, input-output-errors, big-picture]
references:
  - https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap03.html
  - https://man7.org/linux/man-pages/man7/ascii.7.html
  - https://man7.org/linux/man-pages/man7/utf-8.7.html
---
## Goal
By the end you can say what a byte is, what makes a file a text file, and why `cat` on a picture or a program prints gibberish. You can also look at the bytes of a file yourself.

## The idea in plain words
A computer stores everything as numbers. The smallest piece a file is made of is a **byte**: a number from 0 to 255. A file is a long row of bytes.

The bytes do not say what they mean. A program decides how to read them. The same bytes can be a word to one program and a picture to another.

A **text file** is a file whose bytes stand for characters, arranged in lines. The standard puts it this way: characters organized into lines, with no zero bytes in them. Notes, source code and log files are text files.

A **binary file** is a file whose bytes are not meant to be read as characters. Photos, programs and compressed archives are binary files. The standard does not draw a hard line between the two. It only says that many tools give meaningful results only on text.

An **encoding** is the rule that says which number stands for which character. In the old code called **ASCII**, the letter `A` is the number 65. In **UTF-8**, which is the common way to write text today, plain ASCII characters keep the same numbers, and other characters, such as `é`, take two or more bytes.

## A picture
Read the picture from top to bottom. The bytes in the file never change. Only the way they are read changes.

```diagram
   bytes in the file         104      105       10
                              |        |        |
   read as text              'h'      'i'    new line
                              |        |        |
   shown by the terminal       h        i    (next line)

   read as numbers           104      105       10
   (what a tool like od shows)
```

Diagram: three bytes with the numbers 104, 105 and 10. Read as text, 104 is the letter h, 105 is the letter i, and 10 is a new line, so a terminal shows hi and moves to the next line. Read as numbers, the same bytes are just 104, 105 and 10, which is what a tool such as od prints.

## Another way to see it
Think of a row of beads. A text file is a row of beads strung to spell words. A binary file is a row strung to make a pattern for something else, such as a picture or a machine. You can read the words aloud, but reading the machine pattern aloud gives noise.

| Kind | What the bytes mean | Examples |
| --- | --- | --- |
| Text file | Characters in lines | Notes, source code, logs |
| Binary file | Anything else | Photos, programs, archives |

What decides the kind is the content, not the name. A file called `notes.txt` that holds a picture is still binary.

## Walk through it
Here is what happens, step by step, when you run `cat` on a text file and then on a program.

1. You run `cat greeting.txt`. The file holds the three bytes 104, 105 and 10.
2. `cat` reads the bytes and writes them to its standard output without changing them.
3. The terminal receives the bytes and draws them as characters. You see `hi`.
4. Now you run `cat` on a program file. `cat` does the same thing: it copies the bytes out.
5. The terminal receives bytes that were never meant to be characters. It draws strange symbols and may beep.

`cat` did nothing different in the two cases. The difference was in the bytes.

## Try it
These commands only look at things. Each one works on its own.

See the three bytes of a short text file, written in hexadecimal. Hexadecimal is a way of writing numbers in which 68 means 104. The `od` command prints bytes. Then `tr` removes the spaces and new lines, because different systems lay the output out a little differently, and `echo` adds one new line at the end:

```sh
printf 'hi\n' | od -An -tx1 | tr -d ' \n'; echo
```

You should see this:

```text
68690a
```

That is the three bytes 68, 69 and 0a. They are the letters `h` and `i`, and a new line.

See that the letter `A` is stored as the number 65, as ASCII says. Here the same tool prints the byte as an ordinary number:

```sh
printf 'A' | od -An -tu1 | tr -d ' \n'; echo
```

You should see this:

```text
65
```

See that one character can take more than one byte. The word `café` has four characters, but the letter `é` takes two bytes in UTF-8, so the word has five bytes:

```sh
printf 'caf\303\251' | od -An -tx1 | tr -d ' \n'; echo
```

You should see this:

```text
636166c3a9
```

Look at the first four bytes of a program on your computer. You will see numbers and not words. The numbers differ between systems:

```sh
head -c 4 /bin/sh | od -An -tx1
```

## Common confusions
- **A text file is a file that ends in `.txt`.** It is not. The name is only a label. What makes a file text is what is inside it.
- **Binary means broken or secret.** It does not. It only means the bytes are not meant to be read as characters. Compiled programs are binary files, but a script, which is a program written as lines of text, is a text file.
- **A character is always one byte.** Not always. In UTF-8, plain English letters take one byte, but many other characters take two or more.
- **`cat` works on any file.** It copies bytes from any file. The result is only readable when the file is text.
- **Documents from word processors are text.** Usually they are not. They hold formatting and often other data in a binary form, even when the words inside are readable.

## Check yourself
1. What is a byte? Answer: a number from 0 to 255. A file is a row of them.
2. What makes a file a text file? Answer: its bytes stand for characters organized in lines, and not for something else.
3. Why does `cat` on a program print strange symbols? Answer: the bytes were never meant to be characters, but the terminal tries to draw them as if they were.
4. How many bytes does the word `café` take in UTF-8? Answer: five, because `é` takes two bytes.

## Where next
The next guide, `explain learn input-output-errors`, shows how programs read and write streams of these bytes.

To work with text files, look at `explain cat`, `explain head` and `explain wc`.
