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
| Licensing | Choose and add code/content licenses in stage 5, before public release |
| Current status | Stages 1–3 implemented locally; stage 4 iteration next; no public release/deployment |

## Implementation evidence

The shared format is schema version 1, parsed once in Go. Four original Markdown entries cover grep, sed, find, and the cd builtin. The Go CLI embeds the authored files; Astro consumes a generated normalized export through its content collection. The site provides command/tag pages, shareable implementation views, static search, copy buttons, and responsive keyboard-accessible layouts.

Run the commands in [development and verification](development.md) to reproduce validation. Go tests cover content errors, selection, output, and available real-host examples; Chromium tests compare all implementation pages against CLI JSON and test search, copying, keyboard/mobile behavior, and accessibility. The optional QEMU smoke check validates actual Linux GNU/BusyBox selection and authored examples. Exact tested builds are documented rather than claiming universal version coverage.

## Stage 1: validate the shared format

Implemented the minimal Markdown schema using original content. Provenance and reuse constraints remain documented, and final code/content license selection stays in stage 5. The initial entries for grep, sed, find, and cd expose implementation differences and parent-shell resolution limitations.

The normalized Go export supplies both CLI JSON and Astro pages. Parity tests verify selected sections, code fidelity, HTML, references, revision, and variant URLs before corpus expansion. See the content-format and development documents for the implemented contract.

## Stage 2: build the CLI

Implement content validation, embedded lookup, search, terminal rendering, JSON output, and explicit variants first. Add executable detection with confidence reporting and bounded probes. Test mixed GNU/Apple PATH configurations and BusyBox. Add optional shell integration only after standalone behavior is clear.

Do not start by implementing arbitrary script interpretation, AI inference, automatic package installation, or background telemetry. The first useful release should explain a small set of commands accurately.

## Stage 3: build the static site

Implement Astro catalog/detail/tag pages, explicit implementation filters, static search, copyable examples, and accessible mobile layouts. Link command pages to applicable CLI versions and upstream references. Verify that no server-side dependency or browser OS guess is needed for basic reference use.

## Stage 4: iterate on the shared format, CLI, and static site

Treat stages 1–3 as an iterative loop, not a one-time sequence. Review the same command entries in the CLI and local Astro preview, collect feedback from beginners and experienced users, and use the findings to improve the content format and both interfaces together.

For each iteration:

1. Try representative tasks: discover a command, understand an example, find an option, and select or compare implementations.
2. Record unclear explanations, missing metadata, rendering differences, detection mistakes, search problems, and accessibility issues.
3. Return to stage 1 for schema/content changes, stage 2 for CLI changes, and stage 3 for site changes as needed. Update both consumers when the shared contract changes.
4. Rerun shared fixtures and relevant checks, then repeat the tasks to verify the changes helped without introducing regressions.

Keep a small representative corpus during this stage so format changes remain inexpensive. Review offline CLI use, terminal widths, keyboard/mobile navigation, and GNU/Apple/BusyBox differences. Use local builds and previews; public distribution or website deployment is not required to gather feedback.

Repeat until representative tasks are understandable, both outputs preserve the same explanations and examples, and no known blocking format, usability, detection, or rendering issues remain. Document remaining limitations and stabilize the initial schema/JSON contract before advancing to distribution and publishing.

## Stage 5: license, open contributions, distribute, and publish

Begin only after stage 4's iteration criteria are met. Keep public launch as a deliberate release step:

1. Select an open-source code license and an appropriate content/documentation license. Document which paths each covers, add license texts and attribution notices, and confirm imported material is compatible. Do not assume code and explanation text need the same license.
2. Finalize CONTRIBUTING.md with working setup/validation commands, new-command and new-example instructions, implementation/version coverage requirements, and the accepted contribution licensing terms. Prepare issue/PR templates, maintainer review expectations, and a small set of suitable starter issues. Verify a new contributor can follow the process.
3. Review the complete repository history, release artifacts, links, and configuration for secrets, personal information, private dependencies, and material that cannot be redistributed. Resolve findings before opening access.
4. Change the existing GitHub repository from private to public after the publication review and owner authorization. Enable the intended issue/PR contribution channels and update the README to accurately describe public status and licensing.
5. Publish tested immutable binaries/checksums and clear release notes. Validate curl/manual installation, a Homebrew tap, and mise in clean environments. Resolve command-name collisions without silently replacing other programs or shell functions.
6. Publish the static Astro site to GitHub Pages using a GitHub Actions build/deploy workflow after the repository is public. Configure shell-explain.ninegene.com and HTTPS, verify domain ownership and DNS, and check routes/search/assets. Include matching content revision, installation guidance, license/attribution information, and links to report problems or contribute commands/examples. Prepare Pages/domain/TLS configuration for separate deployment authorization; no paid cloud website hosting is planned.

The authorized implementation scope covers stages 1–3 in the private working tree. It does not authorize a commit/push, the future public visibility change, website deployment, or DNS updates.

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

## Later possibilities

Translations, additional BSD variants, user-supplied content packs, full pipeline explanation using a real parser, and richer portability comparisons can follow demonstrated demand. Keep these optional and do not compromise offline operation or a maintainable shared source.
