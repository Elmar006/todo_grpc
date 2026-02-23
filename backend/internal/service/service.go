package service

import (
	"context"
	"errors"
	"time"

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

	task := &model.Model{
		Title:       title,
		Description: &description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	task, err := s.repo.Create(ctx, title, description)
	if err != nil {
		log.L().Error(err)
		return nil, err
	}

	log.L().Infof("Task created successfully: ID=%d, Title=%q", task.ID, task.Title)
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*model.Model, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.L().Error(err)
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
		log.L().Error(err)
		return nil, err
	}

	log.L().Infof("Tasks listed with filter %q, count=%d", filter, len(tasks))
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *model.Model) error {
	if err := s.repo.Update(ctx, task); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.L().Error(ErrTaskNotFound)
			return ErrTaskNotFound
		}
		log.L().Error(err)
		return err
	}

	log.L().Infof("Task updated successfully: ID=%d", task.ID)
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.L().Error(ErrTaskNotFound)
			return ErrTaskNotFound
		}
		log.L().Error(err)
		return err
	}

	log.L().Infof("Task deleted successfully: ID=%d", id)
	return nil
}
