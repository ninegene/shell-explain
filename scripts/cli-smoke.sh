#!/usr/bin/env bash
# Run the explain CLI smoke matrix against a built binary and capture evidence.
set -euo pipefail

usage() {
	cat <<'HELP'
Run the explain CLI smoke matrix and capture its output as evidence.

Usage: cli-smoke.sh <BINARY> <OUTPUT-DIR>

Arguments:
  <BINARY>      Built explain executable
  <OUTPUT-DIR>  New directory for captured results

Flags:
  -h, --help    Print help

Examples:
  $ cli-smoke.sh bin/explain ./results
HELP
}

run_case() {
	local output="$1" name="$2" binary="$3" status=0
	shift 3
	env -u GOGC -u GOMAXPROCS -u GODEBUG LC_ALL=C TERM=dumb \
		"$binary" "$@" >"$output/$name.stdout" 2>"$output/$name.stderr" || status=$?
	printf '%s\n' "$status" >"$output/$name.exit"
	if [[ "$status" -ne 0 || ! -s "$output/$name.stdout" ]]; then
		printf 'FAILED: %s (exit %s)\n' "$name" "$status" >&2
		return 1
	fi
	printf 'PASS: %s\n' "$name"
}

check_output() {
	local description="$1"
	shift
	if "$@" >/dev/null; then
		printf 'PASS: %s\n' "$description"
	else
		printf 'FAILED: %s\n' "$description" >&2
		return 1
	fi
}

repeat_matches() {
	cmp -s "$1/grep-1.stdout" "$1/grep-2.stdout" && cmp -s "$1/grep-1.stdout" "$1/grep-3.stdout"
}

validate_outputs() {
	local output="$1" revision failures=0
	revision="$(sed -n 's/^explain .* content \([0-9a-f]\{64\}\))$/\1/p' "$output/version.stdout")"
	check_output 'version names a content revision' test -n "$revision" || failures=1
	# shellcheck disable=SC2016 # $revision is a jq variable, not a shell expansion.
	check_output 'JSON matches the version content revision' \
		jq -e --arg revision "$revision" \
		'.command.id == "grep" and .selection.implementation == "gnu" and .content_revision == $revision' \
		"$output/grep-json.stdout" || failures=1
	check_output 'repeated grep lookups match' repeat_matches "$output" || failures=1
	check_output 'list includes grep' grep -Fq 'grep ' "$output/list.stdout" || failures=1
	check_output 'search finds grep' grep -Fq 'grep ' "$output/search.stdout" || failures=1
	check_output 'use-case search finds tar examples' grep -Fq '[tar]' "$output/search-usecase.stdout" || failures=1
	check_output 'guide renders' grep -Fq 'From Keypress to Screen' "$output/guide.stdout" || failures=1
	check_output 'glossary renders' grep -Fq 'Kernel' "$output/glossary.stdout" || failures=1
	return "$failures"
}

run_matrix() {
	local output="$1" binary="$2" failures=0
	run_case "$output" grep-1 "$binary" --no-probe grep || failures=1
	run_case "$output" grep-2 "$binary" --no-probe grep || failures=1
	run_case "$output" grep-3 "$binary" --no-probe grep || failures=1
	run_case "$output" version "$binary" --version || failures=1
	run_case "$output" grep-json "$binary" --no-probe --implementation gnu --format json grep || failures=1
	run_case "$output" list "$binary" list || failures=1
	run_case "$output" search "$binary" search logs || failures=1
	run_case "$output" search-usecase "$binary" search extract archive || failures=1
	run_case "$output" guide "$binary" learn big-picture --no-pager || failures=1
	run_case "$output" glossary "$binary" glossary kernel --no-pager || failures=1
	validate_outputs "$output" || failures=1
	return "$failures"
}

main() {
	if [[ "${1:-}" == -h || "${1:-}" == --help ]]; then
		usage
		return 0
	fi
	if [[ $# -ne 2 ]]; then
		usage >&2
		return 2
	fi
	local binary="$1" output="$2"
	[[ -x "$binary" ]] || {
		printf 'Not an executable file: %s\n' "$binary" >&2
		return 1
	}
	mkdir "$output"
	{
		date -u
		uname -a
		if [[ -r /etc/os-release ]]; then
			cat /etc/os-release
		fi
		printf 'GOGC, GOMAXPROCS and GODEBUG unset for every case.\n'
	} >"$output/environment.txt"
	sha256sum "$binary" >"$output/executable.sha256"
	run_matrix "$output" "$binary"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
	main "$@"
fi
