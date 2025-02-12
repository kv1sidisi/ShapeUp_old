package config

import (
	"fmt"
	"time"
)

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	GRPC       GRPCConfig       `yaml:"grpc"`
	Storage    StorageConfig    `yaml:"storage"`
	GRPCClient GRPCClientConfig `yaml:"grpc_client"`
}

// GRPCConfig configuration for GRPC server.
type GRPCConfig struct {
	Port    int64         `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// GRPCClientConfig configuration for GRPC clients.
type GRPCClientConfig struct {
	SendingServiceAddress string `yaml:"sending_service_address" env-required:"true"`
	JWTServiceAddress     string `yaml:"jwt_service_address" env-required:"true"`
}

// StorageConfig configuration for database.
type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (sc *StorageConfig) GetDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
}
