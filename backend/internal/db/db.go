package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbFile string) error {
	if dbFile == "" {
		dbFile = "./data/todo.db"
	}

	sqlDB, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS task (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_completed ON task(completed);
	CREATE INDEX IF NOT EXISTS idx_created_at ON task(created_at);
	`

	if _, err := sqlDB.Exec(query); err != nil {
		sqlDB.Close()
		return err
	}

	DB = sqlDB
	return nil
}
