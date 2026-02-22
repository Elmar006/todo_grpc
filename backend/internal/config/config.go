package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort int
	DBPath   string
}

func Load() (*Config, error) {
	portStr := os.Getenv("GRPC_PORT")
	port := 50051
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "task.db"
	}

	return &Config{GRPCPort: port, DBPath: dbPath}, nil
}
