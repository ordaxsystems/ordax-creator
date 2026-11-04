# AGENTS.md — OrdaX Creator

> **Execução DEV — Evolução Rápida (regra transversal obrigatória):** aplicar a [política canônica OrdaX](https://github.com/ordaxsystems/ordax-os/blob/main/docs/DEVELOPMENT-WORKFLOW.md#regra-operacional-transversal--evolucao-rapida-dev-2026-11-04) antes de iniciar trabalho. Entregar a alteração no owner correto primeiro e validar com testes proporcionais depois; não ampliar escopo com auditorias, documentação, CI ou workflows sem causalidade. Preservar gates obrigatórios, integridade, segurança e autorização.

Mandatory instructions for agents and maintainers.

## Identity and source authority
- This repository is the sole implementation owner of the independent Creator Core, Windows UI, device adapters, signed-consumer channels, narrow privileged writer, device-target revalidation and readback.
- ordaxsystems/ordax-os remains the sole producer/owner of Base, kernel, boot, official OS artifacts, media/storage contracts, publisher trust, physical-promotion policy and release signer.
- No alternative image builder, artifact catalog, 17-artifact inventory, updater, signature root or destructive authorization mechanism may be introduced. Contract downloads are pinned immutable, read-only inputs, never duplicated authority.
- C6 implementation extraction and OS consumer cutover are complete: OS #2222 removed its 154 embedded Go/module files; OS #2000 integrated independent 17-source bindings. OS-owned Python media producers and physical authorization remain in `ordax-os/tools/creator/`. Signed Stable production and device-specific physical consent are separate unapproved OS gates; see `docs/MIGRATION.md`.

## Reading proportional to the change

For a local Creator edit, read this `AGENTS.md`, the affected source and only the relevant contracts. **Do not read the full Creator and OS documentation set before a minor UI, text or isolated code change. Implement first; then run focused existing checks.**

- **Migration, source locks, imported provenance or OS consumer handoff:** consult `docs/MIGRATION.md`, `docs/PROVENANCE.md` and the exact pinned OS source/contracts affected.
- **Core architecture, ownership or component boundaries:** consult `docs/ARCHITECTURE.md` and relevant parts of `README.md`.
- **Physical writes, media geometry, permission/trust, signing, update or production promotion:** perform the necessary preflights and read the applicable OS documents at the **exact verified revision** (including `ordax-os/AGENTS.md`, `MVP.md`, `docs/CURRENT-STATE.md`, `docs/CREATOR-INSTALLATION.md`, `docs/PROMOTION-GATES.md` and machine-readable contracts as applicable). Never infer device authorization from general DEV permission.

Do not create plans, evidence files, CI workflows or full USB/Windows rebuilds solely to satisfy a reading checklist. Existing safety, trust and release gates remain mandatory when their domain is affected.

## Coordination
- Surface agent may change OS main; the OS Foundation work is paused. Keep this project's work Creator-scoped.
- Create Creator branches and PRs here. Proposed OS modifications require separate Creator-only PRs; never merge those before reconciling with the Surface agent and updated OS main.
- Preserve the merged `ordax-os#2000` 17-source source-binding contract and the single independent Creator Core; reconcile concurrent OS branches before any further cross-repository change.
- No force push, history rewriting of either main, or squashing published OS main ancestry.

## Physical-safety invariants
- Do not enumerate, open for writing, lock, eject, repartition, format or write a real USB or internal disk without separate explicit per-target authorization, regardless of broad software-development permission.
- CI may use synthetic devices, ordinary disposable files and QEMU only. Never promote passing CI to release or real-device approval.
- Reject unknown/untrusted internal system disks, re-enumeration/identity/capacity/serial changes, unsafe paths/symlinks, SHA-256/size or signature mismatches, stale signed release identities, insufficient space, missing confirmation or denied UAC. Flush and independently read back after approved writes.
- Privileged physical writer remains excluded from the ordinary GUI/CLI build unless an explicit publisher and build-tag boundary is passed.
- Private signing keys, credentials, physical data and secrets must not enter Git, Creator payloads, CI artifacts or debug logs.

## Engineering and gates
- Preserve the completed Creator Git subtree ancestry and immutable provenance receipt; do not recreate or reimport the removed OS Go module.
- The target Go module is github.com/ordaxsystems/ordax-creator. Do not create empty parallel Core scaffolding.
- Keep a single shared Core for the present 17 Portable sources, 39-operation application plan, safe Windows UI and future Desktop integration. Retain appchannel, componentchannel, physicalchannel and a single atomic pointer implementation.
- Make builds reproducible before code signing; verify signed channels and Authenticode/RFC3161 only when legitimately configured. A candidate executable is not a published signed release.
- Test Go, Windows cross-compilation, negative security cases, offline signed-cache integrity and disposable media before changing an OS consumer or Creator source pin.
- Preserve the proven zero-duplicate-Go OS boundary; do not delete OS-owned Python producer, signer or consent code. Update `docs/MIGRATION.md` as the single Creator migration handoff and never treat source extraction CI as physical-media or Stable release authorization.
