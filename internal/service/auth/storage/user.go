package storage

import (
	"time"

	"github.com/brewpipes/brewpipesproto/internal/entity"
	bpjwt "github.com/brewpipes/brewpipesproto/internal/jwt"
	"github.com/golang-jwt/jwt/v4"
)

const (
	Day = 24 * time.Hour
	// DefaultAccessTokenTTL is the default TTL for issued access tokens.
	DefaultAccessTokenTTL = 7 * Day

	// DefaultRefreshTokenTTL is the default TTL for issued refresh tokens.
	DefaultRefreshTokenTTL = 30 * Day

	// Issuer is the JWT issuer.
	Issuer = "brewpipes"
)

// A User defines basic information about a user.
type User struct {
	entity.Identifiers

	Username string
	Password string

	entity.Timestamps
}

// GenerateTokens geneates a new access and refresh token for the user.
func (u User) GenerateTokens(secret string) (string, string, error) {
	now := time.Now()
	issuedAt := jwt.NewNumericDate(now)

	accessClaims := &bpjwt.AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   u.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(DefaultAccessTokenTTL)),
			NotBefore: issuedAt,
			IssuedAt:  issuedAt,
		},
		Role: "user",
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &bpjwt.RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   u.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(DefaultRefreshTokenTTL)),
			NotBefore: issuedAt,
			IssuedAt:  issuedAt,
		},
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
