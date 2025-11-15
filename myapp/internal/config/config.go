package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env                       string `yaml:"env"`
	DbType                    string `yaml:"dbType"`
	ServiceStorageSQLLitePath string `yaml:"serviceStorageSQLLitePath"`
	postgeSQLHost             string `yaml:"postgeSQLHost"`
}

func MustInitConfig(rootPath string) Config {
	path := filepath.Join(rootPath, "configs", "config.yaml")

	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return config
}
