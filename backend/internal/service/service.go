package service

import (
	"context"
	"errors"

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
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*model.Model, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *TaskService) ListTask(ctx context.Context, filter string) ([]*model.Model, error) {
	tasks, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *model.Model) error {
	if err := s.repo.Update(ctx, task); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}
