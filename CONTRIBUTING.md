# Contributing

The project is currently in planning and remains private through development and iteration. External contributions will open in stage 5, after code/content licenses and the contribution process are finalized. This document is draft guidance for that process; it does not imply a public contribution channel exists today. Executable implementation should follow the agreed roadmap.

## Planned process for new commands and examples

Before opening contributions, maintainers will replace this draft process with verified setup and validation commands and define accepted contribution licensing terms. The intended workflow is:

1. Search existing commands and issues to avoid duplicates. Open an issue for a new command, implementation variant, or larger change; small example corrections can go directly to a pull request once public contributions open.
2. Fork the public repository and create a branch. Add or update the canonical Markdown entry using the shared schema; do not maintain separate web and CLI explanations.
3. Provide original examples, explain their use cases and important options, and state the tested operating system, implementation, and version. Include upstream references and required attribution for any approved reused material.
4. Run the documented content validation and relevant CLI/web checks. Include an explanation of changes and testing evidence in the pull request, with screenshots or terminal output when useful.
5. Address maintainer feedback on accuracy, portability, accessibility, licensing, and shared rendering. Maintainers review and merge accepted contributions; publishing follows the release process.

The public launch should include issue/PR templates and starter issues for adding examples, commands, or portability corrections. No contributor should need access to private infrastructure or credentials to validate reference content locally.

## Command documentation

Write original explanations with a concrete use case, understandable examples, option meanings, expected behavior, and relevant mistakes. State the command implementation and any tested version/platform. Never imply an example has been tested when it has only been reviewed.

Use plain Markdown according to the [proposed format](docs/content-format.md). Keep examples safe to understand and distinguish read-only inspection from commands that modify data. Use synthetic paths/data. Do not include credentials or ask learners to paste secrets into examples.

Describe shell behavior separately from tool options. Prefer quoted filenames/patterns and explain why quoting matters. Identify portability differences explicitly. Avoid undocumented assumptions about aliases, GNU utilities, or optional packages being installed.

## References and reuse

Link to upstream manuals and useful related projects. Do not copy cheat sheets or manual passages merely because a repository is public. Record source, license, attribution, and changes for any approved imported material. The repository's own code/content licenses must be selected before accepting substantial external content contributions or publishing reusable releases.

## Review

Check relative Markdown links, front matter examples, consistent command spelling, and GNU/Apple/BSD distinctions. Once implemented, content validation must include CLI/web rendering parity, implementation-specific example checks, and accessible output. Do not execute destructive examples against a contributor's normal filesystem.

No test suite, package installation command, or contribution CI exists yet. Add accurate instructions alongside implementation rather than documenting commands that do not work.
