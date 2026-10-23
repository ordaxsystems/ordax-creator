# MIGRATION — Independent Creator extraction

**Owner:** ordaxsystems/ordax-creator  
**Status:** ACTIVE — PRE-CODE / NOT CUT OVER  
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
| C2 | File-by-file producer/consumer classification, adoption of OS #2000 without duplicated 17-source inventory | IN PROGRESS — consumer packages staged, #2000 reconciliation pending |
| C3 | Canonical Creator Go module, import/path rewriting, immutable OS contract resolver and independent packages | PENDING |
| C4 | Go and Windows build/test, fail-closed device/security regressions, deterministic unsigned packages and hashes | PENDING |
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
