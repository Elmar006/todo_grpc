package handler

import (
	"context"
	"errors"
	"time"

	log "github.com/Elmar006/todo_grpc/internal/logger"
	"github.com/Elmar006/todo_grpc/internal/model"
	"github.com/Elmar006/todo_grpc/internal/service"
	todo "github.com/Elmar006/todo_grpc/proto/gen/todoService"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	taskService *service.TaskService
	todo.UnimplementedTodoServiceServer
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTaask(ctx context.Context, req *todo.CreateTaskRequest) (*todo.Task, error) {
	task, err := h.taskService.CreateTask(ctx, req.GetTitle(), req.GetDescription())
	if err != nil {
		if errors.Is(err, service.ErrInvalidData) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		log.L().Errorf("CreateTaask failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	id := task.ID

	return &todo.Task{
		Id:          id,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Completed:   false,
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *todo.GetTaskRequest) (*todo.Task, error) {
	taskModel, err := h.taskService.GetTask(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.L().Errorf("GetTask  failed: %v", err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return convertStruct(taskModel), nil
}

func (h *TaskHandler) ListTask(ctx context.Context, req *todo.ListTasksRequest) (*todo.ListTasksResponse, error) {
	taskModel, err := h.taskService.ListTask(ctx, "")
	if err != nil {
		log.L().Errorf("ListTasks failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	protoTasks := make([]*todo.Task, len(taskModel))
	for i, v := range taskModel {
		protoTasks[i] = convertStruct(v)
	}

	return &todo.ListTasksResponse{Tasks: protoTasks}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *todo.UpdateTaskRequest) (*todo.Task, error) {
	taskModel, err := h.taskService.GetTask(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	if req.Title != nil {
		taskModel.Title = req.GetTitle()
	}
	if req.Description != nil {
		desc := req.GetDescription()
		taskModel.Description = &desc
	}
	if req.Completed != nil {
		taskModel.Completed = req.GetCompleted()
	}
	if err := h.taskService.UpdateTask(ctx, taskModel); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.L().Errorf("Update failed: %v", err)
		return nil, status.Error(codes.Internal, "internal err")
	}

	return convertStruct(taskModel), nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *todo.DeleteTaskRequest) (*todo.DeleteTaskResponse, error) {
	if err := h.taskService.DeleteTask(ctx, req.GetId()); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.L().Errorf("DeleteTask failed: %v", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &todo.DeleteTaskResponse{}, nil
}

func convertStruct(m *model.Model) *todo.Task {
	return &todo.Task{
		Id:          m.ID,
		Title:       m.Title,
		Description: *m.Description,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdateAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}
