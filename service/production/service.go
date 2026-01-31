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
}

type Service struct {
	storage *storage.Client
}

// New creates and initializes a new production service instance.
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
		{Method: http.MethodGet, Path: "/batches", Handler: handler.HandleBatches(s.storage)},
		{Method: http.MethodPost, Path: "/batches", Handler: handler.HandleBatches(s.storage)},
		{Method: http.MethodGet, Path: "/batches/{id}", Handler: handler.HandleBatchByID(s.storage)},
		{Method: http.MethodGet, Path: "/volumes", Handler: handler.HandleVolumes(s.storage)},
		{Method: http.MethodPost, Path: "/volumes", Handler: handler.HandleVolumes(s.storage)},
		{Method: http.MethodGet, Path: "/volumes/{id}", Handler: handler.HandleVolumeByID(s.storage)},
		{Method: http.MethodGet, Path: "/volume-relations", Handler: handler.HandleVolumeRelations(s.storage)},
		{Method: http.MethodPost, Path: "/volume-relations", Handler: handler.HandleVolumeRelations(s.storage)},
		{Method: http.MethodGet, Path: "/volume-relations/{id}", Handler: handler.HandleVolumeRelationByID(s.storage)},
		{Method: http.MethodGet, Path: "/vessels", Handler: handler.HandleVessels(s.storage)},
		{Method: http.MethodPost, Path: "/vessels", Handler: handler.HandleVessels(s.storage)},
		{Method: http.MethodGet, Path: "/vessels/{id}", Handler: handler.HandleVesselByID(s.storage)},
		{Method: http.MethodPost, Path: "/occupancies", Handler: handler.HandleCreateOccupancy(s.storage)},
		{Method: http.MethodGet, Path: "/occupancies/{id}", Handler: handler.HandleOccupancyByID(s.storage)},
		{Method: http.MethodGet, Path: "/occupancies/active", Handler: handler.HandleActiveOccupancy(s.storage)},
		{Method: http.MethodGet, Path: "/transfers", Handler: handler.HandleTransfers(s.storage)},
		{Method: http.MethodPost, Path: "/transfers", Handler: handler.HandleTransfers(s.storage)},
		{Method: http.MethodGet, Path: "/transfers/{id}", Handler: handler.HandleTransferByID(s.storage)},
		{Method: http.MethodGet, Path: "/batch-volumes", Handler: handler.HandleBatchVolumes(s.storage)},
		{Method: http.MethodPost, Path: "/batch-volumes", Handler: handler.HandleBatchVolumes(s.storage)},
		{Method: http.MethodGet, Path: "/batch-volumes/{id}", Handler: handler.HandleBatchVolumeByID(s.storage)},
		{Method: http.MethodGet, Path: "/batch-process-phases", Handler: handler.HandleBatchProcessPhases(s.storage)},
		{Method: http.MethodPost, Path: "/batch-process-phases", Handler: handler.HandleBatchProcessPhases(s.storage)},
		{Method: http.MethodGet, Path: "/batch-process-phases/{id}", Handler: handler.HandleBatchProcessPhaseByID(s.storage)},
		{Method: http.MethodGet, Path: "/batch-relations", Handler: handler.HandleBatchRelations(s.storage)},
		{Method: http.MethodPost, Path: "/batch-relations", Handler: handler.HandleBatchRelations(s.storage)},
		{Method: http.MethodGet, Path: "/batch-relations/{id}", Handler: handler.HandleBatchRelationByID(s.storage)},
		{Method: http.MethodGet, Path: "/additions", Handler: handler.HandleAdditions(s.storage)},
		{Method: http.MethodPost, Path: "/additions", Handler: handler.HandleAdditions(s.storage)},
		{Method: http.MethodGet, Path: "/additions/{id}", Handler: handler.HandleAdditionByID(s.storage)},
		{Method: http.MethodGet, Path: "/measurements", Handler: handler.HandleMeasurements(s.storage)},
		{Method: http.MethodPost, Path: "/measurements", Handler: handler.HandleMeasurements(s.storage)},
		{Method: http.MethodGet, Path: "/measurements/{id}", Handler: handler.HandleMeasurementByID(s.storage)},
	}
}

func (s *Service) Start(ctx context.Context) error {
	slog.Info("production service starting")
	if err := s.storage.Start(ctx); err != nil {
		return fmt.Errorf("starting storage: %w", err)
	}

	return nil
}
