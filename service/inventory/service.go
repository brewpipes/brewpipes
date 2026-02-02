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
	SecretKey   string
}

type Service struct {
	storage   *storage.Client
	secretKey string
}

// New creates and initializes a new inventory service instance.
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
		{Method: http.MethodGet, Path: "/ingredients", Handler: auth(handler.HandleIngredients(s.storage))},
		{Method: http.MethodPost, Path: "/ingredients", Handler: auth(handler.HandleIngredients(s.storage))},
		{Method: http.MethodGet, Path: "/ingredients/{id}", Handler: auth(handler.HandleIngredientByID(s.storage))},
		{Method: http.MethodGet, Path: "/stock-locations", Handler: auth(handler.HandleStockLocations(s.storage))},
		{Method: http.MethodPost, Path: "/stock-locations", Handler: auth(handler.HandleStockLocations(s.storage))},
		{Method: http.MethodGet, Path: "/stock-locations/{id}", Handler: auth(handler.HandleStockLocationByID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-receipts", Handler: auth(handler.HandleInventoryReceipts(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-receipts", Handler: auth(handler.HandleInventoryReceipts(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-receipts/{id}", Handler: auth(handler.HandleInventoryReceiptByID(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lots", Handler: auth(handler.HandleIngredientLots(s.storage))},
		{Method: http.MethodPost, Path: "/ingredient-lots", Handler: auth(handler.HandleIngredientLots(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lots/{id}", Handler: auth(handler.HandleIngredientLotByID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-movements", Handler: auth(handler.HandleInventoryMovements(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-movements", Handler: auth(handler.HandleInventoryMovements(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-movements/{id}", Handler: auth(handler.HandleInventoryMovementByID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-usage", Handler: auth(handler.HandleInventoryUsage(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-usage", Handler: auth(handler.HandleInventoryUsage(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-usage/{id}", Handler: auth(handler.HandleInventoryUsageByID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-adjustments", Handler: auth(handler.HandleInventoryAdjustments(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-adjustments", Handler: auth(handler.HandleInventoryAdjustments(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-adjustments/{id}", Handler: auth(handler.HandleInventoryAdjustmentByID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-transfers", Handler: auth(handler.HandleInventoryTransfers(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-transfers", Handler: auth(handler.HandleInventoryTransfers(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-transfers/{id}", Handler: auth(handler.HandleInventoryTransferByID(s.storage))},
		{Method: http.MethodGet, Path: "/beer-lots", Handler: auth(handler.HandleBeerLots(s.storage))},
		{Method: http.MethodPost, Path: "/beer-lots", Handler: auth(handler.HandleBeerLots(s.storage))},
		{Method: http.MethodGet, Path: "/beer-lots/{id}", Handler: auth(handler.HandleBeerLotByID(s.storage))},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("inventory service starting")
	if s.secretKey == "" {
		return fmt.Errorf("missing BREWPIPES_SECRET_KEY for access token verification")
	}
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
