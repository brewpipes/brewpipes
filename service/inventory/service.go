package inventory

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type Config struct {
	PostgresDSN string
}

type Service struct {
	storage *storage.Client
}

// New creates and initializes a new inventory service instance.
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
		{Method: http.MethodGet, Path: "/ingredients", Handler: handler.HandleIngredients(s.storage)},
		{Method: http.MethodPost, Path: "/ingredients", Handler: handler.HandleIngredients(s.storage)},
		{Method: http.MethodGet, Path: "/ingredients/{id}", Handler: handler.HandleIngredientByID(s.storage)},
		{Method: http.MethodGet, Path: "/stock-locations", Handler: handler.HandleStockLocations(s.storage)},
		{Method: http.MethodPost, Path: "/stock-locations", Handler: handler.HandleStockLocations(s.storage)},
		{Method: http.MethodGet, Path: "/stock-locations/{id}", Handler: handler.HandleStockLocationByID(s.storage)},
		{Method: http.MethodGet, Path: "/inventory-receipts", Handler: handler.HandleInventoryReceipts(s.storage)},
		{Method: http.MethodPost, Path: "/inventory-receipts", Handler: handler.HandleInventoryReceipts(s.storage)},
		{Method: http.MethodGet, Path: "/inventory-receipts/{id}", Handler: handler.HandleInventoryReceiptByID(s.storage)},
		{Method: http.MethodGet, Path: "/ingredient-lots", Handler: handler.HandleIngredientLots(s.storage)},
		{Method: http.MethodPost, Path: "/ingredient-lots", Handler: handler.HandleIngredientLots(s.storage)},
		{Method: http.MethodGet, Path: "/ingredient-lots/{id}", Handler: handler.HandleIngredientLotByID(s.storage)},
		{Method: http.MethodGet, Path: "/inventory-movements", Handler: handler.HandleInventoryMovements(s.storage)},
		{Method: http.MethodPost, Path: "/inventory-movements", Handler: handler.HandleInventoryMovements(s.storage)},
		{Method: http.MethodGet, Path: "/inventory-movements/{id}", Handler: handler.HandleInventoryMovementByID(s.storage)},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("inventory service starting")
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
