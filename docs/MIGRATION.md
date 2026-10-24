# MIGRATION — Independent Creator extraction

**Owner:** ordaxsystems/ordax-creator  
**Status:** ACTIVE — INDEPENDENT CREATOR MAIN, OS CUTOVER PENDING  
**Started:** 2026-10-23  
**Nature:** execution plan and immutable handoff, not a new OS contract or physical-release authorization.

This is the *only active Creator migration plan* in this repository. Update it rather than creating a second roadmap.

## Verified source and target

- New Creator GitHub repository: **accessible, public, admin/push enabled**, repository ID 1431616910, created empty.
- Initial Creator main: 4d32db78183eee586f5aaad930753afabb5400b3 (documentation bootstrap, no executable).
- Audited OS main: f65acc4fb9f8ded376fc9b44b4c7c2068c793265; **recheck at every execution**, never treat the dated SHA as moving main.
- OS source root: tools/creator/; original module github.com/ordaxsystems/ordax-os/tools/creator, Go 1.27, product version 0.6.2-dev.
- OS #2000 (a8c11fb636e2c652aa2c0c24f4fa8ad5ae349deb) is an OPEN/DRAFT six-file single-inventory refactor; preserve. At that head, 12/13 workflows passed: First Physical USB Payload failed at OS Stable Base's pinned Alpine ca-certificates version; the new Core 17-source tests passed. Do not modify Base to make Creator CI green.
- OS #2193 remains a separate Creator-documentation handoff draft; no merge into OS until coordination with Surface agent.
- Existing OS machine-readable code-signing contract has publish_allowed=false; no public signed Creator or physical write permission exists.

## Scope by owner

| Current OS location | Execution owner after cutover | Handling |
| --- | --- | --- |
| tools/creator/core/ | Creator | Preserve history, one Core, exactly 17 artifacts and 39 operations |
| tools/creator/cmd/ | Creator | Migrate real GUI/CLI/helper and rewrite imports |
| tools/creator/host/, platform/ | Creator | Thin OS-specific adapters; keep RAW writer gated |
| tools/creator/appchannel, componentchannel, physicalchannel, internal/atomicfile, update | Creator | Retain one signed-channel consumer and atomic-pointer policy |
| tools/creator/proof/, root Python scripts, windows/ | Split by audited producer/consumer role | Move only genuine consumer tests/helpers; never blindly duplicate producer scripts |
| OS bootstrap/, docs/contracts/, tools/release-signing/, OS physical-promotion scripts | OS | Keep one official OS artifact, signing, release and authorization authority |
| OS CI producers / release workflows | OS | Retain producer; replace consumer references only with verified immutable Creator artifact |

## Phases and objective exit criteria

| Phase | Deliverable | State |
| --- | --- | --- |
| C0 | Source owners, migration plan, provenance and guarded independent CI foundation | IN PROGRESS |
| C1 | Exact-commit history extraction from tools/creator/ with Git ancestry, source tree equality and receipt | PASS — import receipt in docs/evidence |
| C2 | File-by-file producer/consumer classification, adoption of OS #2000 without duplicated 17-source inventory | Creator code imported and Core #2000 applied; OS producer #2000 reconciliation still PENDING |
| C3 | Canonical Creator Go module, import/path rewriting, immutable OS contract integration and independent packages | PASS for source/tests; package publisher integration pending |
| C4 | Go and Windows build/test, fail-closed device/security regressions, deterministic unsigned packages and hashes | PASS for Go, Windows-native CI and unsigned Windows ZIP; public signing/hardware approval pending |
| C5 | OS Native helper consumer pinned to independent Creator, merged PR #2197, 5/5 CI | PASS — `ordax-os/main@f07fe50d` |
| C6 | Retire old OS Creator Go source after consumer/publisher migration and full reference proofs | IN PROGRESS — Native, Candidate, Toolkit, Authenticode and official handoff migrated; dev channel, physical source-hash context and remaining producer paths block deletion |

## C1 — history-safe source import

The one-shot GitHub Actions workflow in .github/workflows/creator-source-import.yml is a **migration branch producer only**. It:
1. Reads the public OS repository at exact commit f65acc4fb9f8ded376fc9b44b4c7c2068c793265 and confirms it exists in its history.
2. Computes the source directory tree and git subtree split --prefix=tools/creator at that commit.
3. Requires exact tree SHA equality.
4. Fetches the extracted commit object into Creator and merges it as a real Git parent (no rewriting main).
5. Stages the initial snapshot under temporary source-import/ on migration/creator-history-import and records all SHAs.
6. Pushes this separate branch once, **refusing to overwrite an existing import branch**.

The snapshot under source-import/ is quarantined and MUST NOT be merged unchanged as a second implementation. Subsequent C2/C3 changes select and relocate real Creator-owned files to canonical module root and remove source-import/ before final merge. A source copy or detached SHA comment alone does not preserve ancestry. No OS publisher private key or physical media are imported.

If the GitHub Actions token cannot write to this repository, the action fails closed without alternate secrets; repair workflow permissions using normal repository governance and revalidate.

## C2 — PR #2000 coordination

- Preserve the original OS PR and its review discussions, six changed paths, history and open draft state.
- Reconcile its Creator Core changes to portable_sources.go, portable_sources_test.go, and cmd/ordax-creator/main.go against the live OS source history; never maintain a second 17-artifact inventory.
- Changes to the OS producer workflow, test and physical evidence documentation remain in OS PRs and must not be blindly copied into the Creator executable.
- The upstream Stable Base dependency failure is not a Creator security or Core failure; ask the OS Base owner for reproducible resolution, without relaxing pins or changing OS Foundation.

## C3/C4 — product and proof

- Build standalone Creator with module github.com/ordaxsystems/ordax-creator. No full OS checkout as the runtime source; use only source-pinned, SHA-verified OS contract/asset downloads in CI.
- Keep exactly one host-neutral Core, existing app/component/physical signed channels, one atomic signed-pointer implementation and narrowly elevated Windows adapter. Future OrdaX Desktop consumes this Core, not another flasher.
- Prevent physical write code from appearing in default public binary. No implicit promotion of unsigned CI builds.
- Add Go tests, cross-OS compilation checks, Windows GUI/update/localization checks, signed-source/cache negative tests, device serial/capacity/re-enumeration mismatch and OS-disk exclusion tests, aborted UAC, partial write/readback failures, disposable-only filesystem tests and repeat-build SHA-256/provenance. Never simulate full hardware proof claims.
- Refuse OS contracts not pinned by exact commit, release signature, size/hash, source commit and public trust. No duplicate local machine-readable OS contract, media geometry, signing root or release catalog.

## C5/C6 — integration and cleanup

OS must remain producer of Base, kernel, bootstrap/release and media/physical policy. Integrate only the verified Creator consumer from a SHA-bound package or source revision via a scoped OS draft PR, rechecking the Surface agent's latest OS main. Mark old tools/creator code deletable only after standalone Creator build, Windows packaging, OS workflow integration, references/CI audits and exact-head tests pass. Do not remove partially superseded scripts or change Foundation, Surface, kernel or Buildroot.

## Evidence checklist for completion

Document original OS SHA, original directory tree, extracted split commit+tree, Creator merge parent, classification of source files and imported commit, OS #2000 reconciliation, contract SHA/version pins, test/run URLs at exact heads, Windows unsigned artifact hashes/reproducibility, package/Authenticode status, OS consumer cutover and deletion PRs, and an explicit separate statement of physical authorization. Until then: **NOT DONE, PHYSICAL_WRITE_AUTHORIZED=NO**.

## C2 staging receipt (automated source migration)

- Moved Core, cmd, host, platform, internal, appchannel, componentchannel, physicalchannel, update, proof, go.mod and version.json from immutable history into Creator module root.
- Removed the temporary source-import directory in the branch; files omitted from active Creator stay reachable via source history.
- Explicitly excluded OS release/promotion root Python scripts and windows publisher scripts pending per-file producer audit.
- Converted Go import prefix to github.com/ordaxsystems/ordax-creator, without creating a second Core.
- This branch is NOT cut over; CI tests and canonical contracts integration remain required.

## Source consolidation and independent evidence — 2026-10-23

- History import: original OS `f65acc4fb9f8ded376fc9b44b4c7c2068c793265` tree `cc287d6dafba00b91a1c3e2f86d154ad4a943cd0` equals split commit `d38a92ab75c2ae94319740f3240b785b2ad819f1` tree; actual two-parent ancestry merge `162dbe4a546c95a768c442f6779c4670a031465b`; source snapshot commit `d78287f6b1e3b551d5c41a52765e95eb6f7668af`. [History CI](https://github.com/ordaxsystems/ordax-creator/actions/runs/40035831823): PASS.
- Consumer-only code extraction: 153 Go files, a single Creator root module, one existing Core, app/component/physical signed-consumer packages; removed transitional `source-import/`. No OS producer authorization, release/publisher Python script or Windows publisher script was retained as active Creator source. [Extraction CI](https://github.com/ordaxsystems/ordax-creator/actions/runs/40036089102): PASS.
- Audited source comparison: for the three Creator Core/CLI paths changed by OS PR #2000, source OS `main@f65acc4f` blob IDs **equaled the PR's base**; adopted exact PR head `a8c11fb636e2c652aa2c0c24f4fa8ad5ae349deb` versions and rewrote only the Go module prefix. Creator adoption commit `0680dae87fbf3f7e683fb9c432b46678e344a32c`. OS producer workflow/test/document changes remain in OS PR #2000; do not close or merge its conflicting producer work by assumption.
- Test-path decoupling `f08c2b212fd96ad60360e139ff78fc15fb4f67fe`: explicit `ORDAX_OS_CONTRACTS_DIR` for contract-integration tests. Unset runs offline unit tests; set but missing/invalid fails closed. CI sparse-checks out and verifies exact OS source Git commit and canonical schemas, hashes read-only inputs, then runs **all** Go tests, `go vet`, Win32 GUI/CLI builds and RAW build-tag isolation. [Standalone CI](https://github.com/ordaxsystems/ordax-creator/actions/runs/40036405670): PASS.
- [Unsigned Windows candidate](https://github.com/ordaxsystems/ordax-creator/actions/runs/40036493128): PASS at `f72fae0c2a8450fc114f2c7d686aea8b2285d616`; four unprivileged Windows executables build twice with byte-for-byte identical outputs; fixed-timestamp ZIP plus source hash/provenance. Artifact `ordax-creator-UNSIGNED-READONLY-f72fae0c2a8450fc114f2c7d686aea8b2285d616` excludes physical helper, trust remains UNRESOLVED and public writing/publishing remain disabled.
- Still pending: real Windows GUI/native-platform smoke and installation packaging, final official signer/Authenticode custody, signed distribution channels, OS consumer integration, physical validation with separate authorization, reconciliation/closing producer OS PR #2000, and OS source cleanup after proving equivalence. Do not mark C5/C6 done.


## Creator mainline cutover of implementation source (not OS consumer) — 2026-10-23

- Creator PR #1 integrated as merge commit `525737374c89d92e0712989489d77d45c40784ab`.
- Creator PR #3 integrated as merge commit `d1d32e4858ab7e2f3629fba2ac4dde5e1e1e97ef`, preserving history split as a real Git ancestor. The effective Go packages now live on Creator main, without `source-import/`.
- Full standalone CI against the exact source PR head `e8905207dd7f200363f53fb74ad31941b17ba11e`: [PASS](https://github.com/ordaxsystems/ordax-creator/actions/runs/40036689755), including Go tests/vet, OS immutable read-only contract schemas, development binding UNRESOLVED, Windows GUI/CLI cross-build and RAW backend build-tag isolation.
- Deterministic unsigned Windows development ZIP: [PASS](https://github.com/ordaxsystems/ordax-creator/actions/runs/40036493128). The imported code includes NO producer-side OS release signing or physical authorization tooling.
- Existing migration-only import/stage workflows are historical and are being removed from the active `main` CI route now that source ancestry is recorded; long-lived Go tests and unsigned candidate builds must execute on `main` push and Creator PRs.
- **No OrdaX OS consumer cutover, publisher code signing, Windows hardware GUI smoke, real write authorization or removal of old OS Creator yet.** C5 and C6 remain pending. Active OS PR #2000 must still reconcile its producer-owned changes with current OS main and Surface.

## Native Windows acceptance and OS source dependency cutover — 2026-10-23

- Creator [PR #6](https://github.com/ordaxsystems/ordax-creator/pull/6) merged as `e7432abf30f178a495196b576cb7fd212006d45e` with **3/3 exact-head CI checks PASS**: [standalone Linux Core](https://github.com/ordaxsystems/ordax-creator/actions/runs/40037999005), [actual Windows 2025 runner tests](https://github.com/ordaxsystems/ordax-creator/actions/runs/40037999017), and [unsigned deterministic Windows ZIP](https://github.com/ordaxsystems/ordax-creator/actions/runs/40037999018). Windows native UI test creates actual Win32 window/controls without triggering refresh, update, disk discovery, UAC, physical writer or network. Default physical button stays disabled. CLI smoke invokes read-only inspection `status` and launcher `--launcher-version` only. All signed consumer tests executed on Windows. Windows Go vet has a **narrow, explicitly audited four-callback Win32 LPARAM pointer ABI exception**; other packages use strict vet.
- Single OS contract dependency lock `.github/ordax-os-contracts.lock.json` pins the *external authoritative source revision*, not contract content. Both Linux and Windows jobs verify the full exact OS commit and schema while keeping all source bytes read-only. OS source `tools/creator/` tree at OS `main@70b0f7eb5d5085be0d4f176f5232def37c7a2a24` remained unchanged from the historical extracted tree `cc287d6dafba00b91a1c3e2f86d154ad4a943cd0`.

### C5/C6 verified dependencies that forbid premature OS deletion

1. `ordax-os/tools/native-release-assembly/build.py` (at exact OS `70b0f7eb`) builds `system/bin/ordax-native-install-targets` from `tools/creator/go.mod`, `core/`, `host/linux/` and `cmd/ordax-creator-native-targets/`. This is a read-only Native installation helper built into the **OS release**; its final publisher remains OS. After extraction, the **single** implementation stays in Creator while the OS must build it from an exact Creator source commit without mirroring code.
2. Canonical OS `docs/contracts/native-release-assembly.json` and `docs/contracts/build-autonomy.json` contain the old source input paths; `docs/contracts/native-installation.json` records the helper and native source discovery function paths. Repoint contracts only as part of one coordinated OS consumer cutover PR with updated tests and proven deterministic build.
3. `ordax-os/tools/creator/physical_promotion.py` is **OS physical authorization/publisher authority** and currently hashes Creator Go source paths inside the OS authorization context. Move/reshape its source authority only with a fresh canonical owner-accepted release policy; never infer that an existing authorization hash is still valid after a path/commit cutover.
4. `ordax-os/tests/test_portable_usb_v2_contract.py` references `tools/creator/proof/portable_usb_v2.py` and the old Core CLI; OS `release-pipeline.yml` and `full-bootstrap-media-proof.yml` have old-path filters. Classify test-by-test rather than deleting, duplicating OS contracts, or silently skipping the producer CI.
5. `ordax-os#2000` remains producer-side and must reconcile its OS workflow/test/doc changes independently. Its OS Base pinned Alpine dependency failure does not authorize relaxing signatures, bypassing physical promotion or editing Foundation.

**Current-state superseding the historical draft evidence below:** Creator Native helper source build and read-only dependencies were proven in Creator PR #7; OS [PR #2197](https://github.com/ordaxsystems/ordax-os/pull/2197) subsequently **merged** at `f07fe50d414e883cc88b11988008e799b17366be` with all 5/5 CI gates passed. **C5 PASS; C6 remains IN PROGRESS.** Physical write remains unauthorized.


## C5 OS Native helper consumer draft — 2026-10-23

- Proposed, **not merged**, at [ordax-os PR #2197](https://github.com/ordaxsystems/ordax-os/pull/2197), branch `integration/creator-native-external-source-20261023`, exact PR head `082d5afc374f2e667b53c8003a12076b9438410e`, based on OS `f4119971e8e73b49b6e0c9310b64240881e4d514`.
- The OS's existing canonical `docs/contracts/native-release-assembly.json#creator_source` alone pins `ordaxsystems/ordax-creator@8adf641fb2cd9a3ac8351f23f912ce29ee2b70c9` with full Git tree `312d18bb8e5cb97a216a54697724bbc8893d33d0`. It records Go module and helper package, `physical_apply_authorized=false` and `publisher_authority=ordaxsystems/ordax-os`. No external contract clone, parallel release owner, or floating `main` source exists.
- OS Release Pipeline obtains a read-only external Creator checkout into an isolated temporary runner directory and re-verifies source SHA/tree/clean status before both assemblies. The OS Native assembly recipe independently checks Git-root, exact commit, tree, source cleanliness, original module and helper entrypoint. Missing or modified source fails before any output is created.
- OS Native `system/.ordax/native-release-tools.json` tracks the external helper repository, commit and tree **in addition to** existing SHA-256, size, mode, role, OS source commit and physical authorization state. OS final deterministic `system.tar` bundler, signing/physical promotion owners and separate runtime-component verifier are **unchanged**.
- 12 synthetic source integrity negative tests added in the OS repo. [Exact-head OS Release Pipeline run #40039904876](https://github.com/ordaxsystems/ordax-os/actions/runs/40039904876) PASS, including full Native assembly, two identical helper builds, two identical Native provenance files, `system.tar` deterministic equality and signed envelope compatibility. [OS Surface Consumer run #40039904975](https://github.com/ordaxsystems/ordax-os/actions/runs/40039904975) PASS; [Native Disposable Storage run #40039904834](https://github.com/ordaxsystems/ordax-os/actions/runs/40039904834) PASS. Other broad OS CI checks must still be rechecked before coordination/merge.
- OS `tools/creator/physical_promotion.py`, old Creator producer and physical authorization-context hashes, #2000 and unrelated Surface/Foundation/Buildroot code **remain untouched**. Do not remove historical `tools/creator/` until all active publisher consumers have migrated with exact-source-hash reauthorization. No hardware write, release publish, signing keys or public physical authorization.
- **Exit gate:** obtain explicit integration coordination with the agent modifying Surface/OS Native release assembly, verify current OS main/other PRs, complete all relevant exact-head CI, and only then consider merge of OS #2197. Mark C5 done only after actual merge/proven current OS consumer. C6 and full Creator product signing/physical proof remain separate.


## C5 complete; C6 consumer audit (live status, 2026-10-23)

The dated C5 draft section above is retained as an evidence record, not current execution state.

**Merged OS consumer cutover:** [ordax-os PR #2197](https://github.com/ordaxsystems/ordax-os/pull/2197) merged to `main@f07fe50d414e883cc88b11988008e799b17366be`. All **5/5 PR-head checks PASS**: Release Pipeline 40039904876, Surface External Source Consumer 40039904975, Native Disposable Storage 40039904834, Foundation 40039904815, Public Site 40039904925. It continues to pin the Creator code revision `8adf641fb2cd9a3ac8351f23f912ce29ee2b70c9` and tree `312d18bb8e5cb97a216a54697724bbc8893d33d0`, with no old source copy or publisher transfer. The subsequent `main` CI is a separate post-merge verification.

**C6 in progress:** [ordax-os draft PR #2199](https://github.com/ordaxsystems/ordax-os/pull/2199), exact head `dc8fcf6c2be7be768f48a208d6958dfcbe5ce475`, extracts a reusable OS checkout/verification action from the release pipeline and repoints the OS-owned **unsigned, non-destructive Creator Candidate** to the exact external source. This CI validation is distinct from Creator's own standalone code build, since it verifies consumption of authoritative OS contracts. This change is **not merged** at the time of this record. Its CI has approved Creator Candidate 40040816780 and Release Pipeline 40040816790; Foundation 40040816774 was still in progress as this snapshot was written. Do not treat passing branches as integration.

**OS old-source dependency categories (audit of OS `main@f07fe50d`):**

- **Already transferred:** Linux Native read-only helper in the OS release assembly to pinned independent Creator, while preserving the existing OS signing/release chain.
- **Consumer build pending transition:** OS Creator Candidate (draft #2199); other candidate/toolkit jobs including `.github/workflows/windows-prototype-toolkit.yml`.
- **Developer rollout/publishing — DO NOT RETIRE or redirect blindly:** `creator-dev-channel.yml` currently builds from `tools/creator` and carries `contents:write` / development release upload. A safe cutover must preserve update channel URL/manifest compatibility, source identity, non-destructive policy and release ownership; no automatic publication by a source-only PR.
- **Official consumer, signing and physical publisher:** `creator-official-candidate.yml`, `creator-code-signing-boundary.yml`, `creator-physical-publisher.yml`, `physical-write-promotion.yml`, `first-physical-usb-payload.yml`, `full-bootstrap-media-proof.yml`, `tools/release-signing/` and `tools/publisher/` have separate OS-owned signatures/authorization contracts that **must not** be copied into Creator or deleted.
- **Physical authorization trust context:** `tools/creator/physical_promotion.py` hashes `tools/creator/go.*`, Core, Windows host, physical channel and raw writer source. Moving the code changes the authorization context and requires a safe **new source+revision-bound policy** with proof of fail-closed default; old authorization hashes cannot be reused.
- **OS-owned Python build/proof scripts and test references:** `tools/creator/assemble.py`, `proof/`, several producer Python scripts and many `tests/test_*.py` paths remain operational. Classify as OS release producer, Creator code consumer, or historical evidence. No physical-media access to retire dependencies; do not suppress tests or duplicate the Go module.
- **Producer PR #2000:** Core/CLI changes already adopted in independent Creator; workflow/test/docs work remains open in `ordax-os#2000`. Reconcile only its producer-owned changes against current OS main and new Creator source lock; no duplicate Core source inventory.

**C6 exit gate:** before deleting `ordax-os/tools/creator/`, replace all active old-module imports, CI builds, runtime source paths, OS publisher source-hash references, contract paths and test dependencies with exact independently owned Creator code or correctly separated OS producer modules; prove complete CI (including unsigned Windows, physical fail-closed and signed-release contracts), check no unresolved references, and preserve user-facing channel compatibility. Physical hardware validation requires separate explicit device-specific authorization.

## C6 OS non-destructive candidate source cutover merged — 2026-10-23

- [ordax-os PR #2199](https://github.com/ordaxsystems/ordax-os/pull/2199) **MERGED** into `main@7677226def5fd027d37a84a8c7981c3c90d49b10`. Exact PR head `dc8fcf6c2be7be768f48a208d6958dfcbe5ce475` received **3/3 PASS**: [Creator Candidate](https://github.com/ordaxsystems/ordax-os/actions/runs/40040816780), [Release Pipeline Contract](https://github.com/ordaxsystems/ordax-os/actions/runs/40040816790) and [Foundation Contract](https://github.com/ordaxsystems/ordax-os/actions/runs/40040816774). Independent post-merge `main` checks must be evaluated separately.
- One reusable OS consumer action, `.github/actions/resolve-creator-source/action.yml`, reads **only** the canonical existing OS pin `docs/contracts/native-release-assembly.json#creator_source`. It obtains a clean read-only checkout of exact Creator Git SHA/tree and refuses wrong repo/module, altered tree, a floating branch, missing source or promotion authority. Both the OS Native Release Pipeline and OS `Creator Candidate` use the action, no duplicate Go source or pin.
- OS `Creator Candidate` runs the independent Creator Go tests, OS manifest fail-closed tests, Windows raw-backend compile-only isolation and Windows/Linux unsigned CLI candidate build from that checkout. The non-destructive candidate provenance identifies Creator code SHA/repository/tree separately from the OS consumer commit. **This is not the official signed Windows release and does not activate physical writes.**
- Significant still-active historical source consumers: `windows-prototype-toolkit.yml` builds the old Go module and OS signer tooling, `creator-dev-channel.yml` retains write-permitted development update publication from the OS copy, `creator-code-signing-boundary.yml` and `creator-official-candidate.yml` still compile old paths, and physical publisher/authorization context hashes old Go files. Preserve their safety and existing publisher responsibility during migration; no source removal or channel endpoint change without complete proof.
- Surface owners have been notified in OS PRs #2195 and #2198; their PRs are managed separately. The remaining C6 phases must migrate actual consumers, reconcile OS PR #2000 producer changes and preserve physical trust, rather than delete live files to satisfy a duplicate count.

## C6 Windows Prototype Toolkit consumer migrated — 2026-10-23

- [OS PR #2200](https://github.com/ordaxsystems/ordax-os/pull/2200) **MERGED** as `ordax-os/main@e01f4a92d43278fa66ebc68847469dfbdbc47c11`. Exact merged PR head `4151d23b9b18ce619d0f532d8cb65f83a061b42e` passed **4/4 workflows**: [Windows Prototype Toolkit](https://github.com/ordaxsystems/ordax-os/actions/runs/40041985452) (both native Windows and Ubuntu jobs PASS), [Creator Candidate](https://github.com/ordaxsystems/ordax-os/actions/runs/40041985411) PASS, [Release Pipeline Contract](https://github.com/ordaxsystems/ordax-os/actions/runs/40041985404) PASS, [Foundation Contract](https://github.com/ordaxsystems/ordax-os/actions/runs/40041985438) PASS. A new post-merge main CI cycle was independently started and must not be confused with PR-head success.
- Windows toolkit and Native Windows Go tests now compile/execute the **exact same pinned standalone Creator source** as other read-only OS consumers using one `OS/.github/actions/resolve-creator-source` and the single existing `native-release-assembly.json#creator_source` lock. They do **not** compile the duplicated old `tools/creator` Go module.
- Existing OS release signer/trust ceremony scripts are still built and packaged from the OS repository, with no publisher key included in CI. Toolkit provenance explicitly records two owners: OS `source_commit/source_repository` for publisher/trust contracts and `creator_source_commit/creator_source_repository/creator_source_tree` for source implementation. Physical apply and public writer remain **blocked**.
- The first Windows job proved a cross-platform CRLF defect (Python to Git Bash injected CR into the clone URL). The final PR fixes the **single shared resolver** by stripping CR and validating owner/hex SHA/tree before Git; both OS Windows and Linux consumers passed on the corrected exact head. Do not add Windows-specific source locks, checkout fallbacks or disable integrity checking.
- **Not yet migrated:** `.github/workflows/creator-dev-channel.yml` (writes development GitHub Release assets and must retain URL/manifest backward compatibility), `creator-official-candidate.yml`, `creator-code-signing-boundary.yml`, the physical publisher/authorization-source context, old OS-owned Python producer scripts/proofs, and conflicting producer changes in OS PR #2000. Old Go files remain tracked until all their live references have safely cut over; do not delete or bypass checks.

## C6 Windows Authenticode boundary consumer merged — 2026-10-23

- [OS PR #2201](https://github.com/ordaxsystems/ordax-os/pull/2201) **MERGED** to `ordax-os/main@0d2011e6f704ba7ebdfaba5435cf77adf378e175` at exact approved PR head `8d18508b303d8b1e9346be08b2e6762c6f3b5ff2`. Both PR-head checks **PASS**: [native Windows Creator Code Signing Boundary](https://github.com/ordaxsystems/ordax-os/actions/runs/40043640272) and [Foundation Contract](https://github.com/ordaxsystems/ordax-os/actions/runs/40043640248). The independent post-merge OS-main workflows are separate and are not counted as proof until complete.
- The **unsigned** Windows launcher and Creator GUI built for Authenticode negative tests now use the **single** immutable Creator implementation source declared in `docs/contracts/native-release-assembly.json#creator_source`. The shared resolver verifies Git owner, exact commit and tree, clean checkout and module identity; the Windows job checks source commit/tree again before running the Creator appchannel, componentchannel, physicalchannel tests and building the two executables.
- OS-owned `tools/release-signing` Go tests/binaries and **historical OS-tracked** `tools/creator/windows/Verify-OrdaXCreatorSignature.ps1` / `Finalize-OrdaXCreatorRelease.ps1` remain in place. They are publisher/production-trust tools that are **not** part of the standalone Creator Go code and were not silently transferred to the public consumer repo. These old PS1 paths must be rehomed under OS-owned publisher authority in a later coordinated phase before the full `tools/creator/` subtree can be deleted.
- The exact Windows test still proves the production Authenticode policy remains **unconfigured**, unsigned/unlisted executables are rejected, publisher finalizer cannot generate manifest/envelope before valid signing, and no release is published. **No certificate, private signing key, production Authenticode signature, signed public artifact or physical-media access was involved**.
- **C6 blockers still active:** `creator-official-candidate.yml` with OS publisher handoff; `creator-dev-channel.yml` with existing rolling public development download URL; OS publisher/trust PowerShell scripts still located under legacy Creator paths; physical source-hash/authorization gate; remaining Python producer/test references; OS draft PR #2000. Do not delete old OS sources, falsify source identities or bypass signed/physical safety gates to complete the cleanup.

## C6 official unsigned publisher handoff source cutover merged — 2026-10-23

- [OS PR #2202](https://github.com/ordaxsystems/ordax-os/pull/2202) **MERGED** to `ordax-os/main@fe462f09cd70baeadb515de4a605fef40fbbc769`, source PR head `3dd302dfd9550bd53e63addeaf0b023717d0dd3c`. All **4/4 exact PR-head workflows PASS**: [Creator Official Candidate 40048288274](https://github.com/ordaxsystems/ordax-os/actions/runs/40048288274), [Physical USB Backend Safety 40048288414](https://github.com/ordaxsystems/ordax-os/actions/runs/40048288414), [Creator Disposable Media 40048288288](https://github.com/ordaxsystems/ordax-os/actions/runs/40048288288), [Foundation Contract 40048288386](https://github.com/ordaxsystems/ordax-os/actions/runs/40048288386). Post-merge main CI must be inspected separately.
- OS-owned `tools/creator/windows/Build-OrdaXCreatorPublisherHandoff.ps1` no longer compiles the old embedded Go module. It now **requires** an external Creator checkout; before creating an output folder, it independently verifies existing canonical `docs/contracts/native-release-assembly.json#creator_source` (repository, module, publisher owner, physical safety, lowercase pinned commit and tree), the Git checkout root and cleanliness, and exact launcher/app entrypoints.
- The **same** existing `.github/actions/resolve-creator-source` obtains immutable Creator code in the Ubuntu trust gate and both Windows handoff jobs. OS source and Creator code source now have **separate provenance identities**, with OS as release authority and Creator as implementation owner. No duplicate lock, Go Core copy or floating `main` checkout.
- A Windows disposable wrong-source negative test confirms an invalid checkout is rejected **before output generation**. Publisher handoff preserves canonical versus disposable public trust, signed-release manifest/Ed25519 order, OS signer binaries/PS1 scripts and explicit `authenticode_signed=false` / `publishable_to_end_users=false`. This is **not** a public signed Creator release and grants **no hardware write permission**.
- **Remaining C6 blockers:** OS `creator-dev-channel.yml` still compiles the old Go module and performs actual GitHub Release uploads with `contents: write`. That migration must keep exact old channel URLs, schema compatibility and update consumer behavior; its source-only tests must not accidentally publish. The OS physical promotion context still hashes old `tools/creator` Go files, with a separate authorization protocol; active OS Python media/release proof consumers and draft OS PR #2000 must be reconciled. Publisher-owned OS PowerShell scripts still reside under `tools/creator/windows/`, requiring an explicit owner-preserving rehome before deleting the old tree. Never delete or bypass any existing guard before validated cutover.
