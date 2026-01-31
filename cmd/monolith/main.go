package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/service/inventory"
	"github.com/brewpipes/brewpipes/service/procurement"
	"github.com/brewpipes/brewpipes/service/production"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent identity service application.

	// Initialize services.
	// identitySvc, err := identity.NewService(ctx, &identity.Config{
	// 	PostgresDSN: os.Getenv("POSTGRES_DSN"),
	// 	SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	// })
	// if err != nil {
	// 	return fmt.Errorf("initializing identity service: %w", err)
	// }

	productionSvc, err := production.New(ctx, production.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return fmt.Errorf("initializing production service: %w", err)
	}

	inventorySvc, err := inventory.New(ctx, inventory.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return fmt.Errorf("initializing inventory service: %w", err)
	}

	procurementSvc, err := procurement.New(ctx, procurement.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return fmt.Errorf("initializing procurement service: %w", err)
	}

	return cmd.RunServices(ctx, productionSvc, inventorySvc, procurementSvc)
}
