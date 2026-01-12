package database

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(from, to string) error {
	// m, err := migrate.New(from, to)
	m, err := migrate.New("file:///migrations", "pgx5://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=enable")
	if err != nil {
		return fmt.Errorf("creating migration instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Error("error closing DB migration source", "error", srcErr)
		}
		if dbErr != nil {
			slog.Error("error closing DB migration database", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		return fmt.Errorf("applying DB migrations: %w", err)
	}

	return nil
}
