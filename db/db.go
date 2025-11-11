package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func Open(appDir string) (*sql.DB, error) {
	dataDir := filepath.Join(appDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "snip.sqlite")
	db, err := sql.Open("sqlite3", dbPath+"?journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	return db, nil
}
