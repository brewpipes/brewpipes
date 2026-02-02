package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brewpipes/brewpipes/internal/jwt"
	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/identity/handler"
	"github.com/brewpipes/brewpipes/service/identity/storage"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthStore struct {
	user                 storage.User
	getUserErr           error
	getUserByUsernameErr error

	refreshToken       storage.RefreshToken
	getRefreshTokenErr error
	consumeErr         error

	createdTokens []storage.RefreshToken
	consumed      bool
}

func (s *fakeAuthStore) GetUser(ctx context.Context, id uuid.UUID) (storage.User, error) {
	if s.getUserErr != nil {
		return storage.User{}, s.getUserErr
	}
	return s.user, nil
}

func (s *fakeAuthStore) GetUserByUsername(ctx context.Context, username string) (storage.User, error) {
	if s.getUserByUsernameErr != nil {
		return storage.User{}, s.getUserByUsernameErr
	}
	return s.user, nil
}

func (s *fakeAuthStore) CreateRefreshToken(ctx context.Context, token storage.RefreshToken) (storage.RefreshToken, error) {
	s.createdTokens = append(s.createdTokens, token)
	return token, nil
}

func (s *fakeAuthStore) GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (storage.RefreshToken, error) {
	if s.getRefreshTokenErr != nil {
		return storage.RefreshToken{}, s.getRefreshTokenErr
	}
	if s.refreshToken.TokenID != tokenID {
		return storage.RefreshToken{}, service.ErrNotFound
	}
	return s.refreshToken, nil
}

func (s *fakeAuthStore) ConsumeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	if s.consumeErr != nil {
		return s.consumeErr
	}
	s.consumed = true
	return nil
}

func TestHandleLoginSuccess(t *testing.T) {
	secret := "test-secret"
	userUUID := uuid.Must(uuid.NewV4())

	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	store := &fakeAuthStore{
		user: storage.User{
			Username: "brewer",
			Password: string(hashed),
		},
	}
	store.user.UUID = userUUID

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"brewer","password":"password"}`))

	handler.HandleLogin(store, store, secret).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatalf("expected tokens in response")
	}
	if len(store.createdTokens) != 1 {
		t.Fatalf("expected refresh token to be stored")
	}

	decoded, err := jwt.DecodeRefreshToken(resp.RefreshToken, secret)
	if err != nil {
		t.Fatalf("decode refresh token: %v", err)
	}
	refreshID, err := uuid.FromString(decoded.Claims.ID)
	if err != nil {
		t.Fatalf("parse refresh id: %v", err)
	}
	if store.createdTokens[0].TokenID != refreshID {
		t.Fatalf("stored refresh token id mismatch")
	}
}

func TestHandleRefreshSuccess(t *testing.T) {
	secret := "test-secret"
	userUUID := uuid.Must(uuid.NewV4())

	user := storage.User{
		Username: "brewer",
	}
	user.UUID = userUUID

	_, refresh, err := user.GenerateTokens(secret)
	if err != nil {
		t.Fatalf("generate tokens: %v", err)
	}
	decoded, err := jwt.DecodeRefreshToken(refresh, secret)
	if err != nil {
		t.Fatalf("decode refresh token: %v", err)
	}
	refreshID, err := uuid.FromString(decoded.Claims.ID)
	if err != nil {
		t.Fatalf("parse refresh id: %v", err)
	}

	store := &fakeAuthStore{
		user: user,
		refreshToken: storage.RefreshToken{
			TokenID:   refreshID,
			UserUUID:  userUUID,
			ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
		},
	}

	rec := httptest.NewRecorder()
	reqBody := `{"refresh_token":"` + refresh + `"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(reqBody))

	handler.HandleRefresh(store, store, secret).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !store.consumed {
		t.Fatalf("expected refresh token to be consumed")
	}
	if len(store.createdTokens) != 1 {
		t.Fatalf("expected new refresh token to be stored")
	}
}

func TestHandleLogoutSuccess(t *testing.T) {
	secret := "test-secret"
	userUUID := uuid.Must(uuid.NewV4())

	user := storage.User{Username: "brewer"}
	user.UUID = userUUID

	_, refresh, err := user.GenerateTokens(secret)
	if err != nil {
		t.Fatalf("generate tokens: %v", err)
	}

	store := &fakeAuthStore{}

	rec := httptest.NewRecorder()
	reqBody := `{"refresh_token":"` + refresh + `"}`
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(reqBody))

	handler.HandleLogout(store, secret).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if !store.consumed {
		t.Fatalf("expected refresh token to be consumed")
	}
}
