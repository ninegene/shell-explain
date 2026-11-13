# Pre-publication review

Review performed 2026-11-11 against the local repository and GitHub metadata. This records checks and remaining work; it does not authorize making the repository public or publishing a release.

## Checks performed

- Scanned all local branches and the histories of all 54 PR heads available locally with Gitleaks: 134 commits scanned, no matches. The current tree, GitHub PR discussion data, generated site, and printable strings from the locally built binary also produced no Gitleaks matches. Automated scanning cannot prove that sensitive information is absent.
- Searched commit patches and PR discussions for home-directory paths, email addresses, and links to other workspace repositories. Matches inspected were synthetic example addresses or command names. Commit author and committer addresses use GitHub noreply or GitHub domains. All 54 PR descriptions, 12 conversation comments, 92 reviews and 24 review threads returned by the GitHub API were included; none of those collections were truncated.
- Checked GitHub repository metadata: no open PRs, standalone issues, or releases; the additional remote branch is already contained in main. The referenced workspace repository is public. Repository visibility remains private.
- Checked tracked filenames and historical paths for credential-like filenames, archives and database files; none were found. The Go module file has no replacement or private module dependency, and all 347 locked npm package URLs point to the public npm registry. This does not establish the safety or license compliance of every package.
- Checked 50 local Markdown links: none were missing. Of 398 unique non-example external references and Markdown links tested with HTTP HEAD, 260 returned success; the one 403 response was independently accessible, and one 404 was replaced by the upstream OpenSSH manual. All 136 remaining requests timed out or failed to connect to `www.gnu.org`; they were not verified. The check did not validate URL fragments or all URLs embedded in examples or dependencies.
- A subsequent offline cross-check against installed GNU Info manuals matched the nodes behind 101 of the original Coreutils references, six Diffutils references, one Findutils reference, and one gawk reference. Four SHA-2 links named nodes absent from Coreutils 9.11; they now refer to its shared `sha2 utilities` section. Info node matches do not verify current HTTP availability or anchors. On 2026-11-12, a new HEAD sweep of the 133 current GNU references produced 133 network errors; the GNU host also timed out from the host, the web reader, and a Linux container.
- On 2026-11-12, a fresh Debian Linux ARM64 container installed the locked dependencies and Chromium and passed `just test` and `just check`; a second fresh container passed the development-server test and 645 browser tests using `CHOKIDAR_USEPOLLING=true` for the Colima shared mount. Unavailable host implementations and GNU `factor -h` on coreutils 9.1 were skipped. This did not exercise a public fork or an independent contributor.
- Reviewed the local CLI and site artifacts and Go dependency notices as recorded in [third-party notices](third-party-notices.md). No public release artifacts exist yet. Corrected stale learning-guide status in the README and guide plan.

## Remaining before publication

1. Recheck the current GNU reference URLs and relevant anchors from a network that can reach `www.gnu.org`; investigate any actual failures. Four of the original 136 references were replaced with one shared SHA-2 reference after the offline node check. The link check results are time-sensitive.
2. Arrange an independent contributor trial, including the public fork and PR process when available. Fresh-container and macOS test evidence is in [candidate starter issues](starter-issues.md).
3. Inspect the exact binaries, archives, site output, dependencies and notices when release artifacts are prepared. Complete the provenance and redistribution review for any third-party or AI-assisted material; the checks above do not establish authorship of each passage or certify all dependency licenses.
4. Review repository visibility and public contribution settings with the owner. Public visibility, releases, website deployment and DNS changes require their own authorization.

History and PR discussion checks cover the branches, PR heads and API collections available at the review time; recheck after further changes and before changing visibility. Do not rewrite history or remove discussions without a separate decision.
