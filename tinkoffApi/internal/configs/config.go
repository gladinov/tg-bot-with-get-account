package configs

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

type Configs struct {
	TinkoffApiConfig *investgo.Config
	Config           *Config
}

type Config struct {
	HttpServer      string          `yaml:"http_server"`
	RedisHTTPServer RedisHTTPServer `yaml:"redisHTTP"`
}
type RedisHTTPServer struct {
	Address     string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func MustTinkoffConfigLoad(rootPath string) *investgo.Config {
	const op = "configs.MustTinkoffConfigLoad"
	envPath := filepath.Join(rootPath, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("%s:Could not find any .env files:%s", op, err)
	}

	configPath := filepath.Join(rootPath, os.Getenv("TINKOFF_CONFIG_PATH"))
	if configPath == "" {
		log.Fatalf("%s: CONFIG_PATH is not set", op)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%s:config file does not exist:%s", op, err)
	}

	config, err := investgo.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("op: %s , can't load config", op)
	}

	return &config
}

func MustConfigLoad(rootPath string) *Config {
	const op = "configs.MustConfigLoad"
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

func MustInitConfigs(rootPath string) *Configs {
	var configs Configs
	configs.Config = MustConfigLoad(rootPath)
	configs.TinkoffApiConfig = MustTinkoffConfigLoad(rootPath)
	return &configs
}
