package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// MustLoad tries to load config in given structure.
//
// Panics if there is any errors in parsing config.
func MustLoad(cfg interface{}) {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	MustLoadByPath(path, cfg)
}

// MustLoadByPath tries to load config in given structure from arguments.
//
// Panics if there is any errors in parsing config.
func MustLoadByPath(configPath string, cfg interface{}) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
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
