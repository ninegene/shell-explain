# Contributing

The project has a local shared-content, CLI, and Astro implementation and remains private through stage 4. External contributions are planned for stage 5. This guide does not imply a public contribution channel exists today.

## Setup and checks

Install `just`, Go 1.24 or newer, Node.js 22.23.4 or newer, and npm. From the repository root, install the locked web dependencies and the browser used by the web tests:

```sh
npm ci --prefix web
cd web && npx playwright install chromium && cd ..
```

On Linux, the browser may also need system libraries. For supported distributions, run `cd web && npx playwright install --with-deps chromium && cd ..` instead; installing system packages requires administrator privileges. See [Playwright browser setup](https://playwright.dev/docs/browsers#install-system-dependencies).

Run `just` to list commands. Run `just test` for Go and available host-tool example checks, `just check` for Go vet, formatting and Astro diagnostics, and `just web-test` for the development-server and browser checks. `just web-build` builds the static site locally; `just build` builds `bin/explain`. See [development and verification](docs/development.md) for prerequisites, skips, fixture isolation and platform limitations. Record skipped checks and untested implementations in the PR. If a container or VM shared mount misses file-change events, run `CHOKIDAR_USEPOLLING=true just web-test`; polling uses more resources.

## Contribution process

Once public contributions open:

1. Search existing pages and issues. Open an issue for a new command, implementation, or larger change; small corrections can go directly to a pull request.
2. Fork the repository, clone your fork and create a branch. Make changes to the canonical Markdown in `content/en/`, not generated JSON or separate CLI/web copies. Read the [content format](docs/content-format.md) and use an existing page as a schema example.
3. Explain the use case and options, cite upstream references, set the page's `updated` date when changing claims or examples, and record the operating system, implementation and version tested. Add a safe fixture and expected result when applicable. Follow the format's rules for examples that cannot be run automatically.
4. Run the checks above relevant to the change. Describe the change, evidence, checks, and any skips in the pull request. Use screenshots or terminal output where helpful.
5. Maintainers assess accuracy, portability, accessibility, license provenance, and shared CLI/web rendering. Address review feedback; maintainers merge accepted changes. Do not imply untested versions are supported.

See [candidate starter issues](docs/starter-issues.md) for small contributions. Issue and PR templates are prepared for the public launch; the repository remains private until separately authorized.

## Command documentation

Write original explanations with a concrete use case, understandable examples, option meanings, expected behavior, and relevant mistakes. State the command implementation and any tested version/platform. Never imply an example has been tested when it has only been reviewed.

Use plain Markdown according to the [content format](docs/content-format.md). Keep examples safe to understand and distinguish read-only inspection from commands that modify data. Use synthetic paths/data. Do not include credentials or ask learners to paste secrets into examples.

Describe shell behavior separately from tool options. Prefer quoted filenames/patterns and explain why quoting matters. Identify portability differences explicitly. Avoid undocumented assumptions about aliases, GNU utilities, or optional packages being installed.

## References and reuse

Link to upstream manuals and useful related projects. Do not copy cheat sheets or manual passages merely because a repository is public. Record source, license, attribution, and changes for any approved imported material. Contributions of original software and shell examples are accepted under the [MIT License](LICENSE); original explanatory content and documentation are accepted under [CC BY 4.0](LICENSE-CC-BY-4.0). Submit only material you have the right to contribute under those terms. Third-party material requires compatible permission and its own required notices.

## Review

Check relative Markdown links, front matter, consistent command spelling, and GNU/Apple/BSD distinctions. Content validation includes CLI/web rendering parity, implementation-specific example checks, and accessible output. Do not execute destructive examples against a contributor's normal filesystem.

Follow [development and verification](docs/development.md): run `just test`, `just check`, and `just web-test` after installing web/test-browser dependencies. The latter builds the site and checks its selected sections against CLI output. Add a temporary fixture and expected result when adding an executable example. The `Native Linux AMD64` GitHub Actions workflow runs these checks and a CLI smoke test on every pull request. A trial by an independent contributor remains stage-5 work.
