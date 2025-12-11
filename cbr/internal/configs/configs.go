package configs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env        string `yaml:"env"`
	CbrHost    string `yaml:"cbrHost"`
	CbrAppHost string `yaml:"cbr_app_host"`
	CbrAppPort int    `yaml:"cbr_app_port"`
}

func (c *Config) GetCbrAppServer() string {
	cbrPortStr := strconv.Itoa(c.CbrAppPort)
	return c.CbrAppHost + cbrPortStr
}

func MustInitConfig(rootPath string, configPath string) Config {
	var config Config
	Path := filepath.Join(rootPath, configPath)
	file, err := os.ReadFile(Path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	configWithEnvs, err := InjectEnvs(config)
	if err != nil {
		panic(err)
	}

	return configWithEnvs
}

func InjectEnvs(config Config) (Config, error) {
	cbrPort := os.Getenv("CBR_PORT")
	if cbrPort == "" {
		return Config{}, errors.New("CBR_PORT environment variable is required")
	}

	cbrPortFromConf, err := strconv.Atoi(cbrPort)
	if err != nil {
		return Config{}, errors.New("CBR_PORT environment variable bad format")
	}
	config.CbrAppPort = cbrPortFromConf

	return config, nil
}
