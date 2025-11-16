package configs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	MoexHost   string `yaml:"moexHost"`
	HttpServer string `yaml:"http_server"`
}

func MustLoad(rootPath string) *Config {
	const op = "configs.mustLoad"
	envPath := filepath.Join(rootPath, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("%s:Could not find any .env files:%s", op, err)
	}

	configPath := filepath.Join(rootPath, os.Getenv("CONFIG_PATH"))
	if configPath == "" {
		log.Fatalf("%s: CONFIG_PATH is not set", op)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%s:config file does not exist:%s", op, err)
	}

	var cfg Config

	if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("%s:cannot read config:%s", op, err)
	}

	return &cfg
}
