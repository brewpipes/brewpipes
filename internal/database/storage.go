package database

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BaseClient provides shared storage lifecycle (connect, ping, migrate, close)
// for service-specific storage clients.
type BaseClient struct {
	dsn            string
	db             *pgxpool.Pool
	migrations     embed.FS
	migrationTable string
	serviceName    string
}

// NewBaseClient creates a BaseClient with the given DSN, embedded migrations,
// migration table name, and service name (used for log messages).
func NewBaseClient(dsn string, migrations embed.FS, migrationTable, serviceName string) *BaseClient {
	return &BaseClient{
		dsn:            dsn,
		migrations:     migrations,
		migrationTable: migrationTable,
		serviceName:    serviceName,
	}
}

// Start connects to Postgres, pings the database, runs migrations, and
// registers a goroutine to close the pool when ctx is cancelled.
func (c *BaseClient) Start(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, c.dsn)
	if err != nil {
		return fmt.Errorf("creating DB connection pool: %w", err)
	}
	c.db = pool

	if err := c.db.Ping(ctx); err != nil {
		return fmt.Errorf("pinging Postgres: %w", err)
	}

	if err := Migrate(
		c.migrations,
		"migrations",
		MigrationDSN(c.dsn, c.migrationTable),
	); err != nil {
		return fmt.Errorf("migrating DB: %w", err)
	}

	go func() {
		<-ctx.Done()
		slog.Info("closing " + c.serviceName + " service DB pool")
		c.db.Close()
	}()

	return nil
}

// Close closes the underlying connection pool.
func (c *BaseClient) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

// DB returns the underlying pgxpool.Pool for query execution.
func (c *BaseClient) DB() *pgxpool.Pool {
	return c.db
}
