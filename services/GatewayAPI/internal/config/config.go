package config

import (
	"time"
)

type Config struct {
	Env              string `yaml:"env" env-default:"local"`
	HTTPServer       `yaml:"http_server"`
	GRPCClientConfig `yaml:"grpc_client"`
}

// HTTPServer configuration for HTTP server.
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// GRPCClientConfig configuration for GRPC clients.
type GRPCClientConfig struct {
	UserCreationServiceAddress   string `yaml:"user_creation_service_address" env-required:"true"`
	AuthenticationServiceAddress string `yaml:"authentication_service_address" env-required:"true"`
}
