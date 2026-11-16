"""Targeted non-destructive tests for Owner USB Git transfer."""
import hashlib
import importlib.util
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest import mock

HERE = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location("owner_usb_git_delivery", HERE / "owner_usb_git_delivery.py")
APP = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(APP)


def git(repo, *args):
    proc = subprocess.run(["git", "-C", str(repo), *args],
                          capture_output=True, text=True, check=True, timeout=30)
    return proc.stdout.strip()


class DeliveryTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.source = self.root / "source"
        self.source.mkdir()
        git(self.source, "init", "-b", "main")
        git(self.source, "config", "user.name", "Test")
        git(self.source, "config", "user.email", "test@example.invalid")
        (self.source / "sample.txt").write_text("a\n")
        git(self.source, "add", "sample.txt")
        git(self.source, "commit", "-m", "base")
        self.sha = git(self.source, "rev-parse", "HEAD")
        self.delivery = self.root / "data" / APP.SUBDIR
        self.bundle = self.root / APP.BUNDLE

    def make_bundle(self):
        return APP.build_bundle(self.source, self.bundle, self.sha)

    def test_bundles_only_exact_main_with_matching_sha(self):
        digest, size = self.make_bundle()
        self.assertGreater(size, 0)
        self.assertEqual(digest, APP.sha256_file(self.bundle))
        self.assertEqual(git(self.source, "bundle", "list-heads", self.bundle),
                         self.sha + " refs/heads/main")

    def test_atomic_install_stages_without_changing_git_checkout(self):
        digest, size = self.make_bundle()
        APP.install_bundle(self.delivery, self.bundle, self.sha, digest, size)
        self.assertTrue(APP.existing_current(self.delivery, self.sha))
        self.assertEqual(git(self.source, "rev-parse", "HEAD"), self.sha)
        self.assertFalse((self.delivery / (APP.MANIFEST + ".pending")).exists())
        self.assertFalse((self.delivery / (APP.BUNDLE + ".pending")).exists())

    def test_rejects_tampering_with_existing_bundle(self):
        digest, size = self.make_bundle()
        APP.install_bundle(self.delivery, self.bundle, self.sha, digest, size)
        with (self.delivery / APP.BUNDLE).open("ab") as handle:
            handle.write(b"bad")
        self.assertFalse(APP.existing_current(self.delivery, self.sha))

    def test_rejects_mismatched_digest_without_publishing_manifest(self):
        digest, size = self.make_bundle()
        with self.assertRaisesRegex(APP.DeliveryError, "readback"):
            APP.install_bundle(self.delivery, self.bundle, self.sha, "0" * 64, size)
        self.assertFalse((self.delivery / APP.MANIFEST).exists())

    def test_rejects_invalid_manifest_values(self):
        for value in ("INVALID", "F" * 40):
            with self.assertRaises(APP.DeliveryError):
                APP.manifest_bytes(value, "a" * 64, 1)
        with self.assertRaises(APP.DeliveryError):
            APP.manifest_bytes(self.sha, "a" * 64, APP.MAX_BUNDLE + 1)

    def test_rejects_uncanonical_local_origin(self):
        git(self.source, "remote", "add", "origin", "https://evil.invalid/repo.git")
        with self.assertRaisesRegex(APP.DeliveryError, "canonical"):
            APP.source_commit(self.source)

    def test_safe_volume_gate_does_not_accept_wrong_letters(self):
        for bad in ("E:", "C\\", "EE", "1"):
            with self.assertRaises(APP.DeliveryError):
                APP.checked_volume(bad)

    def test_unsafe_pending_link_is_rejected(self):
        digest, size = self.make_bundle()
        self.delivery.mkdir(parents=True)
        link = self.delivery / (APP.BUNDLE + ".pending")
        try:
            link.symlink_to(self.root / "victim")
        except OSError:
            self.skipTest("Windows symlink creation not permitted")
        with self.assertRaisesRegex(APP.DeliveryError, "unsafe USB pending path"):
            APP.install_bundle(self.delivery, self.bundle, self.sha, digest, size)

    def test_auto_discovery_stages_only_local_serial_match(self):
        with (
            mock.patch.object(APP, "available_ordax_data_letters", return_value=["E", "F"]),
            mock.patch.object(APP, "checked_volume",
                              side_effect=[{"serial": "UNTRUSTED"}, {"serial": "6699E36B"}]) as checked,
        ):
            self.assertEqual(APP.discover_trusted_volume("6699e36b"), "F")
        self.assertEqual(checked.call_count, 2)

    def test_auto_discovery_with_no_trusted_usb_is_read_only(self):
        with mock.patch.object(APP, "available_ordax_data_letters", return_value=[]):
            self.assertIsNone(APP.discover_trusted_volume("6699E36B"))
        with self.assertRaisesRegex(APP.DeliveryError, "malformed"):
            APP.discover_trusted_volume("bad serial; injection")

    def test_auto_discovery_rejects_ambiguous_authorized_serial(self):
        with (
            mock.patch.object(APP, "available_ordax_data_letters", return_value=["E", "F"]),
            mock.patch.object(APP, "checked_volume",
                              side_effect=[{"serial": "6699E36B"}, {"serial": "6699E36B"}]),
        ):
            with self.assertRaisesRegex(APP.DeliveryError, "multiple volumes"):
                APP.discover_trusted_volume("6699E36B")


if __name__ == "__main__":
    unittest.main()
