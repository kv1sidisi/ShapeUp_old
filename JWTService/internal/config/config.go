package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env  string     `yaml:"env"`
	JWT  JWTConfig  `yaml:"jwt"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type JWTConfig struct {
	AccessTokenSecretKey  string `yaml:"access_token_secret_key"`
	RefreshTokenSecretKey string `yaml:"refresh_token_secret_key"`
}

type GRPCConfig struct {
	Port    int64         `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// MustLoad gets config path and panics if there is any errors in parsing config.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath gets config path from arguments and panics if there is any errors in parsing config.
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

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
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
