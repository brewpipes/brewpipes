package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

type RunnableService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Done() <-chan error
}

func RunService(svc RunnableService) error {
	// Establish root cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start service.
	if err := svc.Start(ctx); err != nil {
		return err
	}

	// Wait for service to receive interrupt signal or encounter an error.
	select {
	case <-interrupted():
		slog.Info("application received interrupt signal")
	case err := <-svc.Done():
		if err != nil {
			slog.Error("application encountered a fatal error", "error", err)
		}
	}

	slog.Info("stopping service")

	if err := svc.Stop(context.Background()); err != nil {
		slog.Error("error stopping service", "error", err)
	} else {
		slog.Info("service stopped gracefully")
	}

	slog.Info("application terminated cleanly")
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
