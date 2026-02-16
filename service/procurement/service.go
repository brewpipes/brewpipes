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
func New(cfg Config) *Service {
	return &Service{
		storage:   storage.New(cfg.PostgresDSN),
		secretKey: cfg.SecretKey,
	}
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	auth := service.RequireAccessToken(s.secretKey)
	return []service.HTTPRoute{
		{Method: http.MethodGet, Path: "/suppliers", Handler: auth(handler.HandleSuppliers(s.storage))},
		{Method: http.MethodPost, Path: "/suppliers", Handler: auth(handler.HandleSuppliers(s.storage))},
		{Method: http.MethodGet, Path: "/suppliers/{uuid}", Handler: auth(handler.HandleSupplierByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/suppliers/{uuid}", Handler: auth(handler.HandleSupplierByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-orders", Handler: auth(handler.HandlePurchaseOrders(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-orders", Handler: auth(handler.HandlePurchaseOrders(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-orders/{uuid}", Handler: auth(handler.HandlePurchaseOrderByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/purchase-orders/{uuid}", Handler: auth(handler.HandlePurchaseOrderByUUID(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-order-lines/batch-lookup", Handler: auth(handler.HandleBatchLookupPurchaseOrderLines(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-lines", Handler: auth(handler.HandlePurchaseOrderLines(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-order-lines", Handler: auth(handler.HandlePurchaseOrderLines(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-lines/{uuid}", Handler: auth(handler.HandlePurchaseOrderLineByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/purchase-order-lines/{uuid}", Handler: auth(handler.HandlePurchaseOrderLineByUUID(s.storage))},
		{Method: http.MethodDelete, Path: "/purchase-order-lines/{uuid}", Handler: auth(handler.HandlePurchaseOrderLineByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-fees", Handler: auth(handler.HandlePurchaseOrderFees(s.storage))},
		{Method: http.MethodPost, Path: "/purchase-order-fees", Handler: auth(handler.HandlePurchaseOrderFees(s.storage))},
		{Method: http.MethodGet, Path: "/purchase-order-fees/{uuid}", Handler: auth(handler.HandlePurchaseOrderFeeByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/purchase-order-fees/{uuid}", Handler: auth(handler.HandlePurchaseOrderFeeByUUID(s.storage))},
		{Method: http.MethodDelete, Path: "/purchase-order-fees/{uuid}", Handler: auth(handler.HandlePurchaseOrderFeeByUUID(s.storage))},
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
