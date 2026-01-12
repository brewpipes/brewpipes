package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/internal/service/identity"
	"github.com/brewpipes/brewpipes/internal/service/production"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent identity service application.

	// Initialize services.
	identitySvc, err := identity.NewService(&identity.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return fmt.Errorf("initializing identity service: %w", err)
	}

	productionSvc, err := production.NewService(ctx, &production.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return fmt.Errorf("initializing production service: %w", err)
	}

	return cmd.RunServices(identitySvc, productionSvc)
}
