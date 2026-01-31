package identity

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/identity/handler"
	"github.com/brewpipes/brewpipes/service/identity/storage"
)

type Config struct {
	PostgresDSN string
	SecretKey   string
}

type startCloser interface {
	Start(ctx context.Context) error
	Close() error
}

type Service struct {
	storage startCloser
	routes  []service.HTTPRoute
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	stg, err := storage.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("creating storage client: %w", err)
	}

	return &Service{
		storage: stg,
		// inject dependencies for HTTP handlers
		routes: []service.HTTPRoute{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: handler.HandleLogin(stg, cfg.SecretKey),
			},
		},
	}, nil
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	return s.routes
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("identity service starting")
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
