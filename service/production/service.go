package production

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type Config struct {
	PostgresDSN string
}

type Service struct {
	storage *storage.Client
}

// New creates and initializes a new production service instance.
func New(ctx context.Context, cfg Config) (*Service, error) {
	stg, err := storage.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("creating storage client: %w", err)
	}

	return &Service{
		storage: stg,
	}, nil
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	return []service.HTTPRoute{
		{Method: http.MethodGet, Path: "/volumes", Handler: handler.HandleGetVolumes(s.storage)},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("production service starting")
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
