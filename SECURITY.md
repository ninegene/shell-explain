# Security policy

## Reporting a vulnerability

Please report security problems privately, not in a public issue or pull request.

Use GitHub's private vulnerability reporting: open the [Security tab](https://github.com/ninegene/shell-explain/security) and choose **Report a vulnerability**, or go straight to the [advisory form](https://github.com/ninegene/shell-explain/security/advisories/new). Only the maintainers can see the report.

Please include what you found, the version or commit (`explain --version`), the steps to reproduce it, and what you think the impact is. A minimal example helps more than a long description.

## What this covers

- the `explain` command-line tool and the way it reads files, environment variables and the `PATH`;
- the release archives, their checksums and the installer script;
- the website's hosting configuration under `infra/` and the GitHub Actions workflows.

By design, `explain` never runs the example commands it shows, and the only programs it runs are explicit, bounded version probes of the command you ask about (`--no-probe` turns them off). A report that shows it running anything else, writing outside its own output, or sending data over the network is in scope.

Wrong, outdated or unsafe-to-copy command documentation is not a vulnerability: please open a normal issue or pull request for it.

## What to expect

This is a small, volunteer-maintained project, so replies are best effort. The maintainers will confirm that they received a report, ask for more detail if needed, work on a fix and publish a GitHub security advisory when a fix is released. You will be credited in the advisory unless you prefer not to be.

## Supported versions

Fixes go into the latest release and the `main` branch. Older releases are not patched; upgrade to the latest.

## Secrets

If you find a credential, key or token in this repository, in its history or in a release artifact, report it privately as above, and do not use it.
