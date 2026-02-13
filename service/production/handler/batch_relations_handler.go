package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BatchRelationStore interface {
	CreateBatchRelation(context.Context, storage.BatchRelation) (storage.BatchRelation, error)
	GetBatchRelationByUUID(context.Context, string) (storage.BatchRelation, error)
	ListBatchRelationsByBatchUUID(context.Context, string) ([]storage.BatchRelation, error)
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleBatchRelations handles [GET /batch-relations] and [POST /batch-relations].
func HandleBatchRelations(db BatchRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchUUID := r.URL.Query().Get("batch_uuid")
			if batchUUID == "" {
				http.Error(w, "batch_uuid is required", http.StatusBadRequest)
				return
			}

			relations, err := db.ListBatchRelationsByBatchUUID(r.Context(), batchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error listing batch relations", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBatchRelationsResponse(relations))
		case http.MethodPost:
			var req dto.CreateBatchRelationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve parent batch UUID to internal ID
			parentBatch, err := db.GetBatchByUUID(r.Context(), req.ParentBatchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "parent batch not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving parent batch uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			// Resolve child batch UUID to internal ID
			childBatch, err := db.GetBatchByUUID(r.Context(), req.ChildBatchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "child batch not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving child batch uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			relation := storage.BatchRelation{
				ParentBatchID: parentBatch.ID,
				ChildBatchID:  childBatch.ID,
				RelationType:  req.RelationType,
			}

			// Resolve optional volume UUID to internal ID
			if req.VolumeUUID != nil {
				vol, err := db.GetVolumeByUUID(r.Context(), *req.VolumeUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "volume not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving volume uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				relation.VolumeID = &vol.ID
			}

			created, err := db.CreateBatchRelation(r.Context(), relation)
			if err != nil {
				slog.Error("error creating batch relation", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBatchRelationResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBatchRelationByUUID handles [GET /batch-relations/{uuid}].
func HandleBatchRelationByUUID(db BatchRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		relationUUID := r.PathValue("uuid")
		if relationUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		relation, err := db.GetBatchRelationByUUID(r.Context(), relationUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch relation not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch relation", "error", err, "relation_uuid", relationUUID)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchRelationResponse(relation))
	}
}
