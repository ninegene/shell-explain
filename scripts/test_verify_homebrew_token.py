import base64
import os
import unittest
from unittest.mock import patch

import verify_homebrew_token as verify


class VerifyHomebrewTokenTests(unittest.TestCase):
    @patch.dict(os.environ, {"GH_TOKEN": "test"})
    def test_verifies_matching_release_without_writes(self):
        formula = ('  url "https://github.com/ninegene/shell-explain/archive/refs/tags/v0.1.0.tar.gz"\n'
                   '  sha256 "' + 'a' * 64 + '"\n')
        responses = [{"tag_name": "v0.1.0"}, {"permissions": {"push": True}}, [],
                     {"content": base64.b64encode(formula.encode()).decode()}]
        with patch.object(verify, "api", side_effect=responses) as api, \
             patch.object(verify, "request", return_value=b"archive"), \
             patch.object(verify, "updated_formula", return_value=formula):
            verify.main()
        self.assertTrue(all(len(call.args) == 2 for call in api.call_args_list))

    @patch.dict(os.environ, {"GH_TOKEN": "test"})
    def test_rejects_missing_tap_write_access(self):
        with patch.object(verify, "api", side_effect=[{"tag_name": "v0.1.0"}, {"permissions": {"push": False}}]):
            with self.assertRaisesRegex(ValueError, "push access"):
                verify.main()

    @patch.dict(os.environ, {"GH_TOKEN": "test"})
    def test_rejects_formula_mismatch(self):
        responses = [{"tag_name": "v0.1.0"}, {"permissions": {"push": True}}, [],
                     {"content": base64.b64encode(b"old formula").decode()}]
        with patch.object(verify, "api", side_effect=responses), \
             patch.object(verify, "request", return_value=b"archive"), \
             patch.object(verify, "updated_formula", return_value="new formula"):
            with self.assertRaisesRegex(ValueError, "does not match"):
                verify.main()


if __name__ == "__main__":
    unittest.main()
