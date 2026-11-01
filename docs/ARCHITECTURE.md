# OrdaX Creator — canonical architecture and safety boundaries

**Implementation status:** C5/C6 source extraction and OS consumer cutover are complete (2026-10-24). Independent Creator Core and Windows application live in this repository. Distribution as a publicly signed release and authorization to write a physical USB are **separate, uncompleted OS publisher gates**. See [docs/MIGRATION.md](MIGRATION.md) for the migration execution ledger, [docs/PROVENANCE.md](PROVENANCE.md) for immutable history receipts and [AGENTS.md](../AGENTS.md) for agent constraints.

## Single source of truth by owner

| Concern | Sole authority |
| --- | --- |
| Creator Go Core, safe Windows UI, device adapters, verified-source consumers and narrowly gated physical writer | `ordaxsystems/ordax-creator` |
| OS Base/kernel/boot, production media and storage geometry, physical release authorization, published artifacts, signer/trust roots | `ordaxsystems/ordax-os` |
| Immutable OS source/contract binding used in Creator CI | `.github/ordax-os-contracts.lock.json`, pointing to the exact upstream OS commit (read-only test input, not duplicated contracts) |
| Product version | `version.json` (the exact Creator Git SHA separately identifies each build) |

The OS's verified release consumers bind an immutable Creator source commit and tree; the OS may produce/publish outputs but must not mirror the Go Core, create a second image builder or maintain a parallel 17-source inventory or release signer.

## Runtime boundaries

1. **Core and UI:** `core/`, `cmd/ordax-creator-app/` and `host/` implement the single source of target selection, plans and readback. Windows device identity, system-disk exclusion, serial/capacity revalidation, confirmation and UAC remain separate fail-closed gates.
2. **Official signed consumer channels:** `appchannel/`, `componentchannel/` and `physicalchannel/` verify their purpose-specific signed envelopes, the canonical publisher trust binding, file size and SHA-256, and rollback protections before accepting distribution. A development build with unresolved trust must not silently promote a signed or physical capability.
3. **Development-only updater:** `update/` consumes the OS-owned `creator-dev` manifest and a size/SHA-256-bound ZIP. It supports read-only unprivileged inspection, uses `internal/atomicfile/` for its installation pointer and rejects unsafe archive paths/slot references. This **unsigned** channel is not a substitute for signed release or a hardware-writing authorization.
4. **Physical writer:** the narrow Windows RAW backend is isolated behind an explicit build tag, authorized publisher binding and independent device-specific consent. Ordinary Creator GUI/CLI binaries and CI candidates do not include an authorized RAW writer.
5. **OS integration:** `ordax-os` owns the authoritative media producers, boot proof, signing and promotion policy. Its existing consumers reference the independent Creator source; there is no embedded duplicate Creator Go implementation in the OS after C6.

## Evidence and release gates

- Core/Go vet, immutable read-only contract integration, real Win32 UI smoke, unsafe-target negative tests, repeatable unsigned ZIP builds, and synthetic/disposable file tests are non-destructive engineering evidence.
- CI artifacts must display their Creator commit and explicit **UNSIGNED/READONLY** status. They are not Stable artifacts, signed releases or proof of physical USB writing.
- The official publisher must separately complete source and signature verification, Authenticode/RFC3161 custody, signed channels, target-specific authorization and physical readback proof. The Creator repository must not create another signer, trust root, media producer or consent authority.
- No test, document or merge here grants the right to enumerate, repartition, eject, lock or write a real device. Device-specific physical authorization follows the OS-owned procedure in [AGENTS.md](../AGENTS.md).

**Current status authority:** [docs/MIGRATION.md](MIGRATION.md) records historical cutover evidence; [issue #5](https://github.com/ordaxsystems/ordax-creator/issues/5) tracks only outstanding distribution/physical proof, not a second architecture.
