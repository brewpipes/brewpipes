package jwt

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
)

// A RefreshToken defines a refresh token.
type RefreshToken struct {
	Claims *RefreshClaims
	UserID uuid.UUID
}

// RefreshClaims defines the claims for a refresh token.
type RefreshClaims struct {
	jwt.RegisteredClaims
}

// Valid returns true if the refresh claims are valid.
func (c *RefreshClaims) Valid() error {
	return c.RegisteredClaims.Valid()
}

// DecodeRefreshToken decodes a refresh token.
func DecodeRefreshToken(token, secret string) (*RefreshToken, error) {
	t, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %w", fmt.Errorf("%v", token.Header["alg"]))
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing refresh token: %w", err)
	}

	if !t.Valid {
		slog.Warn("refresh token failed validation")
		return nil, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*RefreshClaims)
	if !ok {
		slog.Warn("refresh token has wrong claims type")
		return nil, errors.New("invalid token")
	}

	uid, err := uuid.FromString(claims.Subject)
	if err != nil {
		slog.Warn("refresh token has invalid subject UUID", "subject", claims.Subject)
		return nil, errors.New("invalid token")
	}

	return &RefreshToken{
		Claims: claims,
		UserID: uid,
	}, nil
}
