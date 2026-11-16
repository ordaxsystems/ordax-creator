# OrdaX Creator

Standalone source of **OrdaX Creator**, using the verified-media preparation Core extracted from [ordax-os](https://github.com/ordaxsystems/ordax-os). The Go module now uses `github.com/ordaxsystems/ordax-creator` and includes the existing Windows Win32 GUI, platform adapters and three signed-consumer channels.

**Migration status: C6 source extraction and OS consumer cutover COMPLETE (2026-10-24).** OS PRs [#2222](https://github.com/ordaxsystems/ordax-os/pull/2222) and [#2000](https://github.com/ordaxsystems/ordax-os/pull/2000) removed the duplicate Go implementation and integrated the independent 17-source bindings. OS-owned Python media producers remain under `ordax-os/tools/creator/` by design. The `ordax-os` remains sole producer/authority for Base, kernel, boot, official releases, media/storage layouts, physical promotion and signing policy. **No signed Stable release or live hardware write is authorized by C6.**

## Relação com o restante do OrdaX

Creator recebe um **candidato de mídia OrdaX já produzido e verificável** pelo owner do OS. Ele não recompila kernel/Base/Surface/Apps nem precisa acompanhar cada alteração dos outros repositórios. Seu trabalho começa no artefato direto: verificar identidade/proveniência, obter consentimento do alvo, gravar e conferir readback.

Isso permite Surface, Desktop, Apps e OS evoluírem em paralelo sem transformar o Creator em gargalo. [Regra canônica de composição](https://github.com/ordaxsystems/ordax-os/blob/main/docs/DEVELOPMENT-WORKFLOW.md#integração-multi-repositório--composição-do-produto-2026-11-13).

## Development and validation

- Core and safe CLI: `go test -count=1 ./...` and `go vet ./...`. Without an explicit OS contract integration root, only OS-contract-dependent integration cases skip; other Core tests run offline.
- Exact upstream contract verification: pass `ORDAX_OS_CONTRACTS_DIR` pointing to the read-only `docs/contracts` of a pinned and verified `ordax-os` Git commit. CI acquires the exact recorded OS revision automatically, without committing contract copies here.
- Safe Windows cross-build: `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildvcs=false -o OrdaX-Creator-App.exe ./cmd/ordax-creator-app`. Default build leaves the physical publisher trust binding **UNRESOLVED**, excludes tagged RAW backend and cannot enable physical apply.
- The Windows development ZIP is **unsigned, non-destructive and review-only**. It contains no elevated physical writer; its CI receipt records identical repeated binaries, SHA-256 and provenance. A dependent Windows runner validates the **packaged ZIP**, its exact commit and binary hashes, then runs only the extracted read-only inspection status and the extracted GUI's `--offline-package-smoke` (native controls, no services or device discovery). A green package test is **not** a signed public release, a full interactive Windows installation or proof that hardware writes are authorized.


Disposable OS media/ESP/QEMU proof producers live only in [`ordax-os/tools/creator/proof/`](https://github.com/ordaxsystems/ordax-os/tree/main/tools/creator/proof); they are not part of the standalone Creator source. The Creator-owned `componentchannel` packager remains here because the OS component candidate workflow consumes it from the pinned independent checkout.

## Single-owner boundaries

Creator alone **owns** Core, UI, adapters, source verification and the narrow physical consumer helper. OrdaX Desktop must reuse this same independent Core when it presents or integrates Creator functionality. The OS alone owns authoritative boot/media contracts, OS artifact producer, release signing and physical promotion. No second image builder, artifact catalog, updater, signature root, 17-source inventory or write-consent mechanism may be added.

See [AGENTS.md](AGENTS.md), [migration execution ledger](docs/MIGRATION.md), [architecture](docs/ARCHITECTURE.md), [provenance](docs/PROVENANCE.md), and the merged [OS PR #2000](https://github.com/ordaxsystems/ordax-os/pull/2000).

**Safety:** never operate on a physical USB/SSD/HDD without separate explicit authorization for that exact target and completed publisher/consumer gates. CI uses disposable files only.

## Versionamento e builds de desenvolvimento

`version.json` é a **única fonte** da versão de produto do Creator. A versão
SemVer muda quando uma revisão do produto é promovida (por exemplo, de
`0.6.2-dev` para `0.6.3-dev`); builds de OS ou execuções repetidas do CI **não**
aumentam automaticamente o SemVer. Cada build tem também o commit exato do
Creator e o commit de origem do OS no manifesto de proveniência. A interface
mostra um identificador curto do commit **do executável Creator** ao lado da
versão para distinguir builds da mesma versão. O ZIP Dev Windows incorpora os
dois valores diretamente em `OrdaX-Creator-App.exe` a partir de `version.json`
e do commit do CI. O teste do ZIP extraído confirma esses valores no executável
sem iniciar os serviços; uma versão do componente `creator-dev` baixada
separadamente não substitui a versão mostrada do GUI. Esses valores são
identificação de desenvolvimento, **não** assinatura/autorização física.

Todos os workflows de candidato e validação nativa incluem `version.json` nos
filtros de disparo: uma mudança na versão não pode passar sem reexecutar os
testes. O OS consome um único commit/tree pinado deste repositório no contrato
`docs/contracts/native-release-assembly.json`; a nova versão **não entra nos
pacotes USB do OS** enquanto esse pin não for atualizado e testado.

O catálogo público de downloads pertence ao publicador oficial do OS, não a
este arquivo. Uma versão Dev ou artefato de Actions não autoriza publicar uma
release Stable no site.

## Owner Development: atualização local verificada

O código Git atualizado **não** substitui o executável aberto. O publicador do OS
entrega apenas mídia autenticada e QEMU-verificada no canal
`ordax-os/creator-owner-media`; o publicador independente deste repositório
gera `creator-owner-prototype` com SHA-256 e os dois commits de origem.

No Desktop de desenvolvimento, `python tools/owner_local_sync.py` sincroniza
a instalação existente em `~/Downloads/OrdaX-Creator-Owner-Prototype/current`.
O sincronizador usa somente o manifesto oficial, verifica SHA-256, progressão
Git em ambos os owners, membros do ZIP, proveniência e o smoke offline do EXE.
Se o Creator estiver aberto ou ainda não houver release, preserva a instalação;
a troca ocorre no próximo ciclo ocioso. A tarefa agendada do Windows pode chamar
esse único script, sem recompilar o OS, sem tocar USB e sem alterar o canal
Stable. Este é um mecanismo local de desenvolvimento; não é assinatura
Authenticode nem atualização silenciosa do aplicativo público.

## Reconexão USB no Creator Windows

Com a interface do Creator aberta, inserção, remoção ou mudança de topologia
USB provoca uma nova descoberta de dispositivos. O refresh reutiliza a
identidade física (serial/índice/token) e não seleciona discos novos. Inserção
também consulta o canal de atualização do próprio Creator. Eventos recebidos
durante uma descoberta são reconciliados ao terminar; uma gravação ativa ou
um resultado terminal não é interrompido.

**Limite atual:** detectar o USB **não é** autorização para formatá-lo.
O fluxo físico existente exige seleção, pré-validação, UAC e readback. O
`ordax-pull` pertence ao OS em execução; o Windows não executa diretamente
o shell Linux do pendrive. O release-agent do OS ja aceita pacotes assinados com --local-bundle. Ainda falta integrar a entrega Portable oficial e a ativacao no boot. Nenhum arquivo de atualizacao sem consumidor sera gravado na particao de dados.

## Entrega Git offline para USB Owner Dev

`tools/owner_usb_git_delivery.py` prepara um transporte incremental Git
para a partição **ORDAX-DATA** (exFAT) do USB de desenvolvimento. O script
somente aceita repositório `ordaxsystems/ordax-os/main` local limpo e
igual à `main` remota, disco físico USB GPT com ESP/Linux/DATA, volume
`ORDAX-DATA` e dispositivo diferente do disco do Windows. Sem `--apply`
faz apenas preflight; com `--apply` copia `ordax-main.bundle` e o
manifesto SHA-256 para `.ordax/updates/owner-git`, com readback e
publicação do manifesto por último. Não formata partições, não abre
`PhysicalDrive`, não altera o kernel e não dispara `ordax-pull`.

No OrdaX OS, o helper proprietário
`bootstrap/dev-base/ordax-owner-usb-git` já consegue verificar/importar
esse transporte sem ativá-lo. **O pacote não é assinado**: ele é de uso
exclusivo Owner Development e não deve autoativar sem comprovação adicional
da origem e handoff do Supervisor. A reconexão no Creator ainda não chama
este script automaticamente, e o pendrive físico não é modificado até o
consumidor estar integrado e validado.
