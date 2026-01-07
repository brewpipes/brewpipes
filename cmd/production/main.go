package main

import (
	"log/slog"
	"os"

	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Entry point for the independent production service application.

	// Initialize service.
	productionSvc, err := production.NewService()
	if err != nil {
		return err
	}

	// Start service.
	if err := productionSvc.Start(); err != nil {
		return err
	}
	defer productionSvc.Stop()

	return nil
}
