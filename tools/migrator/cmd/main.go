package main

import (
	"errors"
	"fmt"
	"migrator/internal/config"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	rootPath := os.Getenv("ROOT_PATH")
	if rootPath == "" {
		panic("ROOT_PATH environment variable is required")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH environment variable is required")
	}

	cnfg := config.MustInitConfig(rootPath, configPath)

	MustMigratePostgres(rootPath, cnfg)

}

func MustMigratePostgres(rootPath string, postgresConfig config.Config) {
	migrationPath := filepath.Join(rootPath, postgresConfig.MigrationsPostgresPath)
	migrationPath = filepath.ToSlash(migrationPath)

	if migrationPath == "" {
		panic("migrations-path is required")
	}

	migrationsURL := "file://" + migrationPath + "/"

	databaseURL, err := postgresConfig.PostgresHost.GetHostToGoMigrate()
	if err != nil {
		panic(err)
	}
	m, err := migrate.New(migrationsURL, databaseURL)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations postgres applied")
}
