package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Environment string `yaml:"environment" env:"ENVIRONMENT" env-required:"true"`
	Port        int    `yaml:"port" env:"PORT" env-default:"42069"`
	DbPath      string `yaml:"db_path" env:"DB_PATH"`
}

func MustLoad() *AppConfig {
	path := fetchConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist at path: " + path)
	} else if err != nil {
		panic("Error checking config file: " + err.Error())
	}

	var cfg AppConfig
	cleanenv.ReadConfig(path, &cfg)

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "Path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path == "" {
		panic("Config path must be provided via --config flag or CONFIG_PATH environment variable")
	}

	return path
}
