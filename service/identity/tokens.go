package identity

import (
	"errors"
	"fmt"
	"log"
	"time"

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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing access token: %v", err)
	}

	if !t.Valid {
		log.Println("token failed validation")
		return nil, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*AccessClaims)
	if !ok {
		log.Println("token has wrong type")
		return nil, errors.New("invalid token")
	}

	uid, err := uuid.FromString(claims.Subject)
	if err != nil {
		log.Println("token has invalid subject UUID")
		return nil, errors.New("invalid token")
	}

	return &AccessToken{
		Claims: claims,
		UserID: uid,
	}, nil
}

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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing refresh token: %v", err)
	}

	if !t.Valid {
		log.Println("token failed validation")
		return nil, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*RefreshClaims)
	if !ok {
		log.Println("token has wrong type")
		return nil, errors.New("invalid token")
	}

	uid, err := uuid.FromString(claims.Subject)
	if err != nil {
		log.Println("token has invalid subject UUID")
		return nil, errors.New("invalid token")
	}

	return &RefreshToken{
		Claims: claims,
		UserID: uid,
	}, nil
}

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
	UUID     uuid.UUID
	Username string
	Password string
}

// GenerateTokens geneates a new access and refresh token for the user.
func (u User) GenerateTokens(secret string) (string, string, error) {
	now := time.Now()
	issuedAt := jwt.NewNumericDate(now)

	accessClaims := &AccessClaims{
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

	refreshClaims := &RefreshClaims{
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
