package production

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func migrateDB() error {
	db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	if err != nil {
		return fmt.Errorf("opening DB connection for migration: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("creating migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("applying DB migrations: %w", err)
	}

	return nil
}
