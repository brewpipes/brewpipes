package production

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type Config struct {
	PostgresDSN string
	SecretKey   string
}

type Service struct {
	storage         *storage.Client
	secretKey       string
	inventoryClient *handler.InventoryClient
}

// New creates and initializes a new production service instance.
func New(cfg Config) *Service {
	return &Service{
		storage:   storage.New(cfg.PostgresDSN),
		secretKey: cfg.SecretKey,
	}
}

func (s *Service) HTTPRoutes() []service.HTTPRoute {
	auth := service.RequireAccessToken(s.secretKey)
	return []service.HTTPRoute{
		{Method: http.MethodGet, Path: "/styles", Handler: auth(handler.HandleStyles(s.storage))},
		{Method: http.MethodPost, Path: "/styles", Handler: auth(handler.HandleStyles(s.storage))},
		{Method: http.MethodGet, Path: "/styles/{uuid}", Handler: auth(handler.HandleStyleByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/recipes", Handler: auth(handler.HandleRecipes(s.storage))},
		{Method: http.MethodPost, Path: "/recipes", Handler: auth(handler.HandleRecipes(s.storage))},
		{Method: http.MethodGet, Path: "/recipes/{uuid}", Handler: auth(handler.HandleRecipeByUUID(s.storage, s.storage))},
		{Method: http.MethodPut, Path: "/recipes/{uuid}", Handler: auth(handler.HandleRecipeByUUID(s.storage, s.storage))},
		{Method: http.MethodPatch, Path: "/recipes/{uuid}", Handler: auth(handler.HandleRecipeByUUID(s.storage, s.storage))},
		{Method: http.MethodDelete, Path: "/recipes/{uuid}", Handler: auth(handler.HandleRecipeByUUID(s.storage, s.storage))},
		{Method: http.MethodGet, Path: "/recipes/{uuid}/ingredients", Handler: auth(handler.HandleRecipeIngredients(s.storage, s.storage))},
		{Method: http.MethodPost, Path: "/recipes/{uuid}/ingredients", Handler: auth(handler.HandleRecipeIngredients(s.storage, s.storage))},
		{Method: http.MethodGet, Path: "/recipes/{uuid}/ingredients/{ingredient_uuid}", Handler: auth(handler.HandleRecipeIngredient(s.storage, s.storage))},
		{Method: http.MethodPatch, Path: "/recipes/{uuid}/ingredients/{ingredient_uuid}", Handler: auth(handler.HandleRecipeIngredient(s.storage, s.storage))},
		{Method: http.MethodDelete, Path: "/recipes/{uuid}/ingredients/{ingredient_uuid}", Handler: auth(handler.HandleRecipeIngredient(s.storage, s.storage))},
		{Method: http.MethodGet, Path: "/batches", Handler: auth(handler.HandleBatches(s.storage))},
		{Method: http.MethodPost, Path: "/batches", Handler: auth(handler.HandleBatches(s.storage))},
		{Method: http.MethodPost, Path: "/batches/import", Handler: auth(handler.HandleBatchImport(s.storage))},
		{Method: http.MethodGet, Path: "/batches/{uuid}", Handler: auth(handler.HandleBatchByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/batches/{uuid}", Handler: auth(handler.HandleBatchByUUID(s.storage))},
		{Method: http.MethodDelete, Path: "/batches/{uuid}", Handler: auth(handler.HandleBatchByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/batches/{uuid}/summary", Handler: auth(handler.HandleBatchSummaryByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/brew-sessions", Handler: auth(handler.HandleBrewSessions(s.storage))},
		{Method: http.MethodPost, Path: "/brew-sessions", Handler: auth(handler.HandleBrewSessions(s.storage))},
		{Method: http.MethodGet, Path: "/brew-sessions/{uuid}", Handler: auth(handler.HandleBrewSessionByUUID(s.storage))},
		{Method: http.MethodPut, Path: "/brew-sessions/{uuid}", Handler: auth(handler.HandleBrewSessionByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/volumes", Handler: auth(handler.HandleVolumes(s.storage))},
		{Method: http.MethodPost, Path: "/volumes", Handler: auth(handler.HandleVolumes(s.storage))},
		{Method: http.MethodGet, Path: "/volumes/{uuid}", Handler: auth(handler.HandleVolumeByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/volume-relations", Handler: auth(handler.HandleVolumeRelations(s.storage))},
		{Method: http.MethodPost, Path: "/volume-relations", Handler: auth(handler.HandleVolumeRelations(s.storage))},
		{Method: http.MethodGet, Path: "/volume-relations/{uuid}", Handler: auth(handler.HandleVolumeRelationByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/vessels", Handler: auth(handler.HandleVessels(s.storage))},
		{Method: http.MethodPost, Path: "/vessels", Handler: auth(handler.HandleVessels(s.storage))},
		{Method: http.MethodGet, Path: "/vessels/{uuid}", Handler: auth(handler.HandleVesselByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/vessels/{uuid}", Handler: auth(handler.HandleVesselByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/occupancies", Handler: auth(handler.HandleOccupancies(s.storage))},
		{Method: http.MethodPost, Path: "/occupancies", Handler: auth(handler.HandleCreateOccupancy(s.storage))},
		{Method: http.MethodGet, Path: "/occupancies/{uuid}", Handler: auth(handler.HandleOccupancyByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/occupancies/{uuid}/close", Handler: auth(handler.HandleCloseOccupancy(s.storage))},
		{Method: http.MethodPatch, Path: "/occupancies/{uuid}/status", Handler: auth(handler.HandleOccupancyStatus(s.storage))},
		{Method: http.MethodGet, Path: "/occupancies/active", Handler: auth(handler.HandleActiveOccupancy(s.storage))},
		{Method: http.MethodGet, Path: "/transfers", Handler: auth(handler.HandleTransfers(s.storage))},
		{Method: http.MethodPost, Path: "/transfers", Handler: auth(handler.HandleTransfers(s.storage))},
		{Method: http.MethodGet, Path: "/transfers/{uuid}", Handler: auth(handler.HandleTransferByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/batch-volumes", Handler: auth(handler.HandleBatchVolumes(s.storage))},
		{Method: http.MethodPost, Path: "/batch-volumes", Handler: auth(handler.HandleBatchVolumes(s.storage))},
		{Method: http.MethodGet, Path: "/batch-volumes/{uuid}", Handler: auth(handler.HandleBatchVolumeByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/batch-process-phases", Handler: auth(handler.HandleBatchProcessPhases(s.storage))},
		{Method: http.MethodPost, Path: "/batch-process-phases", Handler: auth(handler.HandleBatchProcessPhases(s.storage))},
		{Method: http.MethodGet, Path: "/batch-process-phases/{uuid}", Handler: auth(handler.HandleBatchProcessPhaseByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/batch-relations", Handler: auth(handler.HandleBatchRelations(s.storage))},
		{Method: http.MethodPost, Path: "/batch-relations", Handler: auth(handler.HandleBatchRelations(s.storage))},
		{Method: http.MethodGet, Path: "/batch-relations/{uuid}", Handler: auth(handler.HandleBatchRelationByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/additions", Handler: auth(handler.HandleAdditions(s.storage))},
		{Method: http.MethodPost, Path: "/additions", Handler: auth(handler.HandleAdditions(s.storage))},
		{Method: http.MethodGet, Path: "/additions/{uuid}", Handler: auth(handler.HandleAdditionByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/measurements", Handler: auth(handler.HandleMeasurements(s.storage))},
		{Method: http.MethodPost, Path: "/measurements", Handler: auth(handler.HandleMeasurements(s.storage))},
		{Method: http.MethodGet, Path: "/measurements/{uuid}", Handler: auth(handler.HandleMeasurementByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/package-formats", Handler: auth(handler.HandlePackageFormats(s.storage))},
		{Method: http.MethodPost, Path: "/package-formats", Handler: auth(handler.HandlePackageFormats(s.storage))},
		{Method: http.MethodGet, Path: "/package-formats/{uuid}", Handler: auth(handler.HandlePackageFormatByUUID(s.storage))},
		{Method: http.MethodPatch, Path: "/package-formats/{uuid}", Handler: auth(handler.HandlePackageFormatByUUID(s.storage))},
		{Method: http.MethodDelete, Path: "/package-formats/{uuid}", Handler: auth(handler.HandlePackageFormatByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/packaging-runs", Handler: auth(handler.HandlePackagingRuns(s.storage, s.inventoryClient))},
		{Method: http.MethodPost, Path: "/packaging-runs", Handler: auth(handler.HandlePackagingRuns(s.storage, s.inventoryClient))},
		{Method: http.MethodGet, Path: "/packaging-runs/{uuid}", Handler: auth(handler.HandlePackagingRunByUUID(s.storage))},
		{Method: http.MethodDelete, Path: "/packaging-runs/{uuid}", Handler: auth(handler.HandlePackagingRunByUUID(s.storage))},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("production service starting")
	if s.secretKey == "" {
		return fmt.Errorf("missing BREWPIPES_SECRET_KEY for access token verification")
	}
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	// Initialize inventory client for inter-service beer lot creation.
	inventoryURL := os.Getenv("INVENTORY_API_URL")
	if inventoryURL == "" {
		inventoryURL = "http://localhost:8080/api"
	}
	s.inventoryClient = handler.NewInventoryClient(inventoryURL)
	slog.Info("inventory client configured", "inventory_api_url", inventoryURL)

	return nil
}
