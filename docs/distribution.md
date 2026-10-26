# Distribution and release plan

The CLI can be built locally with just build and the Astro site can be previewed locally. No public binaries, installers, package definitions, releases, or hosting exist yet. Distribution commands in this document are proposed future interfaces.

The repository stays private through stages 1–4. Stage 5 includes selecting code/content licenses, preparing external contribution guidance, reviewing publication safety, making the repository public, and publishing the site and release artifacts. Private development builds do not constitute a public release. See the [roadmap](roadmap.md) for the sequence and authorization boundaries.

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

Use **GitHub Pages** for the static Astro site at **shell-explain.ninegene.com**. Build and deploy with GitHub Actions following [Astro's GitHub Pages guide](https://docs.astro.build/en/guides/deploy/github/). No paid cloud provider, application server, or database is planned for website hosting. Domain registration remains a separate cost.

During private stages 1–4, use local previews without publishing the website. In stage 5, make the repository public before enabling Pages. GitHub Free supports public-repository Pages; private-repository Pages requires a qualifying paid GitHub plan. See [Pages availability](https://docs.github.com/en/pages/getting-started-with-github-pages/what-is-github-pages).

The publication checklist is to configure Pages with GitHub Actions as its publishing source, build/upload the static Astro artifact, deploy through the Pages environment with minimal workflow permissions, and verify routes/search/assets. Configure Astro's site URL for the custom domain; do not carry over a repository-name base path intended for a github.io project URL.

Verify domain ownership, configure the repository's custom domain and DNS records, enable HTTPS, and test canonical URLs and redirects. Follow [GitHub's custom-domain guidance](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site/about-custom-domains-and-github-pages). Pages settings, deployment, and DNS updates require separate authorization; this plan does not perform them. Store CLI downloads in GitHub Releases rather than bundling release archives into the website.

The site should show its content revision, link matching CLI releases, and label unreleased installation instructions. Prefer tagged release content for reproducibility; if deploying newer content from the default branch, label the difference clearly. Do not silently claim parity with an older installed binary.
