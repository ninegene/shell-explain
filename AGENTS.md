# AGENTS.md: shell-explain

## Purpose

Command-reference project, private through development and iteration, with a public release planned in stage 5. Repository name: shell-explain. Executable name: explain. Planned static Astro website: shell-explain.ninegene.com. Read README.md and relevant docs before work.

## Current state

Stages 1–3 are implemented locally: original content in content/en, Go parser/CLI in internal and cmd, and an Astro static site in web. Stage 4 is the next iteration phase. Do not create installers, release/deployment workflows, change repository visibility, or publish without authorization.

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

## Publication and repository rules

- Keep the repository private through stage 4. Stage 5 includes choosing code/content licenses, finalizing contribution guidance, changing visibility to public, and publishing the static site. Do not change visibility or deploy merely because these actions appear in the roadmap.
- Never add secrets, personal data, cloud identifiers, machine names, or private project names/details.
- Keep sample paths and values generic; use example.org and synthetic fixtures.
- Check licenses and attribution before copying content/code. License choices remain pending.
- Do not alter bootstrap/reference repositories as a side effect.
- Never commit without an explicit user request. When authorized, use Conventional Commits.
- Follow workspace confirmation requirements for destructive changes, cloud writes, and SSH. Documentation and repository creation do not authorize hosting or DNS changes.
