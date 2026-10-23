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
| OS consumer integration and deletion commits | PENDING |
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
