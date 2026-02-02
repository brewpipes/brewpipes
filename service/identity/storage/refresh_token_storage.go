package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

// RefreshToken represents a refresh token record persisted for rotation and revocation.
type RefreshToken struct {
	ID        int64
	TokenID   uuid.UUID
	UserUUID  uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (c *Client) CreateRefreshToken(ctx context.Context, token RefreshToken) (RefreshToken, error) {
	if err := c.db.QueryRow(ctx, `
		INSERT INTO identity.refresh_token (token_id, user_uuid, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, token_id, user_uuid, expires_at, revoked_at, created_at`,
		token.TokenID,
		token.UserUUID,
		token.ExpiresAt,
	).Scan(
		&token.ID,
		&token.TokenID,
		&token.UserUUID,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	); err != nil {
		return RefreshToken{}, fmt.Errorf("creating refresh token: %w", err)
	}

	return token, nil
}

func (c *Client) GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (RefreshToken, error) {
	var token RefreshToken
	if err := c.db.QueryRow(ctx, `
		SELECT id, token_id, user_uuid, expires_at, revoked_at, created_at
		FROM identity.refresh_token
		WHERE token_id = $1`,
		tokenID,
	).Scan(
		&token.ID,
		&token.TokenID,
		&token.UserUUID,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshToken{}, service.ErrNotFound
		}
		return RefreshToken{}, fmt.Errorf("getting refresh token: %w", err)
	}

	return token, nil
}

func (c *Client) ConsumeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	cmd, err := c.db.Exec(ctx, `
		UPDATE identity.refresh_token
		SET revoked_at = timezone('utc', now())
		WHERE token_id = $1 AND revoked_at IS NULL`,
		tokenID,
	)
	if err != nil {
		return fmt.Errorf("consuming refresh token: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (c *Client) RevokeRefreshTokensForUser(ctx context.Context, userUUID uuid.UUID) error {
	_, err := c.db.Exec(ctx, `
		UPDATE identity.refresh_token
		SET revoked_at = timezone('utc', now())
		WHERE user_uuid = $1 AND revoked_at IS NULL`,
		userUUID,
	)
	if err != nil {
		return fmt.Errorf("revoking refresh tokens: %w", err)
	}

	return nil
}
