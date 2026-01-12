package identity

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/internal/service"
	"github.com/brewpipes/brewpipes/internal/service/identity/handler"
	"github.com/brewpipes/brewpipes/internal/service/identity/storage"
	"github.com/gofrs/uuid/v5"
)

type Database interface {
	GetUser(ctx context.Context, id uuid.UUID) (storage.User, error)
	GetUserByUsername(ctx context.Context, username string) (storage.User, error)
	ListUsers(ctx context.Context) ([]storage.User, error)
}

type Service struct {
	db        Database
	secretKey string
}

func NewService(cfg *Config) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	return []service.HTTPRoute{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: handler.HandleLogin(s.db, s.secretKey),
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
