# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make setup       # install go tools (goimports, golangci-lint) + go mod tidy/download + pre-commit hooks
make mod         # go mod tidy && go mod download
make lint        # golangci-lint run ./...
make format      # go fmt ./... && goimports -w .
make test        # go test -v -race ./...
make build        # builds bin/crm-api-core from main.go (note: entrypoint is cmd/api, verify path before relying on this target)
make clean        # rm bin/ + go clean -testcache
```

Single test: `go test -v -race ./internal/application/... -run TestName`

Local dev stack (Postgres + app w/ hot reload via `air`): `docker-compose up`. App reads env from `.env.local`, config from `resources/application.properties` (override path via `APP_PROPERTY_FILE` env var).

DB migrations (`golang-migrate`, files in `migrations/`) run automatically on app startup — see `internal/infra/repository/database/setup.go`. No manual migrate step needed for local dev.

## Architecture

Clean Architecture, three layers under `internal/`:

- **`internal/domain`** — entities and business rules (`Case`, `Customer`, `Partner`, `Contractor`, `Product`, `Transaction`, `Comment`, `User`, etc). No framework/infra dependencies. Includes domain-level interfaces implemented by outer layers (e.g. `domain.CaseBuilder`).
- **`internal/application`** — use cases / services (`*_service.go`), one per aggregate, orchestrating domain + repositories. Services are composed of other services (e.g. `CaseService` depends on `CustomerService`, `ProductService`, `UserService`, etc — see wiring in `internal/infra/runner.go`), not raw repositories, when cross-aggregate logic is needed.
- **`internal/infra`** — everything framework/tech-specific:
  - `entrypoint/rest` — Gin controllers + DTOs (one controller + one dto file per resource)
  - `entrypoint/middleware` — auth middleware
  - `repository/database` — Postgres repositories (raw SQL via `sqlx`-style DTOs, one file per aggregate)
  - `repository/bucket` — S3-backed attachment storage
  - `config` — env/properties-based `AppConfig` loading (`resources/application.properties` + `*Env`-suffixed fields resolved from OS env vars, e.g. `HostEnv`, `PasswordEnv`)
  - `runner.go` — manual dependency injection: wires config → repositories → services → controllers → routes → Gin router. This is the single place to look to understand how everything connects.

Entry point: `cmd/api/main.go` → `infra.RunApp()`.

Routing: all routes registered in `internal/infra/entrypoint/routes.go`, under prefix `/crm/core/api/v1`. Two route groups — `authGroup` (JWT via `AuthenticationMiddleware`) and `publicGroup` (login, user creation, inbound webhook `/web/message`).

### Partner-specific case ingestion (builder pattern)

`internal/application/builder/` implements `domain.CaseBuilder` per partner (`assurant.go`, `ezze.go`, `luiza_seg_cardif.go`, `default.go`) to parse partner-specific CSV column layouts into `Case`/`Product` domain objects for batch case creation (`BatchCaseService`, `POST /cases/batch`). When adding a new partner integration, add a new builder here rather than branching inside the batch service.

## Tooling notes

- Linter: `golangci-lint` (`.golangci.yml`) — enabled: errcheck, govet, ineffassign, staticcheck, unused, bodyclose, goconst, gocritic, gocyclo (max complexity 15), misspell, nilerr. `make lint` before considering work done.
- `pre-commit` hooks configured via `.pre-commit-config.yaml`, installed by `make setup`.
- Mocks live in `mock_domain`/`mock_application` subpackages, generated (check header of a mock file for the generator used before hand-editing).
