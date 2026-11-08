# Stage 4 command collection checklist

Stage 4's next work is to expand the existing four-command collection. The current CLI and web usability are accepted as sufficient; this checklist tracks planned content, not completed explanations.

## Inventory source and scope

Inventory recorded on 2026-10-30 from the local reference's `config/shell/explain.sh`, its 58 `cheatsheets/*.txt` files, and its 142 `gnu-tools.txt` rows. Sources: [explain.sh](https://github.com/winmetta/bootstrap-dev-env/blob/main/config/shell/explain.sh), [cheat sheets](https://github.com/winmetta/bootstrap-dev-env/tree/main/config/shell/cheatsheets), and [GNU catalog](https://github.com/winmetta/bootstrap-dev-env/blob/main/config/shell/gnu-tools.txt). These supply the candidate inventory, not permission to copy explanations. Author original text and consult upstream documentation while implementing each batch.

Include commands exposed by the dispatcher, tool listing, generated GNU helpers, and named lookup aliases. Internal implementation plumbing in `explain.sh` is not an extra content requirement. Dynamic fallback to any executable on a user's PATH is not a finite inventory. Review the inventory against the reference again before closing stage 4 and record any scope changes.

An unchecked command needs an original or expanded page, applicable implementation coverage, and validation. Existing `grep`, `sed`, and `find` pages are starting points and remain unchecked pending this expansion review; preserve and review the existing `cd` page too. Package/topic helpers and lookup aliases have separate dispositions below.

## Per-command completion criteria

- Explain the purpose, prerequisites, common tasks, common mistakes, and portability limits in the canonical Markdown page.
- Cover commonly used options and practical examples: start with 5–10 distinct everyday use cases where meaningful; small commands may need fewer, complex tools may need more. Do not pad trivial commands or claim an exhaustive manual. Explain option arguments and expected behavior.
- Keep each example in a separate fenced block, and use structured option sections and the existing compact output conventions in AGENTS.md.
- Identify applicable implementations separately from OS: GNU, Apple/macOS, BusyBox, and Bash builtins where relevant. Research FreeBSD/OpenBSD or other implementations when applicable; explicitly mark unsupported/not-applicable variants instead of inventing them. BusyBox applet availability/options vary by build. Modern tools such as jq or Git need their actual upstream implementation, not fabricated GNU/BusyBox variants.
- Record upstream manual links, tested OS/implementation/version, minimum-version assumptions, and whether each example was executed or only reviewed. Common examples must hold for all declared variants. Provide variant examples/options for actual differences.
- Exercise safe examples in temporary fixtures on applicable available environments. Clearly document untested or unavailable environments; required supported behavior needs execution evidence before checking the entry off. Never run destructive, privileged, authenticated cloud, or account-changing examples against real resources.
- Verify lookup, search, aliases, option selection, default examples, `--details`, JSON, and web content agree. Run content validation and relevant tests per batch; run the full Go/check/web suites before stage 4 closes.

## Implementation prerequisites

The current schema accepts only `gnu`, `apple`, `busybox`, and `bash` variants, and IDs/aliases must start with a lowercase letter. Before adding affected pages, extend the shared schema, implementation selection, CLI/web filters, and tests for real upstream variants and command names such as `[`. Define a safe canonical filename/ID versus display/lookup spelling if needed. Do not mislabel Git, curl, Python, or other upstream tools just to satisfy the current schema.

Detection currently probes only grep/sed/find. New pages must work with portable/common content or explicit supported selection first. Add bounded, command-specific detection only with tests and an explicit policy; do not probe every new executable with `--version` or infer implementation from OS alone.

## Suggested order

1. Resolve schema/identity prerequisites; review existing grep, sed, find, and cd.
2. Everyday text/files: ls, cat, head, tail, sort, uniq, cut, tr, wc, cp, mv, mkdir, rm, chmod, stat, date, xargs, awk, tar, gzip.
3. Remaining GNU catalog and Unix networking/process tools.
4. Modern search/data tools, languages, developer tools, containers, and cloud CLIs.
5. Audit aliases, package/topic dispositions, applicable variants, and the complete collection; run final validation.

## GNU catalog commands

`coreutils` is both a package topic and an optional multicall executable; document the distinction. Build-specific commands and deprecated wrappers still need an explicit availability/deprecation explanation. A shared cheat sheet does not make distinct executables interchangeable.

### coreutils

- [ ] `[`
- [ ] `arch`
- [ ] `b2sum`
- [ ] `base32`
- [ ] `base64`
- [ ] `basename`
- [ ] `basenc`
- [ ] `cat`
- [ ] `chcon`
- [ ] `chgrp`
- [ ] `chmod`
- [ ] `chown`
- [ ] `chroot`
- [ ] `cksum`
- [ ] `comm`
- [ ] `coreutils`
- [ ] `cp`
- [ ] `csplit`
- [ ] `cut`
- [ ] `date`
- [ ] `dd`
- [ ] `df`
- [ ] `dir`
- [ ] `dircolors`
- [ ] `dirname`
- [ ] `du`
- [ ] `echo`
- [ ] `env`
- [ ] `expand`
- [ ] `expr`
- [ ] `factor`
- [ ] `false`
- [ ] `fmt`
- [ ] `fold`
- [ ] `groups`
- [ ] `head`
- [ ] `hostid`
- [ ] `hostname`
- [ ] `id`
- [ ] `install`
- [ ] `join`
- [ ] `kill`
- [ ] `link`
- [ ] `ln`
- [ ] `logname`
- [ ] `ls`
- [ ] `md5sum`
- [ ] `mkdir`
- [ ] `mkfifo`
- [ ] `mknod`
- [ ] `mktemp`
- [ ] `mv`
- [ ] `nice`
- [ ] `nl`
- [ ] `nohup`
- [ ] `nproc`
- [ ] `numfmt`
- [ ] `od`
- [ ] `paste`
- [ ] `pathchk`
- [ ] `pinky`
- [ ] `pr`
- [ ] `printenv`
- [ ] `printf`
- [ ] `ptx`
- [ ] `pwd`
- [ ] `readlink`
- [ ] `realpath`
- [ ] `rm`
- [ ] `rmdir`
- [ ] `runcon`
- [ ] `seq`
- [ ] `sha1sum`
- [ ] `sha224sum`
- [ ] `sha256sum`
- [ ] `sha384sum`
- [ ] `sha512sum`
- [ ] `shred`
- [ ] `shuf`
- [ ] `sleep`
- [ ] `sort`
- [ ] `split`
- [ ] `stat`
- [ ] `stdbuf`
- [ ] `stty`
- [ ] `sum`
- [ ] `sync`
- [ ] `tac`
- [ ] `tail`
- [ ] `tee`
- [ ] `test`
- [ ] `timeout`
- [ ] `touch`
- [ ] `tr`
- [ ] `true`
- [ ] `truncate`
- [ ] `tsort`
- [ ] `tty`
- [ ] `uname`
- [ ] `unexpand`
- [ ] `uniq`
- [ ] `unlink`
- [ ] `uptime`
- [ ] `users`
- [ ] `vdir`
- [ ] `wc`
- [ ] `who`
- [ ] `whoami`
- [ ] `yes`

### diffutils

- [ ] `cmp`
- [ ] `diff`
- [ ] `diff3`
- [ ] `sdiff`

### findutils

- [ ] `find` — existing page; expand/review coverage
- [ ] `locate`
- [ ] `updatedb`
- [ ] `xargs`

### gawk

- [ ] `awk`
- [ ] `gawk`
- [ ] `gawkbug`

### grep

- [ ] `egrep`
- [ ] `fgrep`
- [ ] `grep` — existing page; expand/review coverage

### gzip

- [ ] `gunzip`
- [ ] `gzexe`
- [ ] `gzip`
- [ ] `uncompress`
- [ ] `zcat`
- [ ] `zcmp`
- [ ] `zdiff`
- [ ] `zegrep`
- [ ] `zfgrep`
- [ ] `zforce`
- [ ] `zgrep`
- [ ] `zless`
- [ ] `zmore`
- [ ] `znew`

### patch

- [ ] `patch`

### sed

- [ ] `sed` — existing page; expand/review coverage

### tar

- [ ] `rmt`
- [ ] `tar`

### wget

- [ ] `wget`

## Additional commands and tools

Keep distinct commands such as `pkill`, `killall`, `zip`, `unzip`, `ssh-agent`, and `limactl` separate even when the reference routes them to another command’s sheet. Choose canonical IDs/aliases only after checking semantics and implementation selection.

- [ ] `actionlint`
- [ ] `aws`
- [ ] `bash`
- [ ] `bat`
- [ ] `claude`
- [ ] `codex`
- [ ] `colima`
- [ ] `curl`
- [ ] `docker`
- [ ] `docker-compose`
- [ ] `fd`
- [ ] `free`
- [ ] `fzf`
- [ ] `gh`
- [ ] `git`
- [ ] `gitleaks`
- [ ] `go`
- [ ] `gzcat`
- [ ] `jq`
- [ ] `just`
- [ ] `killall`
- [ ] `lefthook`
- [ ] `limactl`
- [ ] `lsof`
- [ ] `markdownlint-cli2`
- [ ] `mise`
- [ ] `node`
- [ ] `npm`
- [ ] `npx`
- [ ] `pgrep`
- [ ] `pkill`
- [ ] `prettier`
- [ ] `ps`
- [ ] `pulumi`
- [ ] `python3`
- [ ] `rg`
- [ ] `rsync`
- [ ] `scp`
- [ ] `sd`
- [ ] `security`
- [ ] `shellcheck`
- [ ] `shfmt`
- [ ] `ssh`
- [ ] `ssh-add`
- [ ] `ssh-agent`
- [ ] `ssh-copy-id`
- [ ] `ssh-keygen`
- [ ] `starship`
- [ ] `top`
- [ ] `tree`
- [ ] `unzip`
- [ ] `uv`
- [ ] `uvx`
- [ ] `vim`
- [ ] `zed`
- [ ] `zip`
- [ ] `zoxide`

## Existing additional entry

- [ ] `cd` — retain Bash builtin coverage and review alongside the expanded collection.

## Lookup names to account for

These normalize reference lookup names toward proposed canonical commands; they are not literal dispatch mappings or assertions that executables are interchangeable. For example, the reference sends `lima` and `limactl` to the colima sheet, and `python3` to the python sheet. Decide whether each needs an alias, a separate page, or a documented topic; verify behavior before closing the row.

| Lookup name | Proposed canonical command | Disposition verified |
| --- | --- | --- |
| `compose` | `docker-compose` | Pending |
| `dh` | `df` | Pending |
| `keychain` | `security` | Pending |
| `lima` | `limactl` | Pending |
| `markdownlint` | `markdownlint-cli2` | Pending |
| `python` | `python3` | Pending |
| `ripgrep` | `rg` | Pending |
| `sshadd` | `ssh-add` | Pending |
| `sshcopyid` | `ssh-copy-id` | Pending |
| `sshkeygen` | `ssh-keygen` | Pending |
| `z` | `zoxide` | Pending |

`python3` and `python` can resolve to different interpreters; record version/lookup limits. `z` is shell integration supplied by zoxide, not the standalone executable. `gzcat`, `gunzip`, and `zcat` need implementation-aware treatment even if they share gzip material.

## Package, shell, and bootstrap topics

Record a disposition for each; do not manufacture executable pages for bootstrap-specific helpers.

- [ ] `findutils`, `diffutils`, and the `coreutils` package overview: explain package membership alongside actual command pages.
- [ ] `aliases`, `git-aliases`, `git-config`: preserve reusable alias/configuration guidance in appropriate command pages; machine-local introspection belongs to optional shell integration or the bootstrap.
- [ ] `command`, `path`: cover reusable command resolution/PATH behavior and relevant builtins; live parent-shell inspection remains outside the standalone binary.
- [ ] `completion`: reusable Bash/zsh completion guidance or an explicit future-topic disposition.
- [ ] `tools`: the bootstrap-specific installed-tool listing is replaced by the standalone catalog/list; no machine-installation claim.
- [ ] `iTerm`: GUI app, not a command spelled `iTerm`; record whether generic terminal guidance or the actual `open` command deserves future coverage. No bootstrap profile/configuration migration.

## Stage 4 exit review

- [ ] Every inventory command is covered and every lookup/topic has an explicit verified disposition; no unchecked required commands remain.
- [ ] Common examples/options and applicable implementation variants are reviewed, with upstream references and truthful execution evidence.
- [ ] Schema/selection changes needed for the collection are documented and tested in both CLI and web.
- [ ] Full content/Go/static-site/parity checks pass; no known blocking accuracy or rendering regressions remain.
- [ ] Update README, development status, and roadmap to mark stage 4 complete. Publication, licensing, installers, and deployment remain stage 5.
