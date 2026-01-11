package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/brewpipes/brewpipesproto/internal/service"
)

type Service interface {
	Name() string
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

	// Start all services.
	for _, svc := range services {
		if err := svc.Start(ctx); err != nil {
			return fmt.Errorf("starting %s: %w", svc.Name(), err)
		}
	}

	// Create HTTP server to serve aggregated routes.
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start HTTP server.
	go func() {
		slog.Info("starting aggregated HTTP server on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for application to receive interrupt signal.
	<-interrupted()
	slog.Info("application received interrupt signal, stopping services")
	cancel()
	slog.Info("application terminated")
	return nil
}

// interrupted returns a channel that is closed when an interrupt signal is received.
func interrupted() <-chan struct{} {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	interrupted := make(chan struct{})
	go func() {
		<-sig
		close(interrupted)
	}()
	return interrupted
}
