"""Portable extractor selection checks; no downloads or QEMU required."""
import importlib.util
from pathlib import Path
import subprocess
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location(
    "linux_smoke", Path(__file__).with_name("linux-smoke.py")
)
smoke = importlib.util.module_from_spec(spec)
spec.loader.exec_module(smoke)


class ExtractorTests(unittest.TestCase):
    def select(self, executables, banners):
        def version(args, **_kwargs):
            return subprocess.CompletedProcess(args, 0, banners[args[0]], "")

        with patch.object(smoke.shutil, "which", side_effect=executables.get), \
             patch.object(smoke.subprocess, "run", side_effect=version):
            return smoke.iso_tar()

    def test_explicit_bsdtar(self):
        self.assertEqual(self.select(
            {"bsdtar": "/tools/bsdtar"},
            {"/tools/bsdtar": "bsdtar 3.8.0 - libarchive 3.8.0"},
        ), "/tools/bsdtar")

    def test_macos_without_bsdtar_and_gnu_first_on_path(self):
        self.assertEqual(self.select(
            {"/usr/bin/tar": "/usr/bin/tar", "tar": "/tools/tar"},
            {"/usr/bin/tar": "bsdtar 3.5.3 - libarchive 3.7.4",
             "/tools/tar": "tar (GNU tar) 1.35"},
        ), "/usr/bin/tar")

    def test_libarchive_tar_on_path(self):
        self.assertEqual(self.select(
            {"tar": "/tools/tar"},
            {"/tools/tar": "bsdtar 3.8.0 - libarchive 3.8.0"},
        ), "/tools/tar")

    def test_reject_gnu_tar(self):
        with self.assertRaisesRegex(SystemExit, "Install bsdtar"):
            self.select({"tar": "/tools/tar"},
                        {"/tools/tar": "tar (GNU tar) 1.35"})


if __name__ == "__main__":
    unittest.main()
