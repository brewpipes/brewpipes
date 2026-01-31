package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, username, password, created_at, updated_at, deleted_at
		FROM identity.user
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.UUID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating users: %w", err)
	}

	return users, nil
}

func (c *Client) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	if err := c.db.QueryRow(ctx, `
		SELECT id, uuid, username, password, created_at, updated_at, deleted_at
		FROM identity.user
		WHERE uuid = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, service.ErrNotFound
		}
		return User{}, fmt.Errorf("getting user: %w", err)
	}

	return user, nil
}

// GetUserByUsername looks up a single user by its username.
func (c *Client) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	if err := c.db.QueryRow(ctx, `
		SELECT id, uuid, username, password, created_at, updated_at, deleted_at
		FROM identity.user
		WHERE username = $1 AND deleted_at IS NULL`,
		username,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, service.ErrNotFound
		}
		return User{}, fmt.Errorf("getting user by username: %w", err)
	}

	return user, nil
}

func (c *Client) CreateUser(ctx context.Context, user User) (User, error) {
	if err := c.db.QueryRow(ctx, `
		INSERT INTO identity.user (username, password)
		VALUES ($1, $2)
		RETURNING id, uuid, username, password, created_at, updated_at, deleted_at`,
		user.Username,
		user.Password,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		return User{}, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}

func (c *Client) UpdateUser(ctx context.Context, user User) (User, error) {
	if err := c.db.QueryRow(ctx, `
		UPDATE identity.user
		SET username = $2,
			password = $3,
			updated_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL
		RETURNING id, uuid, username, password, created_at, updated_at, deleted_at`,
		user.UUID,
		user.Username,
		user.Password,
	).Scan(
		&user.ID,
		&user.UUID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, service.ErrNotFound
		}
		return User{}, fmt.Errorf("updating user: %w", err)
	}

	return user, nil
}

func (c *Client) DeleteUser(ctx context.Context, userUUID uuid.UUID) error {
	cmd, err := c.db.Exec(ctx, `
		UPDATE identity.user
		SET deleted_at = timezone('utc', now()),
			updated_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL`,
		userUUID,
	)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}
