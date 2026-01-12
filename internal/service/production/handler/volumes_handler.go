package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipesproto/internal/service"
	"github.com/brewpipes/brewpipesproto/internal/service/production/storage"
)

type VolumeGetter interface {
	GetVolumes(context.Context) ([]storage.Volume, error)
}

// HandleGetVolumes handles [GET /volumes].
func HandleGetVolumes(db VolumeGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		volumes, err := db.GetVolumes(r.Context())
		if err != nil {
			slog.Error("error getting volumes", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		service.JSON(w, volumes)
	}
}
