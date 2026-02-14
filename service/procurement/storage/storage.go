package storage

import (
	"embed"

	"github.com/brewpipes/brewpipes/internal/database"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Client is the procurement service storage client.
type Client struct {
	*database.BaseClient
}

// New creates a new procurement storage Client.
func New(dsn string) *Client {
	return &Client{
		BaseClient: database.NewBaseClient(dsn, migrations, "procurement_schema_migrations", "procurement"),
	}
}
