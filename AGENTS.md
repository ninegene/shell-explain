# AGENTS.md: shell-explain

## Purpose

Command-reference project, private through development and iteration, with a public release planned in stage 5. Repository name: shell-explain. Executable name: explain. Planned static Astro website: shell-explain.ninegene.com. Read README.md and relevant docs before work.

## Current state

Stages 1–3 are implemented locally: original content in content/en, Go parser/CLI in internal and cmd, and an Astro static site in web. Current CLI/web usability is accepted as sufficient. Stage 4 next expands the command collection using docs/command-collection.md; complete its inventory and exit review before marking the stage done. Broader usability studies are deferred in docs/roadmap.md. Do not create installers, release/deployment workflows, change repository visibility, or publish without authorization.

## Commands and implementation

- `just` / `just help`: show the command guide without building or downloading dependencies.
- `just build`: build bin/explain with embedded Markdown.
- `just content`: validate/export the normalized content model for Astro.
- `just test`: Go tests, including temporary-fixture examples and available host-tool integration checks.
- `just check`: Go vet and Astro diagnostics (run npm ci --prefix web first).
- `just web-build` / `just web-dev`: production static build / local preview.
- `just web-test`: Chromium parity, usability, and accessibility tests; first install the test browser from web with npx playwright install chromium.
- Go is the only Markdown/front-matter parser. Astro consumes its generated, ignored JSON via a content collection. Never hand-edit generated data.
- Keep test claims precise: unit fixtures and an optional real Linux QEMU smoke check cover BusyBox/GNU paths and examples, not every build/version. See docs/development.md; cmd/compat-check is test-only.

## Architecture conventions

- One reviewed Markdown source collection must supply CLI and website explanations.
- Recommended CLI implementation is Go with embedded offline content; website uses Astro static generation and GitHub Pages hosting via GitHub Actions after stage 5's public transition. Use local previews during private iteration; do not add paid cloud hosting for this site.
- Generated data, HTML, search indexes, and embedded bundles are outputs, never a second authoring source.
- Resolve documentation by command implementation/version, not OS alone. Make ambiguous selection visible and support overrides.
- Never execute example/input command lines to explain them. Any supported version probes must be explicit, bounded, command-specific, and documented.
- A standalone binary cannot automatically inspect parent-shell aliases/functions; optional shell integration must remain separate.
- Avoid AI/network dependencies in normal offline command lookup.

## Command content and output format

Follow [docs/content-format.md](docs/content-format.md) when adding or updating a command. Author only the canonical `content/en/<command>.md`; generated JSON and rendered text are not authoring sources.

- **Source structure:** declare metadata and implementation coverage, use level-two variant headings (`Common`, `GNU`, `Apple`, `BusyBox`, `Bash`), and put content in level-three sections. Common requires `Overview`, `Background`, `Examples`, `Common mistakes`, and `Portability`, and every declared implementation view requires one `About <name>` section (history and implementation differences, limited to claims the references or an executed check support; shown with `--details` and on the web, not in the default CLI output). Name variant example sections `<implementation> examples`, such as `GNU examples`, so the compact CLI includes them.
- **Guides and glossary:** beginner guides are `content/en/guides/<id>.md` and the glossary is `content/en/glossary/terms.md`, following the Guides section of [docs/content-format.md](docs/content-format.md) and the plan in [docs/learning-guides.md](docs/learning-guides.md). Guides keep the same rules for references, `updated` dates and claims backed by a source. A new content directory must be added to the `go:embed` directive in `content/embed.go`.
- **Last updated date:** every page has a required `updated: YYYY-MM-DD` front-matter field. Set it to today's date (run `date +%F`) whenever you change that page's claims, examples, or version statements. It is shown on the website and in `--details`.
- **Examples in Markdown:** write a short use-case description followed by a separate fenced `sh` or `bash` block for each example. Keep one independently copyable command or coherent multiline example per block, with explicit closing fences. Do not combine unrelated examples into one block or include shell prompts. Preserve meaningful spaces, quoting, and multiline indentation.
- **Examples on the web:** render the original separate code blocks, each with its own copy button that copies only that example. Do not concatenate examples into a single web block.
- **Examples in the CLI:** default text output shows examples only, with a compact implementation/confidence heading. Render descriptions as `# ` comments, omit Markdown fences and renderer-added indentation, and separate examples with exactly one blank line. Preserve indentation authored inside multiline commands.
- **Detailed output:** `--details`, before or after the command, includes overview, options, mistakes, portability notes, detection evidence, and references. Explicit `--option=-n` lookup shows the requested option directly without requiring `--details`. JSON and the website retain the full selected reference; do not prune JSON based on the text presentation mode.
- **Options in Markdown:** author one level-three `Option: <spellings>` section per documented option, such as `### Option: -i, --ignore-case`, followed by its description. Use structured option spellings for lookup; do not replace these sections with a manually combined output block.
- **Options in output:** group selected options under one `Options` heading. Put each option's spellings on their own line, its description immediately below, and one blank line between entries. The website uses one Options block with a dedicated copy button; CLI text has no fences or added indentation. Preserve source anchors when grouping sections.
- **Description width:** use the shared `catalog.DescriptionWidth` limit of **100 columns**, including the `# ` prefix on CLI example comments. Option descriptions honor narrower CLI widths; `--width=0` retains the compact-description cap. Never split authored commands or indivisible inline literals just to meet the prose limit.
- **Shared rendering:** derive compact `example_text` and `option_text` in Go from the same Markdown. Astro consumes generated content without introducing a second Markdown parser. Keep raw source/code fields intact, and verify exact command text, individual clipboard contents, description wrapping, implementation filtering, and CLI/web parity when changing this format.

## Publication and repository rules

- Keep the repository private through stage 4. Stage 5 includes choosing code/content licenses, finalizing contribution guidance, changing visibility to public, and publishing the static site. Do not change visibility or deploy merely because these actions appear in the roadmap.
- Never add secrets, personal data, cloud identifiers, machine names, or private project names/details.
- Keep sample paths and values generic; use example.org and synthetic fixtures.
- Check licenses and attribution before copying content/code. License choices remain pending.
- Do not alter bootstrap/reference repositories as a side effect.
- Never commit without an explicit user request. When authorized, use Conventional Commits.
- Follow workspace confirmation requirements for destructive changes, cloud writes, and SSH. Documentation and repository creation do not authorize hosting or DNS changes.
