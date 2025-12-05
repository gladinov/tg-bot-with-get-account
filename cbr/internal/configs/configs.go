package configs

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env         string `yaml:"env"`
	CbrHost     string `yaml:"cbrHost"`
	Http_server string `yaml:"http_server"`
}

func MustInitConfig(rootPath string) Config {
	var config Config
	Path := filepath.Join(rootPath, "configs", "local.yaml")
	file, err := os.ReadFile(Path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return config
}
