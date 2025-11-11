package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*
var migrationFiles embed.FS

func RunMigrations(dbPath string) error {

	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("embed source : %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, "sqlite3://"+dbPath)
	if err != nil {
		return fmt.Errorf("create migrator : %w", err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			fmt.Println("Error closing migrator", err)
		}
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up : %w", err)
	}

	return nil
}
