package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/internal/jwt"
	"github.com/brewpipes/brewpipes/internal/service"
	"github.com/brewpipes/brewpipes/internal/service/identity/handler/dto"
	"github.com/brewpipes/brewpipes/internal/service/identity/storage"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserGetter interface {
	GetUser(ctx context.Context, id uuid.UUID) (storage.User, error)
	GetUserByUsername(ctx context.Context, username string) (storage.User, error)
}

// handleLogin handles [POST /login].
func HandleLogin(db UserGetter, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := dto.LoginRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("error decoding JSON:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByUsername(r.Context(), req.Username)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			service.InternalError(w, err.Error())
			return
		}

		if !compareHashedPassword(user.Password, req.Password) {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		access, refresh, err := user.GenerateTokens(secretKey)
		if err != nil {
			service.InternalError(w, err.Error())
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
func HandleRefresh(db UserGetter, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
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
			log.Println("error decoding refresh token:", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if token.Claims.ExpiresAt.Before(now) {
			err := errors.New("refresh token expired")
			log.Println(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		user, err := db.GetUser(r.Context(), token.UserID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		} else if err != nil {
			service.InternalError(w, err.Error())
			return
		}

		access, refresh, err := user.GenerateTokens(secretKey)
		if err != nil {
			service.InternalError(w, err.Error())
			return
		}

		resp := dto.TokenPair{
			AccessToken:  access,
			RefreshToken: refresh,
		}

		service.JSON(w, resp)
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
