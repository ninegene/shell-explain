# Learning guides: teaching the shell from the ground up

Status: design accepted; the format, parser, export and a sample guide are implemented (see the Guides section of [content-format.md](content-format.md)), the `explain learn` and `explain glossary` commands are implemented, and so are the website's `/learn/` pages, glossary page, home-page "Start here" section and "Learn more" links on command pages. Series 1 (five guides), six of the nine Series 2 guides (terminal, shell and command; anatomy of a command; files and folders; the three streams; pipes and redirects; exit codes) and 41 glossary terms are written, and `Try it` commands are run by the tests on the host and in the Linux guest. Not yet built: guides and glossary terms in the website's search, and the other three guides in Series 2 (text and binary files; PATH and the environment; permissions) and Series 3 and 4. This document records what to build and in what order, so each step can be reviewed before it is built. It extends [stage 4](roadmap.md) and must be finished before the public release in stage 5, because the home page will lead with it.

## Why this exists

Command pages answer "how do I use `grep`?" A beginner is stuck one step earlier: what is a terminal, what is a shell, why does `ls` show different things on a Mac and on Ubuntu, and what does "command not found" mean? The command pages assume those answers. The guides supply them.

The same pages also explain a fact that confuses everyone who moves between systems: the same command exists in several versions (GNU, BSD/macOS, BusyBox, mawk and others) because of history, licensing and standards. That history is the reason this project selects documentation by implementation instead of by operating system.

## Audience

- Middle- and high-school students, and non-technical professionals who need to use a terminal.
- Assume no prior knowledge. Do not assume words like kernel, process, stdin or path are known; define each one the first time it appears on every page.
- Keep the experienced user in mind too: every guide must be skimmable, link to the exact command pages, and never talk down.

## Teaching principles

1. **Layers, from input to output.** The organizing idea is a stack: person, keyboard and screen, terminal, shell, program, operating system, hardware, with firmware starting the whole thing at power-on. Every guide says which layer it is about and shows where that layer sits in the stack.
2. **Goal first.** A command is a way to get a result: text on the screen, a file written, an email sent, a web service called. Each guide starts from the result the reader wants and works backwards to the typing.
3. **One idea, several angles.** Repeat important ideas on purpose, in different forms: a plain-words explanation, an everyday analogy, a diagram, a table, a step-by-step trace, a real command with its output, and a list of common mistakes. Duplication across guides is intended; do not remove it to save space.
4. **Show, then name.** Show a small working example before introducing the technical word for what it did.
5. **Visuals carry meaning.** Prefer a diagram or table when it explains faster than a paragraph. Every visual has a caption and a plain-text description, so it works in a terminal, for a screen reader, and without images.
6. **Honest sources.** History and "why" claims are limited to what a cited reference or a command we ran supports, the same rule as the command pages. When sources disagree, say so (the command pages already do this for the origin of `cut`).
7. **Safe by default.** Every "try it" command is read-only or works on a throwaway folder. Never ask a beginner to run something destructive.

## The core mental model

This stack is the picture the whole series builds on. Guide 1 introduces it, and every later guide refers back to it.

```text
 YOU                      type:  ls
   |  keys
   v
 KEYBOARD and SCREEN      the hardware you touch
   |
   v
 TERMINAL                 a window that shows text, passes keys on
   |  text
   v
 SHELL                    reads your line, finds the program, starts it
   |  "run ls"            (bash, zsh, ...)
   v
 PROGRAM (command)        ls, grep, awk ... does one job
   |  asks politely: "list this folder"   (system calls)
   v
 OPERATING SYSTEM         the kernel: shares out CPU, memory, disk
   |
   v
 HARDWARE                 CPU, memory, disk (firmware starts it)
```

A second angle on the same stack, for readers who think in analogies:

| Layer | Restaurant analogy |
| --- | --- |
| You | The customer with an order |
| Terminal | The ordering tablet: shows the menu, sends what you tap |
| Shell | The waiter: understands the order, tells the right cook |
| Program | A specialist cook: one dish done well |
| Kernel | The kitchen manager: shares stoves, fridges and cooks fairly |
| Hardware | The stoves, fridges and counters |

A third angle is the trace of one real command, with the output the reader will see:

```text
1. You press l, s, Enter           the terminal sends them to the shell
2. The shell reads ls               it searches the folders in PATH
3. The shell asks the kernel       "start the program at /bin/ls"
4. The kernel loads the program    into memory, gives it some CPU time
5. ls asks the kernel              "what files are in this folder?"
6. The kernel reads the disk       and hands the names back
7. ls writes the names             to its output
8. The terminal draws them         on your screen
```

## Guide catalog

Titles are written to be inviting, with the plain topic underneath. Each guide has an `id` (the CLI name), a catchy title, and one line saying what the reader can do afterwards. The order is the suggested reading order; each guide also works on its own.

### Series 1: How a computer works, just enough

| id | Title | After reading you can... |
| --- | --- | --- |
| `big-picture` | From Keypress to Screen | Name the layers between your finger and the output |
| `computer-body` | The Computer's Body: CPU, Memory, Disk | Explain what each part does and why disk is slower than memory |
| `who-is-boss` | Firmware, Kernel, Programs: Who Is the Boss? | Tell hardware, firmware, operating system and software apart |
| `what-is-an-os` | What an Operating System Really Does | Say why programs ask the kernel instead of touching hardware |
| `processes` | Programs That Are Awake: Processes | Explain a process, how it starts, and how to see running ones |

### Series 2: Terminal, shell and commands

| id | Title | After reading you can... |
| --- | --- | --- |
| `terminal-shell-command` | Terminal, Shell, Command: Three Things People Mix Up | Say which one a name like "bash" or "iTerm" refers to |
| `anatomy-of-a-command` | Anatomy of a Command | Split a command into name, options and arguments |
| `files-and-folders` | The Map Your Shell Walks: Files and Folders | Read a path, tell absolute from relative, and move around |
| `text-and-binary` | Words or Wiring: Text Files and Binary Files | Explain why you can read one with `cat` and not the other |
| `input-output-errors` | Three Streams Every Program Has | Name input, output and error, and say where each goes |
| `pipes-and-redirects` | Snapping Commands Together | Build a pipeline and redirect output into a file |
| `exit-codes` | How a Command Says Yes or No | Read an exit status and use `&&` and `||` |
| `path-and-environment` | Why "Command Not Found" Happens | Explain PATH and environment variables |
| `permissions` | Who May Read, Write, Run | Read a permission string and explain what `chmod` changes |

### Series 3: Why are there so many versions?

| id | Title | After reading you can... |
| --- | --- | --- |
| `family-tree` | A Short History: Unix, BSD, GNU, Linux | Place each name on a timeline and say who made what |
| `whos-who` | Who's Who: Unix, BSD, FreeBSD, macOS, GNU, Linux, Ubuntu, BusyBox | Say what each name is: a kernel, a toolset, a whole system, or a company's product |
| `same-command-many-versions` | One Command, Many Versions | Explain why `sed -i` differs on a Mac and on Linux |
| `distributions` | Ubuntu, Debian, Alpine, Fedora: What Is a Distribution? | Say what a distribution adds and what a package manager does |
| `which-tool-do-i-have` | Which Tool Do I Have? | Find out the installed implementation in under a minute |

### Series 4: Getting things done (goal first)

Each guide starts from a result and shows two or three different routes to it, with their trade-offs. Written after Series 1 to 3.

| id | Title | Goal |
| --- | --- | --- |
| `find-a-word` | Find a Word in Many Files | Text on screen |
| `count-and-sort` | Count, Sort, Summarize | A small report |
| `rename-and-move` | Tidy Up Files | Files in new places |
| `make-a-backup` | Make a Backup | A new archive file |
| `fetch-from-the-web` | Ask a Web Service | A response from another computer |
| `send-a-message` | Send an Email from the Terminal | An email arrives (with the honest warnings) |

### Glossary: Words You Will See

One page of short definitions for the words beginners meet and are assumed to know: kernel, shell, terminal, process, path, directory, stdin, stdout, stderr, pipe, flag, argument, root, sudo, package, repository, binary, script, environment variable, and others. Each entry is a sentence or two, links to the guide that explains it, and can be looked up from the CLI.

## How a guide is written

Every guide follows this outline. A section may be short, but none is skipped.

1. **Goal.** One or two sentences: what you will be able to do, and the layer this guide is about.
2. **The idea in plain words.** No jargon; any technical word is defined in the same sentence.
3. **A picture.** A diagram with a caption and a text description.
4. **Another way to see it.** An analogy, or a table that compares things side by side.
5. **Walk through it.** A numbered trace of one real example, step by step.
6. **Try it.** Safe commands the reader can run, with the exact output they should see.
7. **Common confusions.** The wrong ideas beginners actually have, and the fix.
8. **Check yourself.** Three or four questions with answers underneath.
9. **Where next.** The next guide, and the command pages that apply.

Style rules, checked in review:

- Short sentences; one idea per sentence; active voice; "you" for the reader.
- Define a term the first time it appears on every page, even if an earlier guide did.
- Prefer a concrete example before a rule. Use the same small set of example files as the command pages.
- A diagram is at most 72 columns wide so it fits a terminal, uses plain characters, and never carries information that is not also in the text.
- A table has at most four columns and short cells, so it reads in a terminal and on a phone.
- Every claim about history, standards or licensing links to a reference in the front matter.

## Content model

Guides are a second kind of authored Markdown, separate from command pages, loaded by the same Go parser.

- **Location:** `content/en/guides/<id>.md`, one file per guide, and `content/en/glossary/terms.md` for the glossary. The command loader globs only `content/en/*.md` (`fs.Glob` does not descend into subdirectories), so neither the guides nor the glossary can be mistaken for a command page or claim a command ID. Nothing may be placed directly in `content/en/` except command pages. Each new kind gets its own loader that globs its own directory: `en/guides/*.md` for guides and `en/glossary/*.md` for the glossary.
- **Front matter** (all required unless noted): `schema_version`, `id` (matches the filename), `title` (the catchy title), `subtitle` (the plain topic), `summary`, `series`, `order` (position in the series), `level` (`beginner` for now), `reading_minutes`, `updated` (the same `YYYY-MM-DD` rule as command pages), `related_commands` (each must be an existing command page), `related_guides` (each must exist), `references` (absolute HTTP(S) links), optional `tags`.
- **Body:** level-two headings are the nine outline sections above, in order; level-three headings are allowed inside them. Required sections are checked by name, like the command pages' required sections.
- **Markdown subset additions.** Guides need two things the command format deliberately rejects today:
  - **Tables**, using the standard pipe syntax, at most four columns and at most 78 columns wide when rendered. Rendered as aligned text in the CLI and as a real table on the web, with horizontal scrolling on narrow screens.
  - **Diagrams**, as a fenced block with the language `diagram`. The block must be followed by a paragraph that starts with `Diagram:` and describes it in words. The web shows it in a labelled, non-copyable block; the CLI prints it unchanged.
  - Everything else stays as today: no raw HTML, no images, no nested headings beyond level three.
- **Embedding:** `content/embed.go` currently embeds `en/*.md`, which skips subdirectories, so the offline CLI would contain no guide or glossary. The directive must become `//go:embed en/*.md en/guides/*.md en/glossary/*.md`, with a test that fails if a content file on disk is missing from the embedded filesystem. This keeps one authored collection supplying both the CLI and the website.
- **Command-page front matter:** command pages gain an optional `related_guides` list (each entry must be an existing guide) that drives the `Learn more:` lines and boxes. The parser rejects unknown keys today, so this needs a parser change, a `docs/content-format.md` update and tests. The schema version stays 1 for now: the format is documented as extensible during stage 4, and the field is optional so no existing page changes. Revisit the version number before the public release.
- **Export:** the exporter adds `guides` and `glossary` next to `entries`, with the same per-section IDs, HTML, and exact code strings, so the web uses the CLI's exported content. Search records include guides and glossary terms.
- **Quality checks in the parser:** required sections present and in order, diagram width at most 72 columns, every diagram followed by its `Diagram:` description, tables at most four columns, links and related names resolve, `updated` valid.

## CLI

The CLI should let a beginner discover and read the guides without leaving the terminal.

| Command | What it does |
| --- | --- |
| `explain learn` | Lists the series and guides: number, catchy title, one-line summary, reading time |
| `explain learn <id or number>` | Shows that guide, with the pager behaviour below |
| `explain learn <id> --section <name>` | Shows one section, such as `Try it` |
| `explain glossary` | Lists all terms with one-line definitions |
| `explain glossary <term>` | Shows one definition, its guide, and related commands |
| `explain search <words>` | Also finds guides and glossary terms, marked with their type |

Reading behaviour:

- Long output goes through the user's pager when standard output is a terminal and a pager exists, with `--no-pager` to turn that off. When output is piped, print everything.
- Every guide ends with the same footer: `Next:`, `Previous:`, the series list, and `Commands:` linking to the related command pages (`explain grep`).
- `--format=json` works for guides and the glossary, with the same fields the web uses.
- Command pages show `Learn more:` lines (from a `related_guides` field) in `--details`, so someone reading about `|` is pointed to the pipes guide.
- Width follows `--width` and the terminal, and diagrams are never wrapped.

## Website

- **Home page.** A "Start here" section above the command search: the four series as cards with their catchy titles and a one-line promise, each linking into its first guide. The command search stays one click away.
- **`/learn/`.** The series and guides in reading order, with reading times.
- **`/learn/<id>/`.** The guide, a table of contents, series navigation with previous and next, related command links, and links from glossary words to their definitions. Diagrams and tables work at mobile widths.
- **`/learn/glossary/`.** All terms with stable anchors.
- **Command pages.** A "Learn more" box when the page lists related guides.
- **Search.** The existing search includes guides and glossary terms.

## Verification

- **Same content, both consumers.** The existing parity test extends to guides: CLI sections, HTML and exact code match the static page.
- **Try-it commands run.** Deterministic "Try it" blocks are executed in the shared fixture on the available implementations and compared with the shown output, like the command examples. Blocks that cannot be deterministic (time, process IDs) are marked and reviewed only.
- **Accessibility and layout.** Playwright checks for diagrams and tables on desktop and mobile widths: no horizontal page overflow, labels present, keyboard navigation, axe clean.
- **Language.** Review checks against the style rules above. A simple automated check reports unusually long sentences and the use of a defined term before its definition on a page, as a warning, not a failure.
- **Sources.** History claims are checked against their references in review, with the date recorded in `updated`.

## Delivery plan

Each step is a reviewed pull request, in this order.

1. **This design.** Review and adjust the catalog, outline and open questions.
2. **Format.** Guide and glossary parsing, tables and diagrams in the Markdown subset, validation, the embed change, the optional `related_guides` field on command pages, export, and tests. One sample guide (`big-picture`) to exercise everything.
3. **CLI and web shell.** `learn` and `glossary`, pager, search, the `/learn/` pages, the home page section, and the parity and accessibility tests.
4. **Content, Series 1 and 2.** Written in batches of three to five guides, each with its diagrams, tables and runnable examples.
5. **Content, Series 3 and the glossary.** History and implementation guides, with references.
6. **Content, Series 4 and `Learn more` links** on the command pages.

## Open questions

1. **Placement.** Treat this as a stage 4 deliverable that must land before the exit review (recommended, since the home page leads with it), or as a separate stage between 4 and 5.
2. **Series names.** The four series titles above are working names; confirm or replace them.
3. **Scope of Series 4.** Whether "send an email" and "ask a web service" belong in the first release, given they depend on tools (`curl`, `mail`) with more variation and risk than the rest.
4. **Pager default.** Page automatically on a terminal (recommended) or only with `--page`.
5. **Quizzes.** The "Check yourself" answers are shown inline (recommended); an interactive quiz is out of scope.

## Out of scope

Interactive terminals or in-browser sandboxes, video, translations, user accounts or progress tracking, and any feature that runs commands for the reader. These can be reconsidered after the public release.
