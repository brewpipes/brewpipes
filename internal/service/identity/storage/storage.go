package storage

import (
	"context"
	"fmt"

	"github.com/brewpipes/brewpipesproto/internal/service"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	db *pgxpool.Pool
}

func NewClient(db *pgxpool.Pool) *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := c.db.Query(ctx, "SELECT id, username, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("querying database: %w", err)
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UUID, &user.Username, &user.Password); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (c *Client) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	rows, err := c.db.Query(ctx, "SELECT id, username, password FROM users WHERE id = $1", id)
	if err != nil {
		return User{}, fmt.Errorf("querying database: %w", err)
	}

	var user User
	if rows.Next() {
		if err := rows.Scan(&user.UUID, &user.Username, &user.Password); err != nil {
			return User{}, fmt.Errorf("scanning row: %w", err)
		}
		return user, nil
	} else {
		return User{}, service.ErrNotFound
	}
}

// GetUserByUsername looks up a single user by its username.
func (c *Client) GetUserByUsername(ctx context.Context, username string) (User, error) {
	rows, err := c.db.Query(ctx, "SELECT id, username, password FROM users WHERE username = $1", username)
	if err != nil {
		return User{}, fmt.Errorf("querying database: %w", err)
	}

	var user User
	if rows.Next() {
		if err := rows.Scan(&user.UUID, &user.Username, &user.Password); err != nil {
			return User{}, fmt.Errorf("scanning row: %w", err)
		}
		return user, nil
	} else {
		return User{}, service.ErrNotFound
	}
}
