package main

import (
	"log/slog"
	"os"

	"github.com/brewpipes/brewpipesproto/internal/service/auth"
	"github.com/brewpipes/brewpipesproto/internal/service/production"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Entry point for the monolith application.

	// Initialize services.
	authSvc, err := auth.NewService()
	if err != nil {
		return err
	}

	productionSvc, err := production.NewService()
	if err != nil {
		return err
	}

	// Start services.
	if err := authSvc.Start(); err != nil {
		return err
	}
	defer authSvc.Stop()

	if err := productionSvc.Start(); err != nil {
		return err
	}
	defer productionSvc.Stop()

	slog.Info("monolith application running")
	slog.Info("monolith application stopping")
	return nil
}
