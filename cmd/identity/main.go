package main

import (
	"os"

	"github.com/brewpipes/brewpipesproto/cmd"
	"github.com/brewpipes/brewpipesproto/internal/service/identity"
)

func main() {
	cmd.Main(run)
}

func run() error {
	// Entry point for the independent identity service application.

	svc, err := identity.NewService(&identity.Config{
		PostgresDSN: os.Getenv("identity_POSTGRES_DSN"),
	})
	if err != nil {
		return err
	}

	return cmd.RunServices(svc)
}
