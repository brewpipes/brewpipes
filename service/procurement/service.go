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
}

type Service struct {
	storage *storage.Client
}

// New creates and initializes a new procurement service instance.
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
		{Method: http.MethodGet, Path: "/suppliers", Handler: handler.HandleSuppliers(s.storage)},
		{Method: http.MethodPost, Path: "/suppliers", Handler: handler.HandleSuppliers(s.storage)},
		{Method: http.MethodGet, Path: "/suppliers/{id}", Handler: handler.HandleSupplierByID(s.storage)},
		{Method: http.MethodGet, Path: "/purchase-orders", Handler: handler.HandlePurchaseOrders(s.storage)},
		{Method: http.MethodPost, Path: "/purchase-orders", Handler: handler.HandlePurchaseOrders(s.storage)},
		{Method: http.MethodGet, Path: "/purchase-orders/{id}", Handler: handler.HandlePurchaseOrderByID(s.storage)},
		{Method: http.MethodGet, Path: "/purchase-order-lines", Handler: handler.HandlePurchaseOrderLines(s.storage)},
		{Method: http.MethodPost, Path: "/purchase-order-lines", Handler: handler.HandlePurchaseOrderLines(s.storage)},
		{Method: http.MethodGet, Path: "/purchase-order-lines/{id}", Handler: handler.HandlePurchaseOrderLineByID(s.storage)},
		{Method: http.MethodGet, Path: "/purchase-order-fees", Handler: handler.HandlePurchaseOrderFees(s.storage)},
		{Method: http.MethodPost, Path: "/purchase-order-fees", Handler: handler.HandlePurchaseOrderFees(s.storage)},
		{Method: http.MethodGet, Path: "/purchase-order-fees/{id}", Handler: handler.HandlePurchaseOrderFeeByID(s.storage)},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("procurement service starting")
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
