# Shell Explain

Practical command explanations for the terminal and the web, maintained in one shared Markdown collection.

- Repository: [ninegene/shell-explain](https://github.com/ninegene/shell-explain), private during development and iteration; planned public release in stage 5.
- Executable: **`explain`**, for example `./bin/explain grep` after a local build.
- Planned website: **shell-explain.ninegene.com**, built with Astro and hosted on GitHub Pages after the public release in stage 5.
- CLI language: Go, with original Markdown embedded for offline use.
- Status: **stages 1–3 implemented locally**. Current CLI and web usability are accepted as sufficient. Stage 4 next expands the command collection, common examples/options, and applicable implementation coverage. No installer, public release, or deployment exists yet.

## Goal

Explain what a command does, show useful examples, explain common options and mistakes, and select documentation appropriate to the installed implementation. Beginners should be able to understand the examples; experienced users should be able to find an unfamiliar option quickly.

The CLI and website must consume the same `.md` files. Neither keeps an independently edited copy of the explanations. Installation should eventually support downloadable binaries, a curl-based installer, Homebrew, and mise.

Use local Astro previews during private development and iteration. After the repository becomes public, GitHub Actions will build and deploy the static site to GitHub Pages. No paid cloud hosting is planned for the website; domain registration remains separate.

## Documentation

| Document | Purpose |
| --- | --- |
| [Product and architecture](docs/architecture.md) | Scope, Go choice, Astro site, and integrations |
| [Shared Markdown format](docs/content-format.md) | Proposed metadata, prose conventions, and rendering contract |
| [Implementation detection](docs/detection.md) | GNU/BSD/Apple/BusyBox selection and shell limitations |
| [Distribution](docs/distribution.md) | Releases, curl, Homebrew, mise, and reproducibility |
| [Related projects and references](docs/references.md) | tldr, ExplainShell, all projects in tldr's similar-projects list, and supporting tools |
| [Command collection checklist](docs/command-collection.md) | Stage 4 inventory, implementation coverage, and completion criteria |
| [Roadmap](docs/roadmap.md) | Decisions, validation, and implementation stages |
| [Development and verification](docs/development.md) | Build commands, local previews, test coverage, and known limits |
| [Contributing](CONTRIBUTING.md) | Documentation contributions and source/licensing expectations |

## Local quick start

Requires [just](https://github.com/casey/just), Go 1.24+, and Node 22.23.4+. Install the command runner with `brew install just` on macOS.

Run `just` or `just help` to see all development commands and setup requirements. Plain `just` only displays help; it does not build or download dependencies.

Build and use the offline CLI:

```sh
just build
./bin/explain grep
./bin/explain grep --details
./bin/explain grep --implementation gnu --option=-n
./bin/explain sed --implementation apple --format json
```

Build or preview the site from the same Markdown files:

```sh
npm ci --prefix web
just web-dev
```

The currently implemented interface supports these forms once the binary is on PATH:

```console
explain grep
explain grep --details
explain ggrep
explain grep --option=-n
explain --implementation gnu grep
explain --platform linux grep
explain --format json grep
explain search logs
explain list
```

`explain` provides information; it does not run the example command or install the command being explained. Full pipeline/script explanation is a possible later feature, not part of the first release.

Text output shows example usages for the selected implementation by default, with a compact implementation/confidence heading. Descriptions use `# ` comments, commands have no added indentation, and a blank line separates examples. In the Markdown source and on the web, each example keeps its own code block and copy button so a single command can be copied. The CLI removes fences and renderer indentation, separates examples with one blank line, and prefixes descriptions with `# `. Detailed option explanations are grouped under `Options`, with each option spelling above its description. Descriptions wrap at 100 columns (including the `# ` prefix on example comments); authored commands and indivisible inline literals remain intact. Use `--details` before or after the command for the full overview, options, mistakes, portability notes, detection evidence, and references. Explicit `--option` lookup shows that option's explanation without requiring `--details`. JSON and the website retain the full selected reference.

## Project boundaries

This project is being developed as a reusable reference tool, with public availability planned after iteration. It is not an interactive lab host. Other applications can eventually install a pinned binary or consume its documented JSON output. Lesson-specific hints, grading, accounts, and progress belong to those applications.

GNU/BSD behavior must be selected per command, not inferred solely from the operating system. The website cannot inspect the visitor's terminal; it must offer explicit implementation/platform selection.

## Development status and licensing

Use `just test` for Go checks and reviewed examples, `just check` for Go vet and Astro diagnostics, `just web-build` for the static build, and `just web-test` for Chromium parity/usability/accessibility tests. See [development](docs/development.md) for setup and limitations. The initial corpus is grep, sed, find, and cd; exact version support is not certified universally.

Research date: **2026-10-23**, America/Los_Angeles. Related-project descriptions come from published documentation, not independent trials. See [references](docs/references.md).

Code and content license selection is pending and belongs to stage 5, after the iteration stage. Before opening the repository and publishing the website, select and add the applicable licenses and finalize the contribution process for new commands and examples. Public visibility alone does not grant reuse permission; do not import third-party content before checking its license and attribution requirements.
