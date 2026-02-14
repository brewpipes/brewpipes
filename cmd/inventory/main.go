package main

import (
	"context"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/service/inventory"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent inventory service application.

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_DSN")
	}

	svc := inventory.New(inventory.Config{
		PostgresDSN: dsn,
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})

	return cmd.RunServices(ctx, nil, svc)
}
