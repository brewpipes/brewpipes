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

type VolumeRelationStore interface {
	CreateVolumeRelation(context.Context, storage.VolumeRelation) (storage.VolumeRelation, error)
	GetVolumeRelationByUUID(context.Context, string) (storage.VolumeRelation, error)
	ListVolumeRelationsByVolumeUUID(context.Context, string) ([]storage.VolumeRelation, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleVolumeRelations handles [GET /volume-relations] and [POST /volume-relations].
func HandleVolumeRelations(db VolumeRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			volumeUUID := r.URL.Query().Get("volume_uuid")
			if volumeUUID == "" {
				http.Error(w, "volume_uuid is required", http.StatusBadRequest)
				return
			}

			relations, err := db.ListVolumeRelationsByVolumeUUID(r.Context(), volumeUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "volume not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error listing volume relations", "error", err)
				return
			}

			service.JSON(w, dto.NewVolumeRelationsResponse(relations))
		case http.MethodPost:
			var req dto.CreateVolumeRelationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve parent volume UUID to internal ID
			parentVol, ok := service.ResolveFK(r.Context(), w, req.ParentVolumeUUID, "parent volume", db.GetVolumeByUUID)
			if !ok {
				return
			}

			// Resolve child volume UUID to internal ID
			childVol, ok := service.ResolveFK(r.Context(), w, req.ChildVolumeUUID, "child volume", db.GetVolumeByUUID)
			if !ok {
				return
			}

			relation := storage.VolumeRelation{
				ParentVolumeID: parentVol.ID,
				ChildVolumeID:  childVol.ID,
				RelationType:   req.RelationType,
				Amount:         req.Amount,
				AmountUnit:     req.AmountUnit,
			}

			created, err := db.CreateVolumeRelation(r.Context(), relation)
			if err != nil {
				service.InternalError(w, "error creating volume relation", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewVolumeRelationResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleVolumeRelationByUUID handles [GET /volume-relations/{uuid}].
func HandleVolumeRelationByUUID(db VolumeRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		relationUUID := r.PathValue("uuid")
		if relationUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		relation, err := db.GetVolumeRelationByUUID(r.Context(), relationUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "volume relation not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting volume relation", "error", err, "relation_uuid", relationUUID)
			return
		}

		service.JSON(w, dto.NewVolumeRelationResponse(relation))
	}
}
