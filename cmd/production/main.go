package main

import (
	"context"
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent production service application.

	svc, err := production.NewService(ctx, &production.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(svc)
}
