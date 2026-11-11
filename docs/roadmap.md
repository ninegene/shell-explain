# Roadmap and decisions

## Recorded decisions

| Topic | Decision/status |
| --- | --- |
| Repository | Private ninegene/shell-explain through stage 4; public in stage 5 |
| Executable | explain |
| Website | Static Astro site intended for shell-explain.ninegene.com |
| Website hosting | GitHub Pages via GitHub Actions after public transition; local previews through stage 4 |
| Source text | Shared Markdown files used by web and CLI |
| Language | Go implemented; pinned pure-Go Markdown/YAML dependencies |
| Normal lookup | Offline, curated content; no command execution |
| Platform support | Initial macOS/Linux ARM64/AMD64 target matrix proposed |
| Detection | Per-command implementation/version evidence plus explicit overrides |
| Distribution | Release binaries, curl installer, Homebrew tap, mise GitHub backend |
| Licensing | MIT for software and shell examples; CC BY 4.0 for original explanatory content |
| Current status | Stages 1–4 implemented (194 command pages, checklist and exit review complete); all planned beginner guides are written; no public release/deployment |

## Implementation evidence

The shared format is schema version 1, parsed once in Go. Four original Markdown entries cover grep, sed, find, and the cd builtin. The Go CLI embeds the authored files; Astro consumes a generated normalized export through its content collection. The site provides command/tag pages, shareable implementation views, static search, copy buttons, and responsive keyboard-accessible layouts.

Run the commands in [development and verification](development.md) to reproduce validation. Go tests cover content errors, selection, output, and available real-host examples; Chromium tests compare all implementation pages against CLI JSON and test search, copying, keyboard/mobile behavior, and accessibility. The optional QEMU smoke check validates actual Linux GNU/BusyBox selection and authored examples. Exact tested builds are documented rather than claiming universal version coverage.

## Stage 1: validate the shared format

Implemented the minimal Markdown schema using original content. Provenance and reuse constraints remain documented. The initial entries for grep, sed, find, and cd expose implementation differences and parent-shell resolution limitations.

The normalized Go export supplies both CLI JSON and Astro pages. Parity tests verify selected sections, code fidelity, HTML, references, revision, and variant URLs before corpus expansion. See the content-format and development documents for the implemented contract.

## Stage 2: build the CLI

Implement content validation, embedded lookup, search, terminal rendering, JSON output, and explicit variants first. Add executable detection with confidence reporting and bounded probes. Test mixed GNU/Apple PATH configurations and BusyBox. Add optional shell integration only after standalone behavior is clear.

Do not start by implementing arbitrary script interpretation, AI inference, automatic package installation, or background telemetry. The first useful release should explain a small set of commands accurately.

## Stage 3: build the static site

Implement Astro catalog/detail/tag pages, explicit implementation filters, static search, copyable examples, and accessible mobile layouts. Link command pages to applicable CLI versions and upstream references. Verify that no server-side dependency or browser OS guess is needed for basic reference use.

## Stage 4: build the commonly used command collection

Status: complete on 2026-11-10. Every command in the checklist has a page with reviewed examples, applicable implementation coverage and recorded evidence, every lookup name and topic has a disposition, and the content, Go, Linux guest and website checks pass. The last guide from the learning-guides proposal, `fetch-from-the-web`, landed on 2026-11-10. The text below is kept as the plan that was followed.

The current shared format, CLI output, and static site usability are good enough to proceed. Do not keep stage 4 open for another broad usability iteration or restrict it to the initial four commands. Remaining non-blocking usability work moves to the future backlog below.

The next step is the [command collection checklist](command-collection.md): inventory all commands exposed by the reference explain helper, its cheat sheets, GNU catalog, tool listing, and lookup aliases. Add original explanations, commonly used examples and options, and applicable implementation variants for each command. GNU, Apple/macOS, BusyBox, shell builtins, and other upstream implementations are distinct coverage targets, not interchangeable OS labels.

Work in reviewed batches, beginning with everyday commands and the existing grep/sed/find/cd pages. Extend the shared schema/selection only where the collection requires it, including upstream tools and symbolic command spellings. Keep the established Markdown/code-copy/compact-output format and use local builds/previews. Record upstream references, implementation/version coverage, and actual test evidence; do not claim unsupported variants or copy the reference corpus without compatible licensing.

Stage 4 is complete when every required command in the checklist has reviewed common use cases/options and applicable implementation coverage, lookup aliases and non-command topics have explicit dispositions, both consumers expose the same content, and content/Go/web checks pass with no known blocking accuracy or rendering regressions. Completion of this collection is the stage's exit gate; the deferred usability tasks are not extra gates. Mark the stage complete only after the inventory and exit review are checked off.

### Proposed addition to stage 4: learning guides

The command pages assume readers already know what a terminal, shell, process or pipe is, and why one command has several versions. [Learning guides](learning-guides.md) propose a beginner series (operating systems, the terminal and shell, input and output, and the history behind GNU, BSD, macOS, Linux and BusyBox) with diagrams, tables and runnable examples, shown on the home page and through `explain learn` in the CLI. The planned guides are now implemented.

## Stage 5: license, open contributions, distribute, and publish

Begin only after stage 4's command collection and exit review are complete. Keep public launch as a deliberate release step:

1. License software and shell examples under MIT and original explanatory content under CC BY 4.0. Review third-party material and dependency notices, and include applicable notices in releases; see [third-party notices](third-party-notices.md) for the inventory and remaining artifact checks.
2. Prepare contribution guidance, issue/PR templates, review expectations, and candidate starter issues. Setup and checks passed in the existing development environment; verify them in a fresh contributor environment before opening contributions; see [contributing](../CONTRIBUTING.md).
3. Review the complete repository history, release artifacts, links, and configuration for secrets, personal information, private dependencies, and material that cannot be redistributed. Resolve findings before opening access.
4. Change the existing GitHub repository from private to public after the publication review and owner authorization. Enable the intended issue/PR contribution channels and update the README to accurately describe public status and licensing.
5. Publish tested immutable binaries/checksums and clear release notes. Validate curl/manual installation, a Homebrew tap, and mise in clean environments. Resolve command-name collisions without silently replacing other programs or shell functions.
6. Publish the static Astro site to GitHub Pages using a GitHub Actions build/deploy workflow after the repository is public. Configure shell-explain.ninegene.com and HTTPS, verify domain ownership and DNS, and check routes/search/assets. Include matching content revision, installation guidance, license/attribution information, and links to report problems or contribute commands/examples. Prepare Pages/domain/TLS configuration for separate deployment authorization; no paid cloud website hosting is planned.

Roadmap tasks do not themselves authorize a commit/push, the future public visibility change, website deployment, or DNS updates. Keep those actions subject to explicit authorization.

## Acceptance criteria

- One authored Markdown source supplies matching explanations to both outputs.
- A downloaded binary works offline on the declared support matrix.
- Mixed-tool installations do not silently receive an OS-only implementation guess.
- Missing/ambiguous coverage is visible and overrideable.
- No explanation path executes user input or sample commands.
- JSON is stable, versioned, and free of ANSI escapes.
- Examples and detection rules have targeted implementation-specific checks.
- Installer/package artifacts are versioned and verified; no unsupported claims of registry availability.
- Licenses, provenance, and attribution are documented before reusing third-party material.
- Before public release, code/content licenses are selected, contribution instructions are verified, and publication review is complete.
- Public repository content includes no secrets, personal information, or private project details.

## Future usability backlog (not stage 4 blockers)

- Recruit beginners and experienced users for structured task sessions; record a feedback log for discovering commands, understanding/copying examples, looking up options, and selecting/comparing implementations.
- Broaden manual offline/narrow-terminal, keyboard/mobile, assistive-technology, and mixed-installation usability testing beyond current automated checks.
- Revisit search, explanations, metadata, implementation-comparison UX, and detection based on observed feedback. Return to stages 1–3 as needed and verify improvements in both consumers.
- Formalize stronger public schema/JSON compatibility and migration guarantees as consumers emerge. Keep current schema versioning and regression checks during expansion.

These replace the unfinished broad iteration tasks from the previous stage 4. Schedule them after the command collection or from real user feedback; fix any blocking correctness/accessibility regressions when discovered.

## Later possibilities

Translations, additional BSD variants, user-supplied content packs, full pipeline explanation using a real parser, and richer portability comparisons can follow demonstrated demand. Keep these optional and do not compromise offline operation or a maintainable shared source.
