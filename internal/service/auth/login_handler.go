package auth

import (
	"log/slog"
	"net/http"
)

// handleLogin handles [POST /login].
func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Login endpoint\n")); err != nil {
		slog.Error("error writing login response", "error", err)
	}
}
