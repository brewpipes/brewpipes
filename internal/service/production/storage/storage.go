package storage

import (
	"context"
	"fmt"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	DB *pgxpool.Pool
}

func (c *Client) Ping(ctx context.Context) error {
	if err := c.DB.Ping(ctx); err != nil {
		return fmt.Errorf("pinging Postgres: %w", err)
	}
	return nil
}

func (c *Client) Migrate() error {
	if err := database.Migrate(
		"file://./db/migrations",
		"pgx5://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=enable",
	); err != nil {
		return fmt.Errorf("migrating DB: %w", err)
	}

	return nil
}
