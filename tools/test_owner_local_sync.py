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

    @unittest.skipUnless(__import__("os").name == "nt", "Windows-only native process snapshot")
    def test_creator_running_check_is_native_and_never_spawns_tasklist(self):
        with mock.patch.object(SYNC.subprocess, "run",
                               side_effect=AssertionError("unattended task must not spawn tasklist")):
            result = SYNC.creator_process_active()
        self.assertIs(type(result), bool)

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
        # Old packages remain readable; the new digest is strictly checked
        # when supplied by the canonical independent Creator publisher.
        current = dict(self.manifest, os_media_payload_sha256="f" * 64)
        self.assertIs(SYNC.validate_manifest(current), current)
        for key, replacement in [
            ("download_url", "https://example.com/Creator.zip"),
            ("automatic_in_app_update", True),
            ("sha256", "not-a-digest"),
            ("os_media_payload_sha256", "not-a-digest"),
            ("os_media_payload_sha256", None),
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
        observed = []
        actual_copy = SYNC.shutil.copyfileobj
        def capture_copy(source, destination, length):
            observed.append((Path(destination.name).name, length))
            return actual_copy(source, destination, length=length)
        with mock.patch.object(SYNC.shutil, "copyfileobj", side_effect=capture_copy):
            SYNC.verify_extracted(stage, archive, self.manifest, smoke=False)
        self.assertTrue((stage / "OrdaX-Creator/OrdaX-Creator.exe").is_file())
        self.assertIn(("OrdaX-Creator.exe.part", 8 * 1024 * 1024), observed)
        self.assertFalse((stage / "OrdaX-Creator/OrdaX-Creator.exe.part").exists())
        with mock.patch.object(SYNC.shutil, "copyfileobj",
                               side_effect=OSError("interrupted archive extraction")):
            with self.assertRaisesRegex(OSError, "interrupted archive"):
                SYNC.verify_extracted(self.root / "interrupted-stage", archive, self.manifest, smoke=False)
        self.assertFalse((self.root / "interrupted-stage/OrdaX-Creator/OrdaX-Creator.exe.part").exists())
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

    def test_local_git_lineage_fast_path_never_calls_remote_compare(self):
        import subprocess
        repo = self.root / "git-source"
        subprocess.run(["git", "init", "-b", "main", str(repo)], check=True,
                       capture_output=True)
        def run(*args):
            return subprocess.run(["git", "-C", str(repo), *args], check=True,
                                  capture_output=True, text=True).stdout.strip()
        run("config", "user.email", "ci@example.invalid")
        run("config", "user.name", "CI")
        run("remote", "add", "origin", "https://github.com/ordaxsystems/ordax-os.git")
        (repo / "first.txt").write_text("first", encoding="utf-8")
        run("add", "first.txt")
        run("commit", "-m", "first")
        old = run("rev-parse", "HEAD")
        (repo / "second.txt").write_text("second", encoding="utf-8")
        run("add", "second.txt")
        run("commit", "-m", "second")
        new = run("rev-parse", "HEAD")
        with (
            mock.patch.object(SYNC, "local_repo_path", return_value=repo),
            mock.patch.object(SYNC, "request_json") as remote,
        ):
            self.assertEqual(SYNC.source_relation("ordax-os", old, new), 1)
            remote.assert_not_called()
            with self.assertRaises(SYNC.SyncError):
                SYNC.source_relation("ordax-os", new, old)
            remote.assert_not_called()

    def test_noncanonical_local_source_falls_back_to_remote_authority(self):
        with (
            mock.patch.object(SYNC, "local_repo_path", return_value=self.root / "missing"),
            mock.patch.object(SYNC, "request_json", return_value={
                "status": "ahead", "ahead_by": 2, "behind_by": 0
            }) as remote,
        ):
            self.assertEqual(SYNC.source_relation("ordax-os", "a"*40, "b"*40), 1)
            remote.assert_called_once()

    def test_creator_gui_reuses_single_usb_transport_without_app_replacement(self):
        receipt = {"status": "CURRENT", "source_commit": "a" * 40}
        with (
            mock.patch.object(SYNC.sys, "argv", [
                "owner_local_sync.py", "--creator-ui-stage",
                "--stage-owner-dev-usb-serial", "6699E36B",
                "--trusted-disk-number", "1",
            ]),
            mock.patch.object(SYNC, "stage_owner_usb_dev", return_value=receipt) as stage,
            mock.patch.object(SYNC, "sync") as update_app,
            mock.patch("sys.stdout", new_callable=io.StringIO) as stdout,
        ):
            self.assertEqual(SYNC.main(), 0)
        stage.assert_called_once_with(
            "6699E36B", from_creator_ui=True, expected_disk_number=1,
        )
        update_app.assert_not_called()
        self.assertEqual(json.loads(stdout.getvalue()), receipt)

    def test_creator_gui_rejects_missing_disk_without_staging(self):
        with (
            mock.patch.object(SYNC.sys, "argv", [
                "owner_local_sync.py", "--creator-ui-stage",
                "--stage-owner-dev-usb-serial", "6699E36B",
            ]),
            mock.patch.object(SYNC, "stage_owner_usb_dev") as stage,
            mock.patch.object(SYNC, "sync") as update_app,
            mock.patch("sys.stderr", new_callable=io.StringIO) as stderr,
        ):
            self.assertEqual(SYNC.main(), 2)
        stage.assert_not_called()
        update_app.assert_not_called()
        self.assertIn("missing selected physical USB", stderr.getvalue())

    def test_creator_ui_retries_only_existing_writer_mutex_collision(self):
        busy = SYNC.subprocess.CompletedProcess(
            ["python"], 2, "", "OWNER_USB_GIT_DELIVERY=BLOCKED:Owner USB transport already running")
        ready = SYNC.subprocess.CompletedProcess(["python"], 0, '{"status": "CURRENT"}', "")
        with (
            mock.patch.object(SYNC.subprocess, "run", side_effect=[busy, ready]) as run,
            mock.patch.object(SYNC.time, "sleep") as sleep,
        ):
            actual = SYNC.run_owner_usb_transport(["python", "existing_writer.py"],
                                                  from_creator_ui=True)
        self.assertIs(actual, ready)
        self.assertEqual(run.call_count, 2)
        sleep.assert_called_once()

    def test_scheduled_task_defers_on_existing_writer_mutex(self):
        busy = SYNC.subprocess.CompletedProcess(
            ["python"], 2, "", "OWNER_USB_GIT_DELIVERY=BLOCKED:Owner USB transport already running")
        with (
            mock.patch.object(SYNC.subprocess, "run", return_value=busy) as run,
            mock.patch.object(SYNC.time, "sleep") as sleep,
        ):
            self.assertIsNone(SYNC.run_owner_usb_transport(["python"], from_creator_ui=False))
        run.assert_called_once()
        sleep.assert_not_called()

    def test_invalid_usb_never_retried_as_transient_mutex_collision(self):
        failed = SYNC.subprocess.CompletedProcess(
            ["python"], 2, "", "OWNER_USB_GIT_DELIVERY=BLOCKED:USB serial changed")
        with (
            mock.patch.object(SYNC.subprocess, "run", return_value=failed) as run,
            mock.patch.object(SYNC.time, "sleep") as sleep,
        ):
            self.assertIs(SYNC.run_owner_usb_transport(["python"], from_creator_ui=True), failed)
        run.assert_called_once()
        sleep.assert_not_called()

    @unittest.skipUnless(__import__("os").name == "nt", "Windows native cache lock")
    def test_git_main_stage_lock_serializes_without_removing_lockfile(self):
        with tempfile.TemporaryDirectory() as temp, mock.patch.dict(
                SYNC.os.environ, {"LOCALAPPDATA": temp}):
            with SYNC.owner_usb_stage_gate(0.0) as first:
                self.assertTrue(first)
                with SYNC.owner_usb_stage_gate(0.0) as second:
                    self.assertFalse(second)
            with SYNC.owner_usb_stage_gate(0.0) as third:
                self.assertTrue(third)
            self.assertTrue((Path(temp) / "OrdaX-Creator" / "owner-main-delivery.lock").exists())

    def test_trusted_usb_staging_defers_during_creator_physical_session(self):
        with mock.patch.object(SYNC, "creator_process_active", return_value=True):
            self.assertEqual(SYNC.stage_owner_usb_dev("6699E36B"),
                             "DEFERRED_CREATOR_RUNNING")
    def test_usb_is_prepared_before_creator_release_network_check(self):
        events = []
        with (
            mock.patch.object(SYNC.sys, "argv",
                              ["owner_local_sync.py", "--stage-owner-dev-usb-serial", "6699E36B"]),
            mock.patch.object(SYNC, "stage_owner_usb_dev",
                              side_effect=lambda serial: events.append(("usb", serial)) or "CURRENT"),
            mock.patch.object(SYNC, "sync",
                              side_effect=lambda *args, **kwargs: events.append(("creator", None)) or "CURRENT"),
        ):
            self.assertEqual(SYNC.main(), 0)
        self.assertEqual(events, [("usb", "6699E36B"), ("creator", None)])


if __name__ == "__main__":
    unittest.main()
