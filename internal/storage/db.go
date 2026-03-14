// Package storage manages the local SQLite database for tgcli.
package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Pure-Go SQLite driver, no CGO required.
)

// InitDB opens (or creates) the SQLite database at path and runs schema
// migrations. The returned *sql.DB is safe for concurrent use.
func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Single-writer optimizations for a local CLI tool.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set journal_mode: %w", err)
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set busy_timeout: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

// migrate applies all schema migrations in order.
func migrate(db *sql.DB) error {
	// FTS5 virtual table for full-text search on messages.
	// UNINDEXED columns are stored but not tokenized by the FTS engine.
	const createMessages = `
		CREATE VIRTUAL TABLE IF NOT EXISTS messages USING fts5(
			id UNINDEXED,
			chat_id UNINDEXED,
			sender,
			text,
			timestamp UNINDEXED
		);`

	if _, err := db.Exec(createMessages); err != nil {
		return fmt.Errorf("create messages table: %w", err)
	}
	return nil
}
