# Release artifact and provenance review

This is a preparation checklist, not release approval. Version `v0.1.0` has been released; repeat the inspection for each exact release candidate and keep its evidence with the release record. The [third-party inventory](third-party-notices.md) is a starting point, not an exhaustive certification.

## Local candidate reviewed on 2026-11-15

At commit `087a446c7a49ef09a6d638587851bfe94aa59dfe`, `just build` and `just web-build` produced local macOS ARM64 and static-site candidates. `go version -m bin/explain` identified Go 1.27.3, goldmark v1.8.6 and go.yaml.in/yaml/v3 v3.0.5. The checked-in Go runtime LICENSE/PATENTS and both module LICENSE files and yaml NOTICE matched the installed build inputs byte for byte. The checked-in Apache-2.0 text accompanies the yaml notice. The candidate binary was not a release archive and its toolchain settings do not establish compatibility with other release targets.

The generated site contained 737 HTML files, one favicon, one CSS file and four JavaScript files. Inspection found one distinct inline catalog script and external scripts corresponding to the site's copy-button source and generated imports. This is an inspection of that build, not proof that future bundles contain no third-party code. At that time, no archive contents or released website were available to inspect.

## Candidate review on 2026-11-28

At commit `d2878a0eb1ce0ecbe586ef002703b83c330f10c0`, the [release dry run](https://github.com/ninegene/shell-explain/actions/runs/46118292392) built `0.0.0-dispatch` archives using Go 1.27.3 and passed native smoke tests for macOS and Linux on AMD64 and ARM64. The four downloaded archives passed `shasum -a 256 -c checksums.txt`; each contained 29 notice and documentation files matching the source checkout. The 21 embedded Go module versions in each binary appeared in the CLI notice inventory. This verifies those dry-run files, not the future `v0.2.0` archives or their attestations.

| Dry-run archive suffix | SHA-256 |
| --- | --- |
| `darwin_amd64` | `feec457b1530bb65fe4335268af1363e516a432d6785c30332dc892b403ceabd` |
| `darwin_arm64` | `3de6730bffe991157be316e9e1cdfcad5b16e27ab190e9bc491ee9a16e77454a` |
| `linux_amd64` | `bef3de17e06a4ae632886d07fd04df2ebf02fe751135039a8472492cdd5b96d8` |
| `linux_arm64` | `41d38c3cf2a6473c6c2f0df1d7b1702a57067c470fae8e8eeb704a187e1d4cdd` |

Separately, four local `v0.2.0`-labeled archives were cross-compiled in an ARM64 container with Go 1.27.3 from that commit's code and uncommitted release-documentation changes. Their checksums passed, and all 29 packaged notice and documentation files matched the working tree at build time. The macOS ARM64 archive passed the CLI smoke matrix natively; the Linux ARM64 archive passed version, lookup, search and guide checks in an ARM64 container. The local AMD64 archives were not executed on native AMD64. These are development candidates, not the final tagged artifacts, and their contents and checksums can change.

An isolated Linux ARM64 build of that commit used Go 1.27.3, Node 26.13.1 and `npm ci` against `web/package-lock.json`. Its 738 HTML pages, one authored SVG, one CSS file and five JavaScript files were inspected. The SVG matched `web/public/favicon.svg`; the scripts corresponded to the project's catalog and copy-button sources and their generated imports. No third-party client code or font assets were identified in the inspected output. The generated About page included the source link, MIT and CC BY 4.0 references. This finding applies only to that build; review each future output and any redistributed third-party material separately.

## Exact artifact checklist

For each binary target and published site build:

1. Record the source revision, build command, toolchain and dependencies, target platform, artifact names and checksums. Inspect the actual archives and site output, including inline scripts, CSS, fonts, images, generated assets, and external references. Distinguish build-only dependencies from material actually redistributed; investigate uncertain bundled code instead of inferring inclusion from a lockfile.
2. Compare the exact Go module graph and toolchain licenses and notices against the files packaged with each archive. Include the project MIT and CC BY 4.0 texts, content attribution and source link, the applicable Go and dependency license/copyright texts, patent grant and any applicable NOTICE text. Review Apache-licensed modifications and notices when relevant. Confirm these files are actually inside every archive, not just present in the source checkout.
3. Inspect site output and build inputs for third-party materials and notices. Include the required copyright, license, attribution, modification, and source notices for material actually distributed. Ensure the deployed site links to the source, content license and attribution and describes its public status accurately. Recheck native dependency obligations separately if native build tools are redistributed.
4. Test each artifact on its claimed platform, verify the checksums against the exact files to be released, and document skipped targets and remaining compatibility questions. Review the publishing and integrity process before release; checksums from the same source are not independent proof of authenticity.

## Native Linux AMD64 verification in CI

The candidate at source revision `e63fcc22e4c989ffb5c0f6591e641d365368cce2` passed a nine-case CLI matrix under Rosetta with normal GC; native Linux AMD64 was not verified for it. Its local archive was `shell-explain-candidate-e63fcc22e4c9-linux-amd64.tar.gz` (SHA-256 `bc587eaab6eaebe6581db3a0e5bec5b0b56095e7c24f14795e0dd5dd44d81430`; executable SHA-256 `f3f2e1f62f64483176e760425a522ae0afb42fe2fa57da1d8f37436a33bb296a`). On 2026-11-16 a manual workflow could not download that archive from a private draft release: the read-only workflow token received HTTP 403, so the archive never ran.

Native checks now build from source instead. The `Native Linux AMD64` workflow runs on GitHub's `ubuntu-26.04` x64 runner for pull requests, pushes to `main` and manual dispatch, with read-only repository access. It builds `explain` with `-trimpath`, records `go version -m` and the executable's SHA-256, runs [scripts/cli-smoke.sh](../scripts/cli-smoke.sh) (the same nine cases, with `GOGC`, `GOMAXPROCS` and `GODEBUG` cleared), then the Go tests, `go vet`, and website formatting and Astro checks. It also runs the Playwright browser tests; the dev-server watcher test that `just web-test` adds is not run in CI for now, because it failed on the runner. It uploads the smoke-test output as evidence and never publishes anything.

A passing run verifies the binary CI built from that commit, not a separately built archive. Before a release, test the exact archives to be released on each claimed platform (item 4 above), preferably by building them in the same CI run that tests them.

## Provenance and AI-assisted material

The owner reports that much of the content was AI-generated. A link to an upstream manual and an executable example support factual review; they do not prove independent authorship or permission to redistribute wording, code, images, or other expressive material. No per-passage provenance or complete similarity audit has been established by this review.

For further provenance work, review available generation and editing records, source materials, model/service terms applicable when used, and any identifiable human contributions. Record what can be established, known third-party reuse, the applicable permission and attribution, and any missing evidence without inventing authorship. Investigate close matches or uncertain imported material, and replace, obtain permission for, attribute, or otherwise resolve it as appropriate. Keep unresolved cases visible for an owner decision. Review the copyright and licensing implications with qualified counsel if needed; a CC license applies only to rights the licensor can grant, and the copyright status of AI-generated material depends on the applicable law and human contribution. See [Creative Commons' guidance on generative AI](https://creativecommons.org/2023/08/18/understanding-cc-licenses-and-generative-ai/).

The owner marked the [content provenance review](content-provenance-review.md) complete with its uncertainties documented. Revisit those questions if new evidence or additional third-party material emerges. Do not imply that an artifact inventory, automated comparison, or a license label alone establishes rights in every passage.
