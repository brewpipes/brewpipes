# BrewPipes Agent Guide

This guide is for agentic coding tools working in this repo.
It captures commands and conventions observed in the current codebase.

## Project quick facts
- Module path: `github.com/brewpipes/brewpipes`
- Language/runtime: Go `1.25.5` (see `go.mod`).
- Primary services: `service/identity`, `service/production`.
- Entrypoints: `cmd/monolith`, `cmd/identity`, `cmd/production`.
- HTTP server: aggregated by `cmd.RunServices` on `:8080`.
- Database: Postgres via `pgx/v5`, migrations per service.

## Repo layout
- `cmd/` contains service entrypoints and shared lifecycle helpers.
- `service/<name>/` contains service logic, HTTP handlers, and storage.
- `service/<name>/storage/migrations/` holds SQL migrations.
- `internal/` holds shared packages like `database` and `jwt`.
- `docker-compose.yml` starts Postgres + monolith container.

## Environment variables
- `POSTGRES_DSN` is required by all services.
- Local DSN used in Makefile:
  `postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable`
- `BREWPIPES_SECRET_KEY` is required for identity/JWT issuance and for access token verification in production, inventory, and procurement services.

## Setup and run
- Start Postgres: `make start-postgres`
- Stop and delete volumes: `make nuke`
- Connect with psql: `make psql`
- Run monolith: `make run`
- Run monolith in background: `make run-bg`
- Run identity service:
  `POSTGRES_DSN=... BREWPIPES_SECRET_KEY=... go run ./cmd/identity`
- Run production service:
  `POSTGRES_DSN=... go run ./cmd/production`
- Compose stack (monolith + postgres): `docker compose up`

## Build commands
- Build monolith: `go build ./cmd/monolith`
- Build identity: `go build ./cmd/identity`
- Build production: `go build ./cmd/production`
- Build all packages: `go build ./...`

## Test commands
- Run all tests: `go test ./...`
- Run package tests: `go test ./service/production/handler`
- Run a single test:
  `go test ./service/production/handler -run TestVolumesHandler`
- Run with verbose output:
  `go test ./service/production/handler -run TestVolumesHandler -v`
- Disable caching (when debugging):
  `go test ./service/production/handler -run TestVolumesHandler -count=1`

## Lint and formatting
- Formatting is standard `gofmt`.
- Format all packages: `go fmt ./...`
- Or format the repo: `gofmt -w .`
- Vetting: `go vet ./...`
- No configured golangci-lint or staticcheck in this repo.

## Code style and conventions
- Use `gofmt`-compatible layout (tabs, aligned struct fields).
- Keep imports in three blocks: standard library, blank line, third-party/local.
- Prefer `log/slog` for structured logs in services/handlers.
- Some legacy helpers use `log.Println`; keep style consistent within a file.
- File names are lower `snake_case.go` (e.g., `volumes_handler.go`).
- Package names are lower-case and match folder names.
- Exported identifiers use PascalCase; unexported use lowerCamelCase.
- Public types and constants are documented with a short comment.

## Error handling
- Wrap errors with context using `fmt.Errorf("action: %w", err)`.
- Use sentinel error `service.ErrNotFound` for missing DB records.
- Prefer returning errors to callers; avoid panics.
- For HTTP handlers:
  - Use `service.JSON(w, payload)` for JSON responses.
  - Use `service.InternalError(w, err)` for 500s with correlation IDs.
  - Use `http.Error` for validation/authorization failures.
- Log server errors with enough context (method, path, or params).

## HTTP handlers
- Handlers typically return `http.HandlerFunc`.
- Accept dependencies via interfaces (e.g., `VolumeGetter`, `UserGetter`).
- Use `r.Context()` for request-scoped calls.
- Keep handler logic thin; delegate storage/logic to service packages.
- Wire routes in `HTTPRoutes()` on the service type.
- All API routes are prefixed with `/api` by `cmd.RunServices` (e.g., `/batches` becomes `/api/batches`).

## Services and lifecycle
- Service entrypoints use `cmd.Main(run)` to standardize exit codes.
- `cmd.RunServices` aggregates routes, prefixes them with `/api`, and runs a single HTTP server.
- The monolith also serves embedded frontend assets at non-API routes via `service/www`.
- `Start(ctx)` should connect to dependencies and run migrations.
- If a service allocates resources, clean up on context cancellation.
- Background goroutines should respect `ctx.Done()`.

## Storage and database
- Storage clients are in `service/<name>/storage` and use `pgxpool`.
- `Start(ctx)` validates connectivity via `Ping`.
- Migrations are run via `internal/database.Migrate`.
- Migrations live at `service/<name>/storage/migrations`.
- SQL is typically embedded as raw string literals in storage methods.
- Use parameterized queries to avoid SQL injection.

## Migrations
- File naming uses timestamp prefixes and `*.up.sql`/`*.down.sql`.
- Keep schema changes in the service-specific migrations directory.
- Add matching up/down files for every change.

## Entities and DTOs
- DB entities embed `entity.Identifiers` and `entity.Timestamps`.
- DTOs live under `service/<name>/handler/dto`.
- Use JSON tags on DTO fields; keep request validation in `Validate()`.
- Keep response types separate from storage models when fields differ.

## Auth and JWT
- JWT helpers live in `internal/jwt`.
- Tokens are HMAC (HS256) and use `github.com/golang-jwt/jwt/v4`.
- Identity service requires `BREWPIPES_SECRET_KEY` to sign tokens.
- Production, inventory, and procurement services require `Authorization: Bearer <access token>` on all routes and will not start without `BREWPIPES_SECRET_KEY`.
- Access/refresh TTLs are defined in `service/identity/storage/user.go`.

## Logging
- Use `slog.Info`/`slog.Error` with key/value pairs.
- Include error details with `"error", err` when possible.
- `service.InternalError` emits a correlation id under the `cid` key.

## Testing style
- Tests use external package naming (`package handler_test`).
- Prefer table-driven tests for multiple cases.
- Use `httptest.NewRecorder` and `httptest.NewRequest` for handlers.
- Keep test helpers near the test file.

## Agent workflow tips
- When adding endpoints, update handler + service routes + tests.
- When adding DB fields, update migrations, storage scans, and DTOs.
- Keep new packages under `service/<name>` or `internal`.
- Avoid introducing new dependencies unless necessary.
- Follow existing directory structure to reduce merge conflicts.

## Cursor/Copilot rules
- No `.cursor/rules`, `.cursorrules`, or `.github/copilot-instructions.md` found.

## Notes
- The repo currently has a single test file.
- Dockerfile in `cmd/monolith/Dockerfile` is a multi-stage build (Node → Go → Alpine runtime).
- The monolith embeds frontend assets from `service/www/dist/` using `//go:embed`.
- Makefile provides the canonical run and Postgres helpers.
- Go version should match `go.mod` and `cmd/monolith/Dockerfile`.
- When changing migrations, ensure both services can start cleanly.
- When adding a new service, add `cmd/<name>/main.go` and wire routes.
- Prefer `service.JSON` to ensure JSON content-type is set.
- Prefer running `go test ./...` before shipping changes.
