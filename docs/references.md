# Related projects and references

Research date: **2026-10-23**, America/Los_Angeles. References are for comparison and implementation research; inclusion is not endorsement, a maintenance guarantee, or permission to copy content. No tool was installed or independently benchmarked during this planning step.

## Command explanations and example libraries

| Project | Links | Relevance |
| --- | --- | --- |
| tldr-pages | [Repository](https://github.com/tldr-pages/tldr), [site](https://tldr.sh/), [web client](https://tldr.inbrowser.app/), [similar projects](https://github.com/tldr-pages/tldr#similar-projects) | Concise, example-focused Markdown reference; key comparison for content authoring |
| ExplainShell | [Repository](https://github.com/idank/explainshell), [web tool](https://explainshell.com/) | Parses command lines and matches arguments/operators to documentation |
| explain.sh | [Repository](https://github.com/benjamine/explain.sh), [project page](https://benjamine.github.io/explain.sh/) | CLI frontend to explainshell.com |
| explainshell-cli | [Repository](https://github.com/tamnd/explainshell-cli) | Go library/client exploration; reviewed README describes the CLI as an early scaffold, not a finished explainer |
| Shell-AI | [Repository](https://github.com/Deltik/shell-ai) | AI-backed explain subcommand and optional explain alias; different cost/offline tradeoff |
| helpshell | [Repository](https://github.com/HelgeSverre/helpshell) | Searchable local command/config catalog; candidate reference for shell-specific discovery |

The standalone curated-reference goal differs from explaining arbitrary pipelines and from generating model answers. Normal command lookup should not require a hosted service. License compatibility and the licenses of underlying manual/content sources must be reviewed independently of a tool's implementation license.

## Complete list from tldr's similar-projects section

The following eight entries were listed in [tldr's similar-projects section](https://github.com/tldr-pages/tldr#similar-projects) at research time. Descriptions are brief summaries; the list may change upstream.

| Project | Links | Main idea to investigate |
| --- | --- | --- |
| cheat.sh | [Service](https://cheat.sh/), [repository](https://github.com/chubin/cheat.sh) | Unified access to multiple cheat-sheet sources |
| devhints | [Site](https://devhints.io/) | Broader programming reference and cheat-sheet presentation |
| eg | [Repository](https://github.com/srsudar/eg) | Explained CLI examples with personal additions |
| kb | [Repository](https://github.com/gnebbia/kb) | Minimal personal command-line knowledge-base management |
| navi | [Repository](https://github.com/denisidoro/navi) | Interactive command/example discovery |
| Cheat | [Repository](https://github.com/cheat/cheat) | Creating and viewing command-line cheat sheets |
| Command Line Interface Pages | [Organization](https://github.com/command-line-interface-pages) | Standardized help for commands, directories, and configuration |
| bropages | [Repository](https://github.com/pombadev/bropages) | Community examples/voting; tldr labels this project deprecated |

Interactive execution/completion features in some projects are not part of this tool's initial scope. Compare their content formats, browsing, and contribution workflows without automatically adopting command execution.

## Primary manuals and implementation references

- [GNU grep manual](https://www.gnu.org/software/grep/manual/): primary reference for GNU grep examples and options; validate current content during authoring.
- [GNU Bash manual](https://www.gnu.org/software/bash/manual/): shell syntax/builtin reference to consult during implementation.
- [Homebrew grep formula](https://formulae.brew.sh/formula/grep): prefixed GNU binaries and PATH selection on macOS.
- [mvdan/sh](https://github.com/mvdan/sh): Go shell parser candidate for a later full-command feature.
- [ShellCheck](https://github.com/koalaman/shellcheck): shell analysis reference and potential example-validation aid; not a command-explanation database.

The manual links above are authoring starting points, not proof that any proposed example has been tested. Obtain Apple/BSD/BusyBox references for each variant when coverage is added rather than applying GNU documentation universally.

## Implementation and distribution documentation

- [Go build/environment variables](https://pkg.go.dev/cmd/go#hdr-Environment_variables) and [Go embed](https://pkg.go.dev/embed).
- [Rust platform support](https://doc.rust-lang.org/rustc/platform-support.html), for the language comparison.
- [Astro content collections](https://docs.astro.build/en/guides/content-collections/) and [Markdown support](https://docs.astro.build/en/guides/markdown-content/).
- [GoReleaser](https://goreleaser.com/), a candidate release-automation tool.
- [Homebrew tap maintenance](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap), [adding software](https://docs.brew.sh/Adding-Software-to-Homebrew), and [formula acceptance](https://docs.brew.sh/Acceptable-Formulae).
- [mise GitHub backend](https://mise.jdx.dev/dev-tools/backends/github.html) and [registry contributions](https://mise.jdx.dev/contributing/registry.html).

## Evidence and licensing limits

Related-project links are references, not dependencies. Some descriptions were checked directly; the complete similar-project inventory comes from tldr's published list. Capabilities, activity, package-manager policy, and licenses can change. Check the exact release and relevant licenses before implementation or redistribution. Do not infer compatibility or permission from the absence of a paywall.
