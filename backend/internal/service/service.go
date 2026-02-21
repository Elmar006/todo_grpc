package service

import (
	"context"

	//log "github.com/Elmar006/todo_grpc/internal/logger"
	"github.com/Elmar006/todo_grpc/internal/model"
	//"github.com/Elmar006/todo_grpc/internal/repository"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Model) (int64, error)
	GetByID(ctx context.Context, id string) (*model.Model, error)
	List(ctx context.Context, filter string) ([]*model.Model, error)
	Update(ctx context.Context, task *model.Model) error
	Delete(ctx context.Context, id string) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*model.Model, error) {
	return nil, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*model.Model, error) {
	return nil, nil
}

func (s *TaskService) ListTask(ctx context.Context, filter string) ([]*model.Model, error) {
	return nil, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *model.Model) error {
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return nil
}
