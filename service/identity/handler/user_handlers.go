package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/identity/handler/dto"
	"github.com/brewpipes/brewpipes/service/identity/storage"
	"github.com/gofrs/uuid/v5"
)

type UserStore interface {
	ListUsers(ctx context.Context) ([]storage.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (storage.User, error)
	CreateUser(ctx context.Context, user storage.User) (storage.User, error)
	UpdateUser(ctx context.Context, user storage.User) (storage.User, error)
	DeleteUser(ctx context.Context, userUUID uuid.UUID) error
	RevokeRefreshTokensForUser(ctx context.Context, userUUID uuid.UUID) error
}

// HandleUsers handles [GET /users] and [POST /users].
func HandleUsers(db UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			users, err := db.ListUsers(r.Context())
			if err != nil {
				slog.Error("error listing users", "error", err)
				service.InternalError(w, "error listing users", "error", err)
				return
			}

			service.JSON(w, dto.NewUsersResponse(users))
		case http.MethodPost:
			var req dto.CreateUserRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			hashed, err := hashPassword(req.Password)
			if err != nil {
				slog.Error("error hashing password", "error", err)
				service.InternalError(w, "error hashing password", "error", err)
				return
			}

			user := storage.User{
				Username: req.Username,
				Password: hashed,
			}

			created, err := db.CreateUser(r.Context(), user)
			if err != nil {
				slog.Error("error creating user", "error", err)
				service.InternalError(w, "error creating user", "error", err)
				return
			}

			service.JSON(w, dto.NewUserResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleUserByUUID handles [GET /users/{uuid}], [PUT /users/{uuid}], and [DELETE /users/{uuid}].
func HandleUserByUUID(db UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidValue := r.PathValue("uuid")
		if uuidValue == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}
		userUUID, err := parseUUIDParam(uuidValue)
		if err != nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			user, err := db.GetUser(r.Context(), userUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting user", "error", err)
				service.InternalError(w, "error getting user", "error", err)
				return
			}

			service.JSON(w, dto.NewUserResponse(user))
		case http.MethodPut:
			var req dto.UpdateUserRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user, err := db.GetUser(r.Context(), userUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting user", "error", err)
				service.InternalError(w, "error getting user", "error", err)
				return
			}

			if req.Username != nil {
				user.Username = *req.Username
			}
			if req.Password != nil {
				hashed, err := hashPassword(*req.Password)
				if err != nil {
					slog.Error("error hashing password", "error", err)
					service.InternalError(w, "error hashing password", "error", err)
					return
				}
				user.Password = hashed
			}

			updated, err := db.UpdateUser(r.Context(), user)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating user", "error", err)
				service.InternalError(w, "error updating user", "error", err)
				return
			}

			service.JSON(w, dto.NewUserResponse(updated))
		case http.MethodDelete:
			if err := db.DeleteUser(r.Context(), userUUID); err != nil {
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "user not found", http.StatusNotFound)
					return
				}
				slog.Error("error deleting user", "error", err)
				service.InternalError(w, "error deleting user", "error", err)
				return
			}
			if err := db.RevokeRefreshTokensForUser(r.Context(), userUUID); err != nil {
				slog.Error("error revoking refresh tokens", "error", err)
				service.InternalError(w, "error revoking refresh tokens", "error", err)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	}
}
