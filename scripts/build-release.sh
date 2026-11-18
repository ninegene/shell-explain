#!/usr/bin/env bash
# Build the release archives and checksums for the explain CLI.
set -euo pipefail

DEFAULT_TARGETS="darwin/arm64 darwin/amd64 linux/arm64 linux/amd64"

usage() {
	cat <<'HELP'
Build release archives and a checksums file for the explain CLI.

Usage: build-release.sh [flags] <VERSION> [OUTPUT-DIR]

Arguments:
  <VERSION>     Release version such as v0.1.0 (a leading v is required)
  [OUTPUT-DIR]  Directory for the archives (default: dist)

Flags:
  --targets "OS/ARCH ..."  Targets to build (default: darwin/arm64 darwin/amd64
                           linux/arm64 linux/amd64)
  -h, --help               Print help

Environment:
  BUILD_RELEASE_SKIP_NOTICE_CHECK  Set to skip the Go-toolchain notice check (tests only)

Each archive is explain_<version>_<os>_<arch>.tar.gz with the executable at the
top level next to the licenses and notices. The binaries are pure Go
(CGO_ENABLED=0) and the archives are reproducible for a given commit: file times
come from the commit, and GNU tar plus gzip -n give a stable byte stream.

Examples:
  $ build-release.sh v0.1.0
  $ build-release.sh --targets linux/amd64 v0.1.0 /tmp/out
HELP
}

die() {
	printf 'build-release: %s\n' "$*" >&2
	exit 1
}

# Print the GNU tar command: gtar on macOS (Homebrew), tar elsewhere.
find_gnu_tar() {
	local candidate
	for candidate in gtar tar; do
		if command -v "$candidate" >/dev/null 2>&1 && "$candidate" --version 2>/dev/null | grep -q 'GNU tar'; then
			printf '%s\n' "$candidate"
			return 0
		fi
	done
	return 1
}

# Print a SHA-256 of a file: only the hex digest.
sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	else
		shasum -a 256 "$1" | cut -d' ' -f1
	fi
}

validate_version() {
	[[ "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$ ]] ||
		die "version must look like v1.2.3 (got '$1')"
}

# The notices record which Go toolchain's license files are shipped; refuse to
# build with a different toolchain until the notices are updated.
check_toolchain_notice() {
	# The script tests skip this check: they test the mechanics, not a release.
	[[ -z "${BUILD_RELEASE_SKIP_NOTICE_CHECK:-}" ]] || return 0
	local go_version
	go_version="$(go env GOVERSION)"
	grep -q "Go ${go_version#go}" docs/third-party-notices.md ||
		die "docs/third-party-notices.md does not mention ${go_version}; update the notices and third-party/go-runtime first"
}

# Copy the files that go into every archive next to the executable.
stage_extras() {
	local stage="$1"
	cp LICENSE LICENSE-CC-BY-4.0 README.md "$stage/"
	mkdir -p "$stage/docs"
	cp docs/third-party-notices.md "$stage/docs/"
	cp -R third-party "$stage/third-party"
}

build_target() {
	local version="$1" outdir="$2" target="$3" epoch="$4" tar_cmd="$5"
	local goos="${target%/*}" goarch="${target#*/}"
	local name="explain_${version#v}_${goos}_${goarch}"
	local stage
	stage="$(mktemp -d)"
	CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build -trimpath \
		-ldflags "-s -w -X main.version=${version#v}" -o "$stage/explain" ./cmd/explain
	stage_extras "$stage"
	# Fixed order, owner, times and modes keep the archive reproducible.
	"$tar_cmd" --sort=name --owner=0 --group=0 --numeric-owner --mtime="@$epoch" \
		--mode='u+rwX,go+rX,go-w' -cf - -C "$stage" . | gzip -n -9 >"$outdir/$name.tar.gz"
	rm -rf "$stage"
	printf 'built %s\n' "$outdir/$name.tar.gz"
}

write_checksums() {
	local outdir="$1" file
	: >"$outdir/checksums.txt"
	for file in "$outdir"/explain_*.tar.gz; do
		printf '%s  %s\n' "$(sha256_of "$file")" "$(basename "$file")" >>"$outdir/checksums.txt"
	done
}

main() {
	local targets="$DEFAULT_TARGETS" positional=()
	while [[ $# -gt 0 ]]; do
		case "$1" in
		-h | --help)
			usage
			return 0
			;;
		--targets)
			[[ $# -ge 2 ]] || die "--targets needs a value"
			targets="$2"
			shift 2
			;;
		-*) die "unknown flag $1 (try --help)" ;;
		*)
			positional+=("$1")
			shift
			;;
		esac
	done
	[[ ${#positional[@]} -ge 1 && ${#positional[@]} -le 2 ]] || {
		usage >&2
		return 2
	}
	local version="${positional[0]}" outdir="${positional[1]:-dist}" tar_cmd epoch target
	validate_version "$version"
	tar_cmd="$(find_gnu_tar)" || die "GNU tar is required (on macOS: brew install gnu-tar)"
	cd "$(git rev-parse --show-toplevel)"
	check_toolchain_notice
	epoch="$(git log -1 --format=%ct)"
	rm -rf "$outdir"
	mkdir -p "$outdir"
	for target in $targets; do
		[[ "$target" == */* ]] || die "target must look like linux/amd64 (got '$target')"
		build_target "$version" "$outdir" "$target" "$epoch" "$tar_cmd"
	done
	write_checksums "$outdir"
	printf 'wrote %s\n' "$outdir/checksums.txt"
}

main "$@"
