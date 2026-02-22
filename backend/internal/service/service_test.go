package service

import (
	"context"
	"errors"

	//"errors"
	"testing"
	//"time"

	"github.com/Elmar006/todo_grpc/internal/model"
	"github.com/Elmar006/todo_grpc/internal/repository"
	//"github.com/Elmar006/todo_grpc/internal/repository"
)

type fakeRepo struct {
	createFunc  func(ctx context.Context, title, description string) (int64, error)
	getByIdFunc func(ctx context.Context, id string) (*model.Model, error)
	listFunc    func(ctx context.Context, filter string) ([]*model.Model, error)
	updateFunc  func(ctx context.Context, task *model.Model) error
	deleteFunc  func(ctx context.Context, id string) error
}

func (f *fakeRepo) Create(ctx context.Context, title, description string) (int64, error) {
	return f.createFunc(ctx, title, description)
}

func (f *fakeRepo) GetByID(ctx context.Context, id string) (*model.Model, error) {
	return f.getByIdFunc(ctx, id)
}

func (f *fakeRepo) List(ctx context.Context, filter string) ([]*model.Model, error) {
	return f.listFunc(ctx, filter)
}

func (f *fakeRepo) Update(ctx context.Context, task *model.Model) error {
	return f.updateFunc(ctx, task)
}

func (f *fakeRepo) Delete(ctx context.Context, id string) error {
	return f.deleteFunc(ctx, id)
}

func TestCreateTaskCorrected(t *testing.T) {
	taskCheck := &fakeRepo{
		createFunc: func(ctx context.Context, title, description string) (int64, error) {
			if title != "Test Task" {
				t.Errorf("Expected title 'Test Task', got %q", title)
			}
			if description != "Test Desc" {
				t.Errorf("Expected description 'Test Desc', got %q", description)
			}
			var id int64 = 123
			return id, nil
		},
	}

	service := NewTaskService(taskCheck)
	id, err := service.CreateTask(context.Background(), "Test Task", "Test Desc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 123 {
		t.Errorf("expected ID 123, got %d", id)
	}
}

func TestCreateTaskInvalid(t *testing.T) {
	taskCheck := &fakeRepo{}
	service := NewTaskService(taskCheck)

	_, err := service.CreateTask(context.Background(), "", "desc")
	if err != nil {
		if !errors.Is(err, ErrInvalidData) {
			t.Errorf("expected ErrInvalidData, got %v", err)
		}
	}
}

func TestGetTaskCorrected(t *testing.T) {
	testTaskRequest := &model.Model{
		ID:    "123",
		Title: "Test Title",
	}

	taskCheck := &fakeRepo{
		getByIdFunc: func(ctx context.Context, id string) (*model.Model, error) {
			if id != "123" {
				t.Errorf("Expected id: 123, got %s", id)
			}

			return testTaskRequest, nil
		},
	}

	service := NewTaskService(taskCheck)
	task, err := service.GetTask(context.Background(), "123")
	if err != nil {
		t.Fatal(err)
	}
	if task != testTaskRequest {
		t.Error("Expected the returned task to match")
	}
}

func TestGetTaskNotFound(t *testing.T) {
	taskCheck := &fakeRepo{
		getByIdFunc: func(ctx context.Context, id string) (*model.Model, error) {
			return nil, nil
		},
	}

	service := NewTaskService(taskCheck)
	task, err := service.GetTask(context.Background(), "123")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
	if task != nil {
		t.Error("Expected nil task")
	}
}

func TestGetTaskRepoError(t *testing.T) {
	repoErr := errors.New("Data Base Error")
	taskCheck := &fakeRepo{
		getByIdFunc: func(ctx context.Context, id string) (*model.Model, error) {
			return nil, repoErr
		},
	}

	service := NewTaskService(taskCheck)
	_, err := service.GetTask(context.Background(), "123")
	if err != repoErr {
		t.Errorf("expected repo error %v, got %v", repoErr, err)
	}
}
func TestListTask(t *testing.T) {
	tasksTest := []*model.Model{
		{ID: "1"},
		{ID: "2"},
	}
	taskCheck := &fakeRepo{
		listFunc: func(ctx context.Context, filter string) ([]*model.Model, error) {
			if filter != "test" {
				t.Errorf("Expect filter 'test', got %s", filter)
			}

			return tasksTest, nil
		},
	}

	service := NewTaskService(taskCheck)
	tasks, err := service.ListTask(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestUpdateTaskSuccess(t *testing.T) {
	called := false
	taskCheck := &fakeRepo{
		updateFunc: func(ctx context.Context, task *model.Model) error {
			called = true
			if task.ID != "123" {
				t.Errorf("expected id 123, got %s", task.ID)
			}

			return nil
		},
	}

	service := NewTaskService(taskCheck)
	if err := service.UpdateTask(context.Background(), &model.Model{ID: "123"}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("expected Update to be called")
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	taskCheck := &fakeRepo{
		updateFunc: func(ctx context.Context, task *model.Model) error {
			return repository.ErrNotFound
		},
	}

	service := NewTaskService(taskCheck)
	if err := service.UpdateTask(context.Background(), &model.Model{ID: "123"}); err != nil {
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	}
}

func TestDeleteTaskSuccess(t *testing.T) {
	called := false
	taskCheck := &fakeRepo{
		deleteFunc: func(ctx context.Context, id string) error {
			called = true
			if id != "123" {
				t.Errorf("expected id 123, got %s", id)
			}

			return nil
		},
	}

	service := NewTaskService(taskCheck)
	if err := service.DeleteTask(context.Background(), "123"); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Expected Delete to be called")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	taskCheck := &fakeRepo{
		deleteFunc: func(ctx context.Context, id string) error {
			return repository.ErrNotFound
		},
	}

	service := NewTaskService(taskCheck)
	if err := service.DeleteTask(context.Background(), "123"); err != nil {
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	}
}
