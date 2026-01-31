package database

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(from, to string) error {
	m, err := migrate.New(from, to)
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

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("applying DB migrations: %w", err)
	}

	return nil
}

func MigrationDSN(dsn, table string) string {
	migrationDSN := strings.Replace(dsn, "postgres://", "pgx5://", 1)
	if table == "" {
		return migrationDSN
	}

	separator := "?"
	if strings.Contains(migrationDSN, "?") {
		separator = "&"
	}

	return migrationDSN + separator + "x-migrations-table=" + table
}
