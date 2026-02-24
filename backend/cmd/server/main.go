package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Elmar006/todo_grpc/internal/config"
	"github.com/Elmar006/todo_grpc/internal/db"
	"github.com/Elmar006/todo_grpc/internal/grpc/handler"
	"github.com/Elmar006/todo_grpc/internal/logger"
	"github.com/Elmar006/todo_grpc/internal/repository"
	"github.com/Elmar006/todo_grpc/internal/service"
	todo "github.com/Elmar006/todo_grpc/proto/gen/todoService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := logger.L()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := db.Init(cfg.DBPath); err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	defer db.DB.Close()

	repo := &repository.RepositoryDB{DB: db.DB}
	taskService := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(taskService)

	grpcServer := grpc.NewServer()
	todo.RegisterTodoServiceServer(grpcServer, taskHandler)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info("Graceful shutdown initiated")
		grpcServer.GracefulStop()
		cancel()
	}()

	log.Infof("gRPC server listening on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		select {
		case <-ctx.Done():
			log.Info("gRPC server stopped gracefully")
		default:
			log.Fatalf("Failed to serve: %v", err)
		}
	}
}
