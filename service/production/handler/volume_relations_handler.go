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
	GetVolumeRelation(context.Context, int64) (storage.VolumeRelation, error)
	ListVolumeRelations(context.Context, int64) ([]storage.VolumeRelation, error)
}

// HandleVolumeRelations handles [GET /volume-relations] and [POST /volume-relations].
func HandleVolumeRelations(db VolumeRelationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			volumeValue := r.URL.Query().Get("volume_id")
			if volumeValue == "" {
				http.Error(w, "volume_id is required", http.StatusBadRequest)
				return
			}
			volumeID, err := parseInt64Param(volumeValue)
			if err != nil {
				http.Error(w, "invalid volume_id", http.StatusBadRequest)
				return
			}

			relations, err := db.ListVolumeRelations(r.Context(), volumeID)
			if err != nil {
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

			relation := storage.VolumeRelation{
				ParentVolumeID: req.ParentVolumeID,
				ChildVolumeID:  req.ChildVolumeID,
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

// HandleVolumeRelationByID handles [GET /volume-relations/{id}].
func HandleVolumeRelationByID(db VolumeRelationStore) http.HandlerFunc {
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

		relation, err := db.GetVolumeRelation(r.Context(), relationID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "volume relation not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting volume relation", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewVolumeRelationResponse(relation))
	}
}
