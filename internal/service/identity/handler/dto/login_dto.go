package dto

import "errors"

// LoginRequest is the request body [POST /login].
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates the request.
func (r LoginRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	} else if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// RefreshRequest is the request body [POST /refresh].
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Validate validates the request.
func (r RefreshRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}

	return nil
}

// TokenPair is the DTO for a pair of access tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
