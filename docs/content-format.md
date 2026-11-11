# Shared Markdown content format

Schema version 1 is implemented in internal/catalog and used by the CLI and Astro export. It is an initial contract that stage 4 may extend for the command collection, not a frozen public API.

## Authoring source

Keep shell explanations in `.md` files, with YAML front matter and a restricted Markdown body. English is the initial language. Use one canonical file per command entry, for example `content/en/grep.md`; its filename must match its ID. Front matter holds structured metadata, not a second copy of prose. The body holds the explanation, examples, and portability details shown by both consumers.

Metadata example (the grep entry provides a working full page):

```yaml
---
schema_version: 1
id: grep
title: grep
summary: Search text for lines that match a pattern.
aliases: [ggrep]
tags: [text, search, logs]
implementations: [gnu, apple, busybox]
updated: 2026-11-08
references:
  - https://www.gnu.org/software/grep/manual/
---
```

All fields shown above are required except aliases, which may be empty. `updated` is the date, written `YYYY-MM-DD` and checked as a real calendar date, when the page's content was last reviewed or changed. Set it by hand to the current date whenever you edit a page's claims, examples, or version statements, because the Background and About sections quote tool versions and history that go stale. The date is shown on the website page and in `--details` output, and is exported in JSON; it is part of the page source, so changing it changes the content revision like any other edit, and it is never compared with the clock, so builds stay deterministic. IDs, tags, and implementation IDs use lowercase letters/digits/hyphens starting with a letter; titles equal the command ID. Aliases follow the same rule, plus the registered symbolic spellings in `symbolicNames` (currently `[`, for `test`). A symbolic spelling is only ever an alias: the canonical ID, filename, and web URL stay `test`, and the CLI looks it up directly (`explain '['`). Unknown metadata keys, duplicate YAML keys, duplicate names, unsupported schema versions, and unsafe reference URLs are rejected. References and inline links must use absolute HTTP(S) URLs without userinfo. Plain scalar metadata must not contain control characters.

An implementation list declares authored coverage; it does not prove every implementation supports every option. Supported IDs come from one registry, `catalog.Implementations`, which the parser, CLI help, JSON export (`implementations`), and website labels all read: gnu, apple, busybox, mawk, procps, bash, zsh, and upstream. `bash` and `zsh` are shell builtins; an entry that declares only shells is treated as a builtin, which no PATH probe can describe. `upstream` is for tools with one real implementation (jq, git, curl, and similar). It cannot be combined with other implementations and is selected without detection, so it is never presented as GNU, BusyBox, or an OS guess. FreeBSD/OpenBSD need future content and validation rules before becoming selectable variants; add a registry row with its heading when real content exists. Aliases are lookup names, not permission to assume the installed implementation; `ggrep` still needs selection evidence.

Each page must contain Common sections named Overview, Background, Examples, Common mistakes, and Portability, and each declared implementation view must contain an About section. History and implementation background are kept separate from the examples: the Common section named `Background` (where the command came from and how its implementations relate), and one `About <name>` section at the start of each implementation view (what that implementation is, where it ships, and how it differs, for example `About GNU grep`). `Background` is valid only in Common, `About …` only in an implementation view, and a view may have at most one. Both are required: a page without `Background`, or a declared view without an `About` section, is rejected. State only history and differences that the page's references or an executed check support; leave out anything unverified. Other sections can explain use cases or options. References are rendered from metadata in both interfaces. Repeated content belongs in Common; variant-specific content uses explicit headings.

## Stage 4 conventions for the command collection

- **Builtin and program pages.** A name that is both a shell builtin and a program (`echo`, `printf`, `pwd`, `true`, `false`, `test`, `kill`) declares `bash` next to the program views. The Common examples must give the same result in the Bash builtin and in every program view: the Go test and the Linux check run them once under Bash and once per executable through a shell function that calls that executable, so the builtin does not answer instead. Differences go in the views.
- **Single-project tools** (`curl`, `git`, `jq` and similar) declare only `upstream`. Their examples run on the tool found on `PATH`, in an empty folder with a clean home folder, and the test skips the page when the tool is not installed.
- **`procps`** is the implementation label for the process tools of most Linux systems (ps, pgrep, pkill, top, free), which are neither GNU coreutils nor BusyBox.
- **Reviewed-only examples.** An example that needs the network, an account, a daemon, a terminal, administrator rights or SELinux, or that would signal processes by a common name, is recorded in `internal/examplecheck` as `notRun`. The page says so and the evidence table in docs/command-collection.md counts it as reviewed only. A network example is checked by hand against the live service on the date the page gives. Examples that call a cloud service, a deployment or an account are never run.
- **Companion programs.** Where a command runs another by name (GNU `diff3` calls `diff`, `locate` and `updatedb` belong together, `colima` starts `limactl`), `examplecheck.Companions` puts the matching program on the example path.
- **Names that are not options.** An option section must start with a dash. Subcommands (`go build`, `npm install`) and operands (`dd if=`) go in a plain section such as "Subcommands" or "Operands".
- **Evidence sentences.** The Portability section names the implementation and version the Common examples ran on, and says when they did not run.

## Variant sections

Level-two headings identify views: `## Common`, `## GNU`, `## Apple`, `## BusyBox`, `## Bash`, `## Zsh`, or `## Upstream`. Level-three headings identify sections within a view. Every declared implementation needs its matching view, and undeclared/unknown/duplicate views are rejected. Text belongs inside a level-three section; duplicate or empty sections are rejected.

An option section heading declares its spellings explicitly, for example `### Option: -i, --in-place`. Each spelling must be unique within its implementation view. Lookup uses these declarations, not prose substring matching. If an option takes a value, explain the value in the section body; the current corpus does not need an extra structured value type yet. The exporter assigns deterministic per-page section IDs used for browser anchors and JSON parity.

Compose the common section with the chosen implementation section in both CLI and web. Common examples must actually work across the implementations the page claims. If a command differs too much, keep the common overview brief and author explicit per-implementation examples.

The first corpus is implementation-level, not a certified option-by-version database. Every output includes a version-policy warning. Do not infer exact support from a banner version. Introduce structured version gates and associated fixtures before adding examples that require newer versions than the entry's baseline; do not silently claim compatibility.

## Markdown subset

Within sections, the accepted CommonMark subset is paragraphs, ordered/unordered lists, emphasis, inline code, fenced code, and links. Fenced code languages are sh, bash, console, and text. Goldmark recognizes fence lengths, indentation, and list containers; the format additionally requires an explicit matching closing fence. Section boundaries come from top-level Markdown headings, so heading-like text inside code remains literal. Raw HTML, images, indented code, nested headings, tables, blockquotes, autolinks, MDX/JSX, and other nodes are rejected. Project planning documents are not subject to this explanation-content restriction. Preserve code exactly while wrapping prose for terminal width. Inline code follows CommonMark code-span normalization (line breaks become spaces and delimiter-padding rules apply); its internal spaces and tabs are then preserved, and the span stays unbroken even when it exceeds terminal width. Adjacent punctuation stays attached.

Every example needs a purpose and a short explanation of important options/arguments. Include expected results where deterministic. Explain placeholders such as `FILE` consistently. Use generic fixtures and quote patterns/paths appropriately. Avoid shell prompts in code that users are expected to copy.

Use `explain grep --option=-n` for exact option lookup. Content for a different implementation must not leak into that result. The full reference combines Common content with the selected implementation; the common view excludes all implementation-specific sections.

## Rendering contract

CLI JSON, `--details` text output, and the website must agree on selected content IDs, order, code text, links, warnings, and variant applicability. Default CLI text shows only sections titled `Examples` or `<implementation> examples`, retaining their authored use-case descriptions as `# ` comments and exact code without renderer-added indentation. The Go parser derives `example_text` from the same Markdown; the CLI uses it directly without Markdown fences, with one blank line between examples. Each example must remain in its own authored Markdown code block; the website renders those original blocks separately with a copy button for each. Options are grouped into a separate block headed `Options`: option spellings, then description text, then a blank line before the next option. The parser derives `option_text` at a 100-column prose width; CLI option descriptions may use a smaller requested width. Example comments also use 100 columns including their `# ` prefix. Commands and indivisible literals are never split. References remain outside both blocks. Original section anchors remain available when option sections are grouped. Non-example sections (including `Background` and `About …`) and references require `--details`, and the website shows them as ordinary page sections; explicit `--option` lookup shows the requested section directly. HTML and terminal styling can differ. Browser pages can expose all variants, while CLI defaults to one selected view.

The content exporter emits schema_version, revision, and entries. Each entry contains metadata, version_policy, and sections. Each section contains id, implementation, title, options, markdown, html, exact code strings, and optional generated example_text for example sections and option_text for option sections. Astro uses the same exported HTML rather than a second Markdown parser. Search records derive from title/ID, summary, aliases, and tags.

CLI command JSON contains schema_version, content_revision, command metadata, selection evidence/warnings, and the selected sections. Invalid requests return a structured error when JSON format is selected. No ANSI escapes appear in JSON. Missing content is an error, not generated prose pretending to be a manual.

Run just content or the exporter with --source content to validate working-tree files. Tests reject schema errors, duplicate IDs/aliases/headings/options, unsupported Markdown, unsafe URLs, control characters, and missing variant coverage. Local documentation links are checked separately; external reference availability is not a network dependency of the offline build. Browser tests compare every implementation's selected CLI sections, HTML, and exact code with the corresponding static page. Reviewed examples run in temporary fixtures against available host tools; see development.md for precise limits.

## Guides and the glossary

Beginner guides and the glossary are a second kind of authored Markdown, parsed by the same Go code and exported next to the command pages. They follow the plan in [learning guides](learning-guides.md). Command pages keep every rule above, including the rejection of tables and headings; only guides and the glossary get the additions below.

**Where they live.** Guides are `content/en/guides/<id>.md`, one file per guide, and the glossary is `content/en/glossary/terms.md`. The command loader reads only `content/en/*.md`, so these files never become command pages. The embed directive in `content/embed.go` lists each directory explicitly because it does not descend into subdirectories; adding a new content directory means adding it there, and `go test ./content` fails if a Markdown file under `content/en` is missing from the embedded copy.

**Guide front matter.** All fields are required except `tags`, which may be empty or omitted:

```yaml
---
schema_version: 1
id: big-picture
title: From Keypress to Screen
subtitle: The big picture of what happens when you type a command
summary: Follow one command from your fingers to the text on your screen.
series: how-a-computer-works
order: 1
level: beginner
reading_minutes: 5
updated: 2026-11-08
tags: [basics, overview]
related_commands: [ls, cat]
related_guides: []
references:
  - https://docs.kernel.org/
---
```

The `id` matches the filename. `series` is one of the registered series in `catalog.Series` (exported as `series`), and `order` is the guide's position in it; two guides cannot share a position. `level` is `beginner` for now. `updated` follows the same `YYYY-MM-DD` rule and hand-maintenance as command pages. Every name in `related_commands` must be an existing command page or alias, and every name in `related_guides` an existing other guide; these are checked when the whole collection loads. Unknown keys are rejected.

**Guide sections.** A guide has exactly these level-two sections, in this order: Goal, The idea in plain words, A picture, Another way to see it, Walk through it, Try it, Common confusions, Check yourself, Where next. Level-three headings are allowed inside a section; level one and levels four to six are not.

**Markdown additions for guides.** Everything the command subset allows, plus:

- Pipe tables, with at most four columns. A table must render within 78 columns, because the CLI prints it as aligned text without wrapping cells.
- Diagrams, written as a fenced block with the language `diagram`. A diagram uses plain printable ASCII, is at most 72 columns wide, and must be followed by a paragraph that starts with `Diagram:` and describes it in words. The CLI prints a diagram exactly as written, with no added indent, and the web page shows it in a block that is not copyable.
- Level-three headings.

**Try it commands.** The `Try it` section of a guide holds commands for the reader to run, one `sh` block each. The tests run every block in a temporary folder on the host (`go test ./internal/catalog`) and in the Linux guest (`scripts/linux-smoke.py`). A block must succeed. If a fenced `text` block follows it directly, the command must print exactly that text; otherwise only success is checked, because things like disk sizes and process numbers differ on every computer. A block may only read, or write inside the current folder, so these words are rejected anywhere in it: `sudo`, `su`, `rm`, `rmdir`, `kill`, `killall`, `pkill`, `shutdown`, `reboot`, `halt`, `poweroff`, `mkfs`, `fdisk`, `dd`, `chown`, `chgrp`, `ln`, `wget`, `ssh`, `scp`, `curl` unless its line reads only a `file://` URL (no other URL scheme, no path outside the current folder, and `$PWD` the only variable), a redirect into a system path, and `mv` on a line that has a word starting with a slash or a tilde, `..` anywhere, a backtick, a command substitution, or a variable other than the loop variable of a `for` line, or in a block where any `cd` is not to a plain folder name inside the current folder, so a bare `cd` and `cd -` are rejected (the guide on renaming files needs `mv`, but only on files the block made in the current folder). The check works by word, so reword a sentence rather than trying to get past it. Use only commands that exist on macOS, GNU/Linux and BusyBox. The exception is `curl`, which a guest without it skips and reports. Commands that need the internet go in a `bash` block, which the tests do not run, with a note of the date they were checked by hand.

**Glossary.** `terms.md` has front matter with `schema_version`, `id: glossary`, `title`, `summary` and `updated`. Each level-two section is one term: the heading is the term and the body its definition. The anchor is the heading in lowercase with hyphens, so `Operating system` is `operating-system`, and anchors must be unique. Keep terms in alphabetical order, ignoring case. A final paragraph written `See: guide-id, other-guide` links the term to guides, which must exist.

**Command pages** may add an optional `related_guides` list of guide IDs, which must exist. It is exported as `related_guides`, and the CLI and web will use it for "Learn more" links.

**Export and revision.** The exporter adds `series`, `guides` and `glossary` next to `entries`, with the same per-section IDs, HTML and exact code strings. Guide and glossary files are part of the content revision, so any edit changes it.

## Content ownership

Write original examples and explanations. Link upstream manuals for detailed authority. Borrowing project ideas does not permit copying their text. If reuse is approved, preserve provenance, attribution, licensing, and modification notices as required. The project licenses original prose under CC BY 4.0 and shell examples under MIT; check compatibility and attribution before importing material or publishing releases.
