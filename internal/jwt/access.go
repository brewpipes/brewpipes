// Package jwt contains types and functions for manipulating JSON Web Tokens.
package jwt

import (
	"errors"
	"fmt"
	"log"

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
	if c.Role == "" && !(c.Role == "admin" || c.Role == "user") {
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
