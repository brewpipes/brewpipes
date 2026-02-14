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
	// Support both DATABASE_URL and POSTGRES_DSN for flexibility.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_DSN")
	}

	// Initialize services.
	identitySvc, err := identity.NewService(ctx, &identity.Config{
		PostgresDSN: dsn,
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return fmt.Errorf("initializing identity service: %w", err)
	}

	productionSvc := production.New(production.Config{
		PostgresDSN: dsn,
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})

	inventorySvc := inventory.New(inventory.Config{
		PostgresDSN: dsn,
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})

	procurementSvc := procurement.New(procurement.Config{
		PostgresDSN: dsn,
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})

	return cmd.RunServices(ctx, www.Handler(), identitySvc, productionSvc, inventorySvc, procurementSvc)
}
