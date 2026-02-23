package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Elmar006/todo_grpc/internal/config"
	"github.com/Elmar006/todo_grpc/internal/db"
	"github.com/Elmar006/todo_grpc/internal/grpc/handler"
	"github.com/Elmar006/todo_grpc/internal/repository"
	"github.com/Elmar006/todo_grpc/internal/service"
	todo "github.com/Elmar006/todo_grpc/proto/gen/todoService"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn, err := db.Init(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	defer dbConn.Close()

	repo := &repository.RepositoryDB{DB: dbConn}
	taskService := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(taskService)
	grpcServer := grpc.NewServer()
	todo.RegisterTodoServiceServer(grpcServer, taskHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
