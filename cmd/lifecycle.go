package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/brewpipes/brewpipes/internal/service"
)

type Service interface {
	HTTPRoutes() []service.HTTPRoute
	Start(ctx context.Context) error
}

func RunServices(services ...Service) error {
	// Establish root cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Aggregate HTTP routes from all services.
	mux := http.NewServeMux()
	for _, svc := range services {
		for _, route := range svc.HTTPRoutes() {
			mux.Handle(route.Path, route.Handler)
		}
	}

	// Create HTTP server to serve aggregated routes.
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start all services.
	for _, svc := range services {
		if err := svc.Start(ctx); err != nil {
			return fmt.Errorf("starting service: %w", err)
		}
	}

	// Start HTTP server.
	go func() {
		slog.Info("starting aggregated HTTP server on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server exited with error", "error", err)
		}
	}()

	// Wait for application to receive interrupt signal.
	<-interrupted()
	slog.Info("application received interrupt signal, stopping services")
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("error while shutting down HTTP server", "error", err)
	}

	cancel()

	slog.Info("application terminated")
	return nil
}

// interrupted returns a channel that is closed when an interrupt signal is received.
func interrupted() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	return sig
}
