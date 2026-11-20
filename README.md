# Shell Explain

Practical command explanations for the terminal and the web, maintained in one shared Markdown collection.

- Repository: [ninegene/shell-explain](https://github.com/ninegene/shell-explain), public since 2026-11-17. Software and shell examples are MIT-licensed ([LICENSE](LICENSE)); original explanatory content is CC BY 4.0 ([LICENSE-CC-BY-4.0](LICENSE-CC-BY-4.0)). Issues and pull requests are welcome; see [CONTRIBUTING.md](CONTRIBUTING.md). Security reports: [SECURITY.md](SECURITY.md).
- Executable: **`explain`**, for example `./bin/explain grep` after a local build.
- Website: **shell-explain.ninegene.com**, built with Astro and hosted on AWS (S3 and CloudFront, defined in Pulumi under [infra](infra/README.md)).
- CLI language: Go, with original Markdown embedded for offline use.
- Status: **stages 1–4 implemented locally**. Current CLI and web usability are accepted as sufficient. The command collection (194 pages: GNU, macOS, BusyBox, procps, Bash builtins and single-project tools) is complete, and all planned beginner guides are written. The website is deployed, and the `explain` CLI is released as binaries (see [Install](#install)).

## Goal

Explain what a command does, show useful examples, explain common options and mistakes, and select documentation appropriate to the installed implementation. Beginners should be able to understand the examples; experienced users should be able to find an unfamiliar option quickly.

The CLI and website must consume the same `.md` files. Neither keeps an independently edited copy of the explanations. The CLI installs from release binaries, a curl-based installer, Homebrew, or mise.

Use local Astro previews during development and iteration. A GitHub Actions workflow builds the static site and deploys it to S3 and CloudFront once the infrastructure is applied; hosting costs a few cents a month and does not require a public repository. Domain registration remains separate.

## Documentation

| Document | Purpose |
| --- | --- |
| [Product and architecture](docs/architecture.md) | Scope, Go choice, Astro site, and integrations |
| [Shared Markdown format](docs/content-format.md) | Implemented metadata, prose conventions, and rendering contract |
| [Implementation detection](docs/detection.md) | GNU/BSD/Apple/BusyBox selection and shell limitations |
| [Distribution](docs/distribution.md) | Releases, curl, Homebrew, mise, and reproducibility |
| [Related projects and references](docs/references.md) | tldr, ExplainShell, all projects in tldr's similar-projects list, and supporting tools |
| [Command collection checklist](docs/command-collection.md) | Stage 4 inventory, implementation coverage, and completion criteria |
| [Learning guides](docs/learning-guides.md) | Beginner guide series on the terminal, shell, operating systems and tool history |
| [Roadmap](docs/roadmap.md) | Decisions, validation, and implementation stages |
| [Development and verification](docs/development.md) | Build commands, local previews, test coverage, and known limits |
| [Contributing](CONTRIBUTING.md) | Documentation contributions and source/licensing expectations |

## Install

Pick one. Binaries are published for macOS and Linux on ARM64 and AMD64.

```sh
brew install ninegene/tap/shell-explain
```

```sh
mise use -g github:ninegene/shell-explain
```

```sh
curl -fsSL https://github.com/ninegene/shell-explain/releases/latest/download/install.sh | sh
```

The install script verifies the SHA-256 checksum, installs into `~/.local/bin`, never uses sudo and refuses to replace an existing `explain` without `--force`; add `--dry-run` to preview. Then run `explain --version` and `explain grep`. Each [release](https://github.com/ninegene/shell-explain/releases) also has the archives, `checksums.txt` and a build provenance attestation (`gh attestation verify`). The same steps are on the [website](https://shell-explain.ninegene.com/install/); packaging details are in [Distribution](docs/distribution.md).

## Local quick start

Requires [just](https://github.com/casey/just), Go 1.26+, and Node 22.23.4+. Install the command runner with `brew install just` on macOS.

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

The site has the command reference at `/commands/`, the beginner guides at `/learn/` (a "Start here" section on the home page links to them), and the glossary at `/learn/glossary/`.

The currently implemented interface supports these forms once the binary is on PATH:

```console
explain
explain browse extract archive
explain grep
explain grep --details
explain ggrep
explain grep --option=-n
explain --implementation gnu grep
explain --platform linux grep
explain --format json grep
explain search logs
explain search extract archive
explain list
explain learn
explain learn big-picture
explain glossary kernel
```

`explain learn` lists the beginner guides, and `explain learn GUIDE` (a name or a number) reads one, optionally a single `--section`. `explain glossary` lists short definitions of words beginners meet, and `explain glossary WORD` reads one. `explain search` also finds guides and glossary words, and matches what you want to do against the descriptions and commands of every authored example (up to 10 shown; `--limit N` shows more). On a terminal, text output opens in your pager (`$PAGER`, or `less`, which prints short output directly); use `--no-pager` or set `EXPLAIN_NO_PAGER=1` to print it directly. JSON is never paged.

### Interactive browser

On a terminal, `explain` with no arguments (or `explain browse [WORDS]`) opens a browser: a search box, a list of commands and a preview of the selected one. The two column titles (`Commands i/n`, `Examples i/n`) show your position, and the focused pane's title is bold. Typing filters the list the way `explain search` does, by command metadata or by a matching example, and the preview shows the matching examples (with a hint to press `enter` for the command's full reference). With no search, or a command that matched only by name, it shows the same text `explain COMMAND` prints. It uses no network, starts no process to detect the implementation (like `--no-probe`), and never runs an example.

| Key | Action |
| --- | --- |
| type | Filter the list |
| `↑` `↓`, `ctrl+p` `ctrl+n` | Move through the list |
| `tab` | Focus the preview; `tab` or `esc` returns |
| `↑` `↓` (or `←` `→`, `p` `n`) in the preview | Select the previous or next example; the preview scrolls to keep it in view |
| `j` `k` in the preview | Scroll a line |
| `ctrl+y`, or `y` in the preview | Copy the selected example's code |
| `enter` | Toggle the full reference (like `--details`); the preview title then reads `Details · Examples` |
| `pgup` `pgdn` | Scroll the preview |
| `esc` | Clear the search; on an empty search, quit |
| `ctrl+c` | Quit |

Copying uses `pbcopy` on macOS and `wl-copy`, `xclip` or `xsel` on Linux, and falls back to the terminal's OSC 52 clipboard sequence when none is installed.

When you run `explain` with no arguments on a terminal, it shows the usage (in the pager), then asks whether to keep showing it. Answer `y` (or just press enter) and you are asked again next time; answer `n` and later runs open the browser directly. Only `n` is kept in `explain/config` in your user config directory (on macOS, `~/Library/Application Support`; on Linux, `$XDG_CONFIG_HOME` or `~/.config`). Delete that file to be asked again; `explain --help` always shows the usage. When standard input or output is not a terminal, `explain` with no arguments just prints the usage, and `explain browse` exits with an error pointing to `explain search`.

`explain` provides information; it does not run the example command or install the command being explained. Full pipeline/script explanation is a possible later feature, not part of the first release.

Text output shows example usages for the selected implementation by default, under a heading with the command's summary (and the chosen implementation when there is a choice). A short command with a short description takes one line, `command  # description`, with the comments aligned; anything longer keeps its description as `# ` comment lines above the command. Commands are never changed or wrapped. On a terminal, command names, flags, quoted strings and comments get different colors. The output ends with links to the command's page on the website and to the upstream manuals. In the Markdown source and on the web, each example keeps its own code block and copy button so a single command can be copied; the terminal layout is only a different presentation of the same examples, using a short label from the code fence (`short="..."`) when there is one. `--details` and the website keep the full descriptions. Detailed option explanations are grouped under `Options`, with each option spelling above its description. Descriptions wrap at 100 columns (including the `# ` prefix on example comments); authored commands and indivisible inline literals remain intact. Use `--details` before or after the command for the full overview, options, mistakes, portability notes, detection evidence, and references. Explicit `--option` lookup shows that option's explanation without requiring `--details`. JSON and the website retain the full selected reference.

## Project boundaries

This project is being developed as a reusable reference tool, now public; release binaries and installers are available (see Install). It is not an interactive lab host. Other applications can eventually install a pinned binary or consume its documented JSON output. Lesson-specific hints, grading, accounts, and progress belong to those applications.

GNU/BSD behavior must be selected per command, not inferred solely from the operating system. The website cannot inspect the visitor's terminal; it must offer explicit implementation/platform selection.

## Development status and licensing

Use `just test` for Go checks and reviewed examples, `just check` for Go vet and Astro diagnostics, `just web-build` for the static build, and `just web-test` for Chromium parity/usability/accessibility tests. See [development](docs/development.md) for setup and limitations. The initial corpus is grep, sed, find, and cd; exact version support is not certified universally.

Research date: **2026-10-23**, America/Los_Angeles. Related-project descriptions come from published documentation, not independent trials. See [references](docs/references.md).

The original software, including scripts, tests, and shell examples in documentation, is licensed under the [MIT License](LICENSE). Original documentation and explanatory content, including command pages, guides, and glossary prose, is licensed under [CC BY 4.0](LICENSE-CC-BY-4.0). Shell examples within that content are MIT-licensed. Attribute content to “Shell Explain contributors” and link to the source and license, and indicate changes when sharing it. These licenses apply to rights held by the contributors; third-party materials retain their own licenses and notices.

The CLI embeds the content, so binary distributions must include the applicable MIT and CC BY notices. See [third-party notices](docs/third-party-notices.md) for the dependency inventory and remaining release-artifact checks.
