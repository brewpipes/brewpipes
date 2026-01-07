package production

import (
	"context"
	"log/slog"
	"net/http"
)

type Service struct {
	done     chan error
	srv      *http.Server
	httpDone chan struct{}
}

func NewService(cfg *Config) (*Service, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK\n")); err != nil {
			slog.Error("error writing health response", "error", err)
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	slog.Info("production service initialized")
	return &Service{
		done:     make(chan error),
		srv:      srv,
		httpDone: make(chan struct{}),
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("production service starting")
	go s.run(ctx)
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	slog.Info("production service stopping")

	slog.Info("shutting down production service HTTP server")
	if err := s.srv.Shutdown(context.Background()); err != nil { // err is always non-nil
		if err == context.DeadlineExceeded {
			slog.Error("production service HTTP server shutdown timed out")
		} else {
			slog.Error("error shutting down production service HTTP server", "error", err)
		}

		close(s.httpDone)
	}

	close(s.done)
	slog.Info("production service stopped")
	return nil
}

func (s *Service) Done() <-chan error {
	return s.done
}

func (s *Service) run(ctx context.Context) {
	// Shutdown server on context cancellation.
	go func() {
		<-ctx.Done()
		s.Stop(context.Background())
	}()

	slog.Info("production service HTTP server starting on :8080")

	if err := s.srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			slog.Error("production service HTTP server error", "error", err)
		}
	}

	<-s.httpDone
	slog.Info("production service HTTP server stopped")
}
