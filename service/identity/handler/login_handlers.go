package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/internal/jwt"
	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/identity/handler/dto"
	"github.com/brewpipes/brewpipes/service/identity/storage"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserGetter interface {
	GetUser(ctx context.Context, id uuid.UUID) (storage.User, error)
	GetUserByUsername(ctx context.Context, username string) (storage.User, error)
}

type RefreshTokenStore interface {
	CreateRefreshToken(ctx context.Context, token storage.RefreshToken) (storage.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (storage.RefreshToken, error)
	ConsumeRefreshToken(ctx context.Context, tokenID uuid.UUID) error
}

// handleLogin handles [POST /login].
func HandleLogin(db UserGetter, tokenStore RefreshTokenStore, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.LoginRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByUsername(r.Context(), req.Username)
		if err != nil && !errors.Is(err, service.ErrNotFound) {
			slog.Error("error getting user", "error", err)
			service.InternalError(w, "error getting user", "error", err)
			return
		}

		if err != nil || !compareHashedPassword(user.Password, req.Password) {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		access, refresh, err := user.GenerateTokens(secretKey)
		if err != nil {
			slog.Error("error generating tokens", "error", err)
			service.InternalError(w, "error generating tokens", "error", err)
			return
		}

		refreshToken, err := jwt.DecodeRefreshToken(refresh, secretKey)
		if err != nil || refreshToken.Claims == nil || refreshToken.Claims.ID == "" {
			slog.Error("error decoding refresh token", "error", err)
			service.InternalError(w, "error decoding refresh token", "error", err)
			return
		}
		refreshID, err := uuid.FromString(refreshToken.Claims.ID)
		if err != nil {
			slog.Error("invalid refresh token id", "error", err)
			service.InternalError(w, "invalid refresh token id", "error", err)
			return
		}
		if refreshToken.Claims.ExpiresAt == nil {
			service.InternalError(w, "refresh token missing expiry")
			return
		}

		_, err = tokenStore.CreateRefreshToken(r.Context(), storage.RefreshToken{
			TokenID:   refreshID,
			UserUUID:  user.UUID,
			ExpiresAt: refreshToken.Claims.ExpiresAt.Time,
		})
		if err != nil {
			slog.Error("error storing refresh token", "error", err)
			service.InternalError(w, "error storing refresh token", "error", err)
			return
		}

		resp := dto.TokenPair{
			AccessToken:  access,
			RefreshToken: refresh,
		}

		service.JSON(w, resp)
	}
}

// HandleRefresh handles [POST /refresh].
func HandleRefresh(db UserGetter, tokenStore RefreshTokenStore, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.RefreshRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "malformed request", http.StatusBadRequest)
			return
		} else if err := req.Validate(); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		token, err := jwt.DecodeRefreshToken(req.RefreshToken, secretKey)
		if err != nil {
			slog.Error("error decoding refresh token", "error", err)
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}
		if token.Claims == nil || token.Claims.ID == "" {
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}
		tokenID, err := uuid.FromString(token.Claims.ID)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}

		stored, err := tokenStore.GetRefreshToken(r.Context(), tokenID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		} else if err != nil {
			slog.Error("error loading refresh token", "error", err)
			service.InternalError(w, "error loading refresh token", "error", err)
			return
		}

		now := time.Now().UTC()
		if stored.RevokedAt != nil || stored.ExpiresAt.Before(now) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		if stored.UserUUID != token.UserID {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		if err := tokenStore.ConsumeRefreshToken(r.Context(), tokenID); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "invalid refresh token", http.StatusUnauthorized)
				return
			}
			slog.Error("error consuming refresh token", "error", err)
			service.InternalError(w, "error consuming refresh token", "error", err)
			return
		}

		user, err := db.GetUser(r.Context(), token.UserID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		} else if err != nil {
			slog.Error("error getting user", "error", err)
			service.InternalError(w, "error getting user", "error", err)
			return
		}

		access, refresh, err := user.GenerateTokens(secretKey)
		if err != nil {
			slog.Error("error generating tokens", "error", err)
			service.InternalError(w, "error generating tokens", "error", err)
			return
		}
		refreshToken, err := jwt.DecodeRefreshToken(refresh, secretKey)
		if err != nil || refreshToken.Claims == nil || refreshToken.Claims.ID == "" {
			slog.Error("error decoding refresh token", "error", err)
			service.InternalError(w, "error decoding refresh token", "error", err)
			return
		}
		refreshID, err := uuid.FromString(refreshToken.Claims.ID)
		if err != nil {
			slog.Error("invalid refresh token id", "error", err)
			service.InternalError(w, "invalid refresh token id", "error", err)
			return
		}
		if refreshToken.Claims.ExpiresAt == nil {
			service.InternalError(w, "refresh token missing expiry")
			return
		}

		_, err = tokenStore.CreateRefreshToken(r.Context(), storage.RefreshToken{
			TokenID:   refreshID,
			UserUUID:  user.UUID,
			ExpiresAt: refreshToken.Claims.ExpiresAt.Time,
		})
		if err != nil {
			slog.Error("error storing refresh token", "error", err)
			service.InternalError(w, "error storing refresh token", "error", err)
			return
		}

		resp := dto.TokenPair{
			AccessToken:  access,
			RefreshToken: refresh,
		}

		service.JSON(w, resp)
	}
}

// HandleLogout handles [POST /logout].
func HandleLogout(tokenStore RefreshTokenStore, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.LogoutRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "malformed request", http.StatusBadRequest)
			return
		} else if err := req.Validate(); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		token, err := jwt.DecodeRefreshToken(req.RefreshToken, secretKey)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}
		if token.Claims == nil || token.Claims.ID == "" {
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}
		tokenID, err := uuid.FromString(token.Claims.ID)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusBadRequest)
			return
		}

		if err := tokenStore.ConsumeRefreshToken(r.Context(), tokenID); err != nil && !errors.Is(err, service.ErrNotFound) {
			slog.Error("error revoking refresh token", "error", err)
			service.InternalError(w, "error revoking refresh token", "error", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func hashPassword(pw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func compareHashedPassword(hashed, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw)) == nil
}
