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
	GetBrewSessionByUUID(context.Context, string) (storage.BrewSession, error)
	ListBrewSessionsByBatchUUID(context.Context, string) ([]storage.BrewSession, error)
	UpdateBrewSession(context.Context, int64, storage.BrewSession) (storage.BrewSession, error)
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
	GetVesselByUUID(context.Context, string) (storage.Vessel, error)
}

// resolveBrewSessionFKs resolves optional UUID fields to internal IDs for brew session create/update.
func resolveBrewSessionFKs(ctx context.Context, db BrewSessionStore, w http.ResponseWriter,
	batchUUID, wortVolumeUUID, mashVesselUUID, boilVesselUUID *string,
	session *storage.BrewSession) bool {

	if batch, ok := service.ResolveFKOptional(ctx, w, batchUUID, "batch", db.GetBatchByUUID); !ok {
		return false
	} else if batchUUID != nil {
		session.BatchID = &batch.ID
	}
	if vol, ok := service.ResolveFKOptional(ctx, w, wortVolumeUUID, "wort volume", db.GetVolumeByUUID); !ok {
		return false
	} else if wortVolumeUUID != nil {
		session.WortVolumeID = &vol.ID
	}
	if vessel, ok := service.ResolveFKOptional(ctx, w, mashVesselUUID, "mash vessel", db.GetVesselByUUID); !ok {
		return false
	} else if mashVesselUUID != nil {
		session.MashVesselID = &vessel.ID
	}
	if vessel, ok := service.ResolveFKOptional(ctx, w, boilVesselUUID, "boil vessel", db.GetVesselByUUID); !ok {
		return false
	} else if boilVesselUUID != nil {
		session.BoilVesselID = &vessel.ID
	}
	return true
}

// HandleBrewSessions handles [GET /brew-sessions] and [POST /brew-sessions].
func HandleBrewSessions(db BrewSessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchUUID := r.URL.Query().Get("batch_uuid")
			if batchUUID == "" {
				http.Error(w, "batch_uuid is required", http.StatusBadRequest)
				return
			}

			sessions, err := db.ListBrewSessionsByBatchUUID(r.Context(), batchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error listing brew sessions", "error", err)
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
				BrewedAt: brewedAt,
				Notes:    req.Notes,
			}

			if !resolveBrewSessionFKs(r.Context(), db, w, req.BatchUUID, req.WortVolumeUUID, req.MashVesselUUID, req.BoilVesselUUID, &session) {
				return
			}

			created, err := db.CreateBrewSession(r.Context(), session)
			if err != nil {
				service.InternalError(w, "error creating brew session", "error", err)
				return
			}

			slog.Info("brew session created", "brew_session_uuid", created.UUID, "batch_uuid", req.BatchUUID)

			service.JSONCreated(w, dto.NewBrewSessionResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleBrewSessionByUUID handles [GET /brew-sessions/{uuid}] and [PUT /brew-sessions/{uuid}].
func HandleBrewSessionByUUID(db BrewSessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.PathValue("uuid")
		if sessionUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			session, err := db.GetBrewSessionByUUID(r.Context(), sessionUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "brew session not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting brew session", "error", err, "session_uuid", sessionUUID)
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

			// Resolve UUID to get internal ID for update
			existing, err := db.GetBrewSessionByUUID(r.Context(), sessionUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "brew session not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting brew session for update", "error", err)
				return
			}

			brewedAt := time.Time{}
			if req.BrewedAt != nil {
				brewedAt = *req.BrewedAt
			}

			session := storage.BrewSession{
				BrewedAt: brewedAt,
				Notes:    req.Notes,
			}

			if !resolveBrewSessionFKs(r.Context(), db, w, req.BatchUUID, req.WortVolumeUUID, req.MashVesselUUID, req.BoilVesselUUID, &session) {
				return
			}

			updated, err := db.UpdateBrewSession(r.Context(), existing.ID, session)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "brew session not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error updating brew session", "error", err)
				return
			}

			slog.Info("brew session updated", "brew_session_uuid", sessionUUID)

			service.JSON(w, dto.NewBrewSessionResponse(updated))
		default:
			service.MethodNotAllowed(w)
		}
	}
}
