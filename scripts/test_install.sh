#!/usr/bin/env bash
# Test scripts/install.sh against a local fixture release (no network).
set -euo pipefail

usage() {
	cat <<'HELP'
Test the installer against a release built from this checkout.

Usage: test_install.sh

Flags:
  -h, --help  Print help
HELP
}

fail() {
	printf 'FAIL: %s\n' "$*" >&2
	exit 1
}

main() {
	case "${1:-}" in
	-h | --help)
		usage
		return 0
		;;
	esac
	cd "$(git rev-parse --show-toplevel)"
	export BUILD_RELEASE_SKIP_NOTICE_CHECK=1
	local work target out
	work="$(mktemp -d)"
	trap 'rm -rf "${work:-}"' EXIT
	target="$(go env GOOS)/$(go env GOARCH)"
	mkdir -p "$work/base/v9.8.7"
	bash scripts/build-release.sh --targets "$target" v9.8.7 "$work/base/v9.8.7" >/dev/null
	export EXPLAIN_INSTALL_BASE_URL="file://$work/base"
	local install=(sh scripts/install.sh --version v9.8.7 --dir "$work/bin")

	out="$("${install[@]}" --dry-run)"
	[[ "$out" == *"would download"* && ! -e "$work/bin/explain" ]] || fail "dry run wrote something or printed nothing"

	"${install[@]}" >/dev/null 2>&1 || fail "install failed"
	[[ "$("$work/bin/explain" --version)" == "explain 9.8.7 "* ]] || fail "installed binary reports the wrong version"
	[[ -x "$work/bin/explain" ]] || fail "installed file is not executable"

	if "${install[@]}" >/dev/null 2>&1; then fail "an existing explain was replaced without --force"; fi
	"${install[@]}" --force >/dev/null 2>&1 || fail "--force did not replace the file"

	# A tampered checksum must stop the install before anything is written.
	sed -i.bak 's/^[0-9a-f]/0/' "$work/base/v9.8.7/checksums.txt"
	rm -f "$work/bin/explain"
	if "${install[@]}" >/dev/null 2>&1; then fail "a bad checksum was accepted"; fi
	[[ ! -e "$work/bin/explain" ]] || fail "a file was installed despite a bad checksum"

	if sh scripts/install.sh --version 1.2.3 --dir "$work/bin" >/dev/null 2>&1; then fail "a version without v was accepted"; fi
	if sh scripts/install.sh --bogus >/dev/null 2>&1; then fail "an unknown flag was accepted"; fi
	printf 'ok: installer\n'
}

main "$@"
