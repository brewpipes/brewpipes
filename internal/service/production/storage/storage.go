package storage

import (
	"context"
	"fmt"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	dsn string
	db  *pgxpool.Pool
}

func NewClient(ctx context.Context, dsn string) (*Client, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("creating DB connection pool: %w", err)
	}

	return &Client{
		dsn: dsn,
		db:  pool,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	if err := c.db.Ping(ctx); err != nil {
		return fmt.Errorf("pinging Postgres: %w", err)
	}

	if err := database.Migrate("file://./db/migrations", c.dsn); err != nil {
		return fmt.Errorf("migrating DB: %w", err)
	}

	return nil
}
