package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UserStorageConfig struct {
	//Env                       string `yaml:"env"`
	DbType                 string `yaml:"dbType"`
	StorageSQLLitePath     string `yaml:"storageSQLLitePath"`
	MigrationsSqllitePath  string `yaml:"migrationsSqllitePath"`
	MigrationsPostgresPath string `yaml:"migrationsPostgresPath"`
	PostgresHost           `yaml:"postgresHost"`
}

func MustInitStorageConfig(rootPath string) UserStorageConfig {
	path := filepath.Join(rootPath, "configs", "userStorageConfig.yaml")

	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config UserStorageConfig
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return config
}
