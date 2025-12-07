package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StorageSQLLitePath     string       `yaml:"storageSQLLitePath"`
	MigrationsSqllitePath  string       `yaml:"migrationsSqllitePath"`
	MigrationsPostgresPath string       `yaml:"migrationsPostgresPath"`
	PostgresHost           PostgresHost `yaml:"postgresHost"`
}

type PostgresHost struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
	SslMode  string `yaml:"sslmode"`
}

func (p *PostgresHost) GetStringHost() (string, error) {
	if p.Host == "" {
		return "", errors.New("empty host in config")
	}
	if p.User == "" {
		return "", errors.New("empty user in config")
	}
	if p.Password == "" {
		return "", errors.New("empty password in config")
	}
	if p.Dbname == "" {
		return "", errors.New("empty dbname in config")
	}
	if p.Port == 0 {
		return "", errors.New("empty port in config")
	}
	host := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s",
		p.Host,
		p.User,
		p.Password,
		p.Dbname,
		p.Port,
		p.SslMode)
	return host, nil
}

func (p *PostgresHost) GetHostToGoMigrate() (string, error) {
	if p.Host == "" {
		return "", errors.New("empty host in config")
	}
	if p.User == "" {
		return "", errors.New("empty user in config")
	}
	if p.Password == "" {
		return "", errors.New("empty password in config")
	}
	if p.Dbname == "" {
		return "", errors.New("empty dbname in config")
	}
	if p.Port == 0 {
		return "", errors.New("empty port in config")
	}
	host := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.Dbname, p.SslMode)
	return host, nil
}

func MustInitConfig(rootPath string, configPath string) Config {
	path := filepath.Join(rootPath, configPath)

	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	configWithEnvs, err := InjectEnv(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return configWithEnvs
}

func InjectEnv(config Config) (Config, error) {
	requiredEnv := []string{"POSTGRES_PASSWORD", "POSTGRES_USER", "POSTGRES_HOST"}
	envValues := make(map[string]string)
	for _, key := range requiredEnv {
		value := os.Getenv(key)
		if value == "" {
			return Config{}, fmt.Errorf("%s environment variable is required", key)
		}
		envValues[key] = value
	}
	config.PostgresHost.Password = envValues["POSTGRES_PASSWORD"]
	config.PostgresHost.User = envValues["POSTGRES_USER"]
	config.PostgresHost.Host = envValues["POSTGRES_HOST"]
	return config, nil
}
