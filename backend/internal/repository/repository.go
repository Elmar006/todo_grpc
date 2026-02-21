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

var (
	ErrNotFound    = errors.New("Failed: RowsAffected count = 0")
	ErrInvalidData = errors.New("Invalid data")
)

func (r *RepositoryDB) Create(ctx context.Context, task *model.Model) (int64, error) {
	if task.Title == "" {
		log.L().Errorf("Title is required. Err: %v", ErrInvalidData)
		return 0, ErrInvalidData
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

	log.L().Info("The data was successfully added to the DB")
	return id, nil
}

func (r *RepositoryDB) GetByID(ctx context.Context, id string) (*model.Model, error) {
	model := &model.Model{}
	query := "SELECT id, title, description, completed, created_at, updated_at FROM task WHERE id = ?"

	err := r.QueryRowContext(ctx, query, id).Scan(&model.ID, &model.Title, &model.Description, &model.Completed, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.L().Info("Task not found")
			return nil, nil
		}

		log.L().Errorf("Failed while retrieving data by ID using the GetByID method. Err: %v", err)
		return nil, err
	}

	log.L().Info("ID data was successfully received")
	return model, nil
}

func (r *RepositoryDB) List(ctx context.Context, filter string) ([]*model.Model, error) {
	tasks := []*model.Model{}
	query := "SELECT * FROM task WHERE title LIKE ? OR description LIKE ? ORDER BY created_at DESC"

	rows, err := r.QueryContext(ctx, query, "%"+filter+"%", "%"+filter+"%")
	if err != nil {
		log.L().Errorf("Failed while retrieving data by search filter using the List method. Err: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := model.Model{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			log.L().Errorf("Failed when scanning data in the List method. Err: %v", err)
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		log.L().Errorf("Error during rows iteration in List: %v", err)
		return nil, err
	}

	log.L().Info("Filter data has been successfully received.")
	return tasks, nil
}

func (r *RepositoryDB) Update(ctx context.Context, task *model.Model) error {
	task.UpdatedAt = time.Now()
	query := "UPDATE task SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?"

	res, err := r.ExecContext(ctx, query, task.Title, task.Description, task.Completed, task.UpdatedAt, task.ID)
	if err != nil {
		log.L().Errorf("Error updating data in the table. Err: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		log.L().Errorf("The data in the column has not been updated. Err: %v", ErrNotFound)
		return ErrNotFound
	}

	log.L().Info("The data in the column has been updated.")
	return nil
}

func (r *RepositoryDB) DeleteByID(ctx context.Context, id string) error {
	query := "DELETE FROM task WHERE id = ?"

	res, err := r.ExecContext(ctx, query, id)
	if err != nil {
		log.L().Errorf("Failed when deleting data. Err: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		log.L().Errorf("The data in the column has not been delete. Err: %v", ErrNotFound)
		return ErrNotFound
	}

	log.L().Info("The data in the column has been successfully deleted")
	return nil
}
