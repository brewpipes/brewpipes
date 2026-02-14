// Package jwt contains types and functions for manipulating JSON Web Tokens.
package jwt

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
)

// An AccessToken defines an access token.
type AccessToken struct {
	Claims *AccessClaims
	UserID uuid.UUID
}

// AccessClaims defines the claims for an access token.
type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// Valid returns true if the access claims are valid.
func (c *AccessClaims) Valid() error {
	if c.Role != "admin" && c.Role != "user" {
		return errors.New("invalid role")
	}

	return c.RegisteredClaims.Valid()
}

// DecodeAccessToken decodes an access token.
func DecodeAccessToken(token, secret string) (*AccessToken, error) {
	t, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %w", fmt.Errorf("%v", token.Header["alg"]))
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing access token: %w", err)
	}

	if !t.Valid {
		slog.Warn("access token failed validation")
		return nil, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*AccessClaims)
	if !ok {
		slog.Warn("access token has wrong claims type")
		return nil, errors.New("invalid token")
	}

	uid, err := uuid.FromString(claims.Subject)
	if err != nil {
		slog.Warn("access token has invalid subject UUID", "subject", claims.Subject)
		return nil, errors.New("invalid token")
	}

	return &AccessToken{
		Claims: claims,
		UserID: uid,
	}, nil
}
