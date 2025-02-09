package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
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

// fetchConfigPath fetches config path from command line flag or environment variable.
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
