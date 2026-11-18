#!/bin/sh
# Install the explain CLI from a GitHub release.
set -eu

REPO="ninegene/shell-explain"
BASE_URL="${EXPLAIN_INSTALL_BASE_URL:-https://github.com/$REPO/releases/download}"
LATEST_URL="${EXPLAIN_INSTALL_LATEST_URL:-https://github.com/$REPO/releases/latest}"

usage() {
	cat <<'HELP'
Install the explain CLI from a GitHub release of ninegene/shell-explain.

Usage: install.sh [flags]

Flags:
  --version VERSION  Release to install, such as v0.1.0 (default: the latest)
  --dir DIR          Directory to install into (default: ~/.local/bin)
  --force            Replace an explain that is already in DIR
  --dry-run          Show what would happen; download and write nothing
  -h, --help         Print help

The script finds your OS and CPU, downloads the matching archive and the
release's checksums.txt, checks the SHA-256 and installs only the explain
executable. It never uses sudo and never edits your shell startup files.
To uninstall, delete the explain file it installed.

Examples:
  $ install.sh
  $ install.sh --version v0.1.0 --dir ~/bin
  $ install.sh --dry-run
HELP
}

die() {
	printf 'install: %s\n' "$*" >&2
	exit 1
}

# Print the OS name used in release archive names.
detect_os() {
	case "$(uname -s)" in
	Darwin) printf 'darwin\n' ;;
	Linux) printf 'linux\n' ;;
	*) die "unsupported OS $(uname -s); releases cover macOS and Linux" ;;
	esac
}

# Print the CPU name used in release archive names.
detect_arch() {
	case "$(uname -m)" in
	arm64 | aarch64) printf 'arm64\n' ;;
	x86_64 | amd64) printf 'amd64\n' ;;
	*) die "unsupported CPU $(uname -m); releases cover arm64 and amd64" ;;
	esac
}

# Print the tag of the latest release by following the /releases/latest redirect.
latest_version() {
	url="$(curl -fsSL -o /dev/null -w '%{url_effective}' "$LATEST_URL")" ||
		die "could not look up the latest release at $LATEST_URL"
	tag="${url##*/}"
	case "$tag" in
	v[0-9]*) printf '%s\n' "$tag" ;;
	*) die "could not find a release tag in $url (no release yet? pass --version)" ;;
	esac
}

# Print the SHA-256 of a file.
sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | cut -d' ' -f1
	else
		die "need sha256sum or shasum to check the download"
	fi
}

verify_checksum() {
	archive_path="$1"
	archive_name="$2"
	checksums_path="$3"
	expected="$(awk -v n="$archive_name" '$2 == n { print $1 }' "$checksums_path")"
	[ -n "$expected" ] || die "checksums.txt has no entry for $archive_name"
	actual="$(sha256_of "$archive_path")"
	[ "$expected" = "$actual" ] ||
		die "checksum mismatch for $archive_name (expected $expected, got $actual); nothing was installed"
}

# Warn when another explain already comes first on PATH.
warn_about_path() {
	dir="$1"
	found="$(command -v explain 2>/dev/null || true)"
	if [ -n "$found" ] && [ "$found" != "$dir/explain" ]; then
		printf 'note: "explain" currently resolves to %s, which may hide the one installed in %s.\n' "$found" "$dir" >&2
		printf '      An alias or shell function named explain also takes precedence; check "type explain".\n' >&2
	fi
	case ":$PATH:" in
	*":$dir:"*) ;;
	*) printf 'note: %s is not on your PATH; add it to run "explain" directly.\n' "$dir" >&2 ;;
	esac
}

main() {
	version=""
	dir="$HOME/.local/bin"
	force=0
	dry_run=0
	while [ $# -gt 0 ]; do
		case "$1" in
		--version)
			[ $# -ge 2 ] || die "--version needs a value"
			version="$2"
			shift 2
			;;
		--dir)
			[ $# -ge 2 ] || die "--dir needs a value"
			dir="$2"
			shift 2
			;;
		--force)
			force=1
			shift
			;;
		--dry-run)
			dry_run=1
			shift
			;;
		-h | --help)
			usage
			return 0
			;;
		*) die "unknown flag $1 (try --help)" ;;
		esac
	done

	command -v curl >/dev/null 2>&1 || die "curl is required"
	os="$(detect_os)"
	arch="$(detect_arch)"
	[ -n "$version" ] || version="$(latest_version)"
	case "$version" in
	v[0-9]*.[0-9]*.[0-9]*) ;;
	*) die "version must look like v1.2.3 (got '$version')" ;;
	esac
	archive="explain_${version#v}_${os}_${arch}.tar.gz"
	archive_url="$BASE_URL/$version/$archive"
	checksums_url="$BASE_URL/$version/checksums.txt"

	if [ -e "$dir/explain" ] && [ "$force" -eq 0 ]; then
		die "$dir/explain already exists; use --force to replace it"
	fi
	if [ "$dry_run" -eq 1 ]; then
		printf 'would download %s\n' "$archive_url"
		printf 'would download %s\n' "$checksums_url"
		printf 'would verify the SHA-256 and install %s/explain\n' "$dir"
		return 0
	fi

	work="$(mktemp -d)"
	trap 'rm -rf "$work"' EXIT INT TERM
	printf 'downloading %s\n' "$archive_url"
	curl -fsSL -o "$work/$archive" "$archive_url" || die "could not download $archive_url"
	curl -fsSL -o "$work/checksums.txt" "$checksums_url" || die "could not download $checksums_url"
	verify_checksum "$work/$archive" "$archive" "$work/checksums.txt"
	tar -xzf "$work/$archive" -C "$work" ./explain || die "the archive has no explain executable"
	mkdir -p "$dir"
	install -m 0755 "$work/explain" "$dir/explain"
	printf 'installed %s (%s)\n' "$dir/explain" "$("$dir/explain" --version | cut -d' ' -f1-2)"
	warn_about_path "$dir"
}

main "$@"
