package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	dsn string
	db  *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Client, error) {
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

	// use the migrations from the "migrations" directory at this level
	if err := database.Migrate(
		"file://internal/service/production/storage/migrations",
		strings.Replace(c.dsn, "postgres://", "pgx5://", 1),
	); err != nil {
		return fmt.Errorf("migrating DB: %w", err)
	}

	return nil
}
