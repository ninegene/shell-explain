# Candidate starter issues

These are proposed issues for maintainers to review and publish after contributions open; they are not existing GitHub issues.

- **Clarify a confusing command example:** pick a documented example, identify the reader's goal, and improve its explanation without changing the command. Cite any new factual claims and run the content and Go checks.
- **Add an example to an existing command:** propose a common use case with an upstream reference, implementation and tested version. Add a safe fixture and expected result where applicable; verify CLI and web rendering.
- **Improve a portability note:** compare two documented implementations against their upstream manuals and available tools, then state the difference and evidence without claiming untested version coverage.

See [contributing](../CONTRIBUTING.md) for scope, validation, and contribution terms. Maintainers should select concrete pages and confirm the work remains open before creating issues.

## Preparation check (2026-11-11)

In the existing macOS development environment, `npm ci --prefix web`, Playwright Chromium setup, `just test`, `just check`, and `just web-test` passed; the browser suite passed 645 tests. The uv project example accepts different generated layouts while checking the common project files. Host-tool checks may skip unavailable implementations. This was not a fresh-machine or external-contributor trial.

On 2026-11-12, a fresh Debian Linux ARM64 container with Go 1.27.3, Node 24.22.1 and just 1.59.0 installed the locked dependencies and Chromium. `just test` and `just check` passed; a second fresh container passed the development-server test and all 645 browser tests with `CHOKIDAR_USEPOLLING=true` for the Colima shared mount. Host-tool checks skipped unavailable implementations and GNU `factor -h`, which coreutils 9.1 does not support. This was not an independent human contributor or a public fork/PR trial.
