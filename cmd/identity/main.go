package main

import (
	"context"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/internal/service/identity"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent identity service application.

	svc, err := identity.NewService(&identity.Config{
		PostgresDSN: os.Getenv("identity_POSTGRES_DSN"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(ctx, svc)
}
