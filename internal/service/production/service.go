package production

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipesproto/internal/service"
)

type Service struct {
}

func NewService(cfg *Config) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Name() string {
	return "auth service"
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	return []service.HTTPRoute{
		{
			Method:  http.MethodGet,
			Path:    "/volumes",
			Handler: http.HandlerFunc(s.handleGetVolumes),
		},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("auth service starting")
	go s.run(ctx)
	return nil
}

func (s *Service) run(ctx context.Context) {
	// stop service on context cancellation
	go func() {
		<-ctx.Done()
		s.stop()
	}()
}

func (s *Service) stop() {
	slog.Info("auth service stopping")
}
