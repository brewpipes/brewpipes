package production

import (
	"log/slog"
	"net/http"
)

// handleLogin handles [GET /volumes].
func (s *Service) handleGetVolumes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("[]")); err != nil {
		slog.Error("error writing login response", "error", err)
	}
}
