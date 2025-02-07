package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	JWT        JWTConfig        `yaml:"jwt"`
	GRPC       GRPCConfig       `yaml:"grpc"`
	Storage    StorageConfig    `yaml:"storage"`
	GRPCClient GRPCClientConfig `yaml:"grpc_client"`
}

// JWTConfig structure represents information from config about jwt
type JWTConfig struct {
	AccessSecret  string `yaml:"access_token_secret_key" env-required:"true"`
	RefreshSecret string `yaml:"refresh_token_secret_key" env-required:"true"`
}

// GRPCConfig structure represents information from config to configure grpc server.
type GRPCConfig struct {
	Port    int64         `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPCClientConfig struct {
	SendingServiceAddress string `yaml:"sending_service_address" env-required:"true"`
	JWTServiceAddress     string `yaml:"jwt_service_address" env-required:"true"`
}

// StorageConfig structure represents information from config to connect to database.
type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
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
