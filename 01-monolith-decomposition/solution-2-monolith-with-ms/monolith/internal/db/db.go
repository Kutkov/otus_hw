package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens an SQLite database at the provided path and configures sane defaults.
func Open(databasePath string) (*sql.DB, error) {
	// Add busy timeout to reduce "database is locked" errors on concurrent access
	dsn := fmt.Sprintf("%s?_fk=1&_busy_timeout=5000", databasePath)
	database, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)
	database.SetConnMaxLifetime(0)
	return database, nil
}

// Migrate ensures all required tables exist.
func Migrate(database *sql.DB) error {
	// Users table stores profile and hashed password
	_, err := database.Exec(`
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL,
    birthdate TEXT NOT NULL,
    biography TEXT,
    city TEXT,
    password_hash BLOB NOT NULL,
    created_at TEXT NOT NULL
);
`)
	if err != nil {
		return err
	}

	// Tokens table stores issued bearer tokens
	_, err = database.Exec(`
CREATE TABLE IF NOT EXISTS tokens (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
`)
	if err != nil {
		return err
	}

	return nil
}

// NowISO returns current time in RFC3339 layout suitable for storage in TEXT columns.
func NowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}
