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
| C5 | Creator-specific OS integration draft PR and Surface reconciliation; no premature merge | PENDING |
| C6 | Eliminate OS legacy implementation only after standalone and OS integration proofs and reference audit | PENDING |

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

**Next exit gate:** independently prove the *Linux Native read-only helper* from Creator Core is reproducible, has a precise source SHA and does not link Windows/physical writer. Only then design/validate an OS-specific PR to consume that pinned source. **C5/C6 remain PENDING** until tests and producer coordination finish. Physical write remains unauthorized.

