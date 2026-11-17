# Distribution and release plan

The CLI can be built locally with just build and the Astro site can be previewed locally. No public binaries, installers, package definitions, releases, or hosting exist yet. Distribution commands in this document are proposed future interfaces.

The repository was private through stages 1–4 and became public on 2026-11-17. Stage 5 continues with contribution guidance and publishing release artifacts. Development builds do not constitute a public release. Binary distributions must include the MIT and CC BY 4.0 license texts and applicable attribution for embedded content and third-party dependencies. See [third-party notices](third-party-notices.md) for the inventory and release checks, and the [roadmap](roadmap.md) for the sequence and authorization boundaries.

## Release artifacts

Start with macOS and Linux on ARM64 and AMD64. Publish versioned GitHub release archives containing the executable named `explain`, relevant notices, and documentation. Include SHA-256 checksums and later verifiable provenance/signatures. Test Linux compatibility explicitly; avoid accidental dependencies on a build host's libc or C libraries.

Embed the Markdown-derived reference bundle so normal lookup works offline without Python, Node, a database, or an AI API. Release the content revision alongside the binary. Keep unknown/missing content responses useful. A future content-update mechanism must be explicit, verified, and compatible with the binary schema; it is not required for version one.

Release automation such as [GoReleaser](https://goreleaser.com/) is a candidate. Review the current free feature set and package-manager policies before relying on any publisher. A release tool supporting a packaging route does not establish that Homebrew's official repositories will accept it.

## curl-based installer

Provide an HTTPS installer and a manual archive-download alternative. Detect OS/CPU, select an immutable versioned artifact, verify its digest, and install to a user-writable destination such as `~/.local/bin`. Support version and destination flags, useful errors, and a dry-run mode. Avoid sudo and automatic shell startup edits by default.

Downloading an installer to inspect before running should be documented alongside any one-line pipe-to-shell shortcut. Checksums obtained from the same compromised source do not independently establish authenticity; signed release metadata/provenance improves the trust story.

Do not overwrite an existing command or managed installation silently. Explain PATH conflicts, including a pre-existing explain alias/function. Uninstallation should remove only files owned by this tool and should be documented.

## Homebrew

Start with an upstream-maintained tap and an open-source CLI formula. Proposed invocation: `brew install ninegene/tap/shell-explain`; the package name is shell-explain, while the installed executable is explain. A separate public tap repository would require its own creation later.

Use immutable sources, declared build dependencies, meaningful tests, and checksums. Official Homebrew inclusion is a separate acceptance process; it is not required for initial availability. Follow current policy rather than assuming a generated binary-only cask is appropriate for an open-source CLI. References: [tap maintenance](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap), [adding software](https://docs.brew.sh/Adding-Software-to-Homebrew), and [acceptable formulae](https://docs.brew.sh/Acceptable-Formulae).

## mise

Proposed invocation: `mise use -g github:ninegene/shell-explain`. Use consistent release asset naming so the backend selects the right executable/architecture. Verify real installation from a clean environment before advertising the command.

A central registry entry is only a convenience shorthand and may require adoption evidence. Current guidance supports using the GitHub backend without one. Avoid designing new integration around the deprecated ubi backend. References: [GitHub backend](https://mise.jdx.dev/dev-tools/backends/github.html) and [registry policy](https://mise.jdx.dev/contributing/registry.html).

## Website publishing

Host the static Astro site at **shell-explain.ninegene.com** on AWS: a private S3 bucket behind CloudFront with an ACM certificate, defined in Pulumi under `infra/pulumi` and deployed by GitHub Actions through an OIDC role (no stored keys). It reuses the account, Pulumi state backend and `ninegene.com` Route 53 zone of the personal site; see [infra/README.md](../infra/README.md). It costs a few cents a month. No application server or database is planned. This replaces the earlier GitHub Pages plan, which would have required making the repository public first.

The hosting does not depend on repository visibility. Publishing the site and making the repository public were separate decisions, each with its own authorization.

The publication checklist is to apply the Pulumi stack, configure the `production` GitHub environment, run the deploy workflow (minimal permissions, OIDC), and verify routes, search and assets on the live host. Astro's site URL is the custom domain; do not carry over a repository-name base path.

Test canonical URLs, the trailing-slash redirect and the 404 page after the first deploy. Applying the stack, deployment and DNS updates require separate authorization; this documentation does not perform them. Store CLI downloads in GitHub Releases rather than bundling release archives into the website.

The site should show its content revision, link matching CLI releases, and label unreleased installation instructions. Prefer tagged release content for reproducibility; if deploying newer content from the default branch, label the difference clearly. Do not silently claim parity with an older installed binary.
