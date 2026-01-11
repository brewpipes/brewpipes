package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/auth"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent production service application.

	svc, err := auth.NewService(&auth.Config{
		PostgresDSN:      os.Getenv("PRODUCTION_POSTGRES_DSN"),
		PostgresPassword: os.Getenv("PRODUCTION_POSTGRES_PASSWORD"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(svc)
}
