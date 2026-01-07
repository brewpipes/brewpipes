package main

import (
	"log/slog"
	"os"

	"github.com/brewpipes/brewpipesproto/internal/service/auth"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Entry point for the independent auth service application.

	// Initialize service.
	authSvc, err := auth.NewService()
	if err != nil {
		return err
	}

	// Start service.
	if err := authSvc.Start(); err != nil {
		return err
	}
	defer authSvc.Stop()

	return nil
}
