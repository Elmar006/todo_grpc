package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort int
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

	return &Config{
		GRPCPort: port,
	}, nil
}
