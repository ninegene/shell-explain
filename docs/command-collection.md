# Stage 4 command collection checklist

Stage 4's next work is to expand the existing four-command collection. The current CLI and web usability are accepted as sufficient; this checklist tracks planned content, not completed explanations.

## Inventory source and scope

Inventory recorded on 2026-10-30 from the local reference's `config/shell/explain.sh`, its 58 `cheatsheets/*.txt` files, and its 142 `gnu-tools.txt` rows. Sources: [explain.sh](https://github.com/winmetta/bootstrap-dev-env/blob/main/config/shell/explain.sh), [cheat sheets](https://github.com/winmetta/bootstrap-dev-env/tree/main/config/shell/cheatsheets), and [GNU catalog](https://github.com/winmetta/bootstrap-dev-env/blob/main/config/shell/gnu-tools.txt). These supply the candidate inventory, not permission to copy explanations. Author original text and consult upstream documentation while implementing each batch.

Include commands exposed by the dispatcher, tool listing, generated GNU helpers, and named lookup aliases. Internal implementation plumbing in `explain.sh` is not an extra content requirement. Dynamic fallback to any executable on a user's PATH is not a finite inventory. Review the inventory against the reference again before closing stage 4 and record any scope changes.

An unchecked command needs an original or expanded page, applicable implementation coverage, and validation. The existing `grep`, `sed`, `find`, and `cd` pages were expanded and reviewed first (checked below). Package/topic helpers and lookup aliases have separate dispositions below.

## Per-command completion criteria

- Explain the purpose, prerequisites, common tasks, common mistakes, and portability limits in the canonical Markdown page, with a `Background` section and an `About` section for each implementation view (required; see [content format](content-format.md)).
- Cover commonly used options and practical examples: start with 5–10 distinct everyday use cases where meaningful; small commands may need fewer, complex tools may need more. Do not pad trivial commands or claim an exhaustive manual. Explain option arguments and expected behavior.
- Keep each example in a separate fenced block, and use structured option sections and the existing compact output conventions in AGENTS.md.
- Identify applicable implementations separately from OS: GNU, Apple/macOS, BusyBox, and Bash builtins where relevant. Research FreeBSD/OpenBSD or other implementations when applicable; explicitly mark unsupported/not-applicable variants instead of inventing them. BusyBox applet availability/options vary by build. Modern tools such as jq or Git need their actual upstream implementation, not fabricated GNU/BusyBox variants.
- Record upstream manual links, tested OS/implementation/version, minimum-version assumptions, and whether each example was executed or only reviewed. Common examples must hold for all declared variants. Provide variant examples/options for actual differences.
- Exercise safe examples in temporary fixtures on applicable available environments. Clearly document untested or unavailable environments; required supported behavior needs execution evidence before checking the entry off. Never run destructive, privileged, authenticated cloud, or account-changing examples against real resources.
- Verify lookup, search, aliases, option selection, default examples, `--details`, JSON, and web content agree. Run content validation and relevant tests per batch; run the full Go/check/web suites before stage 4 closes.

## Implementation prerequisites

Done in the stage 4 batch 1 schema change, with Go, CLI, detection, and web tests:

- Implementations come from one registry (`catalog.Implementations`): gnu, apple, busybox, bash, zsh, and upstream. Git, curl, jq, and similar tools use `upstream` instead of a mislabeled GNU/BusyBox variant.
- Symbolic lookup names such as `[` are aliases of a safe canonical ID (`test`), so filenames, web routes, and JSON IDs keep the lowercase-letter rule.
- Shell-only entries and single-implementation `upstream` entries are selected without PATH probing.

Still open: detection probes only the commands in `detect.probed` (grep, sed, find, cat, head, tail, wc, sort, uniq, cut, tr, ls, mkdir, cp, mv, rm, chmod, xargs, gzip, tar, stat, date; see [detection](detection.md)). New pages must work with portable/common content or explicit supported selection first. Add bounded, command-specific detection only with tests and an explicit policy; do not probe every new executable with `--version` or infer implementation from OS alone.

## Completed-page evidence

Reviewed and executed on 2026-11-08 (macOS Homebrew coreutils 9.11 was exercised through its `g`-prefixed names). Every authored example has a recorded expected result in `internal/examplecheck` (a test enforces this), and the same fixture and expectations run in both runners below. Evidence covers the listed builds, not every version or BusyBox configuration.

| Page | Examples | Apple/macOS 26.6.2 | GNU (Homebrew) | BusyBox 1.37.0 / GNU (Alpine 3.24.3, QEMU) | Bash |
| --- | --- | --- | --- | --- | --- |
| `grep` | 11 common, 1 GNU | system grep 2.6.0-FreeBSD, executed | grep 3.12, executed | BusyBox 1.37.0 and GNU 3.12, executed | n/a |
| `sed` | 10 common, backup edits per variant | system sed, executed | sed 4.10, executed | BusyBox 1.37.0 and GNU 4.9, executed | n/a |
| `find` | 8 common, 1 GNU | system find, executed | findutils 4.10.0, executed | BusyBox 1.37.0 and GNU 4.10.0, executed | n/a |
| `cat` | 5 common | system `/bin/cat`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `head` | 5 common, 1 GNU | system head, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `tail` | 5 common run, 1 reviewed only (`tail -f` never ends) | system tail, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `wc` | 6 common, 1 GNU | system wc, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sort` | 6 common | system sort (2.3-Apple), executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `uniq` | 6 common | system uniq, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `cut` | 6 common, 1 GNU | system cut, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `tr` | 7 common | system tr, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `ls` | 8 common, 1 GNU | system `/bin/ls`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `mkdir` | 5 common, 1 GNU | system mkdir, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `cp` | 5 common, 1 GNU | system cp, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `mv` | 5 common, 1 GNU | system mv, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `rm` | 4 common, 1 GNU | system rm, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `chmod` | 6 common, 1 GNU | system `/bin/chmod`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `xargs` | 6 common, 1 GNU | system xargs, executed | findutils 4.10.0, executed | BusyBox 1.37.0 and findutils 4.10.0, executed | n/a |
| `gzip` | 5 common, 1 GNU | system gzip (Apple gzip 479), executed | gzip 1.14, executed (with Homebrew gzip first on PATH, because a version-manager wrapper script is never probed) | BusyBox 1.37.0 and Alpine GNU gzip, executed | n/a |
| `tar` | 6 common, 1 GNU | system bsdtar 3.5.3, executed | GNU tar, executed | BusyBox 1.37.0 and Alpine GNU tar, executed | n/a |
| `stat` | 2 common; 3 GNU, 3 Apple, 2 BusyBox | system `/usr/bin/stat`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `date` | 5 common; 2 GNU, 2 Apple, 1 BusyBox | system `/bin/date`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `cd` | 5 common, 2 Bash | n/a | n/a | n/a | Bash 5.2.37, executed on macOS; zsh/dash not covered |

Expected output is compared exactly, except directory-traversal order, which is compared as a set of lines, and column padding (`wc`, `uniq -c`, `cat -n`), which is compared with blanks collapsed. Examples that never end are marked reviewed only. Current-time examples (`date`) check only the shape of the output, and `stat` has no portable report format, so its per-implementation examples check the chosen fields. `cd` examples check success only, since the resulting directory is a temporary path. Checks: Go tests for the macOS/Homebrew rows and Bash (`just test`), `python3 -B scripts/linux-smoke.py` for the Linux rows, then `just check` and `just web-test`. Upstream manuals are listed in each page's references.

## Suggested order

1. Resolve schema/identity prerequisites; review existing grep, sed, find, and cd.
2. Everyday text/files: ls, cat, head, tail, sort, uniq, cut, tr, wc, cp, mv, mkdir, rm, chmod, stat, date, xargs, awk, tar, gzip.
3. Remaining GNU catalog and Unix networking/process tools.
4. Modern search/data tools, languages, developer tools, containers, and cloud CLIs.
5. Audit aliases, package/topic dispositions, applicable variants, and the complete collection; run final validation.

## GNU catalog commands

`coreutils` is both a package topic and an optional multicall executable; document the distinction. Build-specific commands and deprecated wrappers still need an explicit availability/deprecation explanation. A shared cheat sheet does not make distinct executables interchangeable.

### coreutils

- [ ] `[` (alias of `test`; add the `test` page)
- [ ] `arch`
- [ ] `b2sum`
- [ ] `base32`
- [ ] `base64`
- [ ] `basename`
- [ ] `basenc`
- [x] `cat` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `chcon`
- [ ] `chgrp`
- [x] `chmod` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)
- [ ] `chown`
- [ ] `chroot`
- [ ] `cksum`
- [ ] `comm`
- [ ] `coreutils`
- [x] `cp` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [ ] `csplit`
- [x] `cut` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `date` (stage 4 batch 2d, 2026-11-08; see completed-page evidence)
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
- [x] `head` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `hostid`
- [ ] `hostname`
- [ ] `id`
- [ ] `install`
- [ ] `join`
- [ ] `kill`
- [ ] `link`
- [ ] `ln`
- [ ] `logname`
- [x] `ls` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [ ] `md5sum`
- [x] `mkdir` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [ ] `mkfifo`
- [ ] `mknod`
- [ ] `mktemp`
- [x] `mv` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
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
- [x] `rm` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
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
- [x] `sort` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `split`
- [x] `stat` (stage 4 batch 2d, 2026-11-08; see completed-page evidence)
- [ ] `stdbuf`
- [ ] `stty`
- [ ] `sum`
- [ ] `sync`
- [ ] `tac`
- [x] `tail` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `tee`
- [ ] `test`
- [ ] `timeout`
- [ ] `touch`
- [x] `tr` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `true`
- [ ] `truncate`
- [ ] `tsort`
- [ ] `tty`
- [ ] `uname`
- [ ] `unexpand`
- [x] `uniq` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `unlink`
- [ ] `uptime`
- [ ] `users`
- [ ] `vdir`
- [x] `wc` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [ ] `who`
- [ ] `whoami`
- [ ] `yes`

### diffutils

- [ ] `cmp`
- [ ] `diff`
- [ ] `diff3`
- [ ] `sdiff`

### findutils

- [x] `find` — expanded and reviewed 2026-11-08 (see completed-page evidence)
- [ ] `locate`
- [ ] `updatedb`
- [x] `xargs` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)

### gawk

- [ ] `awk`
- [ ] `gawk`
- [ ] `gawkbug`

### grep

- [ ] `egrep`
- [ ] `fgrep`
- [x] `grep` — expanded and reviewed 2026-11-08 (see completed-page evidence)

### gzip

- [ ] `gunzip`
- [ ] `gzexe`
- [x] `gzip` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)
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

- [x] `sed` — expanded and reviewed 2026-11-08 (see completed-page evidence)

### tar

- [ ] `rmt`
- [x] `tar` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)

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

- [x] `cd` — Bash builtin coverage reviewed and expanded 2026-11-08 (see completed-page evidence)

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
