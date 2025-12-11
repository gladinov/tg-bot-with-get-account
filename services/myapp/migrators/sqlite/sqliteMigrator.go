package sqliteMigrator

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MustMigrateSqllite(rootPath string, storagePath string, migrationPath string) {
	migrationPath = filepath.Join(rootPath, migrationPath)
	storagePath = filepath.ToSlash(storagePath)
	migrationPath = filepath.ToSlash(migrationPath)

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationPath == "" {
		panic("migrations-path is required")
	}

	migrationsURL := "file://" + migrationPath + "/"
	databaseURL := fmt.Sprintf("sqlite3://%s", storagePath)
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

	fmt.Println("migrations sqlLite applied")
}
