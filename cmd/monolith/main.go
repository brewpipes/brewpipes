package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/identity"
	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent identity service application.
	identityCfg := &identity.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	}

	productionCfg := &production.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	}

	// Initialize services.
	identitySvc, err := identity.NewService(identityCfg)
	if err != nil {
		return err
	}

	productionSvc, err := production.NewService(productionCfg)
	if err != nil {
		return err
	}

	return cmd.RunServices(identitySvc, productionSvc)
}
