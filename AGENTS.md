# AGENTS.md — OrdaX Creator

Mandatory instructions for agents and maintainers.

## Identity and source authority
- This repository will be the sole implementation owner of the independent Creator Core, Windows UI, device adapters, official signed-consumer channels, narrow privileged writer, device-target revalidation and readback.
- ordaxsystems/ordax-os remains the sole producer/owner of Base, kernel, boot, official OS artifacts, media/storage contracts, publisher trust, physical-promotion policy and release signer.
- No alternative image builder, artifact catalog, 17-artifact inventory, updater, signature root or destructive authorization mechanism may be introduced. Contract downloads are pinned immutable, read-only inputs, never duplicated authority.
- During migration, the original implemented Creator remains in ordax-os/tools/creator. This new repository does not become a functional Creator until source history, independent tests and OS consumer cutover are proven.

## Mandatory reading
1. README.md
2. docs/MIGRATION.md
3. docs/ARCHITECTURE.md
4. docs/PROVENANCE.md
5. At an exact OS commit: ordax-os/AGENTS.md, MVP.md, docs/README.md, docs/CURRENT-STATE.md, docs/CREATOR-INSTALLATION.md, docs/PROMOTION-GATES.md and the relevant machine-readable OS contracts.

## Coordination
- Surface agent may change OS main; the OS Foundation work is paused. Keep this project's work Creator-scoped.
- Create Creator branches and PRs here. Proposed OS modifications require separate Creator-only PRs; never merge those before reconciling with the Surface agent and updated OS main.
- Preserve open ordax-os#2000 and its exact 17-source single-Core refactor; do not override competing agents' branches.
- No force push, history rewriting of either main, or squashing published OS main ancestry.

## Physical-safety invariants
- Do not enumerate, open for writing, lock, eject, repartition, format or write a real USB or internal disk without separate explicit per-target authorization, regardless of broad software-development permission.
- CI may use synthetic devices, ordinary disposable files and QEMU only. Never promote passing CI to release or real-device approval.
- Reject unknown/untrusted internal system disks, re-enumeration/identity/capacity/serial changes, unsafe paths/symlinks, SHA-256/size or signature mismatches, stale signed release identities, insufficient space, missing confirmation or denied UAC. Flush and independently read back after approved writes.
- Privileged physical writer remains excluded from the ordinary GUI/CLI build unless an explicit publisher and build-tag boundary is passed.
- Private signing keys, credentials, physical data and secrets must not enter Git, Creator payloads, CI artifacts or debug logs.

## Engineering and gates
- Preserve existing Creator Go implementation/history before moving it; use a filtered subtree import with a true Git ancestry link and an immutable provenance receipt.
- The target Go module is github.com/ordaxsystems/ordax-creator. Do not create empty parallel Core scaffolding.
- Keep a single shared Core for the present 17 Portable sources, 39-operation application plan, safe Windows UI and future Desktop integration. Retain appchannel, componentchannel, physicalchannel and a single atomic pointer implementation.
- Make builds reproducible before code signing; verify signed channels and Authenticode/RFC3161 only when legitimately configured. A candidate executable is not a published signed release.
- Test Go, Windows cross-compilation, negative security cases, offline signed-cache integrity and disposable media before any OS cutover.
- Never delete old OS code until standalone CI, contract equivalence, OS integration and reference audits prove replacement. Update docs/MIGRATION.md as the single Creator migration handoff and never claim unverified success.
