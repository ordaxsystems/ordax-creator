#!/usr/bin/env python3
"""Developer-only Windows USB transport for the canonical OrdaX Git source.

This is not a boot updater, release signer or raw-disk writer. Its only
permitted destination is ORDAX-DATA on a physically validated USB disk.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile

REPO_URL = "https://github.com/ordaxsystems/ordax-os.git"
SCHEMA = "ordax.owner-dev-usb-git/1"
SHA = re.compile(r"[0-9a-f]{40}\Z")
DIGEST = re.compile(r"[0-9a-f]{64}\Z")
MAX_BUNDLE = 1_000_000_000
SUBDIR = Path(".ordax") / "updates" / "owner-git"
BUNDLE = "ordax-main.bundle"
MANIFEST = "owner-git-bundle.txt"


class DeliveryError(RuntimeError):
    pass


def git(repo: Path, *args: str, timeout=90) -> str:
    try:
        result = subprocess.run(["git", "-C", str(repo), *args], text=True,
                                capture_output=True, check=True, timeout=timeout)
    except (subprocess.CalledProcessError, FileNotFoundError, subprocess.TimeoutExpired) as exc:
        raise DeliveryError("Git source or bundle verification failed") from exc
    return result.stdout.strip()


def source_commit(repo: Path) -> str:
    if repo.is_symlink() or not (repo / ".git").is_dir():
        raise DeliveryError("canonical OS checkout missing or linked")
    if git(repo, "remote", "get-url", "origin") != REPO_URL:
        raise DeliveryError("OS Git origin is not canonical")
    if git(repo, "branch", "--show-current") != "main":
        raise DeliveryError("OS Git checkout must be on main")
    if git(repo, "status", "--porcelain=v1"):
        raise DeliveryError("OS Git checkout contains pending edits")
    head = git(repo, "rev-parse", "HEAD")
    if not SHA.fullmatch(head):
        raise DeliveryError("OS Git HEAD is malformed")
    remote = git(repo, "ls-remote", "--heads", "origin", "refs/heads/main", timeout=40)
    if remote != head + "\trefs/heads/main":
        raise DeliveryError("OS checkout differs from current remote main")
    return head


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def manifest_bytes(commit: str, sha256: str, size: int) -> bytes:
    if not SHA.fullmatch(commit) or not DIGEST.fullmatch(sha256) or not 0 < size <= MAX_BUNDLE:
        raise DeliveryError("invalid offline Git provenance")
    return ("schema=" + SCHEMA + "\nrepository=ordaxsystems/ordax-os\n"
            + "source_commit=" + commit + "\n"
            + "bundle_sha256=" + sha256 + "\n"
            + "bundle_bytes=" + str(size) + "\n").encode("ascii")


def existing_current(target: Path, commit: str) -> bool:
    manifest = target / MANIFEST
    bundle = target / BUNDLE
    if not manifest.is_file() or not bundle.is_file() or manifest.is_symlink() or bundle.is_symlink():
        return False
    try:
        lines = manifest.read_bytes().splitlines(keepends=True)
        if len(lines) != 5 or not all(line.endswith(b"\n") and b"\r" not in line for line in lines):
            return False
        size = bundle.stat().st_size
        digest = sha256_file(bundle)
        return manifest.read_bytes() == manifest_bytes(commit, digest, size)
    except (OSError, UnicodeError, DeliveryError):
        return False


def build_bundle(repo: Path, destination: Path, commit: str) -> tuple[str, int]:
    git(repo, "bundle", "create", str(destination), "main", timeout=180)
    heads = git(repo, "bundle", "list-heads", str(destination))
    if heads != commit + " refs/heads/main":
        raise DeliveryError("bundle is not bound to one exact main revision")
    git(repo, "bundle", "verify", str(destination))
    size = destination.stat().st_size
    if size <= 0 or size > MAX_BUNDLE:
        raise DeliveryError("bundle exceeds bounded ORDAX-DATA budget")
    return sha256_file(destination), size


def checked_volume(letter: str) -> dict:
    if not re.fullmatch(r"[A-Z]", letter):
        raise DeliveryError("USB destination requires one drive letter")
    powershell = """
$ErrorActionPreference = 'Stop'
$p=Get-Partition -DriveLetter '{letter}'
$d=Get-Disk -Number $p.DiskNumber
$v=Get-Volume -DriveLetter '{letter}'
$s=Get-Partition -DriveLetter ($env:SystemDrive.Substring(0,1))
$types=@(Get-Partition -DiskNumber $d.Number | ForEach-Object {{ $_.GptType.ToString().Trim('{{}}').ToLowerInvariant() }})
[pscustomobject]@{{ disk=$d.Number; serial=$d.SerialNumber; bus=$d.BusType.ToString();
  style=$d.PartitionStyle.ToString(); label=$v.FileSystemLabel; fs=$v.FileSystem;
  systemDisk=$s.DiskNumber; partition=$p.PartitionNumber; types=$types }} | ConvertTo-Json -Compress
""".format(letter=letter)
    try:
        raw = subprocess.run(["powershell.exe", "-NoProfile", "-NonInteractive", "-Command", powershell],
                             text=True, capture_output=True, timeout=25, check=True).stdout
        value = json.loads(raw)
    except (subprocess.CalledProcessError, FileNotFoundError, subprocess.TimeoutExpired, ValueError) as exc:
        raise DeliveryError("physical USB identity could not be verified") from exc
    required = {
        "c12a7328-f81f-11d2-ba4b-00a0c93ec93b",
        "0fc63daf-8483-4772-8e79-3d69d8477de4",
        "ebd0a0a2-b9e5-4433-87c0-68b6b72699c7",
    }
    if (value.get("bus") != "USB" or value.get("style") != "GPT" or
        value.get("fs") != "exFAT" or value.get("label") != "ORDAX-DATA" or
        value.get("disk") == value.get("systemDisk") or
        not str(value.get("serial", "")).strip() or
        not required.issubset(set(value.get("types", [])))):
        raise DeliveryError("target is not an authorized OrdaX development USB data volume")
    return value


def ensure_real_ancestors(volume: Path, target: Path) -> None:
    current = volume
    for part in target.relative_to(volume).parts:
        current = current / part
        if current.exists() and (current.is_symlink() or
                                 (hasattr(current, "is_junction") and current.is_junction()) or
                                 not current.is_dir()):
            raise DeliveryError("USB update path is linked or not a directory")


def install_bundle(target: Path, bundle: Path, commit: str, digest: str, size: int) -> None:
    # Save manifest last: any interruption before that leaves mismatched data
    # that the OS verifier rejects. No raw block device is opened or formatted.
    if shutil.disk_usage(target.anchor).free < 2 * size + 32 * 1024 * 1024:
        raise DeliveryError("insufficient ORDAX-DATA free space for atomic transfer")
    target.mkdir(parents=True, exist_ok=True)
    temp_bundle = target / (BUNDLE + ".pending")
    temp_manifest = target / (MANIFEST + ".pending")
    for pending in (temp_bundle, temp_manifest):
        if pending.is_symlink() or (hasattr(pending, "is_junction") and pending.is_junction()):
            raise DeliveryError("unsafe USB pending path")
    try:
        with bundle.open("rb") as source, temp_bundle.open("wb") as dest:
            shutil.copyfileobj(source, dest, length=1024 * 1024)
            dest.flush()
            os.fsync(dest.fileno())
        if temp_bundle.stat().st_size != size or sha256_file(temp_bundle) != digest:
            raise DeliveryError("USB readback of Git bundle failed")
        with temp_manifest.open("wb") as dest:
            dest.write(manifest_bytes(commit, digest, size))
            dest.flush()
            os.fsync(dest.fileno())
        os.replace(temp_bundle, target / BUNDLE)
        os.replace(temp_manifest, target / MANIFEST)
        if not existing_current(target, commit):
            raise DeliveryError("final USB bundle readback failed")
    finally:
        temp_bundle.unlink(missing_ok=True)
        temp_manifest.unlink(missing_ok=True)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--volume", required=True, help="ORDAX-DATA drive letter, e.g. E")
    parser.add_argument("--apply", action="store_true", help="stage only on verified USB data volume")
    args = parser.parse_args()
    try:
        letter = args.volume.strip().upper().rstrip(":\\")
        before = checked_volume(letter)
        volume = Path(letter + ":\\")
        target = volume / SUBDIR
        ensure_real_ancestors(volume, target)
        commit = source_commit(args.repo.resolve())
        if existing_current(target, commit):
            print(json.dumps({"status": "CURRENT", "source_commit": commit}))
            return 0
        if not args.apply:
            print(json.dumps({"status": "AVAILABLE", "source_commit": commit, "action": "staging-not-authorized"}))
            return 0
        with tempfile.TemporaryDirectory(prefix="ordax-usb-bundle-") as temp:
            bundle = Path(temp) / BUNDLE
            digest, size = build_bundle(args.repo.resolve(), bundle, commit)
            if checked_volume(letter) != before:
                raise DeliveryError("physical USB identity changed during packaging")
            install_bundle(target, bundle, commit, digest, size)
            if checked_volume(letter) != before:
                raise DeliveryError("physical USB identity changed after transfer")
            print(json.dumps({"status": "STAGED_NOT_ACTIVATED", "source_commit": commit,
                              "sha256": digest, "bytes": size}))
        return 0
    except DeliveryError as exc:
        print("OWNER_USB_GIT_DELIVERY=BLOCKED:" + str(exc), file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
