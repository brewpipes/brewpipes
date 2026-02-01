package production

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type Config struct {
	PostgresDSN string
	SecretKey   string
}

type Service struct {
	storage   *storage.Client
	secretKey string
}

// New creates and initializes a new production service instance.
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
		{Method: http.MethodGet, Path: "/styles", Handler: auth(handler.HandleStyles(s.storage))},
		{Method: http.MethodPost, Path: "/styles", Handler: auth(handler.HandleStyles(s.storage))},
		{Method: http.MethodGet, Path: "/styles/{id}", Handler: auth(handler.HandleStyleByID(s.storage))},
		{Method: http.MethodGet, Path: "/recipes", Handler: auth(handler.HandleRecipes(s.storage))},
		{Method: http.MethodPost, Path: "/recipes", Handler: auth(handler.HandleRecipes(s.storage))},
		{Method: http.MethodGet, Path: "/recipes/{id}", Handler: auth(handler.HandleRecipeByID(s.storage))},
		{Method: http.MethodPut, Path: "/recipes/{id}", Handler: auth(handler.HandleRecipeByID(s.storage))},
		{Method: http.MethodGet, Path: "/batches", Handler: auth(handler.HandleBatches(s.storage))},
		{Method: http.MethodPost, Path: "/batches", Handler: auth(handler.HandleBatches(s.storage))},
		{Method: http.MethodPost, Path: "/batches/import", Handler: auth(handler.HandleBatchImport(s.storage))},
		{Method: http.MethodGet, Path: "/batches/{id}", Handler: auth(handler.HandleBatchByID(s.storage))},
		{Method: http.MethodGet, Path: "/batches/{id}/summary", Handler: auth(handler.HandleBatchSummary(s.storage))},
		{Method: http.MethodGet, Path: "/brew-sessions", Handler: auth(handler.HandleBrewSessions(s.storage))},
		{Method: http.MethodPost, Path: "/brew-sessions", Handler: auth(handler.HandleBrewSessions(s.storage))},
		{Method: http.MethodGet, Path: "/brew-sessions/{id}", Handler: auth(handler.HandleBrewSessionByID(s.storage))},
		{Method: http.MethodPut, Path: "/brew-sessions/{id}", Handler: auth(handler.HandleBrewSessionByID(s.storage))},
		{Method: http.MethodGet, Path: "/volumes", Handler: auth(handler.HandleVolumes(s.storage))},
		{Method: http.MethodPost, Path: "/volumes", Handler: auth(handler.HandleVolumes(s.storage))},
		{Method: http.MethodGet, Path: "/volumes/{id}", Handler: auth(handler.HandleVolumeByID(s.storage))},
		{Method: http.MethodGet, Path: "/volume-relations", Handler: auth(handler.HandleVolumeRelations(s.storage))},
		{Method: http.MethodPost, Path: "/volume-relations", Handler: auth(handler.HandleVolumeRelations(s.storage))},
		{Method: http.MethodGet, Path: "/volume-relations/{id}", Handler: auth(handler.HandleVolumeRelationByID(s.storage))},
		{Method: http.MethodGet, Path: "/vessels", Handler: auth(handler.HandleVessels(s.storage))},
		{Method: http.MethodPost, Path: "/vessels", Handler: auth(handler.HandleVessels(s.storage))},
		{Method: http.MethodGet, Path: "/vessels/uuid/{uuid}", Handler: auth(handler.HandleVesselByUUID(s.storage))},
		{Method: http.MethodGet, Path: "/vessels/{id}", Handler: auth(handler.HandleVesselByID(s.storage))},
		{Method: http.MethodGet, Path: "/occupancies", Handler: auth(handler.HandleOccupancies(s.storage))},
		{Method: http.MethodPost, Path: "/occupancies", Handler: auth(handler.HandleCreateOccupancy(s.storage))},
		{Method: http.MethodGet, Path: "/occupancies/{id}", Handler: auth(handler.HandleOccupancyByID(s.storage))},
		{Method: http.MethodPatch, Path: "/occupancies/{id}/status", Handler: auth(handler.HandleOccupancyStatus(s.storage))},
		{Method: http.MethodGet, Path: "/occupancies/active", Handler: auth(handler.HandleActiveOccupancy(s.storage))},
		{Method: http.MethodGet, Path: "/transfers", Handler: auth(handler.HandleTransfers(s.storage))},
		{Method: http.MethodPost, Path: "/transfers", Handler: auth(handler.HandleTransfers(s.storage))},
		{Method: http.MethodGet, Path: "/transfers/{id}", Handler: auth(handler.HandleTransferByID(s.storage))},
		{Method: http.MethodGet, Path: "/batch-volumes", Handler: auth(handler.HandleBatchVolumes(s.storage))},
		{Method: http.MethodPost, Path: "/batch-volumes", Handler: auth(handler.HandleBatchVolumes(s.storage))},
		{Method: http.MethodGet, Path: "/batch-volumes/{id}", Handler: auth(handler.HandleBatchVolumeByID(s.storage))},
		{Method: http.MethodGet, Path: "/batch-process-phases", Handler: auth(handler.HandleBatchProcessPhases(s.storage))},
		{Method: http.MethodPost, Path: "/batch-process-phases", Handler: auth(handler.HandleBatchProcessPhases(s.storage))},
		{Method: http.MethodGet, Path: "/batch-process-phases/{id}", Handler: auth(handler.HandleBatchProcessPhaseByID(s.storage))},
		{Method: http.MethodGet, Path: "/batch-relations", Handler: auth(handler.HandleBatchRelations(s.storage))},
		{Method: http.MethodPost, Path: "/batch-relations", Handler: auth(handler.HandleBatchRelations(s.storage))},
		{Method: http.MethodGet, Path: "/batch-relations/{id}", Handler: auth(handler.HandleBatchRelationByID(s.storage))},
		{Method: http.MethodGet, Path: "/additions", Handler: auth(handler.HandleAdditions(s.storage))},
		{Method: http.MethodPost, Path: "/additions", Handler: auth(handler.HandleAdditions(s.storage))},
		{Method: http.MethodGet, Path: "/additions/{id}", Handler: auth(handler.HandleAdditionByID(s.storage))},
		{Method: http.MethodGet, Path: "/measurements", Handler: auth(handler.HandleMeasurements(s.storage))},
		{Method: http.MethodPost, Path: "/measurements", Handler: auth(handler.HandleMeasurements(s.storage))},
		{Method: http.MethodGet, Path: "/measurements/{id}", Handler: auth(handler.HandleMeasurementByID(s.storage))},
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

	return nil
}
