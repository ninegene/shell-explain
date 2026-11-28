#!/usr/bin/env python3
"""Open a tap PR for the latest published stable shell-explain release."""

import base64
import hashlib
import json
import os
import re
import sys
import urllib.error
import urllib.parse
import urllib.request

API = "https://api.github.com"
SOURCE = "ninegene/shell-explain"
TAP = "ninegene/homebrew-tap"
FORMULA = "Formula/shell-explain.rb"
TAG = re.compile(r"v(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\Z")


def request(url, token=None, data=None, method=None):
    headers = {"Accept": "application/vnd.github+json", "User-Agent": "shell-explain-release"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    payload = None if data is None else json.dumps(data).encode()
    if payload:
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=payload, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=60) as response:
        return response.read()


def api(path, token, data=None, method=None):
    return json.loads(request(API + path, token, data, method))


def updated_formula(formula, tag, digest):
    url = f"https://github.com/{SOURCE}/archive/refs/tags/{tag}.tar.gz"
    formula, urls = re.subn(
        rf'(?m)^  url "https://github\.com/{SOURCE}/archive/refs/tags/v[^"/]+\.tar\.gz"$',
        f'  url "{url}"', formula,
    )
    formula, hashes = re.subn(r'(?m)^  sha256 "[0-9a-f]{64}"$', f'  sha256 "{digest}"', formula)
    if urls != 1 or hashes != 1:
        raise ValueError("formula does not contain exactly one expected URL and checksum")
    return formula


def main():
    tag, token = os.environ.get("RELEASE_TAG", ""), os.environ.get("GH_TOKEN", "")
    if not TAG.fullmatch(tag):
        raise ValueError("release tag must be a stable version such as v1.2.3")
    if not token:
        raise ValueError("HOMEBREW_TAP_TOKEN is required; see docs/distribution.md")
    latest = api(f"/repos/{SOURCE}/releases/latest", token)
    if latest["tag_name"] != tag:
        print(f"Skipping {tag}: latest stable release is {latest['tag_name']}")
        return

    branch = f"automation/shell-explain-{tag}"
    branch_ref = urllib.parse.quote(branch, safe="")
    pulls = api(f"/repos/{TAP}/pulls?state=open&head=ninegene:{branch_ref}", token)
    if pulls:
        print(f"Update PR already exists: {pulls[0]['html_url']}")
        return
    source_url = f"https://github.com/{SOURCE}/archive/refs/tags/{tag}.tar.gz"
    digest = hashlib.sha256(request(source_url)).hexdigest()
    current = api(f"/repos/{TAP}/contents/{FORMULA}?ref=main", token)
    original = base64.b64decode(current["content"]).decode("utf-8")
    current_tag = re.search(r'(?m)^  url "https://github\.com/ninegene/shell-explain/archive/refs/tags/(v[0-9]+\.[0-9]+\.[0-9]+)\.tar\.gz"$', original)
    if not current_tag or not TAG.fullmatch(current_tag.group(1)):
        raise ValueError("formula does not contain a recognized stable release tag")
    if tuple(map(int, current_tag.group(1)[1:].split("."))) > tuple(map(int, tag[1:].split("."))):
        raise ValueError("latest GitHub release is older than the Homebrew formula; resolve the release mismatch")
    updated = updated_formula(original, tag, digest)
    if updated == original:
        print(f"Formula already points to {tag}")
        return
    base = api(f"/repos/{TAP}/git/ref/heads/main", token)["object"]["sha"]
    try:
        api(f"/repos/{TAP}/git/refs", token, {"ref": f"refs/heads/{branch}", "sha": base})
    except urllib.error.HTTPError as exc:
        if exc.code != 422:
            raise
        # A retry may reuse only our existing branch; never overwrite its contents.
        api(f"/repos/{TAP}/git/ref/heads/{branch_ref}", token)
    existing = api(f"/repos/{TAP}/contents/{FORMULA}?ref={branch_ref}", token)
    branch_formula = base64.b64decode(existing["content"]).decode("utf-8")
    if branch_formula != original and branch_formula != updated:
        raise ValueError("existing automation branch contains a different formula; inspect it manually")
    if branch_formula != updated:
        api(f"/repos/{TAP}/contents/{FORMULA}", token, {
            "message": f"chore: update shell-explain to {tag}",
            "content": base64.b64encode(updated.encode()).decode(),
            "sha": existing["sha"], "branch": branch,
        }, "PUT")
    comparison = api(f"/repos/{TAP}/compare/main...{branch_ref}", token)
    if [item["filename"] for item in comparison["files"]] != [FORMULA]:
        raise ValueError("automation branch contains unexpected changes; inspect it manually")
    pr = api(f"/repos/{TAP}/pulls", token, {
        "title": f"chore: update shell-explain to {tag} [{branch}]",
        "head": branch, "base": "main",
        "body": f"## Why?\n\nMake the published {tag} release available through Homebrew.\n\n"
                f"## Changelog\n\n- chore: update the source archive URL and SHA-256 for {tag}.\n\n"
                "## Validation\n\n- SHA-256 calculated from the tagged source archive. "
                "Review and test the formula before merging.\n",
    })
    print(f"Opened {pr['html_url']}")


if __name__ == "__main__":
    try:
        main()
    except (ValueError, KeyError, OSError, urllib.error.HTTPError) as exc:
        print(f"Homebrew update failed: {exc}", file=sys.stderr)
        sys.exit(1)
