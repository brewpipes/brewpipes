package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/service/identity"
	"github.com/brewpipes/brewpipes/service/inventory"
	"github.com/brewpipes/brewpipes/service/procurement"
	"github.com/brewpipes/brewpipes/service/production"
	"github.com/brewpipes/brewpipes/service/www"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Initialize services.
	identitySvc, err := identity.NewService(ctx, &identity.Config{
		PostgresDSN: postgresDSN(),
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return fmt.Errorf("initializing identity service: %w", err)
	}

	productionSvc, err := production.New(ctx, production.Config{
		PostgresDSN: postgresDSN(),
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return fmt.Errorf("initializing production service: %w", err)
	}

	inventorySvc, err := inventory.New(ctx, inventory.Config{
		PostgresDSN: postgresDSN(),
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return fmt.Errorf("initializing inventory service: %w", err)
	}

	procurementSvc, err := procurement.New(ctx, procurement.Config{
		PostgresDSN: postgresDSN(),
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return fmt.Errorf("initializing procurement service: %w", err)
	}

	return cmd.RunServices(ctx, www.Handler(), identitySvc, productionSvc, inventorySvc, procurementSvc)
}

func postgresDSN() string {
	params := ""
	if os.Getenv("ENVIRONMENT") == "local" {
		params = "?sslmode=disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		params)
}
