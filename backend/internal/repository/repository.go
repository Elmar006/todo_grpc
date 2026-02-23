package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/Elmar006/todo_grpc/internal/logger"
	"github.com/Elmar006/todo_grpc/internal/model"
)

type RepositoryDB struct {
	*sql.DB
}

var (
	ErrNotFound    = errors.New("failed: rows affected count = 0")
	ErrInvalidData = errors.New("invalid data")
)

func (r *RepositoryDB) Create(ctx context.Context, title, description string) (*model.Model, error) {
	if title == "" {
		log.L().Errorf("Title is required. Err: %v", ErrInvalidData)
		return nil, ErrInvalidData
	}

	query := `INSERT INTO task (title, description) VALUES (?, ?)`
	res, err := r.ExecContext(ctx, query, title, description)
	if err != nil {
		log.L().Errorf("Failed to insert data: %v", err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.L().Errorf("Failed to get last insert ID: %v", err)
		return nil, err
	}

	task := &model.Model{
		ID:          id,
		Title:       title,
		Description: &description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	log.L().Info("Task created successfully")
	return task, nil
}

func (r *RepositoryDB) GetByID(ctx context.Context, id int64) (*model.Model, error) {
	task := &model.Model{}
	query := `SELECT id, title, description, completed, created_at, updated_at FROM task WHERE id = ?`

	err := r.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Description,
		&task.Completed, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.L().Info("Task not found")
			return nil, nil
		}
		log.L().Errorf("Failed to get task by ID: %v", err)
		return nil, err
	}

	log.L().Infof("Task retrieved: ID=%d", id)
	return task, nil
}

func (r *RepositoryDB) List(ctx context.Context, filter string) ([]*model.Model, error) {
	tasks := []*model.Model{}
	query := `SELECT id, title, description, completed, created_at, updated_at 
	          FROM task 
	          WHERE title LIKE ? OR description LIKE ? 
	          ORDER BY created_at DESC`

	rows, err := r.QueryContext(ctx, query, "%"+filter+"%", "%"+filter+"%")
	if err != nil {
		log.L().Errorf("Failed to list tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := model.Model{}
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description,
			&task.Completed, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			log.L().Errorf("Failed to scan task: %v", err)
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		log.L().Errorf("Error during rows iteration: %v", err)
		return nil, err
	}

	log.L().Infof("Listed %d tasks with filter: %q", len(tasks), filter)
	return tasks, nil
}

func (r *RepositoryDB) Update(ctx context.Context, task *model.Model) error {
	task.UpdatedAt = time.Now()
	query := `UPDATE task SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?`

	res, err := r.ExecContext(ctx, query,
		task.Title, task.Description, task.Completed, task.UpdatedAt, task.ID,
	)
	if err != nil {
		log.L().Errorf("Failed to update task: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		log.L().Errorf("Task not found for update, ID=%d", task.ID)
		return ErrNotFound
	}

	log.L().Infof("Task updated: ID=%d", task.ID)
	return nil
}

func (r *RepositoryDB) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM task WHERE id = ?`

	res, err := r.ExecContext(ctx, query, id)
	if err != nil {
		log.L().Errorf("Failed to delete task: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		log.L().Errorf("Task not found for deletion, ID=%d", id)
		return ErrNotFound
	}

	log.L().Infof("Task deleted: ID=%d", id)
	return nil
}
