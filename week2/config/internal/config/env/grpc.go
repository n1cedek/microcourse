package env

import (
	"errors"
	"microservices_course/config/internal/config"
	"net"
	"os"
)

var _ config.GPRCConfig = (*grpcConfig)(nil)

const (
	grpcPortEnv = "GRPC_PORT"
	grpcHostEnv = "GRPC_HOST"
)

type grpcConfig struct {
	host string
	port string
}

func NewGrpcConfig() (*grpcConfig, error) {
	port := os.Getenv(grpcPortEnv)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	host := os.Getenv(grpcHostEnv)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	},nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
