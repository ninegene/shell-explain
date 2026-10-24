# Shell Explain

Practical command explanations for the terminal and the web, maintained in one shared Markdown collection.

- Repository: [ninegene/shell-explain](https://github.com/ninegene/shell-explain), private during development and iteration; planned public release in stage 5.
- Planned executable: **`explain`**, for example `explain grep`.
- Planned website: **shell-explain.ninegene.com**, built as a static Astro site.
- Recommended CLI language: Go, with documentation embedded for offline use.
- Status: **documentation only**. No CLI, installer, website implementation, release, or deployment exists yet.

## Goal

Explain what a command does, show useful examples, explain common options and mistakes, and select documentation appropriate to the installed implementation. Beginners should be able to understand the examples; experienced users should be able to find an unfamiliar option quickly.

The CLI and website must consume the same `.md` files. Neither keeps an independently edited copy of the explanations. Installation should eventually support downloadable binaries, a curl-based installer, Homebrew, and mise.

## Documentation

| Document | Purpose |
| --- | --- |
| [Product and architecture](docs/architecture.md) | Scope, Go choice, Astro site, and integrations |
| [Shared Markdown format](docs/content-format.md) | Proposed metadata, prose conventions, and rendering contract |
| [Implementation detection](docs/detection.md) | GNU/BSD/Apple/BusyBox selection and shell limitations |
| [Distribution](docs/distribution.md) | Releases, curl, Homebrew, mise, and reproducibility |
| [Related projects and references](docs/references.md) | tldr, ExplainShell, all projects in tldr's similar-projects list, and supporting tools |
| [Roadmap](docs/roadmap.md) | Decisions, validation, and implementation stages |
| [Contributing](CONTRIBUTING.md) | Documentation contributions and source/licensing expectations |

## Proposed usage

These commands describe the intended interface; they are not available yet.

```console
explain grep
explain ggrep
explain grep --option=-n
explain --implementation gnu grep
explain --platform linux grep
explain --format json grep
explain search logs
explain list
```

`explain` provides information; it does not run the example command or install the command being explained. Full pipeline/script explanation is a possible later feature, not part of the first release.

## Project boundaries

This project is being developed as a reusable reference tool, with public availability planned after iteration. It is not an interactive lab host. Other applications can eventually install a pinned binary or consume its documented JSON output. Lesson-specific hints, grading, accounts, and progress belong to those applications.

GNU/BSD behavior must be selected per command, not inferred solely from the operating system. The website cannot inspect the visitor's terminal; it must offer explicit implementation/platform selection.

## Development status and licensing

There are no build/test commands yet. Validate documentation links, formatting, public-safe content, and consistency of proposed behavior. Add executable commands only when their implementation exists.

Research date: **2026-10-23**, America/Los_Angeles. Related-project descriptions come from published documentation, not independent trials. See [references](docs/references.md).

Code and content license selection is pending and belongs to stage 5, after the iteration stage. Before opening the repository and publishing the website, select and add the applicable licenses and finalize the contribution process for new commands and examples. Public visibility alone does not grant reuse permission; do not import third-party content before checking its license and attribution requirements.
