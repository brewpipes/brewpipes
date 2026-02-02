package main

import (
	"context"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/service/production"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent production service application.

	svc, err := production.New(ctx, production.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(ctx, nil, svc)
}
