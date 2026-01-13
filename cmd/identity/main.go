package main

import (
	"context"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/service/identity"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent identity service application.

	svc, err := identity.NewService(ctx, &identity.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(ctx, svc)
}
