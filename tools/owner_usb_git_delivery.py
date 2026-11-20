#!/usr/bin/env python3
"""Developer-only Windows USB transport for the canonical OrdaX Git source.

This is not a boot updater, release signer or raw-disk writer. Its only
permitted destination is ORDAX-DATA on a physically validated USB disk.
"""
from __future__ import annotations

import argparse
from contextlib import contextmanager
import ctypes
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


@contextmanager
def exclusive_transport():
    """The GUI and the existing Windows task share one Owner USB writer."""
    if os.name != "nt":
        yield
        return
    import msvcrt
    root = Path.home() / "AppData" / "Local" / "OrdaX-Creator"
    root.mkdir(mode=0o700, parents=True, exist_ok=True)
    lock = root / "owner-usb-git.lock"
    if lock.is_symlink() or (hasattr(lock, "is_junction") and lock.is_junction()):
        raise DeliveryError("unsafe USB transport lock path")
    with lock.open("a+b") as handle:
        # LockFile permits a byte-range beyond EOF. Do not read or modify
        # the locked byte: Windows denies even same-process reads while it is
        # held by another descriptor. Only the kernel owns this mutex.
        handle.seek(0)
        try:
            msvcrt.locking(handle.fileno(), msvcrt.LK_NBLCK, 1)
        except OSError as exc:
            raise DeliveryError("Owner USB transport already running") from exc
        try:
            yield
        finally:
            handle.seek(0)
            msvcrt.locking(handle.fileno(), msvcrt.LK_UNLCK, 1)


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
    """Verify the exact Owner USB through read-only native Windows disk APIs.

    Storage CIM Get-Partition/Get-Disk may block on removable disks. The
    kernel layout and volume extents are authoritative and require no UAC.
    """
    if not re.fullmatch(r"[A-Z]", letter):
        raise DeliveryError("USB destination requires one drive letter")
    if os.name != "nt":
        raise DeliveryError("Owner USB detection requires Windows")
    import struct
    import uuid

    kernel = ctypes.WinDLL("kernel32", use_last_error=True)
    kernel.CreateFileW.argtypes = [
        ctypes.c_wchar_p, ctypes.c_uint32, ctypes.c_uint32, ctypes.c_void_p,
        ctypes.c_uint32, ctypes.c_uint32, ctypes.c_void_p,
    ]
    kernel.CreateFileW.restype = ctypes.c_void_p
    kernel.DeviceIoControl.argtypes = [
        ctypes.c_void_p, ctypes.c_uint32, ctypes.c_void_p, ctypes.c_uint32,
        ctypes.c_void_p, ctypes.c_uint32, ctypes.POINTER(ctypes.c_uint32),
        ctypes.c_void_p,
    ]
    kernel.CloseHandle.argtypes = [ctypes.c_void_p]

    def read_ioctl(path: str, command: int, payload: bytes = b"") -> bytes:
        handle = kernel.CreateFileW(path, 0, 3, None, 3, 0, None)
        if handle == ctypes.c_void_p(-1).value:
            raise OSError(ctypes.get_last_error(), "disk metadata inaccessible")
        try:
            output = ctypes.create_string_buffer(8192)
            used = ctypes.c_uint32()
            request = ctypes.create_string_buffer(payload) if payload else None
            if not kernel.DeviceIoControl(
                ctypes.c_void_p(handle), command, request, len(payload),
                output, len(output), ctypes.byref(used), None,
            ):
                raise OSError(ctypes.get_last_error(), "disk metadata query failed")
            return output.raw[:used.value]
        finally:
            kernel.CloseHandle(ctypes.c_void_p(handle))

    def extent(drive: str) -> tuple[int, int, int]:
        raw = read_ioctl(rf"\\.\{drive}:", 0x560000)
        if len(raw) < 32 or struct.unpack_from("<I", raw)[0] != 1:
            raise DeliveryError("USB data volume requires one physical disk extent")
        return struct.unpack_from("<I4xqq", raw, 8)

    try:
        label = ctypes.create_unicode_buffer(256)
        fs = ctypes.create_unicode_buffer(256)
        serial = ctypes.c_uint32()
        maximum = ctypes.c_uint32()
        flags = ctypes.c_uint32()
        if not kernel.GetVolumeInformationW(
            letter + ":\\", label, len(label), ctypes.byref(serial),
            ctypes.byref(maximum), ctypes.byref(flags), fs, len(fs),
        ):
            raise OSError(ctypes.get_last_error(), "USB volume unreadable")
        disk, start, size = extent(letter)
        system_drive = (os.environ.get("SystemDrive") or "C:")[0].upper()
        system_disk, _, _ = extent(system_drive)
        layout = read_ioctl(rf"\\.\PhysicalDrive{disk}", 0x70050)
        if len(layout) < 48:
            raise DeliveryError("USB GPT layout missing")
        style, count = struct.unpack_from("<II", layout)
        if style != 1 or count != 3 or len(layout) < 48 + 144 * count:
            raise DeliveryError("USB does not contain the canonical GPT layout")
        parts = []
        for index in range(count):
            offset = 48 + 144 * index
            partition_style, first, length, number = struct.unpack_from(
                "<I4xqqI", layout, offset,
            )
            guid = str(uuid.UUID(bytes_le=layout[offset + 32:offset + 48]))
            if partition_style != 1 or number != index + 1 or first < 0 or length <= 0:
                raise DeliveryError("USB GPT partition metadata invalid")
            parts.append((number, first, length, guid))

        # STORAGE_DEVICE_DESCRIPTOR: get the physical USB serial and bus
        # from IOCTL_STORAGE_QUERY_PROPERTY (12-byte aligned query).
        descriptor = read_ioctl(
            rf"\\.\PhysicalDrive{disk}", 0x2d1400,
            struct.pack("<III", 0, 0, 0),
        )
        if len(descriptor) < 36:
            raise DeliveryError("USB storage descriptor missing")
        serial_offset = struct.unpack_from("<I", descriptor, 24)[0]
        bus_type = struct.unpack_from("<I", descriptor, 28)[0]
        if not 36 <= serial_offset < len(descriptor):
            raise DeliveryError("USB device serial is absent")
        physical_serial = descriptor[serial_offset:].split(b"\0", 1)[0].decode(
            "ascii", "strict",
        ).strip()
        device = {
            "Index": disk,
            "InterfaceType": "USB" if bus_type == 7 else "OTHER",
            "SerialNumber": physical_serial,
        }
    except (OSError, ValueError, TypeError, IndexError, struct.error) as exc:
        raise DeliveryError("physical USB identity could not be verified") from exc

    required = {
        "c12a7328-f81f-11d2-ba4b-00a0c93ec93b",
        "0fc63daf-8483-4772-8e79-3d69d8477de4",
        "ebd0a0a2-b9e5-4433-87c0-68b6b72699c7",
    }
    types = [part[3] for part in parts]
    if (
        {part[3] for part in parts} != required
        or parts[2][:3] != (3, start, size)
        or parts[2][3] != "ebd0a0a2-b9e5-4433-87c0-68b6b72699c7"
        or disk == system_disk
        or not isinstance(device, dict)
        or device.get("Index") != disk
        or device.get("InterfaceType") != "USB"
        or not str(device.get("SerialNumber") or "").strip()
        or label.value != "ORDAX-DATA" or fs.value != "exFAT"
    ):
        raise DeliveryError("target is not an authorized OrdaX development USB data volume")
    return {
        "disk": disk, "serial": str(device["SerialNumber"]).strip(),
        "bus": "USB", "style": "GPT", "label": label.value, "fs": fs.value,
        "systemDisk": system_disk, "partition": 3, "types": types,
    }


def enforce_selected_disk(identity: dict, expected: int | None) -> None:
    if expected is not None and (type(expected) is not int or expected < 0 or
                                 type(identity.get("disk")) is not int or identity["disk"] != expected):
        raise DeliveryError("USB disk number differs from the explicit Creator selection")


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


def available_ordax_data_letters() -> list[str]:
    """Read Windows volume labels cheaply; checked_volume owns disk identity."""
    if os.name != "nt":
        raise DeliveryError("automatic volume discovery requires Windows")
    from ctypes import wintypes
    kernel = ctypes.WinDLL("kernel32", use_last_error=True)
    logical = kernel.GetLogicalDrives
    logical.restype = wintypes.DWORD
    volume_info = kernel.GetVolumeInformationW
    volume_info.argtypes = [
        wintypes.LPCWSTR, wintypes.LPWSTR, wintypes.DWORD,
        ctypes.POINTER(wintypes.DWORD), ctypes.POINTER(wintypes.DWORD),
        ctypes.POINTER(wintypes.DWORD), wintypes.LPWSTR, wintypes.DWORD]
    volume_info.restype = wintypes.BOOL
    drive_mask = logical()
    if not drive_mask:
        raise DeliveryError("cannot list Windows logical drives")
    matches = []
    for index in range(26):
        if not (drive_mask & (1 << index)):
            continue
        drive = chr(65 + index)
        path = drive + ":\\"
        label = ctypes.create_unicode_buffer(256)
        filesystem = ctypes.create_unicode_buffer(256)
        volume_serial = wintypes.DWORD()
        max_component = wintypes.DWORD()
        flags = wintypes.DWORD()
        if volume_info(path, label, len(label), ctypes.byref(volume_serial),
                       ctypes.byref(max_component), ctypes.byref(flags),
                       filesystem, len(filesystem)):
            if label.value == "ORDAX-DATA" and filesystem.value.lower() == "exfat":
                matches.append(drive)
    return matches


def discover_trusted_volume(expected_serial: str) -> str | None:
    """Find one locally authorized USB, never by drive letter or label alone."""
    if not re.fullmatch(r"[A-Za-z0-9._-]{4,64}", expected_serial):
        raise DeliveryError("authorized physical USB serial is malformed")
    matches = []
    for letter in available_ordax_data_letters():
        if not isinstance(letter, str) or not re.fullmatch(r"[A-Z]", letter):
            raise DeliveryError("malformed USB drive letter")
        try:
            current = checked_volume(letter)
        except DeliveryError:
            continue
        if str(current["serial"]).strip().casefold() == expected_serial.casefold():
            matches.append(letter)
    if len(matches) > 1:
        raise DeliveryError("multiple volumes claim the authorized USB serial")
    return matches[0] if matches else None


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo", type=Path, required=True)
    destination = parser.add_mutually_exclusive_group(required=True)
    destination.add_argument("--volume", help="ORDAX-DATA drive letter, e.g. E")
    destination.add_argument("--auto-volume", action="store_true",
                             help="discover only the locally authorized physical USB")
    parser.add_argument("--trusted-serial", help="USB serial selected locally, never in Git")
    parser.add_argument("--trusted-disk-number", type=int,
                        help="also require selected Windows physical disk number")
    parser.add_argument("--apply", action="store_true", help="stage only on verified USB data volume")
    args = parser.parse_args()
    try:
        with exclusive_transport():
            return deliver(args)
    except DeliveryError as exc:
        print("OWNER_USB_GIT_DELIVERY=BLOCKED:" + str(exc), file=sys.stderr)
        return 2


def deliver(args: argparse.Namespace) -> int:
    try:
        if args.apply and not args.trusted_serial:
            raise DeliveryError("staging requires explicit locally authorized USB serial")
        if args.auto_volume:
            if not args.trusted_serial:
                raise DeliveryError("automatic USB discovery requires trusted serial")
            letter = discover_trusted_volume(args.trusted_serial)
            if letter is None:
                print(json.dumps({"status": "SKIP_TRUSTED_USB_NOT_PRESENT"}))
                return 0
        else:
            letter = args.volume.strip().upper().rstrip(":\\")
        before = checked_volume(letter)
        enforce_selected_disk(before, args.trusted_disk_number)
        if args.trusted_serial and str(before["serial"]).strip().casefold() != args.trusted_serial.casefold():
            raise DeliveryError("USB serial differs from locally authorized device")
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
