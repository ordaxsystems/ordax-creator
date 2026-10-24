# Creator migration provenance ledger

This ledger records origin independently of chat. It is **NOT** evidence that a source migration has finished.

| Item | Identity / state |
| --- | --- |
| Original repo | ordaxsystems/ordax-os |
| Original implementation | tools/creator/ |
| Audited OS main commit | f65acc4fb9f8ded376fc9b44b4c7c2068c793265 |
| OS PR preserving single inventory | #2000, head a8c11fb636e2c652aa2c0c24f4fa8ad5ae349deb, OPEN/DRAFT |
| Original Creator module | github.com/ordaxsystems/ordax-os/tools/creator |
| Original Creator product version | 0.6.2-dev |
| Target Creator repo | ordaxsystems/ordax-creator (ID 1431616910) |
| Initial Creator main | 4d32db78183eee586f5aaad930753afabb5400b3 |
| Original subtree tree SHA | cc287d6dafba00b91a1c3e2f86d154ad4a943cd0 — verified |
| Filtered history commit SHA | d38a92ab75c2ae94319740f3240b785b2ad819f1 — verified |
| Creator integrated history import commit | 162dbe4a546c95a768c442f6779c4670a031465b (two parents) and d78287f6b1e3b551d5c41a52765e95eb6f7668af |
| Rewritten standalone module commit | 335b7e75ae0a0825015c1fcd19e7dcef5407dd3d; PR #2000 Core addition 0680dae87fbf3f7e683fb9c432b46678e344a32c |
| Independently verified Creator tests and Windows package | Core CI run 40036689755 PASS; unsigned Windows candidate run 40036493128 PASS |
| OS consumer integration and deletion commits | Native #2197 and Candidate #2199 and Windows toolkit #2200 integrated; legacy deletion PENDING |
| Physical write authorization | NONE |

History rule: a file copy with the old commit mentioned in text is insufficient. The imported Creator branch must contain genuine Git ancestry to the filtered historical Creator subtree and prove byte-identical source tree at the immutable OS source commit. No source or OS contract cutover before this receipt is completed.

Creator module is merged to main at d1d32e4858ab7e2f3629fba2ac4dde5e1e1e97ef. OS producer and physical release authority not cut over; unsigned builds are review-only.

## 2026-10-23 additional exact-head acceptance

- Creator Windows-native development code and CI merged on main via PR #6: `e7432abf30f178a495196b576cb7fd212006d45e`.
- [Linux Go tests/immutable OS contract integration](https://github.com/ordaxsystems/ordax-creator/actions/runs/40037999005): PASS.
- [Windows 2025 native Win32 control smoke, unit/security tests, read-only CLI execution](https://github.com/ordaxsystems/ordax-creator/actions/runs/40037999017): PASS.
- [Unsigned Windows deterministic ZIP](https://github.com/ordaxsystems/ordax-creator/actions/runs/40037999018): PASS.
- OS source Creator tree has not diverged at OS main `70b0f7eb5d5085be0d4f176f5232def37c7a2a24`: `cc287d6dafba00b91a1c3e2f86d154ad4a943cd0`.
- OS Native source helper consumer, OS producer/promoter and public signing remain not cut over or authorized.

## 2026-10-23 external OS Native consumer proposal

- Proposed, **not merged**: `ordax-os#2197@082d5afc374f2e667b53c8003a12076b9438410e`, based on OS `f4119971e8e73b49b6e0c9310b64240881e4d514`.
- Canonical OS external source pin: Creator `8adf641fb2cd9a3ac8351f23f912ce29ee2b70c9`, Git tree `312d18bb8e5cb97a216a54697724bbc8893d33d0`.
- Native release pipeline `40039904876` PASS, Surface external consumer `40039904975` PASS, disposable Native storage `40039904834` PASS. No OS cutover was merged. All physical release authorization and public signing remain unchanged.

## 2026-10-23 C5 merge and C6 preparation

- OS Creator Native helper consumer **merged** via `ordax-os#2197`, merge commit `f07fe50d414e883cc88b11988008e799b17366be`. Its 5 exact-head CI checks all PASS; OS release publisher/signatures and physical authorization remain OS-owned.
- C6 second consumer candidate **draft** `ordax-os#2199@dc8fcf6c2be7be768f48a208d6958dfcbe5ce475`: shared immutable source checkout plus OS Creator Candidate build from `ordax-creator`. No duplicate Go module, no old-code deletion and no publication. Await exact-head/merge evidence.
- Safety: `tools/creator/physical_promotion.py` OS source context still contains legacy Go paths; do **not** remove/move them before re-binding authorization and demonstrating fail-closed producer gate.

## 2026-10-23 C6 unsigned OS candidate migrated

- OS #2199 merged in `main@7677226def5fd027d37a84a8c7981c3c90d49b10` after PR-head Creator Candidate 40040816780, Native Release Pipeline 40040816790 and Foundation Contract 40040816774 — all PASS.
- OS still holds publisher signing, physical authorization and the historically tracked `tools/creator/` pending additional active-workflow cutovers. No official signed/public or physical release was produced by this migration.

## 2026-10-23 Windows toolkit independent Creator source accepted

- `ordax-os#2200` merged as OS `main@e01f4a92d43278fa66ebc68847469dfbdbc47c11`. Exact PR head `4151d23b9b18ce619d0f532d8cb65f83a061b42e` passed all four required CI workflows: Windows Toolkit 40041985452, Creator Candidate 40041985411, Native Release Pipeline 40041985404, Foundation 40041985438.
- Toolkit Creator implementation code is now sourced from pinned `ordax-creator@8adf641fb2cd9a3ac8351f23f912ce29ee2b70c9` (Git tree `312d18bb8e5cb97a216a54697724bbc8893d33d0`), while the OS remains source authority for the release-signing Go module, trust ceremony, system.tar and physical authorization. Both separate source identities appear in CI provenance.
- Common exact-source action now normalizes CRLF when used by native Windows Git Bash, with repository and SHA/tree validation preserved. No Windows code-signing certificate or public physical write permission was granted.

## 2026-10-23 Creator Authenticode test source migrated

- `ordax-os#2201` merged to OS `main@0d2011e6f704ba7ebdfaba5435cf77adf378e175`; exact PR-head CI Windows 40043640272 PASS and Foundation 40043640248 PASS.
- Windows unsigned GUI/launcher source is pinned to independent Creator commit `8adf641fb2cd9a3ac8351f23f912ce29ee2b70c9` and Git tree `312d18bb8e5cb97a216a54697724bbc8893d33d0`. OS retains Authenticode/Ed25519 publication contracts, verification/finalizer scripts, publisher Go tools, trust keys and physical authorization; no public release was made.
