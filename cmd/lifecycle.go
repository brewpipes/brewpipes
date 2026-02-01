package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/brewpipes/brewpipes/service"
)

// Service is anything with HTTP routes that can be started.
type Service interface {
	HTTPRoutes() []service.HTTPRoute
	Start(ctx context.Context) error
}

// RunServices starts the provided services and manages their lifecycle.
// If staticHandler is non-nil, it is registered as a catch-all handler for
// non-API routes to serve static files (e.g., a frontend SPA).
func RunServices(ctx context.Context, staticHandler http.Handler, services ...Service) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Aggregate HTTP routes from all services.
	mux := http.NewServeMux()
	for _, svc := range services {
		for _, route := range svc.HTTPRoutes() {
			pattern := "/api" + route.Path
			if route.Method != "" {
				pattern = route.Method + " /api" + route.Path
			}
			mux.Handle(pattern, route.Handler)
		}
	}

	// Register static file handler as catch-all for non-API routes.
	if staticHandler != nil {
		mux.Handle("/", staticHandler)
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
	slog.Info("application received interrupt signal")

	// first, gracefully shut down HTTP server so that in-flight requests can complete
	slog.Info("stopping HTTP server")
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("error while shutting down HTTP server", "error", err)
	}

	// cancel service contexts
	slog.Info("stopping services")
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
