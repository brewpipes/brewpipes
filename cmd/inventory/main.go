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

	svc, err := inventory.New(ctx, inventory.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(ctx, svc)
}
