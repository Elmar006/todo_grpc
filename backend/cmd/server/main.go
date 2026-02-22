package main

import (
	//"context"
	"fmt"
	"net"
	"os"

	"github.com/Elmar006/todo_grpc/internal/config"
	"github.com/Elmar006/todo_grpc/internal/db"
	log "github.com/Elmar006/todo_grpc/internal/logger"

	//"github.com/Elmar006/todo_grpc/internal/repository"
	//"github.com/Elmar006/todo_grpc/internal/service"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.L().Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	dbCon, err := db.Init(cfg.DBPath)
	if err != nil {
		log.Log.Errorf("Failed to init db: %v", err)
		os.Exit(1)
	}
	defer dbCon.Close()

	//repo := &repository.RepositoryDB{DB: dbCon}
	//s := service.NewTaskService(repo)
	// конект к грпс хэнндлерам с аргументом (s)
	port := fmt.Sprintf(":%d", cfg.GRPCPort)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.L().Errorf("Failed to listen serve on port: %s. Err: %v\n", port, err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	log.L().Infof("Listening gRPC server on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.L().Errorf("Failed to server star. Err: %v", err)
		os.Exit(1)
	}
}
