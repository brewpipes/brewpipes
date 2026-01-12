package production

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brewpipes/brewpipesproto/internal/database"
	"github.com/brewpipes/brewpipesproto/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	pool, err := pgxpool.New(ctx, "postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=enable")
	if err != nil {
		return nil, fmt.Errorf("creating DB connection pool: %w", err)
	}

	return &Service{
		db: pool,
	}, nil
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
	if err := s.db.Ping(ctx); err != nil {
		return fmt.Errorf("pinging DB: %w", err)
	}

	if err := database.Migrate(
		"file://./db/migrations",
		"pgx5://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=enable",
	); err != nil {
		return fmt.Errorf("migrating DB: %w", err)
	}

	return nil
}
