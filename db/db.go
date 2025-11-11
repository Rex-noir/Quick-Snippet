package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // or mattn/go-sqlite3 depending on your choice
)

func GetDBPath(appDir string) string {
	dataDir := filepath.Join(appDir, "data")
	_ = os.MkdirAll(dataDir, 0755) // ensure the directory exists
	return filepath.Join(dataDir, "snip.db")
}

// Open opens the SQLite database, creating the data directory if needed.
func Open(appDir string) (*sql.DB, error) {
	dbPath := GetDBPath(appDir)
	db, err := sql.Open("sqlite", dbPath+"?journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// optional: ensure connection works
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

// Close safely closes the database
func Close(db *sql.DB) {
	if db != nil {
		_ = db.Close()
	}
}
