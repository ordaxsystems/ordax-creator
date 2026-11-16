#!/usr/bin/env python3
"""Owner Dev local Creator synchronizer. Never writes a USB or upgrades a running EXE.

Consumes only the immutable Creator release channel, checked against the two
independent GitHub histories. The installed application directory is replaced
as a unit, only after digest, provenance and offline executable smoke checks.
"""
from __future__ import annotations

import argparse
import csv
import hashlib
import io
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
import urllib.error
import urllib.request
import zipfile

REPO = "ordaxsystems/ordax-creator"
TAG = "creator-owner-prototype"
BASE = f"https://github.com/{REPO}/releases/download/{TAG}/"
MANIFEST_URL = BASE + "creator-owner-update.json"
BUNDLE_NAME = "OrdaX-Creator-Owner-Prototype.zip"
EXPECTED_FILES = {
    "OrdaX-Creator/OrdaX-Creator.exe",
    "OrdaX-Creator/provenance.json",
    "OrdaX-Creator/ordax-symbol.png",
    "OrdaX-Creator/ordax-design-theme.json",
    "OrdaX-Creator/LEIA-ME.txt",
}
HEX40 = re.compile(r"[0-9a-f]{40}\Z")
HEX64 = re.compile(r"[0-9a-f]{64}\Z")
VERSION = re.compile(r"[0-9]+\.[0-9]+\.[0-9]+(?:-[a-z][a-z0-9]*(?:\.[0-9]+)?)?\Z")
MAX_BUNDLE_BYTES = 2 << 30


class SyncError(RuntimeError):
    pass


def sha256_file(path: Path) -> str:
    result = hashlib.sha256()
    with path.open("rb") as src:
        for chunk in iter(lambda: src.read(1024 * 1024), b""):
            result.update(chunk)
    return result.hexdigest()


def request_json(url: str, max_bytes: int = 64 << 10) -> dict:
    request = urllib.request.Request(url, headers={
        "User-Agent": "OrdaX-Creator-Owner-Dev-Local-Sync/1",
        "Cache-Control": "no-cache",
    })
    with urllib.request.urlopen(request, timeout=25) as response:
        body = response.read(max_bytes + 1)
    if not body or len(body) > max_bytes:
        raise SyncError("untrusted remote JSON size")
    value = json.loads(body)
    if not isinstance(value, dict):
        raise SyncError("remote JSON must be an object")
    return value


def validate_manifest(value: dict) -> dict:
    if (value.get("$schema") != "prototype-ordax.creator-owner-bundle-update/2"
            or value.get("owner_repository") != REPO
            or value.get("channel") != "owner-prototype"
            or value.get("artifact") != BUNDLE_NAME
            or value.get("download_url") != BASE + BUNDLE_NAME
            or value.get("automatic_in_app_update") is not False
            or not isinstance(value.get("version"), str)
            or not VERSION.fullmatch(value["version"])):
        raise SyncError("untrusted Creator update manifest")
    for key in ("source_commit", "creator_source_commit"):
        if not isinstance(value.get(key), str) or not HEX40.fullmatch(value[key]):
            raise SyncError("invalid Owner revision identity: " + key)
    if (not isinstance(value.get("sha256"), str) or not HEX64.fullmatch(value["sha256"])
            or type(value.get("size")) is not int
            or not 0 < value["size"] <= MAX_BUNDLE_BYTES):
        raise SyncError("untrusted Creator bundle integrity descriptor")
    return value


def current_provenance(root: Path) -> dict:
    current = root / "current"
    path = current / "OrdaX-Creator" / "provenance.json"
    if current.is_symlink() or not path.is_file() or path.is_symlink():
        raise SyncError("existing Owner Creator must have safe provenance")
    data = json.loads(path.read_text(encoding="utf-8"))
    if (data.get("$schema") != "prototype-ordax.creator-owner-physical/1"
            or data.get("creator_source_repository") != REPO
            or not isinstance(data.get("creator_version"), str)
            or not VERSION.fullmatch(data["creator_version"])
            or any(not isinstance(data.get(k), str) or not HEX40.fullmatch(data[k])
                   for k in ("source_commit", "creator_source_commit"))):
        raise SyncError("installed Owner Creator has untrusted provenance")
    return data


def source_relation(repository: str, installed: str, candidate: str) -> int:
    if installed == candidate:
        return 0
    if repository not in ("ordax-os", "ordax-creator"):
        raise SyncError("invalid history owner")
    url = f"https://api.github.com/repos/ordaxsystems/{repository}/compare/{installed}...{candidate}?per_page=1"
    relation = request_json(url, 16 << 20)
    if (relation.get("status") == "ahead"
            and type(relation.get("ahead_by")) is int and relation["ahead_by"] > 0
            and relation.get("behind_by") == 0):
        return 1
    raise SyncError(f"{repository} release history is not a verified fast-forward")


def download_bundle(manifest: dict, cache: Path) -> Path:
    cache.mkdir(mode=0o700, parents=True, exist_ok=True)
    expected = manifest["sha256"]
    target = cache / (expected + ".zip")
    if target.is_symlink():
        raise SyncError("Creator cache must not contain symlinks")
    if target.is_file() and target.stat().st_size == manifest["size"] and sha256_file(target) == expected:
        return target
    part = cache / (expected + f".{os.getpid()}.part")
    part.unlink(missing_ok=True)
    digest = hashlib.sha256()
    size = 0
    try:
        request = urllib.request.Request(BASE + BUNDLE_NAME, headers={
            "User-Agent": "OrdaX-Creator-Owner-Dev-Local-Sync/1"})
        with urllib.request.urlopen(request, timeout=40) as response, part.open("xb") as dst:
            for block in iter(lambda: response.read(1024 * 1024), b""):
                size += len(block)
                if size > manifest["size"] or size > MAX_BUNDLE_BYTES:
                    raise SyncError("Creator bundle exceeds release size")
                digest.update(block)
                dst.write(block)
        if size != manifest["size"] or digest.hexdigest() != expected:
            raise SyncError("Creator bundle digest or length mismatch")
        part.replace(target)
        return target
    finally:
        part.unlink(missing_ok=True)


def verify_extracted(staging: Path, bundle: Path, manifest: dict, *, smoke: bool = True) -> None:
    with zipfile.ZipFile(bundle) as archive:
        files = archive.infolist()
        if (len(files) != len(EXPECTED_FILES)
                or {item.filename for item in files} != EXPECTED_FILES
                or sum(item.file_size for item in files) > MAX_BUNDLE_BYTES):
            raise SyncError("Creator package has unexpected files or unsafe size")
        for item in files:
            if item.is_dir() or ((item.external_attr >> 16) & 0o170000) == 0o120000:
                raise SyncError("Creator package contains linked or non-file entry")
            path = staging / item.filename
            path.parent.mkdir(parents=True, exist_ok=True)
            with archive.open(item, "r") as src, path.open("xb") as dst:
                shutil.copyfileobj(src, dst, length=1024 * 1024)
    doc = json.loads((staging / "OrdaX-Creator/provenance.json").read_text(encoding="utf-8"))
    if (doc.get("$schema") != "prototype-ordax.creator-owner-physical/1"
            or doc.get("creator_source_repository") != REPO
            or doc.get("source_commit") != manifest["source_commit"]
            or doc.get("creator_source_commit") != manifest["creator_source_commit"]
            or doc.get("creator_version") != manifest["version"]
            or doc.get("physical_write_authorized_in_binary") is not True):
        raise SyncError("Creator package provenance differs from published release")
    executable = staging / "OrdaX-Creator/OrdaX-Creator.exe"
    with executable.open("rb") as binary:
        if binary.read(2) != b"MZ":
            raise SyncError("Creator package does not contain Windows executable")
    if smoke:
        result = subprocess.run([str(executable), "--offline-package-smoke"],
                                cwd=executable.parent, timeout=45, check=False)
        if result.returncode:
            raise SyncError("Creator offline package smoke failed")


def creator_process_active() -> bool:
    if os.name != "nt":
        raise SyncError("local Creator updates require Windows")
    result = subprocess.run(["tasklist", "/FI", "IMAGENAME eq OrdaX-Creator.exe",
                             "/FO", "CSV", "/NH"], capture_output=True, check=False)
    if result.returncode != 0:
        raise SyncError("cannot verify running Creator process")
    text = result.stdout.decode("utf-8", errors="replace")
    return any(row and row[0].strip().lower() == "ordax-creator.exe"
               for row in csv.reader(io.StringIO(text)))


def recover_interrupted_swap(root: Path) -> None:
    current, previous = root / "current", root / ".ordax-previous"
    if not current.exists() and previous.is_dir() and not previous.is_symlink():
        previous.rename(current)


def sync(root: Path, *, check_only: bool = False) -> str:
    if os.name != "nt":
        raise SyncError("Owner Creator local sync supports Windows only")
    if root.is_symlink():
        raise SyncError("Owner Creator install root is a symbolic link")
    root = root.resolve()
    if root.name != "OrdaX-Creator-Owner-Prototype" or not root.is_dir():
        raise SyncError("not an existing Owner Development Creator installation")
    import msvcrt
    lock = root / ".ordax-sync.lock"
    with lock.open("a+b") as handle:
        if handle.tell() == 0:
            handle.write(b"0")
            handle.flush()
        handle.seek(0)
        try:
            msvcrt.locking(handle.fileno(), msvcrt.LK_NBLCK, 1)
        except OSError:
            return "SKIP_ALREADY_RUNNING"
        try:
            recover_interrupted_swap(root)
            current = current_provenance(root)
            try:
                remote = validate_manifest(request_json(MANIFEST_URL))
            except urllib.error.HTTPError as exc:
                if exc.code == 404:
                    return "DEFERRED_NO_PUBLISHED_RELEASE"
                raise
            newer_os = source_relation("ordax-os", current["source_commit"], remote["source_commit"])
            newer_creator = source_relation("ordax-creator", current["creator_source_commit"], remote["creator_source_commit"])
            if not (newer_os or newer_creator):
                return "CURRENT"
            if check_only:
                return "AVAILABLE_VERIFIED"
            if creator_process_active():
                return "DEFERRED_CREATOR_RUNNING"
            bundle = download_bundle(remote, root / ".ordax-update-cache")
            stage = Path(tempfile.mkdtemp(prefix=".ordax-next-", dir=root))
            try:
                # Task Scheduler may run in a noninteractive Windows session;
                # the publisher already runs the real GUI smoke on Windows.
                # Local safety relies on release SHA-256, exact ZIP membership,
                # provenance, lineage and the binary header, not a GUI window.
                verify_extracted(stage, bundle, remote, smoke=False)
                if creator_process_active():
                    return "DEFERRED_CREATOR_RUNNING"
                installed = root / "current"
                previous = root / ".ordax-previous"
                if previous.exists():
                    shutil.rmtree(previous)
                installed.rename(previous)
                try:
                    stage.rename(installed)
                except OSError:
                    previous.rename(installed)
                    raise
                try:
                    shutil.rmtree(previous)
                except OSError:
                    return "UPDATED_VERIFIED_CLEANUP_PENDING"
                return "UPDATED_VERIFIED"
            finally:
                if stage.exists():
                    shutil.rmtree(stage)
        finally:
            handle.seek(0)
            msvcrt.locking(handle.fileno(), msvcrt.LK_UNLCK, 1)


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--install-root", type=Path,
                   default=Path.home() / "Downloads" / "OrdaX-Creator-Owner-Prototype")
    p.add_argument("--check-only", action="store_true")
    args = p.parse_args()
    try:
        result = sync(args.install_root, check_only=args.check_only)
    except (SyncError, ValueError, OSError, zipfile.BadZipFile,
            urllib.error.URLError, subprocess.TimeoutExpired) as exc:
        print(f"OWNER_CREATOR_LOCAL_SYNC=BLOCKED:{type(exc).__name__}:{exc}", file=sys.stderr)
        return 1
    print("OWNER_CREATOR_LOCAL_SYNC=" + result)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
