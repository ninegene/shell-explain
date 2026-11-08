# Command implementation detection

## Principle

Select documentation for the installed command implementation, not merely its host operating system. macOS can run GNU utilities; Linux can run BusyBox; FreeBSD and Apple commands are not interchangeable labels for every behavior.

Homebrew installs GNU grep with a `g` prefix and offers a `gnubin` PATH location for unprefixed names. Consequently, `grep`, `ggrep`, and `/usr/bin/grep` can select different tools on the same machine. See [Homebrew grep](https://formulae.brew.sh/formula/grep).

## Implemented precedence

1. Explicit implementation override selects the requested documentation and labels it as an override.
2. An explicit executable path constrains local detection.
3. Resolve the command from PATH, inspect symlink/path evidence, and apply known command-specific identification rules.
4. Use narrowly scoped version probes only under a documented policy; apply a short timeout and output limit, without shell evaluation.
5. Fall back to a platform-associated default or common portable content, labeled as inferred/unknown.

Platform override and implementation override serve different purposes. `--platform linux` requests portable common content, not a claim that the local command is GNU. `--platform macos` selects Apple content for external tools; `--platform bash` selects Bash builtin documentation. Combining platform and implementation overrides is rejected.

For every selection, expose command, resolved path if available, implementation/version if known, evidence, and confidence. A filename alone is weak evidence. Do not assume every command accepts `--version`; an unknown wrapper should not be executed just to guess its behavior. Provide a no-probe mode, and avoid loading arbitrary shell startup files.

Current probes apply only to native-format executables for the commands in the allowlist `detect.probed` (grep, sed, find, cat, head, tail, wc, sort, uniq, cut, tr), which were checked against GNU, Apple, and BusyBox builds before being listed; any other command is never executed. Probes run --version directly without a shell, use a 400 ms execution timeout, a 200 ms `WaitDelay` for output draining, and an 8 KiB output cap, and supply no input. On Linux and macOS each probe starts in its own process group. Timeout cancellation kills that group, and cleanup also kills remaining group members after the parent exits. `WaitDelay` bounds output waiting even when a descendant escapes the group and holds an inherited output pipe open. Process groups provide cleanup, not a security sandbox; escaped descendants are not guaranteed to be terminated. These deadlines bound execution/output waiting rather than OS process-start latency. Script wrappers are not probed. Apple system paths can provide labeled inference; a BusyBox-named symlink target provides labeled inference. A native executable is still code that runs during a probe; use --no-probe to disable all such execution. The Apple grep banner says GNU compatible and macOS `sort --version` prints `2.3-Apple`, so the classifier checks Apple identity before a GNU banner pattern. macOS system tools live in `/usr/bin` or `/bin` (for example `/bin/cat`); both count as Apple system paths for labeled inference.

## Shell boundaries

A standalone process cannot generally read unexported aliases/functions from its parent shell. Shell builtins may shadow executables. Distinguish executable lookup from actual interactive-shell resolution in the output. Do not claim complete shell awareness by launching a fresh login shell.

Optional shell integration may explicitly pass command-resolution information. It should be a small reviewed adapter; normal use must not require it. During extraction from an existing bootstrap, remove or delegate the old explain function so it does not shadow the binary. Migration is future work and must not alter another repository without authorization.

## Website behavior

The static site cannot inspect a visitor's PATH, installed packages, aliases, or remote shell. Offer an explicit implementation selector and a portable default. Never infer an installed GNU/BSD tool from browser user-agent alone.

## Initial matrix

Validate at least macOS stock utilities, macOS with prefixed GNU tools, macOS with GNU tools first on PATH, Linux GNU tools, and an Alpine/BusyBox environment. Test missing executables, ambiguous wrappers, explicit paths, symlinks, shell builtins, failed probes, and user overrides. Add BSD variants only with a documented test/review environment.

The project explains behavior; it does not install GNU alternatives, rewrite PATH, execute examples, or guarantee command safety. Flag known incompatibilities without silently rewriting the user's command.
