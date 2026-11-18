# AGENTS.md: shell-explain

## Purpose

Command-reference project; the repository has been public since 2026-11-17 and release binaries, a curl installer and install docs (README, website /install/) exist; the first release is v0.1.0. Repository name: shell-explain. Executable name: explain. Static Astro website: shell-explain.ninegene.com. Read README.md and relevant docs before work.

## Current state

Stages 1–4 are implemented locally: original content in content/en, Go parser/CLI in internal and cmd, and an Astro static site in web. Current CLI/web usability is accepted as sufficient. Stage 4 (the command collection in docs/command-collection.md, 194 pages) is complete, and the planned beginner guides are written. Broader usability studies are deferred in docs/roadmap.md. Do not add further release/deployment workflows beyond the release tooling in docs/distribution.md, change repository visibility, or publish without authorization (the website's AWS hosting under `infra/` and `.github/workflows/deploy.yml` is the authorized exception).

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
- Recommended CLI implementation is Go with embedded offline content; website uses Astro static generation, hosted on AWS (S3 + CloudFront, Pulumi in `infra/pulumi`, runbook in `infra/README.md`) and deployed by GitHub Actions. Use local previews during iteration. Applying the Pulumi stack, deploying and DNS changes are cloud writes: confirm each first.
- Generated data, HTML, search indexes, and embedded bundles are outputs, never a second authoring source.
- Resolve documentation by command implementation/version, not OS alone. Make ambiguous selection visible and support overrides.
- Never execute example/input command lines to explain them. Any supported version probes must be explicit, bounded, command-specific, and documented.
- A standalone binary cannot automatically inspect parent-shell aliases/functions; optional shell integration must remain separate.
- Avoid AI/network dependencies in normal offline command lookup.

## Command content and output format

Follow [docs/content-format.md](docs/content-format.md) when adding or updating a command. Author only the canonical `content/en/<command>.md`; generated JSON and rendered text are not authoring sources.

- **Source structure:** declare metadata and implementation coverage, use level-two variant headings (`Common`, `GNU`, `Apple`, `BusyBox`, `Bash`), and put content in level-three sections. Common requires `Overview`, `Background`, `Examples`, `Common mistakes`, and `Portability`, and every declared implementation view requires one `About <name>` section (history and implementation differences, limited to claims the references or an executed check support; shown with `--details` and on the web, not in the default CLI output). Name variant example sections `<implementation> examples`, such as `GNU examples`, so the compact CLI includes them.
- **Guides and glossary:** beginner guides are `content/en/guides/<id>.md` and the glossary is `content/en/glossary/terms.md`, following the Guides section of [docs/content-format.md](docs/content-format.md) and the plan in [docs/learning-guides.md](docs/learning-guides.md). Guides keep the same rules for references, `updated` dates and claims backed by a source, and every `Try it` command is run by the tests and must be safe (see the content format). A new content directory must be added to the `go:embed` directive in `content/embed.go`.
- **Last updated date:** every page has a required `updated: YYYY-MM-DD` front-matter field. Set it to today's date (run `date +%F`) whenever you change that page's claims, examples, or version statements. It is shown on the website and in `--details`.
- **Examples in Markdown:** write a short use-case description followed by a separate fenced `sh` or `bash` block for each example. Keep one independently copyable command or coherent multiline example per block, with explicit closing fences. Do not combine unrelated examples into one block or include shell prompts. Preserve meaningful spaces, quoting, and multiline indentation.
- **Examples on the web:** render the original separate code blocks, each with its own copy button that copies only that example. Do not concatenate examples into a single web block.
- **Examples in the CLI:** default text output shows examples only, under a heading with the command's summary (the implementation and confidence are shown only when the page has more than one implementation). Lay examples out densely from the same authored pairs: a single-line command of at most 46 columns whose description fits becomes `command  # description` with the comments of neighbouring lines aligned; any other example keeps its description as `# ` comment lines above the command, separated by a blank line. Drop the "reviewed only" remark and the trailing colon from descriptions in this layout. Omit Markdown fences and renderer-added indentation, never change or wrap a command, and preserve indentation authored inside multiline commands. Color command names (bold), flags (blue), quoted strings (magenta), operators (dim) and comments and option descriptions (green) when color is on, using only the 8 standard ANSI colors plus bold, dim and underline so every terminal theme, light or dark, maps them to readable values: no 256-color or truecolor values, no cyan, yellow or bright colors, and no bold combined with a color. End the output with the command's website page and its reference URLs, written in full so a terminal can open them. Give each example an optional `short="..."` label on its code fence (plain text, about 50 characters, at most 70, no backticks) so it fits on one line: the default CLI layout uses the label, while the website and `--details` keep the full authored description. The full description stays the reference text; never put information only in the label.
- **Detailed output:** `--details`, before or after the command, includes overview, options, mistakes, portability notes, detection evidence, and references. Explicit `--option=-n` lookup shows the requested option directly without requiring `--details`. JSON and the website retain the full selected reference; do not prune JSON based on the text presentation mode.
- **Options in Markdown:** author one level-three `Option: <spellings>` section per documented option, such as `### Option: -i, --ignore-case`, followed by its description. Use structured option spellings for lookup; do not replace these sections with a manually combined output block.
- **Options in output:** group selected options under one `Options` heading. On the website, put each option's spellings on their own line with its description immediately below and one blank line between entries, in one block with a dedicated copy button. In CLI text, align the spellings and descriptions in two columns (a spelling wider than 22 columns gets its own line), wrap descriptions under their column, color the spellings, and add no fences or indentation. `--details` ends with the website page as `More:` after the references and colors its labels and URLs. Preserve source anchors when grouping sections.
- **Description width:** use the shared `catalog.DescriptionWidth` limit of **100 columns**, including the `# ` prefix on CLI example comments. Option descriptions honor narrower CLI widths; `--width=0` retains the compact-description cap. Never split authored commands or indivisible inline literals just to meet the prose limit.
- **Shared rendering:** derive compact `example_text` and `option_text` in Go from the same Markdown. Astro consumes generated content without introducing a second Markdown parser. Keep raw source/code fields intact, and verify exact command text, individual clipboard contents, description wrapping, implementation filtering, and CLI/web parity when changing this format.

## Publication and repository rules

- The repository is public, and `main` is protected by a ruleset: changes go through pull requests. Do not publish a release, deploy the website, or change DNS or repository visibility merely because these actions appear in the roadmap; each needs explicit authorization.
- Never add secrets, personal data, cloud account IDs or keys, machine names, or private project names/details. Resource names that grant no access (the Pulumi state bucket, KMS alias, site bucket) and GitHub repository IDs may appear in `infra/`; the AWS account ID stays a GitHub environment secret.
- Keep sample paths and values generic; use example.org and synthetic fixtures.
- Check licenses and attribution before copying content/code. Original software and shell examples use MIT; original explanatory content uses CC BY 4.0. See README.md.
- Do not alter bootstrap/reference repositories as a side effect.
- Never commit without an explicit user request. When authorized, use Conventional Commits.
- Follow workspace confirmation requirements for destructive changes, cloud writes, and SSH. Documentation and repository creation do not authorize hosting or DNS changes.
