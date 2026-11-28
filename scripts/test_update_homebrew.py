import base64
import importlib.util
import os
from pathlib import Path
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("update_homebrew", Path(__file__).with_name("update_homebrew.py"))
updater = importlib.util.module_from_spec(spec)
spec.loader.exec_module(updater)

FORMULA = ('class ShellExplain < Formula\n'
           '  url "https://github.com/ninegene/shell-explain/archive/refs/tags/v0.1.0.tar.gz"\n'
           f'  sha256 "{"a" * 64}"\nend\n')


class UpdateHomebrewTests(unittest.TestCase):
    def test_update_preserves_formula_and_replaces_only_release_fields(self):
        result = updater.updated_formula(FORMULA, "v0.2.0", "b" * 64)
        self.assertEqual(result, FORMULA.replace("v0.1.0", "v0.2.0").replace("a" * 64, "b" * 64))

    def test_unexpected_formula_fails_closed(self):
        with self.assertRaises(ValueError):
            updater.updated_formula(FORMULA + '  sha256 "' + 'c' * 64 + '"\n', "v0.2.0", "b" * 64)

    @patch.dict(os.environ, {"RELEASE_TAG": "v0.2.0", "GH_TOKEN": "test"})
    def test_older_release_does_not_open_pr(self):
        with patch.object(updater, "api", return_value={"tag_name": "v0.3.0"}) as api:
            updater.main()
        self.assertEqual(api.call_count, 1)

    @patch.dict(os.environ, {"RELEASE_TAG": "v0.2.0", "GH_TOKEN": "test"})
    def test_existing_pr_is_not_duplicated(self):
        with patch.object(updater, "api", side_effect=[{"tag_name": "v0.2.0"}, [{"html_url": "https://example.org/pr"}]]) as api:
            updater.main()
        self.assertEqual(api.call_count, 2)

    @patch.dict(os.environ, {"RELEASE_TAG": "v0.0.9", "GH_TOKEN": "test"})
    def test_latest_github_release_cannot_downgrade_formula(self):
        content = {"content": base64.b64encode(FORMULA.encode()).decode(), "sha": "file-sha"}
        with patch.object(updater, "api", side_effect=[{"tag_name": "v0.0.9"}, [], content]) as api, patch.object(updater, "request", return_value=b"archive"):
            with self.assertRaisesRegex(ValueError, "older than the Homebrew formula"):
                updater.main()
        self.assertEqual(api.call_count, 3)

    @patch.dict(os.environ, {"RELEASE_TAG": "v0.2.0", "GH_TOKEN": "test"})
    def test_new_release_creates_branch_formula_commit_and_pr(self):
        content = {"content": base64.b64encode(FORMULA.encode()).decode(), "sha": "file-sha"}
        responses = [{"tag_name": "v0.2.0"}, [], content, {"object": {"sha": "base-sha"}},
                     {}, content, {}, {"files": [{"filename": updater.FORMULA}]},
                     {"html_url": "https://example.org/pr"}]
        with patch.object(updater, "api", side_effect=responses) as api, patch.object(updater, "request", return_value=b"archive"):
            updater.main()
        calls = api.call_args_list
        self.assertEqual(calls[4].args[2]["ref"], "refs/heads/automation/shell-explain-v0.2.0")
        self.assertIn("v0.2.0.tar.gz", base64.b64decode(calls[6].args[2]["content"]).decode())
        self.assertEqual(calls[8].args[2]["base"], "main")

    @patch.dict(os.environ, {"RELEASE_TAG": "v0.2.0", "GH_TOKEN": "test"})
    def test_unrelated_branch_changes_do_not_open_pr(self):
        content = {"content": base64.b64encode(FORMULA.encode()).decode(), "sha": "file-sha"}
        responses = [{"tag_name": "v0.2.0"}, [], content, {"object": {"sha": "base-sha"}},
                     {}, content, {}, {"files": [{"filename": updater.FORMULA}, {"filename": "other"}]}]
        with patch.object(updater, "api", side_effect=responses) as api, patch.object(updater, "request", return_value=b"archive"):
            with self.assertRaisesRegex(ValueError, "unexpected changes"):
                updater.main()
        self.assertEqual(api.call_count, 8)


if __name__ == "__main__":
    unittest.main()
