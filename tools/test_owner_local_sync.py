"""Offline contract and fail-closed behavior for Owner Dev local Creator synchronization."""
from __future__ import annotations
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest import mock
import urllib.error
import zipfile

MODULE_PATH = Path(__file__).with_name("owner_local_sync.py")
SPEC = importlib.util.spec_from_file_location("owner_local_sync_tests", MODULE_PATH)
SYNC = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(SYNC)


class LocalSyncTests(unittest.TestCase):
    def setUp(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        self.root = Path(tmp.name) / "OrdaX-Creator-Owner-Prototype"
        self.root.mkdir()
        self.old = {
            "$schema": "prototype-ordax.creator-owner-physical/1",
            "creator_source_repository": SYNC.REPO,
            "creator_version": "0.6.4-dev",
            "source_commit": "a" * 40,
            "creator_source_commit": "b" * 40,
        }
        current = self.root / "current" / "OrdaX-Creator"
        current.mkdir(parents=True)
        (current / "provenance.json").write_text(json.dumps(self.old), encoding="utf-8")
        self.manifest = {
            "$schema": "prototype-ordax.creator-owner-bundle-update/2",
            "owner_repository": SYNC.REPO,
            "channel": "owner-prototype",
            "artifact": SYNC.BUNDLE_NAME,
            "download_url": SYNC.BASE + SYNC.BUNDLE_NAME,
            "automatic_in_app_update": False,
            "version": "0.6.5-dev",
            "source_commit": "c" * 40,
            "creator_source_commit": "d" * 40,
            "sha256": "e" * 64,
            "size": 4096,
        }

    def test_mutable_release_tag_uses_fresh_digest_bound_cache_keys(self):
        url = SYNC.fresh_manifest_url(12345)
        self.assertEqual(url, SYNC.MANIFEST_URL + "?ordax_nocache=12345")
        expected = self.manifest["sha256"]
        self.assertEqual(SYNC.verified_bundle_url(expected),
                         SYNC.BASE + SYNC.BUNDLE_NAME + "?ordax_sha256=" + expected)
        with self.assertRaises(SYNC.SyncError):
            SYNC.verified_bundle_url("not-a-sha")

    def test_manifest_rejects_unsafe_redirect_and_digest(self):
        self.assertIs(SYNC.validate_manifest(self.manifest), self.manifest)
        for key, replacement in [
            ("download_url", "https://example.com/Creator.zip"),
            ("automatic_in_app_update", True),
            ("sha256", "not-a-digest"),
            ("source_commit", "not-a-commit"),
            ("owner_repository", "ordaxsystems/ordax-os"),
        ]:
            value = dict(self.manifest, **{key: replacement})
            with self.assertRaises(SYNC.SyncError):
                SYNC.validate_manifest(value)

    @unittest.skipUnless(__import__("os").name == "nt", "Windows-only local updater")
    def test_missing_release_does_not_mutate_installed_creator(self):
        with mock.patch.object(SYNC, "request_json", side_effect=urllib.error.HTTPError(
            SYNC.MANIFEST_URL, 404, "Missing", None, None)):
            result = SYNC.sync(self.root)
        self.assertEqual(result, "DEFERRED_NO_PUBLISHED_RELEASE")
        self.assertEqual(SYNC.current_provenance(self.root), self.old)

    @unittest.skipUnless(__import__("os").name == "nt", "Windows-only local updater")
    def test_active_creator_defers_without_downloading_or_writing(self):
        with (
            mock.patch.object(SYNC, "request_json", return_value=self.manifest),
            mock.patch.object(SYNC, "source_relation", return_value=1),
            mock.patch.object(SYNC, "creator_process_active", return_value=True),
            mock.patch.object(SYNC, "download_bundle") as download,
        ):
            self.assertEqual(SYNC.sync(self.root), "DEFERRED_CREATOR_RUNNING")
            download.assert_not_called()
        self.assertEqual(SYNC.current_provenance(self.root), self.old)

    def test_zip_membership_is_exact_and_provenance_matches_two_repositories(self):
        contents = {
            "OrdaX-Creator/OrdaX-Creator.exe": b"MZ" + b"0" * 8,
            "OrdaX-Creator/provenance.json": json.dumps({
                **self.old,
                "source_commit": "c" * 40,
                "creator_source_commit": "d" * 40,
                "creator_version": "0.6.5-dev",
                "physical_write_authorized_in_binary": True,
            }).encode(),
            "OrdaX-Creator/ordax-symbol.png": b"icon",
            "OrdaX-Creator/ordax-design-theme.json": b"{}",
            "OrdaX-Creator/LEIA-ME.txt": b"docs",
        }
        archive = self.root / "fixture.zip"
        with zipfile.ZipFile(archive, "w") as z:
            for name, data in contents.items():
                z.writestr(name, data)
        stage = self.root / "stage"
        SYNC.verify_extracted(stage, archive, self.manifest, smoke=False)
        self.assertTrue((stage / "OrdaX-Creator/OrdaX-Creator.exe").is_file())
        with zipfile.ZipFile(archive, "a") as z:
            z.writestr("../evil.txt", b"unsafe")
        with self.assertRaisesRegex(SYNC.SyncError, "unexpected files"):
            SYNC.verify_extracted(self.root / "bad-stage", archive, self.manifest, smoke=False)

    def test_swap_recovery_returns_old_current_after_interruption(self):
        current = self.root / "current"
        previous = self.root / ".ordax-previous"
        current.rename(previous)
        SYNC.recover_interrupted_swap(self.root)
        self.assertEqual(SYNC.current_provenance(self.root), self.old)

    @unittest.skipUnless(__import__("os").name == "nt", "Windows-only local updater")
    def test_headless_verified_update_does_not_launch_gui_smoke(self):
        bundle = self.root / "verified.zip"
        bundle.write_bytes(b"already-verified-by-publisher")
        def staged(stage, archive, manifest, *, smoke):
            self.assertFalse(smoke, "Task Scheduler cannot perform interactive Win32 smoke")
            self.assertEqual(archive, bundle)
            target = stage / "OrdaX-Creator"
            target.mkdir()
            (target / "provenance.json").write_text(json.dumps({
                **self.old,
                "source_commit": manifest["source_commit"],
                "creator_source_commit": manifest["creator_source_commit"],
                "creator_version": manifest["version"],
            }), encoding="utf-8")
        with (
            mock.patch.object(SYNC, "request_json", return_value=self.manifest),
            mock.patch.object(SYNC, "source_relation", return_value=1),
            mock.patch.object(SYNC, "creator_process_active", return_value=False),
            mock.patch.object(SYNC, "download_bundle", return_value=bundle),
            mock.patch.object(SYNC, "verify_extracted", side_effect=staged),
        ):
            self.assertEqual(SYNC.sync(self.root), "UPDATED_VERIFIED")
        self.assertEqual(SYNC.current_provenance(self.root)["creator_version"], "0.6.5-dev")

    def test_trusted_usb_staging_defers_during_creator_physical_session(self):
        with mock.patch.object(SYNC, "creator_process_active", return_value=True):
            self.assertEqual(SYNC.stage_owner_usb_dev("6699E36B"),
                             "DEFERRED_CREATOR_RUNNING")


if __name__ == "__main__":
    unittest.main()
