package db

import (
	"database/sql"
	"os"

	log "github.com/Elmar006/todo_grpc/internal/logger"
)

func Init(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)
	isNotEx := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.L().Errorf("Failed to 'open' dbFile. Err: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.L().Errorf("Failed to connect(Ping) DB. Err: %v", err)
		return nil, err
	}

	if isNotEx {
		_, err := createTableTask(db)
		if err != nil {
			log.L().Errorf("Failed to create DB. Err: %v", err)
			return nil, err
		}

		log.L().Info("Database Create")
	}

	return db, nil
}

func createTableTask(db *sql.DB) (sql.Result, error) {
	query := `
	CREATE TABLE task (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL, 
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_completed ON task(completed);
	CREATE INDEX idx_created_at ON task(created_at);
`

	table, err := db.Exec(query)
	if err != nil {
		log.L().Errorf("Failed to create table. Err: %v", err)
		return nil, err
	}

	log.L().Info("Create DB")
	return table, nil
}
