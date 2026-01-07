package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

type RunnableService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Done() <-chan error
}

func RunServices(services ...RunnableService) error {
	// Establish root cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	startErrsCh := make(chan error, len(services))
	for _, svc := range services {
		wg.Go(func() {
			if err := svc.Start(ctx); err != nil {
				startErrsCh <- err
			}
		})
	}

	// Wait for all services to start or any to return an error.
	wg.Wait()
	close(startErrsCh)
	for err := range startErrsCh {
		if err != nil {
			slog.Error("error starting service", "error", err)
			return err
		}
	}

	// Wait for application to receive interrupt signal.
	<-interrupted()
	slog.Info("application received interrupt signal, stopping services")

	for _, svc := range services {
		if err := svc.Stop(context.Background()); err != nil {
			slog.Error("error stopping service", "error", err)
		} else {
			slog.Info("service stopped gracefully")
		}
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
