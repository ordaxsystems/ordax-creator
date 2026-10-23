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
