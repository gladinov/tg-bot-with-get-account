package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	HTTPServer `yaml:"redis"`
}

type HTTPServer struct {
	Address     string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func MustInitRedisConfig(rootpath string, configPath string) *RedisConfig {
	fpath := filepath.Join(rootpath, configPath)
	file, err := os.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	var redisConfig RedisConfig
	err = yaml.Unmarshal(file, redisConfig)
	if err != nil {
		panic(err)
	}
	return &redisConfig
}
