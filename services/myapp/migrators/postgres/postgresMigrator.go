package postgresMigrator

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"main.go/internal/config"
)

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
