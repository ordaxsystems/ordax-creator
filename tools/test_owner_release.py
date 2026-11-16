"""Fast, media-only contract tests for the independent Creator publisher."""
import gzip
import hashlib
import json
from pathlib import Path
import tempfile
import unittest

from owner_release import (
    BUNDLE, RELEASE_FILES, digest, package, stage, up_to_date, valid_media,
    verify_published,
)

CREATOR_SHA = "a" * 40
CREATOR_TREE = "b" * 40
OS_SHA = "c" * 40
ASSEMBLY_SHA = "d" * 40


def sha(data):
    return hashlib.sha256(data).hexdigest()


class OwnerReleaseTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.path = Path(self.tmp.name)
        self.root = self.path / "creator"
        self.root.mkdir()
        (self.root / "version.json").write_text(
            json.dumps({"$schema": "ordax.creator-product-version/1", "version": "0.6.5-dev"}),
            encoding="utf-8",
        )
        (self.root / "cmd/ordax-creator-app").mkdir(parents=True)
        self.media = self.path / "media"
        self.media.mkdir()
        raw = b"verified-qemu-only-synthetic-seed" * 127
        compressed = gzip.compress(raw, mtime=0)
        (self.media / "ordax-bootstrap-seed.raw.gz").write_bytes(compressed)
        self.symbol = b"\x89PNG\r\n\x1a\n-synthetic-symbol"
        self.theme = b'{"theme":"source-owned-fixture"}'
        (self.media / "ordax-symbol.png").write_bytes(self.symbol)
        (self.media / "ordax-design-theme.json").write_bytes(self.theme)
        self.metadata = {
            "$schema": "prototype-ordax.creator-owner-media/1",
            "channel": "owner-prototype",
            "status": "qemu-accepted-development-media-only",
            "source_repository": "ordaxsystems/ordax-os",
            "source_commit": OS_SHA,
            "assembly_creator_source_commit": ASSEMBLY_SHA,
            "seed_sha256": sha(raw),
            "seed_size": len(raw),
            "compressed_seed_sha256": sha(compressed),
            "compressed_seed_size": len(compressed),
            "symbol_sha256": sha(self.symbol),
            "theme_sha256": sha(self.theme),
            "trust_sha256": "e" * 64,
            "prototype_manifest_sha256": "f" * 64,
            "qemu_proof_sha256": "a" * 64,
            "qemu_two_boot_acceptance": True,
            "physical_hardware_proven": False,
            "canonical_public_release": False,
            "private_key_in_package": False,
            "windows_executable_in_media": False,
        }
        (self.media / "creator-owner-media.json").write_text(
            json.dumps(self.metadata), encoding="utf-8",
        )
        self.bundle = self.path / "bundle"
        self.release = self.path / "release"

    def test_one_executable_and_two_independent_revision_identities(self):
        data = valid_media(self.media, verify_gzip=True)
        self.assertEqual(data["source_commit"], OS_SHA)
        stage(self.media, self.root, self.bundle, CREATOR_SHA, CREATOR_TREE)
        self.assertEqual((self.root / "cmd/ordax-creator-app/owner_payload/provenance.json").read_bytes(),
                         (self.bundle / "provenance.json").read_bytes())
        (self.bundle / "OrdaX-Creator.exe").write_bytes(b"MZfake-compiler-artifact")
        package(self.media, self.root, self.bundle, self.release, CREATOR_SHA, CREATOR_TREE)
        manifest = json.loads((self.release / "creator-owner-update.json").read_text())
        self.assertEqual(manifest["$schema"], "prototype-ordax.creator-owner-bundle-update/2")
        self.assertEqual(manifest["source_commit"], OS_SHA)
        self.assertEqual(manifest["creator_source_commit"], CREATOR_SHA)
        self.assertEqual(manifest["owner_repository"], "ordaxsystems/ordax-creator")
        self.assertFalse(manifest["automatic_in_app_update"])
        self.assertTrue(up_to_date(self.media / "creator-owner-media.json",
                                   self.release / "creator-owner-update.json", self.root, CREATOR_SHA))
        self.assertFalse(up_to_date(self.media / "creator-owner-media.json",
                                    self.release / "creator-owner-update.json", self.root, "1"*40))
        import zipfile
        with zipfile.ZipFile(self.release / BUNDLE) as z:
            names = z.namelist()
        self.assertEqual(len(names), 5)
        self.assertEqual([x for x in names if x.endswith(".exe")],
                         ["OrdaX-Creator/OrdaX-Creator.exe"])
        record = {"tag_name": "creator-owner-prototype", "prerelease": True, "draft": False,
                  "assets": []}
        for name in RELEASE_FILES:
            h, size = digest(self.release / name)
            record["assets"].append({"name": name, "digest": "sha256:" + h,
                                     "size": size, "state": "uploaded"})
        remote = self.path / "release.json"
        remote.write_text(json.dumps(record))
        verify_published(self.release, remote)

    def test_mismatched_raw_sha_rejected(self):
        self.metadata["seed_sha256"] = "0" * 64
        (self.media / "creator-owner-media.json").write_text(json.dumps(self.metadata))
        with self.assertRaisesRegex(ValueError, "decompressed image"):
            valid_media(self.media, verify_gzip=True)

    def test_extra_file_or_prototype_authority_rejected(self):
        self.metadata["canonical_public_release"] = True
        (self.media / "creator-owner-media.json").write_text(json.dumps(self.metadata))
        with self.assertRaisesRegex(ValueError, "not accepted"):
            valid_media(self.media, verify_gzip=False)
        self.metadata["canonical_public_release"] = False
        (self.media / "creator-owner-media.json").write_text(json.dumps(self.metadata))
        (self.media / "ordax-creator-physical-test.exe").write_bytes(b"MZ")
        with self.assertRaisesRegex(ValueError, "unexpected"):
            valid_media(self.media, verify_gzip=False)

    def test_existing_generated_payload_rejected(self):
        stage(self.media, self.root, self.bundle, CREATOR_SHA, CREATOR_TREE)
        with self.assertRaisesRegex(ValueError, "already contains"):
            stage(self.media, self.root, self.bundle, CREATOR_SHA, CREATOR_TREE)


if __name__ == "__main__":
    unittest.main()
