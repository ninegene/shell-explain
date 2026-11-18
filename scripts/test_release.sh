#!/usr/bin/env bash
# Test scripts/build-release.sh by building the host's own target.
set -euo pipefail

usage() {
	cat <<'HELP'
Build the release archive for this machine and check it.

Usage: test_release.sh

Checks the archive contents, the checksum file, the version reported by the
binary, and that building twice gives identical bytes.

Flags:
  -h, --help  Print help
HELP
}

fail() {
	printf 'FAIL: %s\n' "$*" >&2
	exit 1
}

host_target() {
	local os arch
	os="$(go env GOOS)"
	arch="$(go env GOARCH)"
	printf '%s/%s\n' "$os" "$arch"
}

archive_has() {
	grep -qxF "./$2" <<<"$1" || fail "archive is missing $2"
}

sha_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	else
		shasum -a 256 "$1" | cut -d' ' -f1
	fi
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
	local work target name archive listing file
	work="$(mktemp -d)"
	trap 'rm -rf "${work:-}"' EXIT
	target="$(host_target)"
	name="explain_9.8.7_${target%/*}_${target#*/}.tar.gz"

	bash scripts/build-release.sh --targets "$target" v9.8.7 "$work/a" >/dev/null
	bash scripts/build-release.sh --targets "$target" v9.8.7 "$work/b" >/dev/null
	archive="$work/a/$name"
	[[ -f "$archive" ]] || fail "missing $name"

	listing="$(tar -tzf "$archive")"
	for file in explain LICENSE LICENSE-CC-BY-4.0 README.md docs/third-party-notices.md \
		third-party/go-runtime/LICENSE third-party/go-runtime/PATENTS third-party/goldmark/LICENSE \
		third-party/go-yaml/LICENSE third-party/go-yaml/NOTICE third-party/go-yaml/Apache-2.0.txt; do
		archive_has "$listing" "$file"
	done

	[[ "$(cat "$work/a/checksums.txt")" == "$(sha_of "$archive")  $name" ]] ||
		fail "checksums.txt does not match the archive"
	[[ "$(sha_of "$archive")" == "$(sha_of "$work/b/$name")" ]] ||
		fail "two builds of the same commit differ: the archive is not reproducible"

	mkdir "$work/x"
	tar -xzf "$archive" -C "$work/x"
	[[ "$("$work/x/explain" --version)" == "explain 9.8.7 "* ]] ||
		fail "unexpected version output: $("$work/x/explain" --version)"
	"$work/x/explain" ssh >/dev/null || fail "explain ssh failed"
	bash scripts/build-release.sh bad >/dev/null 2>&1 && fail "a version without v was accepted"

	# A relative output directory is resolved against the caller's directory, and
	# unrelated files in it survive.
	mkdir "$work/rel" && : >"$work/rel/keep.txt"
	(cd "$work/rel" && bash "$OLDPWD/scripts/build-release.sh" --targets "$target" v9.8.7 out >/dev/null)
	[[ -f "$work/rel/out/$name" ]] || fail "a relative output directory was not used as given"
	: >"$work/rel/out/keep.txt"
	(cd "$work/rel" && bash "$OLDPWD/scripts/build-release.sh" --targets "$target" v9.8.7 out >/dev/null)
	[[ -f "$work/rel/out/keep.txt" && -f "$work/rel/keep.txt" ]] || fail "building deleted files it did not create"
	printf 'ok: %s\n' "$name"
}

main "$@"
