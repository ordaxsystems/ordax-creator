# OrdaX Creator

## Governança documental multi-repositório

A [entrada canônica da documentação transversal](https://github.com/ordaxsystems/ordax-os/blob/main/docs/README.md) e o [contrato de ownership](https://github.com/ordaxsystems/ordax-os/blob/main/docs/contracts/repository-ownership.json) definem **onde encontrar a autoridade**, não transferem para o OS o código, as APIs locais ou as responsabilidades dos outros owners.

- Normas globais de arquitetura, composição, segurança e integração pertencem aos documentos e contratos canônicos indicados pelo índice do OS. Não duplicar essas regras aqui nem criar um SSOT paralelo.
- Este repositório continua autoridade para seu próprio source, contratos que produz, guias de operação, testes e evidências. Documentos locais não podem alterar silenciosamente contratos de outro owner.
- Em conflito, prevalecem o source/contrato vigente do owner aplicável e o índice de autoridades; planos, handoffs e evidências datadas não substituem estado vivo.
- Alterações cross-repo devem atualizar a autoridade primária e somente os links/consumidores afetados, sem sincronização artificial de versões, rebuilds ou PRs.


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
esse único script, sem recompilar o OS e sem alterar o canal Stable.
A descoberta de um Creator aberto usa diretamente o snapshot de processos do
Windows (Toolhelp32), sem `tasklist` ou subprocesso sem limite em sessão
não interativa. O transporte Git para `ORDAX-DATA` continua sujeito ao serial
autorizado e ao bloqueio de instâncias concorrentes.

Este é um mecanismo local de desenvolvimento; não é assinatura
Authenticode nem atualização silenciosa do aplicativo público.

### Descoberta rápida no executável Owner único

No Creator Owner Windows compilado como `ordax_single_exe`, a descoberta
`targets`/`status` agora usa o próprio executável instalado. Iniciar a
interface **não** extrai/copia outro EXE de centenas de MB nem uma segunda
imagem de Base para `%TEMP%\OrdaX-Creator`. O binário verifica a
proveniência incorporada e a identidade do caminho antes de invocar suas
operações estritamente de leitura.

Somente quando o usuário confirma **Criar OrdaX** é que a mesma
implementação existente materializa o backend isolado e a seed necessários
para o preparo físico. A atualização incremental via Git continua no único
transportador canônico; não requer materializar nem escrever a imagem Base.
Nenhum segundo atualizador, cache de fonte, writer, serviço ou build do OS foi
introduzido. O CI Windows prova a variante `ordax_single_exe` com seed
**sintética**, sem enumerar ou tocar mídia física.

## Reconexão USB no Creator Windows

Com a interface do Creator aberta, inserção, remoção ou mudança de topologia
USB provoca uma nova descoberta de dispositivos. O refresh reutiliza a
identidade física (serial/índice/token) e não seleciona discos novos. Inserção
também consulta o canal de atualização do próprio Creator. Eventos recebidos
durante uma descoberta são reconciliados ao terminar; uma gravação ativa ou
um resultado terminal não é interrompido.

No perfil Owner Dev, ao abrir o Creator com o USB conectado ou receber um
evento de inserção, a descoberta já existente também consulta **em segundo
plano** a identidade autorizada na tarefa Windows canônica. Somente quando
existe **um único** USB físico correspondente ao serial autorizado, sem
escrita física/atualização concorrente nem seleção explícita de outro disco,
a interface inicia a preparação incremental pelo mesmo sincronizador do
botão e da tarefa. USB desconhecido ou ambíguo nunca é atualizado
silenciosamente. O evento não seleciona um disco para a ação destrutiva
**Criar OrdaX**, não formata, não grava a ESP e não ativa o Linux no Windows.

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
da origem e handoff do Supervisor.

O Supervisor Owner Dev já consome este cache em modo somente leitura e só
reutiliza objetos Git quando a `main` HTTPS confirma exatamente o mesmo commit.
A cópia do bundle não equivale à instalação ou à ativação do OS.

### Cache Git independente do checkout do desenvolvedor

GUI e tarefa Owner Dev usam o mesmo cache de entrega existente, adquirido
diretamente da `main` canônica do OS. Nenhum checkout local do desenvolvedor
é necessário ou modificado. O cache mantém objetos próprios para o bundle;
um cache antigo com alternates, shallow ou objetos promisor é substituído
somente após validar o novo checkout. Atualizações subsequentes fazem
fast-forward do mesmo cache; o escritor ainda confirma a main remota antes
de tocar a mídia autorizada.

### Preparação automática na tarefa Owner Dev já existente

O agendador local `OrdaX Creator Dev Release Sync` pode receber o argumento
`--stage-owner-dev-usb-serial SERIAL_LOCAL`. Não cria outra tarefa, serviço
ou atualizador: primeiro prepara o cache USB pelo transportador canônico,
depois sincroniza a instalação do Creator. O mesmo script
`tools/owner_usb_git_delivery.py` é executado com `--auto-volume --apply`.

A busca de volumes usa a API nativa do Windows somente para encontrar a
partição exFAT chamada `ORDAX-DATA`. A identidade continua validada pelo
**único verificador físico existente**, exigindo GPT ESP/Linux/DATA, conexão
USB, disco diferente do sistema e comparação exata com o serial físico
autorizado localmente. Não gravar o serial no repositório público.

Quando o pendrive não está conectado, o resultado é
`SKIP_TRUSTED_USB_NOT_PRESENT`; quando já contém a revisão verificada,
`CURRENT`. Enquanto o Creator está aberto, a **tarefa automática** não disputa o
dispositivo; a ação manual na interface usa a mesma exclusão mútua e exige
seleção explícita. O pacote Git é gerado apenas quando a `main` local coincide
com a remota, e sempre passa por SHA-256 e readback.

**A preparação não altera boot, ESP, ext4, seleção de slot, Git ativo nem
instala automaticamente código não assinado.** A ativação permanece
exclusiva do Supervisor e do `ordax-pull`, após comprovação da identidade
do candidato. Também não transforma um cache de desenvolvimento em release
Stable assinada.

### Botão "Atualizar pendrive" no Creator Owner Development

Ao selecionar explicitamente um USB elegível, o Creator oferece **Atualizar pendrive**
sem formatar a mídia. O botão exige o binário Owner Dev autorizado e a tarefa
local **já existente** `OrdaX Creator Dev Release Sync`, com serial físico
localmente configurado, instância única e checkout canônico de desenvolvimento.
O aplicativo confere também o número do disco selecionado; não aceita uma
segunda unidade com o mesmo rótulo e não escolhe automaticamente um destino.

A interface entra no **mesmo sincronizador** `tools/owner_local_sync.py`
com `--creator-ui-stage`. Esse modo reutiliza o cache exato da `main`,
prepara o Git e, quando o OS comprova uma Base publicada compatível, usa
o **mesmo** transportador `tools/owner_usb_git_delivery.py` para copiar
o candidato binário verificado para `ORDAX-DATA`, sem atualizar o boot.
Não substitui o executável aberto. A tarefa agendada mantém seu
comportamento de sincronizar o Creator quando ele estiver fechado;
o modo de interface não cria serviço, segundo updater ou segundo
checkout de código. O cache de binários `owner-base-delivery` apenas
reutiliza os bytes autenticados pelo adquiridor canônico do OS.
A preparação Git do cache compartilhado fica serializada entre GUI e tarefa
pelo bloqueio nativo `owner-main-delivery.lock`, mantido durante o transporte.
O bloqueio físico `owner-usb-git.lock` já existente permanece como autoridade
única de escrita do USB, inclusive para chamadas diretas ao transportador.
Conflitos momentâneos são aguardados de forma limitada na interface; a tarefa
agendada adia a tentativa para seu próximo ciclo. **Somente** contenção de
mutex permite repetição: serial, disco, origem, SHA ou volume inválidos nunca
são repetidos como se fossem conflitos transitórios.
As entradas devolvem recibos `CURRENT` ou `STAGED_NOT_ACTIVATED` após
conclusão efetiva; uma espera vencida não é anunciada como atualização. A interface aceita um
Creator cuja revisão esteja comprovadamente na ancestralidade da `main`
atual, sem exigir que o EXE já incorpore toda alteração posterior do Git;
a implementação local permanece limpa e exatamente igual à `main` remota. O único escritor valida o GitHub `main`, GPT,
serial, volume, SHA-256 e readback antes de publicar o manifesto. A interface
não executa `ordax-pull`, não abre `PhysicalDrive` nem escreve na ESP/ext4.

### Read-only USB Git transport inspection

The existing canonical transport tool supports a read-only --inspect mode
with --auto-volume, --trusted-serial and --trusted-disk-number. It requires
no OS checkout or network connection, and rechecks the physical device identity
before and after inspection. Local validation covers the exact manifest,
SHA-256, file length and internal Git bundle head.

GIT_BUNDLE_LOCAL_VERIFIED proves only local integrity of the transport stored
in ORDAX-DATA. It does NOT authenticate the remote main, prove the Surface
EROFS was installed, inspect Linux Base/Kernel, or permit activation.
GIT_BUNDLE_NOT_STAGED means no transport was staged; incomplete or tampered
transport is blocked. Creator's post-stage summary separates Git transport
from executable evidence and OS/Supervisor activation.

**Primeira instalação da Base:** continua sendo a ação separada
**Criar OrdaX**, com confirmação de formatação e UAC. O botão incremental
não substitui uma Base antiga que ainda não possui o consumidor do Supervisor.
A ativação/rollback permanece exclusivamente no OS. Fora de um ambiente
Owner Dev configurado, a atualização incremental é bloqueada, sem procurar
Python ou scripts arbitrários nem introduzir um caminho público alternativo.

## Base DEV no Creator Windows — aquisição e transporte sem ativação

A ação de preparo do pendrive conserva o fluxo Git existente e, somente
quando o OS identifica uma Base compatível, usa o **mesmo código**
`ordax-os/system/services/base-update/dev_channel.py` extraído de `main`.
`--resolve-only` compara o SHA solicitado com as tags publicadas e a política
única de impacto. `--acquire-only` recebe os quatro arquivos binários sob
SHA-256 e limites do produtor: `dev-base.json`, `vmlinuz`, `initrd.gz`,
`rootfs.tar`. A aquisição é nativa do Windows; sem WSL, distro auxiliar
nem reimplementação do resolvedor.

O Creator valida o serial, número do disco, GPT, volume e espaço disponível;
sua única escrita é em `ORDAX-DATA` sob o mesmo mutex físico
`owner-usb-git.lock`. `--import-only` do OS publica o candidato imutável em
`.ordax/updates/owner-base/<artifact_sha>` somente depois do readback,
sem sobrescrever uma revisão existente. O recipiente mostra
`STAGED_NOT_ACTIVATED` somente para os arquivos Base efetivamente
copiados e verificados. `AVAILABLE_NOT_STAGED` significa que a Base está
publicada, mas a cópia ainda não foi comprovada ou não coube na mídia.
Os recibos do Git são independentes dos da Base; a tarefa agendada continua
responsável apenas pelo transporte Git não ativador.

**Não há promessa de instalação:** essa transferência não altera ESP, ext4,
EFI, slots A/B, Kernel executado, Surface nem o estado de saúde.
O USB Owner Dev preparado possui apenas 1 GiB de `ORDAX-DATA`, podendo não
comportar os artefatos publicados. O OS precisa consumir o candidato no boot
por seu mecanismo único de update/rollback, e uma Base antiga não recebe
esse consumidor só porque a partição exFAT tem novos arquivos. A missão
`conectar, clicar e inicializar o notebook já atualizado` só pode ser
declarada concluída após prova física de boot, readback da versão ativa,
Wi-Fi e reversão de uma candidata inválida. Não é permitido transformar
transporte em ativação pelo texto da interface.
