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
references:
  - https://www.gnu.org/software/grep/manual/
---
```

All fields shown above are required except aliases, which may be empty. IDs, tags, and implementation IDs use lowercase letters/digits/hyphens starting with a letter; titles equal the command ID. Aliases follow the same rule, plus the registered symbolic spellings in `symbolicNames` (currently `[`, for `test`). A symbolic spelling is only ever an alias: the canonical ID, filename, and web URL stay `test`, and the CLI looks it up directly (`explain '['`). Unknown metadata keys, duplicate YAML keys, duplicate names, unsupported schema versions, and unsafe reference URLs are rejected. References and inline links must use absolute HTTP(S) URLs without userinfo. Plain scalar metadata must not contain control characters.

An implementation list declares authored coverage; it does not prove every implementation supports every option. Supported IDs come from one registry, `catalog.Implementations`, which the parser, CLI help, JSON export (`implementations`), and website labels all read: gnu, apple, busybox, bash, zsh, and upstream. `bash` and `zsh` are shell builtins; an entry that declares only shells is treated as a builtin, which no PATH probe can describe. `upstream` is for tools with one real implementation (jq, git, curl, and similar). It cannot be combined with other implementations and is selected without detection, so it is never presented as GNU, BusyBox, or an OS guess. FreeBSD/OpenBSD need future content and validation rules before becoming selectable variants; add a registry row with its heading when real content exists. Aliases are lookup names, not permission to assume the installed implementation; `ggrep` still needs selection evidence.

Each page must contain Common sections named Overview, Background, Examples, Common mistakes, and Portability, and each declared implementation view must contain an About section. History and implementation background are kept separate from the examples: the Common section named `Background` (where the command came from and how its implementations relate), and one `About <name>` section at the start of each implementation view (what that implementation is, where it ships, and how it differs, for example `About GNU grep`). `Background` is valid only in Common, `About …` only in an implementation view, and a view may have at most one. Both are required: a page without `Background`, or a declared view without an `About` section, is rejected. State only history and differences that the page's references or an executed check support; leave out anything unverified. Other sections can explain use cases or options. References are rendered from metadata in both interfaces. Repeated content belongs in Common; variant-specific content uses explicit headings.

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

## Content ownership

Write original examples and explanations. Link upstream manuals for detailed authority. Borrowing project ideas does not permit copying their text. If reuse is approved, preserve provenance, attribution, licensing, and modification notices as required. Choose code and content licenses before importing a corpus or publishing reusable releases.
