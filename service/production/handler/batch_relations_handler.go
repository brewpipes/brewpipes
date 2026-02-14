package handler

import (
	"context"
	"encoding/json"
	"errors"
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
				service.InternalError(w, "error listing batch relations", "error", err)
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
			parentBatch, ok := service.ResolveFK(r.Context(), w, req.ParentBatchUUID, "parent batch", db.GetBatchByUUID)
			if !ok {
				return
			}

			// Resolve child batch UUID to internal ID
			childBatch, ok := service.ResolveFK(r.Context(), w, req.ChildBatchUUID, "child batch", db.GetBatchByUUID)
			if !ok {
				return
			}

			relation := storage.BatchRelation{
				ParentBatchID: parentBatch.ID,
				ChildBatchID:  childBatch.ID,
				RelationType:  req.RelationType,
			}

			// Resolve optional volume UUID to internal ID
			if vol, ok := service.ResolveFKOptional(r.Context(), w, req.VolumeUUID, "volume", db.GetVolumeByUUID); !ok {
				return
			} else if req.VolumeUUID != nil {
				relation.VolumeID = &vol.ID
			}

			created, err := db.CreateBatchRelation(r.Context(), relation)
			if err != nil {
				service.InternalError(w, "error creating batch relation", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewBatchRelationResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleBatchRelationByUUID handles [GET /batch-relations/{uuid}].
func HandleBatchRelationByUUID(db BatchRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			service.InternalError(w, "error getting batch relation", "error", err, "relation_uuid", relationUUID)
			return
		}

		service.JSON(w, dto.NewBatchRelationResponse(relation))
	}
}
