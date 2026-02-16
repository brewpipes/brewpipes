package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

// PackagingRunStore defines the storage methods needed by packaging run handlers.
type PackagingRunStore interface {
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetOccupancyByUUID(context.Context, string) (storage.Occupancy, error)
	GetPackageFormatByUUID(context.Context, string) (storage.PackageFormat, error)
	CreatePackagingRunWithLines(context.Context, storage.PackagingRun, []storage.PackagingRunLine) (storage.PackagingRun, []storage.PackagingRunLine, error)
	GetPackagingRunByUUID(context.Context, string) (storage.PackagingRun, error)
	ListPackagingRuns(context.Context) ([]storage.PackagingRun, error)
	ListPackagingRunsByBatchUUID(context.Context, string) ([]storage.PackagingRun, error)
	DeletePackagingRun(context.Context, int64) error
	ListPackagingRunLinesByRunID(context.Context, int64) ([]storage.PackagingRunLine, error)
	CloseOccupancy(context.Context, int64, time.Time) error
}

// BeerLotCreator abstracts the inter-service call to the Inventory service for creating beer lots.
type BeerLotCreator interface {
	CreateBeerLot(ctx context.Context, authToken string, req BeerLotRequest) (*BeerLotResponse, error)
}

// BeerLotRequest is the request payload for creating a beer lot via the Inventory service.
type BeerLotRequest struct {
	ProductionBatchUUID string    `json:"production_batch_uuid"`
	PackagingRunUUID    string    `json:"packaging_run_uuid"`
	LotCode             *string   `json:"lot_code,omitempty"`
	PackageFormatName   string    `json:"package_format_name"`
	Container           string    `json:"container"`
	VolumePerUnit       int64     `json:"volume_per_unit"`
	VolumePerUnitUnit   string    `json:"volume_per_unit_unit"`
	Quantity            int       `json:"quantity"`
	PackagedAt          time.Time `json:"packaged_at"`
	StockLocationUUID   string    `json:"stock_location_uuid"`
	Notes               *string   `json:"notes,omitempty"`
}

// BeerLotResponse is the response from the Inventory service after creating a beer lot.
type BeerLotResponse struct {
	UUID string `json:"uuid"`
}

// containerAbbreviation maps container types to short codes for lot code generation.
var containerAbbreviation = map[string]string{
	"keg":     "KEG",
	"can":     "CAN",
	"bottle":  "BTL",
	"cask":    "CSK",
	"growler": "GRL",
	"other":   "OTH",
}

// HandlePackagingRuns handles [GET /packaging-runs] and [POST /packaging-runs].
func HandlePackagingRuns(db PackagingRunStore, invClient BeerLotCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchUUID := r.URL.Query().Get("batch_uuid")

			var runs []storage.PackagingRun
			var err error
			if batchUUID != "" {
				runs, err = db.ListPackagingRunsByBatchUUID(r.Context(), batchUUID)
			} else {
				runs, err = db.ListPackagingRuns(r.Context())
			}
			if err != nil {
				service.InternalError(w, "error listing packaging runs", "error", err)
				return
			}

			// Load lines for each run
			linesByRunID := make(map[int64][]storage.PackagingRunLine, len(runs))
			for _, run := range runs {
				lines, err := db.ListPackagingRunLinesByRunID(r.Context(), run.ID)
				if err != nil {
					service.InternalError(w, "error listing packaging run lines", "error", err, "packaging_run_uuid", run.UUID)
					return
				}
				linesByRunID[run.ID] = lines
			}

			service.JSON(w, dto.NewPackagingRunsResponse(runs, linesByRunID))
		case http.MethodPost:
			var req dto.CreatePackagingRunRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve batch UUID to internal ID
			batch, ok := service.ResolveFK(r.Context(), w, req.BatchUUID, "batch", db.GetBatchByUUID)
			if !ok {
				return
			}

			// Resolve occupancy UUID to internal ID
			occupancy, ok := service.ResolveFK(r.Context(), w, req.OccupancyUUID, "occupancy", db.GetOccupancyByUUID)
			if !ok {
				return
			}

			// Resolve each line's package format UUID to internal ID and collect format details
			type resolvedLine struct {
				format   storage.PackageFormat
				quantity int
			}
			resolvedLines := make([]resolvedLine, 0, len(req.Lines))
			for _, line := range req.Lines {
				format, ok := service.ResolveFK(r.Context(), w, line.PackageFormatUUID, "package format", db.GetPackageFormatByUUID)
				if !ok {
					return
				}
				resolvedLines = append(resolvedLines, resolvedLine{
					format:   format,
					quantity: line.Quantity,
				})
			}

			startedAt := time.Now().UTC()
			if req.StartedAt != nil {
				startedAt = *req.StartedAt
			}

			run := storage.PackagingRun{
				BatchID:     batch.ID,
				OccupancyID: occupancy.ID,
				StartedAt:   startedAt,
				EndedAt:     req.EndedAt,
				LossAmount:  req.LossAmount,
				LossUnit:    req.LossUnit,
				Notes:       req.Notes,
			}

			// Build storage lines from resolved formats
			storageLines := make([]storage.PackagingRunLine, 0, len(resolvedLines))
			for _, rl := range resolvedLines {
				storageLines = append(storageLines, storage.PackagingRunLine{
					PackageFormatID: rl.format.ID,
					Quantity:        rl.quantity,
				})
			}

			// Create packaging run and lines atomically
			created, createdLines, err := db.CreatePackagingRunWithLines(r.Context(), run, storageLines)
			if err != nil {
				service.InternalError(w, "error creating packaging run", "error", err)
				return
			}

			// Close source occupancy if requested (default true)
			closeSource := true
			if req.CloseSource != nil {
				closeSource = *req.CloseSource
			}
			if closeSource {
				outAt := created.StartedAt
				if err := db.CloseOccupancy(r.Context(), occupancy.ID, outAt); err != nil {
					service.InternalError(w, "error closing source occupancy during packaging run",
						"error", err,
						"packaging_run_uuid", created.UUID,
						"occupancy_id", occupancy.ID,
					)
					return
				}
			}

			// Create beer lots in inventory if stock_location_uuid is provided
			if req.StockLocationUUID != nil && invClient != nil {
				// Extract auth token from the inbound request
				authHeader := r.Header.Get("Authorization")
				authToken := strings.TrimPrefix(authHeader, "Bearer ")

				// Determine lot code prefix
				lotCodePrefix := batch.ShortName
				if req.LotCodePrefix != nil && *req.LotCodePrefix != "" {
					lotCodePrefix = *req.LotCodePrefix
				}

				packagedAt := created.StartedAt

				for i, rl := range resolvedLines {
					abbrev := containerAbbreviation[rl.format.Container]
					if abbrev == "" {
						abbrev = "OTH"
					}
					lotCode := fmt.Sprintf("%s-%s-%02d", lotCodePrefix, abbrev, i+1)

					beerLotReq := BeerLotRequest{
						ProductionBatchUUID: req.BatchUUID,
						PackagingRunUUID:    created.UUID.String(),
						LotCode:             &lotCode,
						PackageFormatName:   rl.format.Name,
						Container:           rl.format.Container,
						VolumePerUnit:       rl.format.VolumePerUnit,
						VolumePerUnitUnit:   rl.format.VolumePerUnitUnit,
						Quantity:            rl.quantity,
						PackagedAt:          packagedAt,
						StockLocationUUID:   *req.StockLocationUUID,
					}

					resp, err := invClient.CreateBeerLot(r.Context(), authToken, beerLotReq)
					if err != nil {
						slog.Error("error creating beer lot in inventory (best-effort)",
							"error", err,
							"packaging_run_uuid", created.UUID,
							"line_index", i,
							"lot_code", lotCode,
						)
						continue
					}
					slog.Info("beer lot created in inventory",
						"beer_lot_uuid", resp.UUID,
						"packaging_run_uuid", created.UUID,
						"lot_code", lotCode,
					)
				}
			}

			slog.Info("packaging run created",
				"packaging_run_uuid", created.UUID,
				"batch_uuid", req.BatchUUID,
				"lines", len(createdLines),
			)

			service.JSONCreated(w, dto.NewPackagingRunResponse(created, createdLines))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandlePackagingRunByUUID handles [GET /packaging-runs/{uuid}] and [DELETE /packaging-runs/{uuid}].
func HandlePackagingRunByUUID(db PackagingRunStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runUUID := r.PathValue("uuid")
		if runUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			run, err := db.GetPackagingRunByUUID(r.Context(), runUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "packaging run not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting packaging run", "error", err, "packaging_run_uuid", runUUID)
				return
			}

			lines, err := db.ListPackagingRunLinesByRunID(r.Context(), run.ID)
			if err != nil {
				service.InternalError(w, "error listing packaging run lines", "error", err, "packaging_run_uuid", runUUID)
				return
			}

			service.JSON(w, dto.NewPackagingRunResponse(run, lines))
		case http.MethodDelete:
			run, err := db.GetPackagingRunByUUID(r.Context(), runUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "packaging run not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting packaging run for delete", "error", err, "packaging_run_uuid", runUUID)
				return
			}

			if err := db.DeletePackagingRun(r.Context(), run.ID); err != nil {
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "packaging run not found", http.StatusNotFound)
					return
				}
				service.InternalError(w, "error deleting packaging run", "error", err, "packaging_run_uuid", runUUID)
				return
			}

			slog.Info("packaging run deleted", "packaging_run_uuid", runUUID)

			w.WriteHeader(http.StatusNoContent)
		default:
			service.MethodNotAllowed(w)
		}
	}
}
