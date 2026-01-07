package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/auth"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent auth service application.
	cfg := &auth.Config{
		PostgresDSN:      os.Getenv("AUTH_POSTGRES_DSN"),
		PostgresPassword: os.Getenv("AUTH_POSTGRES_PASSWORD"),
	}

	// Initialize service.
	svc, err := auth.NewService(cfg)
	if err != nil {
		return err
	}

	return cmd.RunServices(svc)
}
