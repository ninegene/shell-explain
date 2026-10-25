# AGENTS.md: shell-explain

## Purpose

Command-reference project, private through development and iteration, with a public release planned in stage 5. Repository name: shell-explain. Executable name: explain. Planned static Astro website: shell-explain.ninegene.com. Read README.md and relevant docs before work.

## Current state

Documentation only. Do not scaffold Go, Astro, installers, workflows, or deployment resources unless asked. No build/test commands exist yet; check local links, whitespace, content consistency, and public-safe material for documentation changes.

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
