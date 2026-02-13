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
				slog.Error("error listing volume relations", "error", err)
				service.InternalError(w, err.Error())
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
			parentVol, err := db.GetVolumeByUUID(r.Context(), req.ParentVolumeUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "parent volume not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving parent volume uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			// Resolve child volume UUID to internal ID
			childVol, err := db.GetVolumeByUUID(r.Context(), req.ChildVolumeUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "child volume not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving child volume uuid", "error", err)
				service.InternalError(w, err.Error())
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
				slog.Error("error creating volume relation", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewVolumeRelationResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleVolumeRelationByUUID handles [GET /volume-relations/{uuid}].
func HandleVolumeRelationByUUID(db VolumeRelationStore) http.HandlerFunc {
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

		relation, err := db.GetVolumeRelationByUUID(r.Context(), relationUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "volume relation not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting volume relation", "error", err, "relation_uuid", relationUUID)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewVolumeRelationResponse(relation))
	}
}
