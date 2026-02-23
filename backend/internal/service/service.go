package service

import (
	"context"
	"errors"

	log "github.com/Elmar006/todo_grpc/internal/logger"
	"github.com/Elmar006/todo_grpc/internal/model"
	"github.com/Elmar006/todo_grpc/internal/repository"
)

var (
	ErrInvalidData  = errors.New("invalid data")
	ErrTaskNotFound = errors.New("task not found")
)

type TaskRepository interface {
	Create(ctx context.Context, title, description string) (*model.Model, error)
	GetByID(ctx context.Context, id int64) (*model.Model, error)
	List(ctx context.Context, filter string) ([]*model.Model, error)
	Update(ctx context.Context, task *model.Model) error
	Delete(ctx context.Context, id int64) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*model.Model, error) {
	if title == "" {
		return nil, ErrInvalidData
	}

	task, err := s.repo.Create(ctx, title, description)
	if err != nil {
		log.L().Errorf("Failed to create task: %v", err)
		return nil, err
	}

	log.L().Infof("Task created successfully: ID=%d, Title=%q", task.ID, task.Title)
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*model.Model, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.L().Errorf("Failed to get task: %v", err)
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	log.L().Infof("Task retrieved: ID=%d", id)
	return task, nil
}

func (s *TaskService) ListTask(ctx context.Context, filter string) ([]*model.Model, error) {
	tasks, err := s.repo.List(ctx, filter)
	if err != nil {
		log.L().Errorf("Failed to list tasks: %v", err)
		return nil, err
	}

	log.L().Infof("Listed %d tasks with filter %q", len(tasks), filter)
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *model.Model) error {
	if err := s.repo.Update(ctx, task); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.L().Errorf("Task not found for update, ID=%d", task.ID)
			return ErrTaskNotFound
		}
		log.L().Errorf("Failed to update task: %v", err)
		return err
	}

	log.L().Infof("Task updated: ID=%d", task.ID)
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.L().Errorf("Task not found for deletion, ID=%d", id)
			return ErrTaskNotFound
		}
		log.L().Errorf("Failed to delete task: %v", err)
		return err
	}

	log.L().Infof("Task deleted: ID=%d", id)
	return nil
}
