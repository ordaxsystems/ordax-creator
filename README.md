# OrdaX Creator

Standalone source of **OrdaX Creator**, using the verified-media preparation Core extracted from [ordax-os](https://github.com/ordaxsystems/ordax-os). The Go module now uses `github.com/ordaxsystems/ordax-creator` and includes the existing Windows Win32 GUI, platform adapters and three signed-consumer channels.

**Migration status: C6 source extraction and OS consumer cutover COMPLETE (2026-10-24).** OS PRs [#2222](https://github.com/ordaxsystems/ordax-os/pull/2222) and [#2000](https://github.com/ordaxsystems/ordax-os/pull/2000) removed the duplicate Go implementation and integrated the independent 17-source bindings. OS-owned Python media producers remain under `ordax-os/tools/creator/` by design. The `ordax-os` remains sole producer/authority for Base, kernel, boot, official releases, media/storage layouts, physical promotion and signing policy. **No signed Stable release or live hardware write is authorized by C6.**

## Development and validation

- Core and safe CLI: `go test -count=1 ./...` and `go vet ./...`. Without an explicit OS contract integration root, only OS-contract-dependent integration cases skip; other Core tests run offline.
- Exact upstream contract verification: pass `ORDAX_OS_CONTRACTS_DIR` pointing to the read-only `docs/contracts` of a pinned and verified `ordax-os` Git commit. CI acquires the exact recorded OS revision automatically, without committing contract copies here.
- Safe Windows cross-build: `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildvcs=false -o OrdaX-Creator-App.exe ./cmd/ordax-creator-app`. Default build leaves the physical publisher trust binding **UNRESOLVED**, excludes tagged RAW backend and cannot enable physical apply.
- The Windows development ZIP is **unsigned, non-destructive and review-only**. It contains no elevated physical writer; its CI receipt records identical repeated binaries, SHA-256 and provenance. It is neither a verified public signed release nor proof that hardware writes are authorized.


Disposable OS media/ESP/QEMU proof producers live only in [`ordax-os/tools/creator/proof/`](https://github.com/ordaxsystems/ordax-os/tree/main/tools/creator/proof); they are not part of the standalone Creator source. The Creator-owned `componentchannel` packager remains here because the OS component candidate workflow consumes it from the pinned independent checkout.

## Single-owner boundaries

Creator alone will own Core, UI, adapters, source verification and narrow physical consumer helper. OrdaX Desktop must reuse the same Core. The OS alone owns authoritative boot/media contracts, OS artifact producer, release signing and physical promotion. No second image builder, artifact catalog, updater, signature root, 17-source inventory or write-consent mechanism may be added.

See [AGENTS.md](AGENTS.md), [migration execution ledger](docs/MIGRATION.md), [architecture](docs/ARCHITECTURE.md), [provenance](docs/PROVENANCE.md), and the merged [OS PR #2000](https://github.com/ordaxsystems/ordax-os/pull/2000).

**Safety:** never operate on a physical USB/SSD/HDD without separate explicit authorization for that exact target and completed publisher/consumer gates. CI uses disposable files only.

## Versionamento e builds de desenvolvimento

`version.json` é a **única fonte** da versão de produto do Creator. A versão
SemVer muda quando uma revisão do produto é promovida (por exemplo, de
`0.6.2-dev` para `0.6.3-dev`); builds de OS ou execuções repetidas do CI **não**
aumentam automaticamente o SemVer. Cada build tem também o commit exato do
Creator e o commit de origem do OS no manifesto de proveniência. A interface
mostra um identificador curto do commit ao lado da versão para distinguir
builds da mesma versão.

Todos os workflows de candidato e validação nativa incluem `version.json` nos
filtros de disparo: uma mudança na versão não pode passar sem reexecutar os
testes. O OS consome um único commit/tree pinado deste repositório no contrato
`docs/contracts/native-release-assembly.json`; a nova versão **não entra nos
pacotes USB do OS** enquanto esse pin não for atualizado e testado.

O catálogo público de downloads pertence ao publicador oficial do OS, não a
este arquivo. Uma versão Dev ou artefato de Actions não autoriza publicar uma
release Stable no site.
