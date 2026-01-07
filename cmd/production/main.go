package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent production service application.
	cfg := &production.Config{
		PostgresDSN:      os.Getenv("PRODUCTION_POSTGRES_DSN"),
		PostgresPassword: os.Getenv("PRODUCTION_POSTGRES_PASSWORD"),
	}

	// Initialize service.
	svc, err := production.NewService(cfg)
	if err != nil {
		return err
	}

	return cmd.RunServices(svc)
}
