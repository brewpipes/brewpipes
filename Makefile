.PHONY: info clean monolith auth-service production-service run-monolith psql

info:
	@echo "  make monolith            Build the monolith application"
	@echo "  make auth-service        Build the authentication service"
	@echo "  make production-service  Build the production service"
	@echo "  make clean               Clean all built binaries"

start-postgres:
	docker compose up -d postgres

nuke:
	docker compose down -v

run-server:
	POSTGRES_DSN=postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable \
	BREWPIPES_SECRET_KEY=dummy \
	go run ./cmd/monolith

run-web:
	nvm use && \
	cd service/www && \
	pnpm dev --force

# connect to the postgres container using psql
psql:
	psql "postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable"
