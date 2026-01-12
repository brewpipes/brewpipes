package production

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brewpipes/brewpipesproto/internal/service"
	"github.com/brewpipes/brewpipesproto/internal/service/production/handler"
	"github.com/brewpipes/brewpipesproto/internal/service/production/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	storage *storage.Client
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("creating DB connection pool: %w", err)
	}

	return &Service{
		storage: &storage.Client{DB: pool},
	}, nil
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	return []service.HTTPRoute{
		{
			Method:  http.MethodGet,
			Path:    "/volumes",
			Handler: handler.HandleGetVolumes(s.storage),
		},
	}
}

func (s *Service) Start(ctx context.Context) error {
	if err := s.storage.Ping(ctx); err != nil {
		return fmt.Errorf("pinging DB: %w", err)
	}

	return nil
}
