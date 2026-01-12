.PHONY: info clean monolith auth-service production-service run-monolith

info:
	@echo "  make monolith            Build the monolith application"
	@echo "  make auth-service        Build the authentication service"
	@echo "  make production-service  Build the production service"
	@echo "  make clean               Clean all built binaries"

clean:
	rm -rf bin/*

monolith:
	go build -o bin/monolith ./cmd/monolith

auth-service:
	go build -o bin/auth-service ./cmd/auth

production-service:
	go build -o bin/production-service ./cmd/production


run-monolith:
	POSTGRES_DSN=postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable \
	go run ./cmd/monolith
