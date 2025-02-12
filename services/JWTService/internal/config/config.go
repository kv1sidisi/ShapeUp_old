package config

import (
	"time"
)

type Config struct {
	Env  string     `yaml:"env"`
	JWT  JWTConfig  `yaml:"jwt"`
	GRPC GRPCConfig `yaml:"grpc"`
}

// JWTConfig configuration for JWT.
type JWTConfig struct {
	AccessTokenSecretKey  string `yaml:"access_token_secret_key"`
	RefreshTokenSecretKey string `yaml:"refresh_token_secret_key"`
}

// GRPCConfig configuration for GRPC server.
type GRPCConfig struct {
	Port    int64         `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
