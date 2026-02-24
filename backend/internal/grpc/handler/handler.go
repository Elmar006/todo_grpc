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

const requestTimeout = 5 * time.Second

type TaskHandler struct {
	taskService *service.TaskService
	todo.UnimplementedTodoServiceServer
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *todo.CreateTaskRequest) (*todo.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	log.L().Infof("CreateTask request: title=%q", req.GetTitle())

	task, err := h.taskService.CreateTask(ctx, req.GetTitle(), req.GetDescription())
	if err != nil {
		if errors.Is(err, service.ErrInvalidData) {
			log.L().Warnf("CreateTask failed: invalid data - title is empty")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			log.L().Errorf("CreateTask timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		log.L().Errorf("CreateTask failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.L().Infof("CreateTask success: id=%d", task.ID)
	return &todo.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: *task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdateAt:    task.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *todo.GetTaskRequest) (*todo.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	log.L().Infof("GetTask request: id=%d", req.GetId())

	taskModel, err := h.taskService.GetTask(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			log.L().Warnf("GetTask not found: id=%d", req.GetId())
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			log.L().Errorf("GetTask timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		log.L().Errorf("GetTask failed: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.L().Infof("GetTask success: id=%d", req.GetId())
	return convertStruct(taskModel), nil
}

func (h *TaskHandler) ListTask(ctx context.Context, req *todo.ListTasksRequest) (*todo.ListTasksResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	log.L().Info("ListTask request")

	taskModel, err := h.taskService.ListTask(ctx, "")
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.L().Errorf("ListTask timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		log.L().Errorf("ListTask failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	protoTasks := make([]*todo.Task, len(taskModel))
	for i, v := range taskModel {
		protoTasks[i] = convertStruct(v)
	}

	log.L().Infof("ListTask success: count=%d", len(protoTasks))
	return &todo.ListTasksResponse{Tasks: protoTasks}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *todo.UpdateTaskRequest) (*todo.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	log.L().Infof("UpdateTask request: id=%d", req.GetId())

	taskModel, err := h.taskService.GetTask(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			log.L().Warnf("UpdateTask not found: id=%d", req.GetId())
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			log.L().Errorf("UpdateTask timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		log.L().Errorf("UpdateTask failed: %v", err)
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
			log.L().Warnf("UpdateTask not found: id=%d", req.GetId())
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			log.L().Errorf("UpdateTask timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		log.L().Errorf("UpdateTask failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.L().Infof("UpdateTask success: id=%d", req.GetId())
	return convertStruct(taskModel), nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *todo.DeleteTaskRequest) (*todo.DeleteTaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	log.L().Infof("DeleteTask request: id=%d", req.GetId())

	if err := h.taskService.DeleteTask(ctx, req.GetId()); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			log.L().Warnf("DeleteTask not found: id=%d", req.GetId())
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			log.L().Errorf("DeleteTask timeout exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "request timeout")
		}
		log.L().Errorf("DeleteTask failed: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.L().Infof("DeleteTask success: id=%d", req.GetId())
	return &todo.DeleteTaskResponse{}, nil
}

func convertStruct(m *model.Model) *todo.Task {
	desc := ""
	if m.Description != nil {
		desc = *m.Description
	}

	return &todo.Task{
		Id:          m.ID,
		Title:       m.Title,
		Description: desc,
		Completed:   m.Completed,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdateAt:    m.UpdatedAt.Format(time.RFC3339),
	}
}
