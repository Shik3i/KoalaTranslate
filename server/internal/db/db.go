package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database at dbPath and creates any missing parent directories.
func InitDB(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		dbPath = "koalatranslate.db"
	}

	// Ensure parent directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failure: %w", err)
	}

	return db, nil
}

// RunMigrations runs schema migrations sequentially.
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS locales (
			project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
			code TEXT NOT NULL,
			name TEXT NOT NULL,
			is_source BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (project_id, code)
		);`,

		`CREATE TABLE IF NOT EXISTS translation_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
			key_name TEXT NOT NULL,
			description TEXT,
			max_chars INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (project_id, key_name)
		);`,

		`CREATE TABLE IF NOT EXISTS translation_values (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key_id INTEGER REFERENCES translation_keys(id) ON DELETE CASCADE,
			locale_code TEXT NOT NULL,
			value TEXT NOT NULL,
			is_dirty BOOLEAN DEFAULT 0,
			last_edited_by TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (key_id, locale_code)
		);`,

		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS imports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
			locale_code TEXT NOT NULL,
			source_type TEXT NOT NULL,
			source_url TEXT,
			imported_by TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS exports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
			locale_code TEXT NOT NULL,
			export_type TEXT NOT NULL,
			pr_number INTEGER,
			pr_url TEXT,
			exported_by TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("failed executing migration %d: %w", i, err)
		}
	}

	return nil
}
