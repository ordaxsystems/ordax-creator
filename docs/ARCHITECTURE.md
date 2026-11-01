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
3. **Development-only updater:** `update/` consumes the OS-owned `creator-dev` manifest and a size/SHA-256-bound ZIP. It supports read-only unprivileged inspection, uses `internal/atomicfile/` for its installation pointer and rejects unsafe archive paths/slot references. Its local receipt records each extracted file's size and SHA-256 and verifies all files before offline reuse; a legacy cache can only be adopted online after byte-for-byte comparison with a newly downloaded verified ZIP. This **unsigned** local integrity receipt is not publisher authentication, a signed release or a hardware-writing authorization.
4. **Physical writer:** the narrow Windows RAW backend is isolated behind an explicit build tag, authorized publisher binding and independent device-specific consent. Ordinary Creator GUI/CLI binaries and CI candidates do not include an authorized RAW writer.
5. **OS integration:** `ordax-os` owns the authoritative media producers, boot proof, signing and promotion policy. Its existing consumers reference the independent Creator source; there is no embedded duplicate Creator Go implementation in the OS after C6.

## Shared cache directory policy

All Creator package consumers (Dev, App, Component and Physical) use `internal/cachepath` as the **single implementation** of the `versions/` cache-root boundary. Offline reads require an existing real directory; installers may create it but then independently verify it is not a symbolic link or a regular file. Each channel retains its own separate signature, artifact-integrity, rollback and authorization policy; the shared path check grants no physical-write authority.

## Windows owner-bundle update presentation

The independent Creator GUI's owner-prototype update check is **read-only**. It validates the OS-published owner-update manifest (channel, exact pinned release URL, version, SHA-256, source commit, declared bundle size and no automatic-install permission) and compares the source lineage before displaying availability. The historical 512 MiB manifest-size ceiling was below the OS-published full owner bundle; the validation bound is now **2 GiB**, without downloading, executing or promoting an artifact. All other checks remain unchanged.

The main GUI version label always identifies the **installed** build and immutable short source commit. An update notice belongs to the dedicated button; the manual instructions show **installed and available** versions/commits side by side. The candidate never becomes the displayed installed version until a separate, legitimate external bundle replacement has occurred. The update remains manual, not a signed Stable release or physical-write authorization. Windows regression tests cover the published-size class, invalid manifests, both locales and installed-version preservation.

The update button reflects four mutually exclusive **check** outcomes rather than an installation attempt: searching (disabled), failed check (retry available), successfully checked with no newer version, or a newer version available for manual replacement. A network or lineage failure never asserts that the installed build is current. A result from an already-started manual check must not open a modal dialog during physical writing or USB discovery; the check result remains on the button for later review. The installed build and physical authorization are unaffected by every check outcome. Regression tests exercise the state policy and real owner-drawn Win32 button captions, using only synthetic state and no network or device I/O.

## Windows refresh UI interlock

The Win32 GUI has one UI-thread-owned refresh-in-progress gate. From the moment a refresh starts until its newly verified snapshot is rendered, another refresh or a physical write cannot start. USB selection and the destructive action remain disabled, including on locale changes and queued native button messages. A locale change continues to show the translated searching state rather than redisplaying stale targets. Queued progress/completion events from an earlier operation cannot reenable the previous device while the new refresh is active.

This presentation-level interlock complements, but **never replaces**, the signed physical channel, device-specific confirmation and independent privileged revalidation. The Windows native smoke uses only synthetic in-memory device data and intentionally starts no discovery, update, UAC or disk I/O.

## Windows physical progress evidence

The Win32 progress adapter treats the elevated helper's `prototype-ordax.creator-physical-progress/1` snapshot as **bounded, untrusted display evidence**, not permission to write or proof of a successful installation. Each snapshot is read through the shared `internal/safefile.ReadBoundedRegular` descriptor with its 4 KiB cap and exact schema. Byte counters must be nonnegative and satisfy `completed_bytes <= total_bytes`; in particular, nonzero completed bytes with zero total are invalid. Unknown totals show the beginning of the current phase rather than a fabricated fraction.

The phase percentage uses an **exact overflow-free 128-bit multiplication/division** of valid 64-bit byte counters, preserving the existing phase ranges and avoiding signed `int64` overflow for large images. Tests cover the maximum signed byte counter with an independent arbitrary-precision reference, near-complete values and invalid snapshots. The privileged writer's independent consent, revalidation and success/failure exit checks are unchanged; visual progress cannot bypass them.

## Windows USB selection across refreshes

Device discovery never implicitly selects a physical USB, even if only one target appears. The user must make an explicit selection before the GUI can expose the destructive confirmation step. A subsequent refresh may retain that explicit choice only when there is exactly one newly discovered safe target with the same non-empty device serial, disk number, physical capacity and lowercase physical confirmation token. A changed drive letter or volume label alone does not erase a provably stable choice. Missing serial, changed identity, duplicate candidates, a fatal refresh error or unready physical channel clears the selection and requires the user to choose again.

This is a **presentation-only** convenience, not authorization: the elevated writer still independently validates a freshly enumerated target, signed sources, consent and the exact physical identity before writing. The existing target policy remains owned by the Windows adapter and Core; the GUI introduces no second token derivation.

## Offline cache fallback and physical readiness

Windows refresh distinguishes **update availability** from **validity of installed content**. If the latest App/Component update cannot be acquired but that channel's existing `Current()` verifier accepts a complete local installation, the GUI shows a localized nonblocking cache notice instead of an error. Development caches are explicitly described as local **development** copies, never signed releases; signed Component caches are described as signed only after their existing signature validation succeeds.

If no verified cache is available, refresh retains a fatal error. Independent physical-channel signature/trust, target enumeration, `PhysicalReady`, user confirmation and privileged-device revalidation remain separate requirements; neither a development cache nor an informational notice confers physical-write authority. Windows-only regression cases cover verified fallback, absent/invalid cache, physical readiness and fatal target failures.

## Windows design assets and first launch

The Windows Creator consumes the theme JSON and symbol PNG exported by the Surface design source of truth, with their matching bundle `provenance.json`. The GUI preflight uses the existing `internal/safefile.ReadBoundedRegular` for all three local files (64 KiB theme, 2 MiB symbol, 64 KiB provenance), rejecting oversized, empty, non-regular or symlinked entries **before** allocating asset contents. The existing SHA-256/provenance checks remain authoritative for the exact bundle identity; this read guard is not a new palette owner, signer or release trust root.

## Immutable signed release identity

`internal/releaseidentity` is the single rollback/source-identity policy for the three signed channels, including the App pending/health-check workflow. Each previously verified source commit binds exactly one signed manifest and its existing positive release sequence. A source commit cannot be silently reused with another signed artifact, bundle, Portable payload, version or release sequence, even if the new envelope bears a valid signature. A different source commit may only advance the authenticated sequence. Unchanged manifests can still be verified and reused without downloading their artifacts; a valid new commit is not impeded. The consumer compares only **already cryptographically verified** manifests. App `PromotePending` additionally requires the **exact signed manifest returned by `AcquirePending` and subjected to the launcher health check**; a matching source commit alone cannot authorize a pending pointer that was replaced with another validly signed manifest. A mismatch aborts before touching `current` or `previous`, preserving the independently verified last-known-good fallback. No new signer, publication permission or physical authorization is introduced.

## Signed Component repair without rollback reset

The Component channel uses its already verified `current-envelope.json` as the **sole authenticated rollback/immutable release-identity baseline**, even when its cached executable bytes are missing or corrupted. A different signed source commit with strictly greater release sequence can install a fresh verified component in a new immutable slot and activate it only after the existing independent second-envelope read, signature verification and installed-file SHA-256 checks. `Current` still requires the actual component to exist and match its signed hash; corrupted executables are **never** returned as runnable. An invalid or missing-signed-evidence substitution cannot silently reset a published sequence: missing current permits initial installation, but malformed, untrusted, or symlinked current envelopes block all rollback-sensitive updates. Same-source rebinding, older source sequences, and sequence reuse remain prohibited by `internal/releaseidentity`. This is recovery for the signed Component channel only; physical authorization and provenance have independent stricter gates.

## Offline signed-evidence reads

The App, Component and Physical channels use `internal/safefile` for one bounded, opened-file-identity-checked policy over cached signed envelopes, App pending/previous envelopes, physical provenance and physical authorization. The same package also owns streamed SHA-256 hashing of installed App/Component/Physical executables and donor Portable sources, with the signed artifact size enforced **before** reading, and a canonical signed-manifest-derived maximum for the Portable receipt. The local file must be a non-symlink regular file with an expected size before opening; after opening it must still be the same file with an allowed size, and only a bounded read is performed. Each caller retains its existing signed-envelope, SHA-256, provenance, rollback and device-specific authorization checks. This closes the gap between bounded network reads and previously unbounded offline `os.ReadFile` operations. The development updater also uses the same `ReadBoundedRegular` for its unsigned `current.json` receipt, replacing the old independent reader, and `DigestInventoryRegular` for all cached Dev file hashes. The shared hashing implementation validates opened-file identity and exact pre-read size; the development-only wrapper preserves valid zero-byte files, unlike the positive-size signed-artifact policy. The unsigned receipt remains strictly a local integrity check, never publisher authentication or physical-write permission. Windows CI runs the shared cache security packages directly.

## Single official App activation path

The signed Creator App has exactly **one** path for a new release to become active: `appchannel.AcquirePending` stages a SHA-256-bound executable only after two identical, independently retrieved signed envelopes; the Windows launcher runs its bounded `--launcher-healthcheck` against that candidate; and `appchannel.PromotePending` atomically updates `current-envelope.json`, preserving a separately verified `previous-envelope.json` for last-known-good recovery. The former `appchannel.AcquireCached` direct-activation API was removed because it could replace the current signed App pointer without this health proof or previous-App preservation. There is no second App activation or publisher authorization policy.

A staged but unpromoted App is never reported by `Current` as active, including when its signed envelope and executable hash are correct. Incomplete second-envelope verification, failed pointer staging, and discarded pending updates cannot activate it. Component and Physical distribution are separate channels and retain their own signed cache activation contracts.

## Launcher health proof and last-known-good startup

The Windows signed-App launcher runs the already verified pending executable's `--launcher-healthcheck` with a 30-second deadline and at most 8 KiB retained combined stdout/stderr. It accepts **only one exact** `ORDAX_CREATOR_APP_HEALTH=PASS` line with a final LF or CRLF, after a zero exit code; duplicate, embedded, truncated or excessive output and hanging checks cannot promote a candidate. Windows native CI tests the launcher behavior directly, in addition to its standalone binary smoke.

A later synchronous failure of `Process.Start`, even after health-check and promotion, falls back exclusively to the separately validated signed `previous-envelope.json`, if it names a **different** source commit. The same fallback covers start failures when the launcher's selected signed current App cannot start. Failed or missing fallback is reported as an error, never fabricated as success. The signed `current-envelope.json` remains unchanged so its authenticated release sequence still governs future updates; this is a runtime last-known-good launch, not an unsigned downgrade or a second update mechanism. It does not claim to detect failures after the child process has successfully started.

## Portable source reuse from a single verified file handle

Portable Physical candidates may reuse any earlier cached 17-source artifact **only** when its actual bytes match the new, separately verified signed manifest's positive size and SHA-256. The shared `internal/safefile.CopyBoundRegular` confirms the donor path is a regular non-symlink file of exactly the signed size, verifies `os.SameFile` on the opened descriptor, and performs **one** streamed copy and SHA-256 calculation on that same handle into a private unpublished candidate directory. A second `os.Open` or an initial full-file prehash is not required. Incorrect size, symlink, malformed input, and opened-file identity changes cannot be accepted. A stable but stale/mismatched donor is deleted from staging and the next signed-hash-matching cached copy (or the official download) is tried; read, copy, and staging I/O errors are fatal. The staging file is synced and closed only after the actual digest is compared to the signed binding, and the completed 17-source set remains subject to the Creator Core's full `VerifyPortableSources` check before installation or any physical operation.

This reduces local donor read I/O for large already-cached payloads without creating a shared mutable source index, accepting unsigned bytes, altering the OS producer, or weakening physical consent, authorization and readback checks.

## ZIP extraction size contract

The Dev, Component and Physical ZIP consumers use `internal/archivebounds` as their **only** expanded-size arithmetic and exact-copy boundary. It rejects unsigned ZIP64 headers whose `uint64` sizes cannot fit the remaining per-channel `int64` extraction budget, **before** writing the entry. Copying is capped to the validated header length plus one byte and requires precisely the declared number of bytes, so truncated or overlong decompression fails closed. The existing channel-specific file inventories, SHA-256 bindings, ZIP names, signature checks and maximum budgets remain authoritative; the helper itself grants no signature trust or physical permission. Negative tests exercise oversized forged central-directory sizes and ZIP64 arithmetic on Linux and native Windows.

## Evidence and release gates

- Core/Go vet, immutable read-only contract integration, real Win32 UI smoke, unsafe-target negative tests, repeatable unsigned ZIP builds, and synthetic/disposable file tests are non-destructive engineering evidence.
- CI artifacts must display their Creator commit and explicit **UNSIGNED/READONLY** status. They are not Stable artifacts, signed releases or proof of physical USB writing.
- The official publisher must separately complete source and signature verification, Authenticode/RFC3161 custody, signed channels, target-specific authorization and physical readback proof. The Creator repository must not create another signer, trust root, media producer or consent authority.
- All three **official signed consumers** (App, Componentes e Físico) must confirm the identical signed envelope on an independent second read before activation; App additionally requires launcher health proof and explicit pending promotion, while Componentes and Físico may activate their own verified cache pointers. Network errors, changed signed bindings, invalid signatures and pointer-write failures return an error; only a separately verified existing `Current` or the App's explicit last-known-good flow may be used offline. The App launcher continues to require its separate health-check and promotion. Each official channel accepts only its purpose-bound canonical GitHub release asset URL (no query, encoded alias or nested path). The physical channel additionally requires one ZIP asset within the `creator-physical` release.
- No test, document or merge here grants the right to enumerate, repartition, eject, lock or write a real device. Device-specific physical authorization follows the OS-owned procedure in [AGENTS.md](../AGENTS.md).

**Current status authority:** [docs/MIGRATION.md](MIGRATION.md) records historical cutover evidence; [issue #5](https://github.com/ordaxsystems/ordax-creator/issues/5) tracks only outstanding distribution/physical proof, not a second architecture.
