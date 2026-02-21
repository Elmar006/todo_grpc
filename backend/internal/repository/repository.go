package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/Elmar006/todo_grpc/internal/logger"
	model "github.com/Elmar006/todo_grpc/internal/model"
)

type RepositoryDB struct {
	*sql.DB
}

var ErrNotFound = errors.New("Failed: RowsAffected count = 0")

func (r *RepositoryDB) Create(ctx context.Context, task *model.Model) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("title is required")
	}

	query := `INSERT INTO task (title, description) VALUES (?, ?)`
	res, err := r.ExecContext(ctx, query, task.Title, task.Description)
	if err != nil {
		log.L().Errorf("Failed to append data on table 'task'. Err: %v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.L().Errorf("Failed when receiving data. Err: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *RepositoryDB) GetByID(ctx context.Context, id string) (*model.Model, error) {
	model := &model.Model{}
	query := "SELECT id, title, description, completed, created_at, updated_at FROM task WHERE id = ?"

	err := r.QueryRowContext(ctx, query, id).Scan(&model.ID, &model.Title, &model.Description, &model.Completed, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return model, nil
}

func (r *RepositoryDB) List(ctx context.Context, filter string) ([]*model.Model, error) {
	tasks := []*model.Model{}
	query := "SELECT * FROM task WHERE title LIKE ? OR description LIKE ? ORDER BY created_at DESC"

	rows, err := r.QueryContext(ctx, query, "%"+filter+"%", "%"+filter+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := model.Model{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tasks, nil
}

func (r *RepositoryDB) Update(ctx context.Context, task *model.Model) error {
	task.UpdatedAt = time.Now()
	query := "UPDATE task SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?"

	res, err := r.ExecContext(ctx, query, task.Title, task.Description, task.Completed, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *RepositoryDB) DeleteByID(ctx context.Context, id string) error {
	query := "DELETE FROM task WHERE id = ?"

	res, err := r.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}

	return nil
}
