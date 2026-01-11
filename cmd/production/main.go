package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent production service application.

	svc, err := production.NewService(&production.Config{
		PostgresDSN: os.Getenv("PRODUCTION_POSTGRES_DSN"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(svc)
}
