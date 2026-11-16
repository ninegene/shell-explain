#!/usr/bin/env bash
# Verify the existing Linux AMD64 candidate on a native Linux AMD64 host.
set -euo pipefail

usage() {
	cat <<'HELP'
Verify the fixed Linux AMD64 candidate and capture smoke-test evidence.

Usage: verify-linux-amd64-candidate.sh <ARCHIVE> <OUTPUT-DIR>

Arguments:
  <ARCHIVE>     Candidate tar.gz archive
  <OUTPUT-DIR>  New directory for captured results

Flags:
  -h, --help    Print help

Examples:
  $ verify-linux-amd64-candidate.sh candidate.tar.gz ./results
HELP
}

verify_file() {
	local file="$1" expected="$2" actual
	actual="$(sha256sum "$file")"
	actual="${actual%% *}"
	if [[ "$actual" != "$expected" ]]; then
		printf 'SHA-256 mismatch for %s\n' "$file" >&2
		return 1
	fi
	printf '%s  %s\n' "$actual" "$file"
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

compare_rosetta() {
	local output="$1" name expected actual
	while read -r name expected; do
		actual="$(sha256sum "$output/$name.stdout")"
		actual="${actual%% *}"
		if [[ "$actual" == "$expected" ]]; then
			printf '%s match\n' "$name"
		else
			printf '%s differs (inspect output; environment can affect results)\n' "$name"
		fi
	done >"$output/rosetta-comparison.txt" <<'HASHES'
grep-1 6c028b2e24a98ed157d06d5bc5c859364291adfc1d51359f06e55d0cff1d9885
grep-2 6c028b2e24a98ed157d06d5bc5c859364291adfc1d51359f06e55d0cff1d9885
grep-3 6c028b2e24a98ed157d06d5bc5c859364291adfc1d51359f06e55d0cff1d9885
version ba293f70695b69564886c86896d4f8f435c96bd24f85bf546065b25356803aea
grep-json fbff940ed2805aea40a2aa3ef7c49dc22b9771bbed47466c476a799ed0a01e76
list e4f31c5f884f38f294b501ce05b98ac7fe92c4ca9a6aa7fa44149ba17b2eb2f1
search 0ff98182c28ae432b6e1675dca5d621ec468121a5915e87050e1642419875b5b
guide 4d6a92099e596d21528f9d07dfcf21c8f4f2db22ef54f0dd614fc5f91728fe01
glossary a6557d261ae10d175d8d5e2bad4feaa9ac67e286d7940909d72915a62472a959
HASHES
}

validate_outputs() {
	local output="$1" failures=0
	jq -e '.command.id == "grep" and .selection.implementation == "gnu" and .content_revision == "8ee4e9006c7b9e79fe2cf87a2fb2e6b66da1279b7aeb0dbf45b234cb61718fe9"' "$output/grep-json.stdout" >"$output/grep-json-valid.txt" || failures=1
	grep -Fq 'explain dev (schema 1, content 8ee4e900' "$output/version.stdout" || failures=1
	grep -Fq 'grep ' "$output/list.stdout" || failures=1
	grep -Fq 'grep ' "$output/search.stdout" || failures=1
	grep -Fq 'From Keypress to Screen' "$output/guide.stdout" || failures=1
	grep -Fq 'Kernel' "$output/glossary.stdout" || failures=1
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
	run_case "$output" guide "$binary" learn big-picture --no-pager || failures=1
	run_case "$output" glossary "$binary" glossary kernel --no-pager || failures=1
	validate_outputs "$output" || failures=1
	compare_rosetta "$output"
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
	local archive="$1" output="$2" binary
	[[ "$(uname -s)" == Linux && "$(uname -m)" == x86_64 ]] || {
		printf 'Requires Linux x86_64\n' >&2
		return 1
	}
	mkdir "$output"
	{
		date -u
		uname -a
		cat /etc/os-release
		printf 'GOGC, GOMAXPROCS and GODEBUG unset for every case.\n'
	} >"$output/environment.txt"
	verify_file "$archive" bc587eaab6eaebe6581db3a0e5bec5b0b56095e7c24f14795e0dd5dd44d81430 >"$output/archive.sha256"
	tar -xzf "$archive" -C "$output"
	binary="$output/shell-explain-candidate-e63fcc22e4c9-linux-amd64/explain"
	verify_file "$binary" f3f2e1f62f64483176e760425a522ae0afb42fe2fa57da1d8f37436a33bb296a >"$output/executable.sha256"
	run_matrix "$output" "$binary"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
	main "$@"
fi
