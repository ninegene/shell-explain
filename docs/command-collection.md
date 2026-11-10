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

- Implementations come from one registry (`catalog.Implementations`): gnu, apple, busybox, mawk, procps, bash, zsh, and upstream. The `procps` row exists because the process tools of most Linux systems (ps, pgrep, pkill, top, free) come from procps-ng, which is neither GNU coreutils nor BusyBox. The `mawk` row exists because the default awk on Debian-based Linux is mawk, not GNU awk; gawk is covered by `gnu` and macOS awk by `apple`. Git, curl, jq, and similar tools use `upstream` instead of a mislabeled GNU/BusyBox variant.
- Symbolic lookup names such as `[` are aliases of a safe canonical ID (`test`), so filenames, web routes, and JSON IDs keep the lowercase-letter rule.
- Shell-only entries and single-implementation `upstream` entries are selected without PATH probing.

Still open: detection probes only the commands in `detect.probed` (grep, sed, find, cat, head, tail, wc, sort, uniq, cut, tr, ls, mkdir, cp, mv, rm, chmod, xargs, gzip, tar, stat, date, awk; see [detection](detection.md)). New pages must work with portable/common content or explicit supported selection first. Add bounded, command-specific detection only with tests and an explicit policy; do not probe every new executable with `--version` or infer implementation from OS alone.

## Completed-page evidence

Reviewed and executed on 2026-11-08 (touch, ln, du, basename, dirname and tee on 2026-11-09; rmdir, mktemp, readlink, realpath, seq, paste, comm, join, nl and fold on 2026-11-10; env, printenv, id, whoami, uname, hostname, df, yes, sleep, expand and unexpand on 2026-11-10; the checksum and encoding tools of batch 3d on 2026-11-10; tac, shuf, od, split, csplit, fmt, factor, expr, nproc, timeout and truncate on 2026-11-10; dd, install, link, unlink, mkfifo, nohup, nice, sync, sum, tsort, pr and pathchk on 2026-11-10; groups, logname, tty, users, who, uptime, hostid, pinky, chown, chgrp, shred, dir, vdir, dircolors and mknod on 2026-11-10; echo, printf, pwd, true, false, test, kill and arch on 2026-11-10; numfmt, ptx, stdbuf, stty, chroot, chcon, runcon and coreutils on 2026-11-10; diff, cmp, diff3, sdiff, patch and zcat on 2026-11-10; the gzip family of batch 3k on 2026-11-10; the upstream tools of batch 3l on 2026-11-10) (macOS Homebrew coreutils 9.11 was exercised through its `g`-prefixed names). Every authored example has a recorded expected result in `internal/examplecheck` (a test enforces this), and the same fixture and expectations run in both runners below. Evidence covers the listed builds, not every version or BusyBox configuration.

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
| `awk` | 8 common; 2 GNU, 1 BusyBox | system `/usr/bin/awk` (awk version 20200816), executed | gawk 5.4.0, executed (with Homebrew gawk first on PATH, because a version-manager wrapper script is never probed) | BusyBox 1.37.0, gawk 5.3.2 and mawk 1.3.4, executed | n/a |
| `touch` | 6 common, 1 GNU | system `/usr/bin/touch`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `ln` | 7 common, 1 GNU | system `/bin/ln`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `du` | 6 common, 1 GNU (sizes checked by shape only) | system `/usr/bin/du`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `basename` | 4 common, 1 GNU | system `/usr/bin/basename`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `dirname` | 5 common, 1 GNU | system `/usr/bin/dirname`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `tee` | 5 common, 1 GNU | system `/usr/bin/tee`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `rmdir` | 4 common, 1 GNU | system `/bin/rmdir`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `mktemp` | 4 common, 1 GNU | system `/usr/bin/mktemp`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `readlink` | 4 common, 1 GNU | system `/usr/bin/readlink`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `realpath` | 4 common, 2 GNU | system `/bin/realpath`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `seq` | 6 common, 1 GNU | system `/usr/bin/seq`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `paste` | 6 common | system `/usr/bin/paste`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `comm` | 4 common, 1 GNU | system `/usr/bin/comm`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `join` | 5 common, 1 GNU | system `/usr/bin/join`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list` in the guest) | n/a |
| `nl` | 5 common, 1 GNU | system `/usr/bin/nl`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `fold` | 4 common, 1 GNU | system `/usr/bin/fold`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `env` | 4 common, 1 GNU (machine-dependent output checked by shape) | system `/usr/bin/env`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `printenv` | 4 common (machine-dependent output checked by shape) | system `/usr/bin/printenv`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `id` | 5 common (machine-dependent output checked by shape) | system `/usr/bin/id`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `whoami` | 3 common (machine-dependent output checked by shape) | system `/usr/bin/whoami`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `uname` | 5 common (machine-dependent output checked by shape) | system `/usr/bin/uname`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `hostname` | 3 common (no GNU view: not in coreutils) | system `/bin/hostname`, executed | n/a | BusyBox 1.37.0, executed | n/a |
| `df` | 5 common (machine-dependent output checked by shape) | system `/bin/df`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `yes` | 4 common, 1 GNU | system `/usr/bin/yes`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sleep` | 4 common, 1 GNU | system `/bin/sleep`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `expand` | 4 common, 1 GNU | system `/usr/bin/expand`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `unexpand` | 4 common | system `/usr/bin/unexpand`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `md5sum` | 5 common | system `/sbin/md5sum`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sha1sum` | 5 common | system `/sbin/sha1sum`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sha256sum` | 5 common | system `/sbin/sha256sum`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sha512sum` | 5 common | system `/sbin/sha512sum`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sha224sum` | 5 common | system `/sbin/sha224sum`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `sha384sum` | 5 common | system `/sbin/sha384sum`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `cksum` | 4 common, 1 GNU | system `/usr/bin/cksum`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `base64` | 4 common, 1 GNU, 1 Apple | system `/usr/bin/base64`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `b2sum` | 3 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `base32` | 3 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `basenc` | 4 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `tac` | 4 common | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `shuf` | 4 common, 2 GNU (random output checked by shape) | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `od` | 5 common, 1 GNU (column spacing ignored) | system `/usr/bin/od`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `split` | 5 common, 1 GNU | system `/usr/bin/split`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `csplit` | 3 common, 1 GNU | system `/usr/bin/csplit`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `fmt` | 3 common | system `/usr/bin/fmt`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `factor` | 4 common, 1 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `expr` | 6 common, 1 GNU | system `/bin/expr`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `nproc` | 2 common, 1 GNU (output checked by shape) | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `timeout` | 4 common, 1 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `truncate` | 4 common, 1 GNU | system `/usr/bin/truncate`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `dd` | 4 common, 1 GNU, 1 Apple | system `/bin/dd`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `install` | 4 common, 1 GNU | system `/usr/bin/install`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `link` | 3 common | system `/bin/link`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `unlink` | 3 common | system `/bin/unlink`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `mkfifo` | 3 common | system `/usr/bin/mkfifo`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `nohup` | 3 common | system `/usr/bin/nohup`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `nice` | 3 common, 1 GNU | system `/usr/bin/nice`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sync` | 3 common, 1 GNU | system `/bin/sync`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `sum` | 3 common (BusyBox prints no file name) | system `/usr/bin/sum`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `tsort` | 3 common | system `/usr/bin/tsort`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `pr` | 3 common (layout compared with blanks collapsed) | system `/usr/bin/pr`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `pathchk` | 3 common | system `/usr/bin/pathchk`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `groups` | 2 common (shape only) | system `/usr/bin/groups`, executed | coreutils 9.11 on macOS only | BusyBox 1.37.0, executed; GNU ran on macOS with Homebrew coreutils only (not installed in the Alpine guest) | n/a |
| `logname` | 2 common (shape only) | system `/usr/bin/logname`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `tty` | 3 common, 1 GNU, 1 Apple | system `/usr/bin/tty`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `users` | 2 common (success only) | system `/usr/bin/users`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `who` | 2 common, 1 GNU, 1 Apple (success only) | system `/usr/bin/who`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `uptime` | 2 common (shape only) | system `/usr/bin/uptime`, executed | coreutils 9.11 on macOS only | BusyBox 1.37.0, executed; GNU ran on macOS with Homebrew coreutils only (not installed in the Alpine guest) | n/a |
| `hostid` | 1 common (shape only) | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `pinky` | 2 GNU (success only) | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `chown` | 4 common (own IDs only) | system `/usr/sbin/chown`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `chgrp` | 3 common (own group only) | system `/usr/bin/chgrp`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `shred` | 3 common | not shipped by macOS (checked on the host) | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `dir` | 2 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `vdir` | 2 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `dircolors` | 3 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `mknod` | 2 common | system `/sbin/mknod`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | n/a |
| `echo` | 6 common, 2 Bash, 1 GNU, 1 Apple | system `/bin/echo`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | Bash 5.2.37 builtin, executed on macOS |
| `printf` | 6 common, 2 Bash | system `/usr/bin/printf`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | Bash 5.2.37 builtin, executed on macOS |
| `pwd` | 4 common | system `/bin/pwd`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | Bash 5.2.37 builtin, executed on macOS |
| `true` | 5 common | system `/usr/bin/true`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | Bash 5.2.37 builtin, executed on macOS |
| `false` | 4 common | system `/usr/bin/false`, executed | coreutils 9.11, executed | BusyBox 1.37.0 and coreutils 9.11, executed | Bash 5.2.37 builtin, executed on macOS |
| `test` | 9 common | system `/bin/test`, executed | coreutils 9.11, executed in the Linux guest only (GNU `test --version` prints no banner, so it is not auto-identified) | BusyBox 1.37.0 and coreutils 9.11, executed | Bash 5.2.37 builtin, executed on macOS |
| `kill` | 4 common | system `/bin/kill`, executed | no GNU view (the coreutils builds used do not include `kill`) | BusyBox 1.37.0, executed | Bash 5.2.37 builtin, executed on macOS |
| `arch` | 2 common (shape only) | system `/usr/bin/arch`, executed | no GNU view (not installed by the coreutils builds used) | BusyBox 1.37.0, executed | n/a |
| `numfmt` | 7 GNU | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `ptx` | 2 GNU (layout compared with blanks collapsed) | not shipped by macOS (checked on the host) | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `stdbuf` | 3 common | system `/usr/bin/stdbuf`, executed | coreutils 9.11, executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `stty` | 4 common, reviewed only (needs a terminal) | system `/bin/stty`, not executed | not executed | BusyBox 1.37.0, not executed | n/a |
| `chroot` | 3 common, reviewed only (needs administrator rights) | system `/usr/sbin/chroot`, not executed | not executed | not executed | n/a |
| `chcon` | 3 GNU, reviewed only (needs SELinux) | not shipped by macOS (checked on the host) | not executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `runcon` | 3 GNU, reviewed only (needs SELinux) | not shipped by macOS (checked on the host) | not executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `coreutils` | 3 GNU (multi-call program) | not shipped by macOS (checked on the host) | coreutils 9.11, executed in the Linux guest only (Homebrew does not install it) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `diff` | 6 common, 1 GNU, 1 Apple | system `/usr/bin/diff` (Apple diff, based on FreeBSD diff), executed | GNU 3.12 from a conda install managed by mise on the host, executed, and Alpine diffutils 3.12 in the guest | BusyBox 1.37.0, executed | n/a |
| `cmp` | 4 common, 1 GNU, 1 Apple | system `/usr/bin/cmp`, executed | GNU 3.12 from a conda install managed by mise on the host, executed, and Alpine diffutils 3.12 in the guest | BusyBox 1.37.0, executed | n/a |
| `diff3` | 2 common | system `/usr/bin/diff3`, executed | GNU 3.12 from a conda install managed by mise on the host, executed, and Alpine diffutils 3.12 in the guest (with the matching `diff` put on the path) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `sdiff` | 3 common (layout compared with blanks collapsed) | system `/usr/bin/sdiff`, executed | GNU 3.12 from a conda install managed by mise on the host, executed, and Alpine diffutils 3.12 in the guest | no BusyBox applet (checked with `busybox --list`) | n/a |
| `patch` | 4 common | system `/usr/bin/patch` (2.0-12u11-Apple), executed | GNU patch 2.8, executed on the host (mise conda) and in the guest (Alpine patch) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `zcat` | 3 common, 1 GNU, 1 Apple | system `/usr/bin/zcat` and `gzcat`, executed | gzip 1.14 on the host and Alpine gzip in the guest (a shell script, so not identified automatically) | BusyBox 1.37.0, executed | n/a |
| `zgrep` | 5 common | system `/usr/bin/zgrep`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `zless` | 1 common executed, 1 reviewed only (needs a terminal) | system `/usr/bin/zless`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `zmore` | 1 common executed, 1 reviewed only (needs a terminal) | system `/usr/bin/zmore`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `zcmp` | 3 common | system `/usr/bin/zcmp`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `zdiff` | 2 common | system `/usr/bin/zdiff`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `zforce` | 2 common | system `/usr/bin/zforce`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `gzexe` | 2 common | system `/usr/bin/gzexe`, executed | GNU gzip scripts, executed in the Linux guest only (the probe does not run scripts, so the host cannot select them) | no BusyBox applet (checked with `busybox --list`) | n/a |
| `znew` | 2 common, reviewed only (needs `compress`) | system `/usr/bin/znew`, not executed | not executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `uncompress` | 2 common, reviewed only (needs `compress`) | system `/usr/bin/uncompress`, not executed | not executed | no BusyBox applet (checked with `busybox --list`) | n/a |
| `ps` | 4 common, 1 Apple, 1 procps | system `/bin/ps`, executed | procps-ng 4.0.6 in the Linux guest (`procps` view) | BusyBox 1.37.0, executed | n/a |
| `pgrep` | 3 common, 1 Apple, 1 procps | system `/usr/bin/pgrep`, executed | procps-ng 4.0.6 in the guest | BusyBox 1.37.0, executed | n/a |
| `pkill` | 3 common | system `/usr/bin/pkill`, executed | procps-ng 4.0.6 in the guest | BusyBox 1.37.0, executed | n/a |
| `free` | 3 common | not shipped by macOS | procps-ng 4.0.6 in the guest | BusyBox 1.37.0, executed | n/a |
| `top` | 1 Apple, 1 procps; the interactive example is reviewed only | system `/usr/bin/top`, executed in batch mode | procps-ng 4.0.6 in the guest | not run (reviewed only) | n/a |
| `killall` | 4 common, reviewed only (a name such as `sleep` would signal every process of that name) | system `/usr/bin/killall`, not executed | not covered (psmisc is another program) | not executed | n/a |
| `cd` | 5 common, 2 Bash | n/a | n/a | n/a | Bash 5.2.37, executed on macOS; zsh/dash not covered |

Expected output is compared exactly, except directory-traversal order, which is compared as a set of lines, and column padding (`wc`, `uniq -c`, `cat -n`), which is compared with blanks collapsed. Examples that never end are marked reviewed only. Current-time examples (`date`) check only the shape of the output, and `stat` has no portable report format, so its per-implementation examples check the chosen fields. `cd` examples check success only, since the resulting directory is a temporary path. Checks: Go tests for the macOS/Homebrew rows and Bash (`just test`), `python3 -B scripts/linux-smoke.py` for the Linux rows, then `just check` and `just web-test`. Upstream manuals are listed in each page's references.

### Upstream tools

Single-project tools are documented with an `upstream` view, so there are no GNU, Apple or BusyBox columns. Their examples run on the tool installed on the host, in an empty temporary folder with a clean home folder, and are skipped when the tool is not installed. Examples that need the network or a signed-in account are reviewed by hand and never run by the tests. They are marked "reviewed only" below, and each was checked once against the live service on the date shown.

| Page | Examples | Version tested | Notes |
| --- | --- | --- | --- |
| `curl` | 2 executed (`file://` URLs), 7 reviewed only (GitHub API) | curl 8.7.1 (macOS system build) | network examples checked by hand on 2026-11-10 against `api.github.com` and `raw.githubusercontent.com` |
| `jq` | 8 executed, 1 reviewed only (GitHub API) | jq 1.8.2 | the network example checked by hand on 2026-11-10 |
| `git` | 8 executed | git 2.56.0 | each example makes its own repository; names and addresses are given with `-c`, so no setup is needed |
| `rg` | 8 executed | ripgrep 15.2.0 | the order of results is not stable, so one example sorts and one is compared as a set |
| `fd` | 5 executed | fd 10.5.0 | |
| `bat` | 4 executed | bat 0.26.1 | numbers are compared with blanks collapsed |
| `tree` | 4 executed | tree v2.3.2 | the version tested counts the starting folder in the directory total |
| `sd` | 4 executed | sd v1.1.0 | |
| `shellcheck` | 4 executed | ShellCheck 0.11.0 | |
| `shfmt` | 4 executed | shfmt v3.14.1 | `-l` exits with status 1 when it lists a file, in this version |
| `prettier` | 4 executed | Prettier 3.9.10 | runs on the host's Node.js |
| `node` | 5 executed | Node.js long-term-support release on the host | |
| `npm` | 3 executed, 2 reviewed only (they download packages) | the npm that comes with that Node.js | |
| `npx` | 1 executed, 2 reviewed only (they download packages) | the npx that comes with that Node.js | |
| `python3` | 4 executed, 1 reviewed only (it listens on a port) | Python 3 on the host | |
| `uv` | 3 executed, 2 reviewed only (they may download) | uv on the host | |
| `uvx` | 1 executed, 2 reviewed only (they download tools) | uvx of that uv | |
| `go` | 5 executed | Go 1.27.3 | no network: each example builds in a temporary module |
| `just` | 4 executed | just on the host | |
| `bash` | 6 executed | GNU Bash 5.2.37 | macOS ships Bash 3.2, which the page warns about but the examples do not run on |
| `zip` | 4 executed | Info-ZIP Zip 3.0 (macOS build) | |
| `unzip` | 5 executed | UnZip 6.00 (macOS build) | |
| `lsof` | 3 executed, 2 reviewed only (they depend on what is running) | lsof 4.91 (macOS build) | |
| `rsync` | 4 executed, 1 reviewed only (needs a remote host) | openrsync (the macOS `/usr/bin/rsync`) | macOS ships openrsync, not the rsync project's program |
| `ssh-keygen` | 4 executed | OpenSSH of the macOS host | throw-away keys in a temporary folder |
| `actionlint` | 3 executed | actionlint 1.7.12 | |
| `markdownlint-cli2` | 3 executed | markdownlint-cli2 0.24.0 | runs on the host's Node.js |
| `gitleaks` | 3 executed | gitleaks 8.30.1 | the made-up token is written in two pieces; gitleaks ignores well-known example keys |
| `lefthook` | 3 executed | lefthook 2.2.2 | |
| `mise` | 3 executed, 2 reviewed only (they install tools) | mise 2026.10.19 | |
| `zoxide` | 3 executed | zoxide 0.10.0 | |
| `starship` | 2 executed | starship 1.27.0 | |
| `ssh` | 4 executed (settings printed with `-G`, no connection), 4 reviewed only (they connect) | OpenSSH of the macOS host | |
| `scp` | 1 executed (local copy), 4 reviewed only (they connect) | OpenSSH of the macOS host | |
| `ssh-add` | 4 executed | OpenSSH of the macOS host | each example runs a private agent for one command on a short `/tmp` socket path and removes the socket file |
| `ssh-agent` | 2 executed, 2 reviewed only (they leave a program running) | OpenSSH of the macOS host | |
| `ssh-copy-id` | 1 executed (usage text), 3 reviewed only (they connect) | ssh-copy-id of the macOS host | |
| `fzf` | 4 executed, 1 reviewed only (needs a terminal) | fzf 0.74.5 | |
| `vim` | 4 executed, 1 reviewed only (opens the editor) | vim on the host | |
| `zed` | 1 executed (version), 3 reviewed only (they open a window) | Zed 1.26.2 | |
| `security` | 5 reviewed only | macOS system tool | none was run, because the examples touch keychains |
| `docker` | 3 executed (client version, help, contexts), 4 reviewed only (they need a daemon) | Docker 29.9.2 (client, no daemon) | |
| `docker-compose` | 4 executed (version and compose-file checks), 2 reviewed only (they start containers) | Docker Compose v5.6.1 | |
| `colima` | 2 executed (version, list), 3 reviewed only (they start a virtual machine) | colima v0.10.3 | `limactl` is put on the example path, because colima starts it |
| `limactl` | 2 executed (version, help), 4 reviewed only | limactl 2.2.1 | |
| `aws` | 4 executed (version, settings, help), 3 reviewed only (they call AWS) | aws-cli 2.38 | clean home folder, no credentials, never run against a real account |
| `gh` | 4 executed (version, help, login state), 3 reviewed only (they call GitHub as the user) | gh 2.103.2 | clean home folder, so no login exists |
| `pulumi` | 3 executed (version, help), 4 reviewed only (they plan or change infrastructure) | Pulumi v3.272.0 | `up`, `destroy` and `refresh` are never run by this project |
| `claude` | 2 executed (version, help), 3 reviewed only (they send text to a service) | Claude Code 2.1.310 | |
| `codex` | 2 executed (version, help), 2 reviewed only | codex-cli 0.164.0 | its warning about helper programs under a temporary folder is hidden |

## Suggested order

1. Resolve schema/identity prerequisites; review existing grep, sed, find, and cd.
2. Everyday text/files: ls, cat, head, tail, sort, uniq, cut, tr, wc, cp, mv, mkdir, rm, chmod, stat, date, xargs, awk, tar, gzip.
3. Remaining GNU catalog and Unix networking/process tools.
4. Modern search/data tools, languages, developer tools, containers, and cloud CLIs.
5. Audit aliases, package/topic dispositions, applicable variants, and the complete collection; run final validation.

## GNU catalog commands

`coreutils` is both a package topic and an optional multicall executable; document the distinction. Build-specific commands and deprecated wrappers still need an explicit availability/deprecation explanation. A shared cheat sheet does not make distinct executables interchangeable.

### coreutils

- [x] `[` (an alias of the `test` page, 2026-11-10; its examples use `test`, and the `[ ]` form ran in the shell builtin)
- [x] `arch` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `b2sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `base32` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `base64` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `basename` (stage 4 batch 3a, 2026-11-09; see completed-page evidence)
- [x] `basenc` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `cat` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `chcon` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `chgrp` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `chmod` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)
- [x] `chown` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `chroot` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `cksum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `comm` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `coreutils` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `cp` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [x] `csplit` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `cut` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `date` (stage 4 batch 2d, 2026-11-08; see completed-page evidence)
- [x] `dd` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `df` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `dir` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `dircolors` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `dirname` (stage 4 batch 3a, 2026-11-09; see completed-page evidence)
- [x] `du` (stage 4 batch 3a, 2026-11-09; see completed-page evidence)
- [x] `echo` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `env` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `expand` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `expr` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `factor` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `false` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `fmt` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `fold` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `groups` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `head` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `hostid` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `hostname` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `id` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `install` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `join` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `kill` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `link` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `ln` (stage 4 batch 3a, 2026-11-09; see completed-page evidence)
- [x] `logname` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `ls` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [x] `md5sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `mkdir` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [x] `mkfifo` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `mknod` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `mktemp` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `mv` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [x] `nice` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `nl` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `nohup` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `nproc` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `numfmt` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `od` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `paste` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `pathchk` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `pinky` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `pr` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `printenv` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `printf` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `ptx` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `pwd` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `readlink` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `realpath` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `rm` (stage 4 batch 2b, 2026-11-08; see completed-page evidence)
- [x] `rmdir` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `runcon` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `seq` (stage 4 batch 3b, 2026-11-10; see completed-page evidence)
- [x] `sha1sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `sha224sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `sha256sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `sha384sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `sha512sum` (stage 4 batch 3d, 2026-11-10; see completed-page evidence)
- [x] `shred` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `shuf` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `sleep` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `sort` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `split` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `stat` (stage 4 batch 2d, 2026-11-08; see completed-page evidence)
- [x] `stdbuf` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `stty` (stage 4 batch 3i, 2026-11-10; see completed-page evidence)
- [x] `sum` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `sync` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `tac` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `tail` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `tee` (stage 4 batch 3a, 2026-11-09; see completed-page evidence)
- [x] `test` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `timeout` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `touch` (stage 4 batch 3a, 2026-11-09; see completed-page evidence)
- [x] `tr` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `true` (stage 4 batch 3h, 2026-11-10; see completed-page evidence)
- [x] `truncate` (stage 4 batch 3e, 2026-11-10; see completed-page evidence)
- [x] `tsort` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `tty` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `uname` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `unexpand` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `uniq` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `unlink` (stage 4 batch 3f, 2026-11-10; see completed-page evidence)
- [x] `uptime` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `users` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `vdir` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `wc` (stage 4 batch 2a, 2026-11-08; see completed-page evidence)
- [x] `who` (stage 4 batch 3g, 2026-11-10; see completed-page evidence)
- [x] `whoami` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)
- [x] `yes` (stage 4 batch 3c, 2026-11-10; see completed-page evidence)

### diffutils

- [x] `cmp` (stage 4 batch 3j, 2026-11-10; see completed-page evidence)
- [x] `diff` (stage 4 batch 3j, 2026-11-10; see completed-page evidence)
- [x] `diff3` (stage 4 batch 3j, 2026-11-10; see completed-page evidence)
- [x] `sdiff` (stage 4 batch 3j, 2026-11-10; see completed-page evidence)

### findutils

- [x] `find` — expanded and reviewed 2026-11-08 (see completed-page evidence)
- [ ] `locate`
- [ ] `updatedb`
- [x] `xargs` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)

### gawk

- [x] `awk` (stage 4 batch 2e, 2026-11-08; see completed-page evidence)
- [x] `gawk` (an alias of the `awk` page, which has a GNU view for gawk; see completed-page evidence)
- [ ] `gawkbug`

### grep

- [x] `egrep` (stage 4 batch 3j, 2026-11-10; an alias of the `grep` page, which explains the obsolescent spelling)
- [x] `fgrep` (stage 4 batch 3j, 2026-11-10; an alias of the `grep` page, which explains the obsolescent spelling)
- [x] `grep` — expanded and reviewed 2026-11-08 (see completed-page evidence)

### gzip

- [x] `gunzip` (stage 4 batch 3j, 2026-11-10; an alias of the `gzip` page, with an example)
- [x] `gzexe` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `gzip` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)
- [x] `uncompress` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `zcat` (stage 4 batch 3j, 2026-11-10; see completed-page evidence)
- [x] `zcmp` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `zdiff` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `zegrep` (stage 4 batch 3k, 2026-11-10; an alias of the `zgrep` page)
- [x] `zfgrep` (stage 4 batch 3k, 2026-11-10; an alias of the `zgrep` page)
- [x] `zforce` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `zgrep` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `zless` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `zmore` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)
- [x] `znew` (stage 4 batch 3k, 2026-11-10; see completed-page evidence)

### patch

- [x] `patch` (stage 4 batch 3j, 2026-11-10; see completed-page evidence)

### sed

- [x] `sed` — expanded and reviewed 2026-11-08 (see completed-page evidence)

### tar

- [ ] `rmt`
- [x] `tar` (stage 4 batch 2c, 2026-11-08; see completed-page evidence)

### wget

- [ ] `wget`

## Additional commands and tools

Keep distinct commands such as `pkill`, `killall`, `zip`, `unzip`, `ssh-agent`, and `limactl` separate even when the reference routes them to another command’s sheet. Choose canonical IDs/aliases only after checking semantics and implementation selection.

- [x] `actionlint` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `aws` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `bash` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `bat` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `claude` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `codex` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `colima` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `curl` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `docker` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `docker-compose` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `fd` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `free` (stage 4 batch 3m, 2026-11-10; see completed-page evidence)
- [x] `fzf` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `gh` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `git` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `gitleaks` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `go` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `gzcat` (stage 4 batch 3k, 2026-11-10; an alias of the `zcat` page)
- [x] `jq` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `just` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `killall` (stage 4 batch 3m, 2026-11-10; see completed-page evidence)
- [x] `lefthook` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `limactl` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `lsof` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `markdownlint-cli2` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `mise` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `node` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `npm` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `npx` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `pgrep` (stage 4 batch 3m, 2026-11-10; see completed-page evidence)
- [x] `pkill` (stage 4 batch 3m, 2026-11-10; see completed-page evidence)
- [x] `prettier` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `ps` (stage 4 batch 3m, 2026-11-10; see completed-page evidence)
- [x] `pulumi` (stage 4 batch 3q, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `python3` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `rg` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `rsync` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `scp` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `sd` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `security` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `shellcheck` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `shfmt` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `ssh` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `ssh-add` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `ssh-agent` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `ssh-copy-id` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `ssh-keygen` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `starship` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `top` (stage 4 batch 3m, 2026-11-10; see completed-page evidence)
- [x] `tree` (stage 4 batch 3l, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `unzip` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `uv` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `uvx` (stage 4 batch 3n, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `vim` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `zed` (stage 4 batch 3p, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `zip` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)
- [x] `zoxide` (stage 4 batch 3o, 2026-11-10; see the upstream tools table in the completed-page evidence)

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
