package procurement

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/procurement/handler"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type Config struct {
	PostgresDSN string
	SecretKey   string
}

type Service struct {
	storage   *storage.Client
	secretKey string
}

// New creates and initializes a new procurement service instance.
func New(ctx context.Context, cfg Config) (*Service, error) {
	stg, err := storage.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("creating storage client: %w", err)
	}

	return &Service{
		storage:   stg,
		secretKey: cfg.SecretKey,
	}, nil
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	auth := service.RequireAccessToken(s.secretKey)
	return []service.HTTPRoute{
		{Method: http.MethodGet, Path: "/suppliers", Handler: auth(handler.HandleSuppliers(s.storage))},
		{Method: http.MethodPost, Path: "/suppliers", Handler: auth(handler.HandleSuppliers(s.storage))},
		{Method: http.MethodGet, Path: "/suppliers/{id}", Handler: auth(handler.HandleSupplierByID(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-orders", Handler: auth(handler.HandlePurchaseOrders(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-orders", Handler: auth(handler.HandlePurchaseOrders(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-orders/{id}", Handler: auth(handler.HandlePurchaseOrderByID(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-lines", Handler: auth(handler.HandlePurchaseOrderLines(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-order-lines", Handler: auth(handler.HandlePurchaseOrderLines(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-lines/{id}", Handler: auth(handler.HandlePurchaseOrderLineByID(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-fees", Handler: auth(handler.HandlePurchaseOrderFees(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-order-fees", Handler: auth(handler.HandlePurchaseOrderFees(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-fees/{id}", Handler: auth(handler.HandlePurchaseOrderFeeByID(s.storage))},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("procurement service starting")
	if s.secretKey == "" {
		return fmt.Errorf("missing BREWPIPES_SECRET_KEY for access token verification")
	}
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
