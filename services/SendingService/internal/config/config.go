package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
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

// MustLoad tries to get config path.
//
// Panics if there is any errors in parsing config.
//
// Returns Config.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath tries to get config path from arguments.
//
// Panics if there is any errors in parsing config.
//
// Returns Config.
func MustLoadByPath(configPath string) *Config {
	//check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist" + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or env variable.
//
// Priority: flag > env > default.
//
// Default value is empty string "".
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
