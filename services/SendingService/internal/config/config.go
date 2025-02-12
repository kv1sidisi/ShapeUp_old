package config

import (
	"time"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	GRPC GRPCConfig `yaml:"grpc"`
	SMTP SMTPConfig `yaml:"smtp"`
}

// GRPCConfig configuration for GRPC server.
type GRPCConfig struct {
	Port    int64         `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// SMTPConfig configuration for SMTP.
type SMTPConfig struct {
	MailRu MailRuConfig `yaml:"mail_ru"`
	YDX    YandexConfig `yaml:"yandex"`
}

// MailRuConfig configuration for MailRu SMTP.
type MailRuConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
}

// YandexConfig configuration for Yandex SMTP.
type YandexConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
}
