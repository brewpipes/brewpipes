package storage

import (
	"context"
	"fmt"
	"log/slog"

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

	if err := database.Migrate(
		"file://service/inventory/storage/migrations",
		database.MigrationDSN(c.dsn, "inventory_schema_migrations"),
	); err != nil {
		return fmt.Errorf("migrating DB: %w", err)
	}

	go func() {
		<-ctx.Done()
		slog.Info("closing inventory service DB pool")
		c.db.Close()
	}()

	return nil
}

func (c *Client) Close() error {
	c.db.Close()
	return nil
}
