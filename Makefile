.PHONY: info clean monolith auth-service production-service run-monolith psql

postgres_dsn = postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable

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
	POSTGRES_DSN="$(postgres_dsn)" \
	BREWPIPES_SECRET_KEY="dummy" \
	go run ./cmd/monolith

run-web:
	cd service/www && \
	pnpm dev --force

build-docker:
	docker build -f cmd/monolith/Dockerfile -t brewpipes .

run-docker:
	docker run -p 8080:8080 \
	-e POSTGRES_DSN="$(postgres_dsn)" \
	-e BREWPIPES_SECRET_KEY="dummy" \
	brewpipes

# connect to the postgres container using psql
psql:
	psql "$(postgres_dsn)"
