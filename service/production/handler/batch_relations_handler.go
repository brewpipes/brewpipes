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
	GetBatchRelation(context.Context, int64) (storage.BatchRelation, error)
	ListBatchRelations(context.Context, int64) ([]storage.BatchRelation, error)
}

// HandleBatchRelations handles [GET /batch-relations] and [POST /batch-relations].
func HandleBatchRelations(db BatchRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchValue := r.URL.Query().Get("batch_id")
			if batchValue == "" {
				http.Error(w, "batch_id is required", http.StatusBadRequest)
				return
			}
			batchID, err := parseInt64Param(batchValue)
			if err != nil {
				http.Error(w, "invalid batch_id", http.StatusBadRequest)
				return
			}

			relations, err := db.ListBatchRelations(r.Context(), batchID)
			if err != nil {
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

			relation := storage.BatchRelation{
				ParentBatchID: req.ParentBatchID,
				ChildBatchID:  req.ChildBatchID,
				RelationType:  req.RelationType,
				VolumeID:      req.VolumeID,
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

// HandleBatchRelationByID handles [GET /batch-relations/{id}].
func HandleBatchRelationByID(db BatchRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		relationID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		relation, err := db.GetBatchRelation(r.Context(), relationID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch relation not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch relation", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchRelationResponse(relation))
	}
}
