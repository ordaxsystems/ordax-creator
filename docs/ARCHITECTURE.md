# OrdaX Creator — Architecture and security boundary

Status: TARGET ARCHITECTURE; not yet migrated or released.

Creator is a standalone Windows-first verified-media preparation application. Its single host-neutral Creator Core is reusable from the future OrdaX Desktop without implementing a second flasher. The pre-existing code in ordax-os/tools/creator is the sole extraction source until verified cutover.

## Execution ownership
- Creator: graphical shell and localization, target discovery/selection, source acquisition, verified signed cache, host-neutral plan and validation, target-specific explicit confirmation, narrow Windows elevated adapter, post-write verification, app/component/physical consumer channels.
- OS: Base, kernel, initramfs, bootstrap, OS release bytes, machine-readable boot/media/storage contracts, release trust anchor, producer/publisher and physical promotion authorization.
- All OS contract/artifact inputs must be immutable, source-revision pinned and SHA-256/size checked. A downloaded file is a verified consumer copy, never a new source of truth.
- Do not duplicate the 17-source inventory, 39-operation planner, media geometry, signature policy, updater, physical authorization, or actual OS release generator.

## Security boundary
1. Acquire and validate the signed release envelope and version/source commitment.
2. Validate exactly the canonical 17 flat regular non-symlink files with immutable size/SHA-256 bindings.
3. Enumerate supported removable targets, exclude system/internal disks and detect device identity/capacity/serial changes.
4. Prepare Core's exact plan and target-specific destructive confirmation; require explicit human consent and Windows elevation independently.
5. Elevated writer reopens/locks verified source bytes and revalidates target before any physical writes.
6. Flush, independently read back 17 artifacts and final geometry/filesystems/capacity; fail closed on partial results.

Code/build tests are NON-DESTRUCTIVE and cannot grant 4-6. Public signing status and current OS promotion gates remain independent.

## Migration layout
- Existing Core, cmd, host, appchannel, componentchannel, physicalchannel, internal/atomicfile and update packages become Creator-owned.
- Split proof scripts and root Python helpers by actual producer/consumer behavior; retain OS publisher and source-authorized signing scripts in OS.
- First import is quarantined in a migration-only branch; selected implementation moves to canonical paths only after classification and Go module/import rewriting.
- The OS consumes a tagged or SHA-bound released Creator capability; no permanent mirrored sources inside the OS.
