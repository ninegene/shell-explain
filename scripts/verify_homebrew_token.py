#!/usr/bin/env python3
"""Read-only verification of the Homebrew token and current formula."""

import base64
import hashlib
import os
import sys

from update_homebrew import FORMULA, SOURCE, TAG, TAP, api, request, updated_formula


def main():
    token = os.environ.get("GH_TOKEN", "")
    if not token:
        raise ValueError("HOMEBREW_TAP_TOKEN is not configured")
    release = api(f"/repos/{SOURCE}/releases/latest", token)
    tag = release["tag_name"]
    if not TAG.fullmatch(tag):
        raise ValueError("latest release does not have a stable version tag")
    tap = api(f"/repos/{TAP}", token)
    if not tap.get("permissions", {}).get("push"):
        raise ValueError("GitHub does not report account push access to the tap")
    api(f"/repos/{TAP}/pulls?state=open&per_page=1", token)
    contents = api(f"/repos/{TAP}/contents/{FORMULA}?ref=main", token)
    formula = base64.b64decode(contents["content"]).decode("utf-8")
    source = f"https://github.com/{SOURCE}/archive/refs/tags/{tag}.tar.gz"
    digest = hashlib.sha256(request(source)).hexdigest()
    if formula != updated_formula(formula, tag, digest):
        raise ValueError("current formula does not match the latest release archive")
    print(f"Authenticated reads and the {tag} formula checksum passed; GitHub reports account push access.")
    print("Token write permissions and release-event PR creation were not exercised.")


if __name__ == "__main__":
    try:
        main()
    except (OSError, ValueError, KeyError) as exc:
        print(f"Homebrew verification failed: {exc}", file=sys.stderr)
        sys.exit(1)
