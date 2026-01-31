# AGENTS.md

This file guides agentic coding tools working in this repository.
Keep changes consistent with existing patterns and Go conventions.

## Project Overview
- Language: Go (module `github.com/brewpipes/brewpipes`)
- Go version: 1.25.5 (see `go.mod`)
- Services: identity, production, monolith aggregation
- Runtime: HTTP on port 8080
- Database: Postgres (pgx v5) with migrations

## Build / Run Commands
- Build monolith: `make monolith`
- Build auth service: `make auth-service`
- Build production service: `make production-service`
- Clean binaries: `make clean`
- Run monolith (local): `make run-monolith`
- Run identity service: `go run ./cmd/identity`
- Run production service: `go run ./cmd/production`
- Run monolith directly: `go run ./cmd/monolith`

## Docker Compose (local dev)
- Start services: `docker compose up`
- Stop services: `docker compose down`
- Postgres connection (local): `make psql`
- Compose sets `POSTGRES_DSN` for `brewpipes` service

## Lint / Format / Test
- Format all Go code: `gofmt -w .`
- Vet (baseline lint): `go vet ./...`
- Run all tests: `go test ./...`
- Run tests with race detector: `go test -race ./...`

## Run a Single Test
- Run one test by name (package scoped):
  `go test ./service/production/handler -run TestVolumesHandler -count=1`
- Run one test by regex (all packages):
  `go test ./... -run TestVolumesHandler -count=1`
- Run a single file's tests:
  `go test ./service/production/handler -run TestVolumesHandler -count=1`

Notes:
- `-run` uses regex, so anchor or fully name tests when needed.
- `-count=1` avoids cached results.

## Environment Variables
- `POSTGRES_DSN`: database connection string
  - Example: `postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable`
- `BREWPIPES_SECRET_KEY`: JWT signing secret for identity service

## Repo Layout
- `cmd/`: binaries entrypoints
  - `cmd/monolith`, `cmd/identity`, `cmd/production`
- `service/`: service implementations and handlers
  - `service/identity`, `service/production`
- `internal/`: shared infra (db migrations, jwt helpers)

## Code Style Guidelines
### Imports
- Use standard Go import grouping:
  1) standard library
  2) local module imports (`github.com/brewpipes/brewpipes/...`)
  3) third-party imports
- Avoid unused imports; keep order stable after `gofmt`.

### Formatting
- Always run `gofmt` on modified Go files.
- Keep lines readable; prefer early returns over deep nesting.
- Use `go fmt`/`gofmt` instead of manual alignment.

### Types and Structs
- Use explicit types for public API surfaces.
- Leverage existing shared types in `internal/database/entity`.
- Prefer `uuid.UUID` for identifiers and `time.Time` for timestamps.

### Naming Conventions
- Exported identifiers use PascalCase; unexported use camelCase.
- Use descriptive handler names: `HandleX` returning `http.HandlerFunc`.
- Use `DTO` types under `service/.../handler/dto` for API shapes.

### Error Handling
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- For HTTP handlers:
  - Use `http.Error` for 4xx responses.
  - Use `service.InternalError` for 500s and log with correlation id.
- When a storage lookup misses, return `service.ErrNotFound`.

### Logging
- Prefer `log/slog` for structured logs in services.
- Use `log` only for simple local errors in JWT helpers.
- Avoid logging sensitive values (passwords, tokens, secret keys).

### HTTP Handlers and JSON
- Handlers return `http.HandlerFunc`.
- Use `service.JSON(w, v)` for JSON responses.
- `service.JSON` sets `Content-Type: application/json`.
- Return on error early; do not continue after writing error responses.

### Database Access
- Use pgx pool (`pgxpool.Pool`) through `storage.Client`.
- Always `defer rows.Close()` after successful `Query`.
- Wrap DB errors with context; avoid leaking internal SQL details to clients.
- Migrations are run in `storage.Client.Start` via `internal/database`.

### JWT and Auth
- Use `internal/jwt` to decode/validate access/refresh tokens.
- Token TTLs and issuer are defined in `service/identity/storage/user.go`.
- `BREWPIPES_SECRET_KEY` must be set for identity service.

### Testing
- Keep tests in `_test.go` files and use `testing` package.
- Prefer table-driven tests for multiple cases.
- When mocking storage, use small interfaces in handler packages.

## Cursor / Copilot Rules
- No Cursor rules found in `.cursor/rules/` or `.cursorrules`.
- No Copilot rules found in `.github/copilot-instructions.md`.

## Additional Notes for Agents
- Avoid introducing new frameworks or dependencies without need.
- Stick to existing HTTP routing style (mux in `cmd.RunServices`).
- Do not change DB schemas without updating migrations.
- Keep Docker workflows simple; `docker compose up` is standard for dev.
