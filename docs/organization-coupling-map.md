# Mapa de acoplamento a um único tenant (`Organization`)

Este documento registra o que hoje está fortemente acoplado à premissa de que o CRM
atende um único cliente. Serve de referência para as próximas fases de migração
para multi-tenancy, depois da fundação introduzida em conjunto com este documento
(entidade `domain.Organization`, tabela `organizations`, coluna `organization_id`
nullable em `users`, `customers`, `partners`, `contractors`, `cases`, `products` e
`queues` — ver `internal/domain/organization.go` e as migrations `000030`/`000031`).

Levantamento feito por exploração em `internal/domain`, `internal/infra` e
`internal/application` em 2026-08-12.

## Resumo executivo

- Antes desta mudança, nenhuma tabela em `migrations/` tinha coluna de escopo
  (tenant/organização/cliente). Uma busca por `tenant`, `organization`, `org_id`,
  `company_id`, `client_id` em todo `internal/` não retornava nenhum resultado.
- O JWT emitido por `AuthService` carrega apenas `user_id`, `session_token` e `exp`
  — nenhuma claim de organização.
- `internal/infra/runner.go` é 100% singleton: uma conexão Postgres, um bucket S3,
  um segredo JWT, e uma única instância de cada repositório/serviço/controller para
  o processo inteiro (`internal/infra/config/setup.go`).
- A identidade de quem está chamando a API está fragmentada em três mecanismos que
  não se cruzam: a claim `user_id` do JWT (que a `AuthenticationMiddleware` sequer
  propaga para os controllers via contexto do Gin), o header `X-Author` (lido só
  em alguns controllers: `case_controller.go`, `case_action_controller.go`,
  `partner_controller.go`), e os campos `created_by`/`updated_by` enviados no corpo
  da requisição (o chamador é confiado sem checagem cruzada contra o JWT).

## Hardcodes de single-tenant mais explícitos

- **`internal/domain/user.go:94-99`** — `UserRole.THAVANNA_ADMIN` embute o nome da
  empresa operadora atual ("Thavanna") como valor de enum, misturando "admin de
  plataforma" com o nome de uma empresa específica.
- **`internal/domain/region.go:3-61`** — `regions` (mapa estado→região) e
  `AcronymForState` são hardcoded para os 26 estados + DF do Brasil, com
  agrupamento arbitrário em 7 regiões numéricas (0-6). Usado por
  `Customer.GetRegion()`, `Partner.GetRegion()` e no fallback de atribuição de
  dono de caso (`internal/application/case_service.go:164-172`, que filtra
  operadores candidatos por `Region`). Assume um único país/uma única taxonomia
  regional para toda a instalação.
- **`internal/domain/customer.go:42-46`** — `DocumentType` (`CPF`, `CNPJ`, `RG`) é
  específico do Brasil; usado tanto por `Customer` quanto por `Partner`.
- **`internal/domain/billing.go:11-15`** — `BillingType` só tem o valor `PIX`
  (sistema de pagamento instantâneo brasileiro).
- **`internal/application/batch_case_service.go:60-81`** — dispatch de ingestão em
  lote via `switch companyName { case "Assurant": ...; case "Cardif": ...; case
  "Ezze Seguros": ...; case "LuizaSeg": ...; default: ... }`, selecionando a
  implementação de `domain.CaseBuilder` (`internal/application/builder/`:
  `assurant.go`, `ezze.go`, `luiza_seg_cardif.go`, `default.go`) por comparação de
  string livre. O mapeamento partner→builder é duplicado e não validado de forma
  cruzada entre esse switch e o `GetCompanyName()` de cada builder.
- **`internal/application/report_service.go:27-32,188`** — `contractorsTemplates`
  é um `map[string]string` hardcoded de nome-do-contractor → arquivo de template
  `.docx` (`"LuizaSeg"`, `"Assurant"`, `"Cardif"`, `"Ezze Seguros"`); gerar
  relatório para qualquer `Contractor` fora dessa lista falha. Há ainda um
  `isAssurant := reportData.Contractor.CompanyName == "Assurant"` especial dentro
  de `replaceImages` que pula o último slot de imagem só para a Assurant.
- **`internal/application/partner_service.go:90-114`** — import de partners via
  CSV assume prefixo telefônico `"+55 "`, país `"Brazil"` e cabeçalhos de coluna em
  português (`"Nome"`, `"Documento"`, `"Telefone"`), além de inferir
  `DocumentType` pelo tamanho da string do documento (11 dígitos = CPF, 14 = CNPJ).
- **`internal/domain/case.go` (`CaseStatus`/`CasePriority`)** e
  **`internal/domain/comment.go` (`CommentType`)** — workflow fixo modelado no
  processo atual de sinistro de seguro (`WAITING_PARTNER`, `REPORT`, `PAYMENT`,
  `RECEIPT`, `COMMENT_RESOLUTION`, `COMMENT_REJECTION`), e não uma state machine
  genérica configurável por organização. Hoje não há validação de transição de
  status em `CaseActionService.ChangeStatus` — qualquer transição é aceita
  incondicionalmente.
- **`internal/application/builder/*.go`** — literais `"insurance"` (categoria) e
  `"csv"` (canal de origem) hardcoded em todo builder, i.e. o caminho de criação
  de casos em lote assume que o domínio de negócio é processamento de sinistro de
  seguro.

## Ambiguidade `Contractor` vs `Partner` vs `Organization`

`Contractor` já se comporta como um "sub-cliente" de fato dentro do tenant único
atual: dita o parsing de CSV na ingestão em lote, o template de relatório, e o
status default de caso por builder (Ezze/LuizaSeg/Cardif criam casos em `DRAFT`;
Assurant não sobrescreve o status). `Partner`, por outro lado, é uma entidade
completamente diferente — o prestador/reparador que executa o serviço, com
endereço de entrega, dados de pagamento (`Billing`), etc. Nenhum dos dois é a
mesma coisa que uma organização-cliente da plataforma (uma empresa que licencia o
CRM). Ao evoluir para multi-tenancy de verdade, é importante manter essa
distinção clara — `Organization` é um nível acima de ambos, não um substituto.

## Infra singleton

- `internal/infra/config/setup.go` — `AppConfig` carrega uma única
  `Database` (host/porta/schema/credenciais), um único `Bucket` S3, e um único
  `SecretJWTKey`, todos lidos uma vez de `resources/application.properties` no
  boot.
- `internal/infra/runner.go` — `RunApp()` monta um grafo de dependências
  singleton: uma conexão `*sqlx.DB`, um cliente S3, uma instância de cada
  repositório/serviço/controller, criados uma única vez no processo. Não há
  fábrica/registro parametrizável por tenant.
- `internal/infra/repository/database/setup.go` — migrations rodam
  automaticamente contra um único schema no startup; não há roteamento por
  schema/banco por tenant.

## Mecanismo reaproveitável para filtro futuro

Os helpers genéricos de WHERE em `internal/infra/repository/database/utils.go`
(`prepareInQuery`, `prepareLikeQuery`, `prepareGreaterEqualQuery`,
`prepareLesserEqualQuery`, `prepareMetadataContainsQuery`) já são o padrão certo
para acoplar um filtro de `organization_id` no futuro — é o mesmo mecanismo usado
pelo filtro de metadata recém-adicionado em `case_controller.go`
(`?metadata[categoria]=X`). Hoje nenhum repositório injeta esse tipo de cláusula
automaticamente; todo filtro é opt-in e definido pelo que o controller monta em
cada `*Filters`.

## Próximos passos sugeridos (fora do escopo desta primeira fase)

1. **Enforcement de isolamento**: transformar `organization_id` de campo opcional
   em filtro obrigatório nos `Search`/`GetByID` de cada repositório, hoje globais
   (nenhum aggregate — `Case`, `User`, `Queue`, `Contractor`, `Partner`, `Product`,
   `Customer` — tem qualquer filtro por organização aplicado automaticamente).
2. **Identidade de organização na autenticação**: adicionar claim de
   `organization_id` ao JWT (`AuthService.CreateToken`) e propagar de fato via
   contexto do Gin (`AuthenticationMiddleware` hoje seta `user_id` no contexto mas
   nenhum controller o lê).
3. **Papéis/permissões por organização**: decidir o destino de
   `UserRole.THAVANNA_ADMIN` — vira um papel de "admin de plataforma" cross-tenant
   explícito, ou é substituído por um mecanismo de permissão por organização?
4. **Regionalização configurável**: extrair `internal/domain/region.go` de um mapa
   hardcoded de estados brasileiros para algo configurável por organização, para
   suportar clientes fora do Brasil.
5. **Generalizar o workflow de caso**: avaliar se `CaseStatus`/`CommentType`
   devem continuar como enums globais fixos ou se cada organização precisa de seu
   próprio workflow.
6. **Revisar o dispatch de builders**: o `switch companyName` em
   `batch_case_service.go` deveria migrar de comparação de string livre para um
   identificador estruturado (possivelmente ligado a `organization_id` e/ou
   `Contractor`), unificando as duas fontes de verdade hoje divergentes (o switch
   e `GetCompanyName()` de cada builder).
