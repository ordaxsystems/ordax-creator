# OrdaX Creator

Official standalone source repository for **OrdaX Creator**, the OrdaX verified-media preparation application and reusable host-neutral Creator Core.

> Migration in progress. This repository is **not yet an operational, signed or public-physical-write-capable Creator**. Its OS implementation currently resides in [ordaxsystems/ordax-os/tools/creator](https://github.com/ordaxsystems/ordax-os/tree/main/tools/creator). Do not assume this new repository contains production-capable executables until the source history and independent CI have been migrated and proved.

## Ownership

- **This repository** will own the Creator Core, Windows UI, removable-device adapters, exact-target verification and consent implementation, post-write readback and verification consumers.
- **[ordax-os](https://github.com/ordaxsystems/ordax-os)** remains the sole authority for Base, kernel, bootstrap, official OS release artifacts, boot/storage/media layouts, release-producer contracts and physical release promotion. The Creator consumes only immutable, verified producer inputs; it must not generate a competing image, artifact catalog, updater, signing root or physical authorization policy.
- The OrdaX Desktop will reuse this Creator Core through an explicit supported integration; it must not implement a second flashing backend.

Read [AGENTS.md](AGENTS.md), [docs/MIGRATION.md](docs/MIGRATION.md) and [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) before making changes. The new repository is in **PREPARATION** status; none of its files are a substitute for the canonical OS machine-readable contracts.

### Safety

**No disk/USB destructive operation, physical release approval, trust rotation or signing key operation is authorized by this repository's creation or its CI.** Use regular disposable files and simulated or QEMU devices only. Hardware writes require separate explicit per-target approval and all publisher/consumer verification gates.

### Source provenance

The source baseline for auditing is `ordax-os/main@f65acc4fb9f8ded376fc9b44b4c7c2068c793265` (2026-10-23). The open [ordax-os#2000](https://github.com/ordaxsystems/ordax-os/pull/2000) contains independently valuable Core source-binding work and must be reconciled, not discarded. This initial README **is not a source import**; the migration ledger must record exact original, split and imported Git commits before any cutover.
