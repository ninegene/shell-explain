---
schema_version: 1
id: sed
title: sed
summary: Transform text one line at a time using editing rules.
aliases: [gsed]
tags: [text, replace, files]
implementations: [gnu, apple, busybox]
references:
  - https://www.gnu.org/software/sed/manual/
  - https://man.freebsd.org/cgi/man.cgi?query=sed&sektion=1
  - https://busybox.net/downloads/BusyBox.html
---
## Common
### Overview
Use `sed` to print transformed text without changing the original file. A substitution rule has the form `s/pattern/replacement/`. Preview the result before choosing an in-place edit.

### Examples
Replace the first occurrence of old on each line:

```sh
sed 's/old/new/' notes.txt
```

Replace every occurrence on each line:

```sh
sed 's/old/new/g' notes.txt
```

Print only lines 2 through 4:

```sh
sed -n '2,4p' notes.txt
```

### Option: -n
Suppress automatic printing. Pair it with a rule such as `p` to explicitly print selected lines.

### Option: -e
Add an editing rule. Repeat `-e` to apply multiple rules in order.

### Common mistakes
The substitution flag `g` changes every match on a line, not the number of files edited. Replacement text treats `&` specially: it expands to the matched text. Preview rules containing punctuation carefully.

### Portability
In-place editing is not uniform across implementations. The examples below create backup files rather than discarding the original. Version-sensitive extensions are outside the common examples.

## GNU
### GNU examples
Edit a disposable copy in place and retain the previous contents as notes.txt.bak:

```sh
sed -i.bak 's/old/new/g' notes.txt
```

### Option: -i, --in-place
Write changes back to the file. GNU accepts `-i` without a backup suffix, but use a backup while learning. An existing backup path can be overwritten; practice in a temporary directory.

## Apple
### Apple examples
Retain a backup while editing a disposable file:

```sh
sed -i '.bak' 's/old/new/g' notes.txt
```

### Option: -i
Apple sed expects an argument for the backup suffix. The empty argument in `-i ''` means no backup; it is not the same syntax as GNU's usual suffix-free form. Use a backup until you have verified the edit.

## BusyBox
### BusyBox examples
Keep a backup while editing a disposable file:

```sh
sed -i.bak 's/old/new/g' notes.txt
```

### Option: -i
Edit in place with an optional attached backup suffix. BusyBox features depend on the build; check its local help before relying on an extension.
