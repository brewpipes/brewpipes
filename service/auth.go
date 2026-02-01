package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/brewpipes/brewpipes/internal/jwt"
	"github.com/gofrs/uuid/v5"
)

type contextKey string

const (
	contextUserIDKey   contextKey = "userID"
	contextUserRoleKey contextKey = "userRole"
)

// UserIDFromContext returns the authenticated user ID when present.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value := ctx.Value(contextUserIDKey)
	if value == nil {
		return uuid.UUID{}, false
	}

	id, ok := value.(uuid.UUID)
	return id, ok
}

// UserRoleFromContext returns the authenticated user role when present.
func UserRoleFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(contextUserRoleKey)
	if value == nil {
		return "", false
	}

	role, ok := value.(string)
	return role, ok
}

// RequireAccessToken validates bearer access tokens and populates context.
func RequireAccessToken(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			rawToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if rawToken == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			accessToken, err := jwt.DecodeAccessToken(rawToken, secret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextUserIDKey, accessToken.UserID)
			ctx = context.WithValue(ctx, contextUserRoleKey, accessToken.Claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
