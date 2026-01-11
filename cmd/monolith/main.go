package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/auth"
	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent auth service application.
	authCfg := &auth.Config{
		PostgresDSN:      os.Getenv("POSTGRES_DSN"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
	}

	productionCfg := &production.Config{
		PostgresDSN:      os.Getenv("POSTGRES_DSN"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
	}

	// Initialize services.
	authSvc, err := auth.NewService(authCfg)
	if err != nil {
		return err
	}

	productionSvc, err := production.NewService(productionCfg)
	if err != nil {
		return err
	}

	return cmd.RunServices(authSvc, productionSvc)
}
