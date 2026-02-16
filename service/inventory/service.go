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
func New(cfg Config) *Service {
	return &Service{
		storage:   storage.New(cfg.PostgresDSN),
		secretKey: cfg.SecretKey,
	}
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	auth := service.RequireAccessToken(s.secretKey)
	return []service.HTTPRoute{
		{Method: http.MethodGet, Path: "/ingredients", Handler: auth(handler.HandleIngredients(s.storage))},
		{Method: http.MethodPost, Path: "/ingredients", Handler: auth(handler.HandleIngredients(s.storage))},
		{Method: http.MethodGet, Path: "/ingredients/{uuid}", Handler: auth(handler.HandleIngredientByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/stock-locations", Handler: auth(handler.HandleStockLocations(s.storage))},
		{Method: http.MethodPost, Path: "/stock-locations", Handler: auth(handler.HandleStockLocations(s.storage))},
		{Method: http.MethodGet, Path: "/stock-locations/{uuid}", Handler: auth(handler.HandleStockLocationByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-receipts", Handler: auth(handler.HandleInventoryReceipts(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-receipts", Handler: auth(handler.HandleInventoryReceipts(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-receipts/{uuid}", Handler: auth(handler.HandleInventoryReceiptByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lots", Handler: auth(handler.HandleIngredientLots(s.storage))},
		{Method: http.MethodPost, Path: "/ingredient-lots", Handler: auth(handler.HandleIngredientLots(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lots/{uuid}", Handler: auth(handler.HandleIngredientLotByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-movements", Handler: auth(handler.HandleInventoryMovements(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-movements", Handler: auth(handler.HandleInventoryMovements(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-movements/{uuid}", Handler: auth(handler.HandleInventoryMovementByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-usage", Handler: auth(handler.HandleInventoryUsage(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-usage", Handler: auth(handler.HandleInventoryUsage(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-usage/batch", Handler: auth(handler.HandleCreateBatchUsage(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-usage/{uuid}", Handler: auth(handler.HandleInventoryUsageByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-adjustments", Handler: auth(handler.HandleInventoryAdjustments(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-adjustments", Handler: auth(handler.HandleInventoryAdjustments(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-adjustments/{uuid}", Handler: auth(handler.HandleInventoryAdjustmentByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-transfers", Handler: auth(handler.HandleInventoryTransfers(s.storage))},
		{Method: http.MethodPost, Path: "/inventory-transfers", Handler: auth(handler.HandleInventoryTransfers(s.storage))},
		{Method: http.MethodGet, Path: "/inventory-transfers/{uuid}", Handler: auth(handler.HandleInventoryTransferByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/beer-lots", Handler: auth(handler.HandleBeerLots(s.storage))},
		{Method: http.MethodPost, Path: "/beer-lots", Handler: auth(handler.HandleBeerLots(s.storage))},
		{Method: http.MethodGet, Path: "/beer-lots/{uuid}", Handler: auth(handler.HandleBeerLotByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/stock-levels", Handler: auth(handler.HandleStockLevels(s.storage))},
		{Method: http.MethodGet, Path: "/beer-lot-stock-levels", Handler: auth(handler.HandleBeerLotStockLevels(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lot-malt-details", Handler: auth(handler.HandleIngredientLotMaltDetails(s.storage))},
		{Method: http.MethodPost, Path: "/ingredient-lot-malt-details", Handler: auth(handler.HandleIngredientLotMaltDetails(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lot-malt-details/{uuid}", Handler: auth(handler.HandleIngredientLotMaltDetailByUUID(s.storage))},
		{Method: http.MethodPut, Path: "/ingredient-lot-malt-details/{uuid}", Handler: auth(handler.HandleIngredientLotMaltDetailByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lot-hop-details", Handler: auth(handler.HandleIngredientLotHopDetails(s.storage))},
		{Method: http.MethodPost, Path: "/ingredient-lot-hop-details", Handler: auth(handler.HandleIngredientLotHopDetails(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lot-hop-details/{uuid}", Handler: auth(handler.HandleIngredientLotHopDetailByUUID(s.storage))},
		{Method: http.MethodPut, Path: "/ingredient-lot-hop-details/{uuid}", Handler: auth(handler.HandleIngredientLotHopDetailByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lot-yeast-details", Handler: auth(handler.HandleIngredientLotYeastDetails(s.storage))},
		{Method: http.MethodPost, Path: "/ingredient-lot-yeast-details", Handler: auth(handler.HandleIngredientLotYeastDetails(s.storage))},
		{Method: http.MethodGet, Path: "/ingredient-lot-yeast-details/{uuid}", Handler: auth(handler.HandleIngredientLotYeastDetailByUUID(s.storage))},
		{Method: http.MethodPut, Path: "/ingredient-lot-yeast-details/{uuid}", Handler: auth(handler.HandleIngredientLotYeastDetailByUUID(s.storage))},
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
