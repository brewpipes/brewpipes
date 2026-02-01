package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BrewSessionStore interface {
	CreateBrewSession(context.Context, storage.BrewSession) (storage.BrewSession, error)
	GetBrewSession(context.Context, int64) (storage.BrewSession, error)
	ListBrewSessionsByBatch(context.Context, int64) ([]storage.BrewSession, error)
	UpdateBrewSession(context.Context, int64, storage.BrewSession) (storage.BrewSession, error)
}

// HandleBrewSessions handles [GET /brew-sessions] and [POST /brew-sessions].
func HandleBrewSessions(db BrewSessionStore) http.HandlerFunc {
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

			sessions, err := db.ListBrewSessionsByBatch(r.Context(), batchID)
			if err != nil {
				slog.Error("error listing brew sessions", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBrewSessionsResponse(sessions))
		case http.MethodPost:
			var req dto.CreateBrewSessionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			brewedAt := time.Time{}
			if req.BrewedAt != nil {
				brewedAt = *req.BrewedAt
			}

			session := storage.BrewSession{
				BatchID:      req.BatchID,
				WortVolumeID: req.WortVolumeID,
				MashVesselID: req.MashVesselID,
				BoilVesselID: req.BoilVesselID,
				BrewedAt:     brewedAt,
				Notes:        req.Notes,
			}

			created, err := db.CreateBrewSession(r.Context(), session)
			if err != nil {
				slog.Error("error creating brew session", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("brew session created", "brew_session_id", created.ID, "batch_id", created.BatchID)

			service.JSON(w, dto.NewBrewSessionResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBrewSessionByID handles [GET /brew-sessions/{id}] and [PUT /brew-sessions/{id}].
func HandleBrewSessionByID(db BrewSessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		sessionID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			session, err := db.GetBrewSession(r.Context(), sessionID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "brew session not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting brew session", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBrewSessionResponse(session))
		case http.MethodPut:
			var req dto.UpdateBrewSessionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			brewedAt := time.Time{}
			if req.BrewedAt != nil {
				brewedAt = *req.BrewedAt
			}

			session := storage.BrewSession{
				BatchID:      req.BatchID,
				WortVolumeID: req.WortVolumeID,
				MashVesselID: req.MashVesselID,
				BoilVesselID: req.BoilVesselID,
				BrewedAt:     brewedAt,
				Notes:        req.Notes,
			}

			updated, err := db.UpdateBrewSession(r.Context(), sessionID, session)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "brew session not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating brew session", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("brew session updated", "brew_session_id", updated.ID)

			service.JSON(w, dto.NewBrewSessionResponse(updated))
		default:
			methodNotAllowed(w)
		}
	}
}
