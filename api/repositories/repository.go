package repositories

import (
	"database/sql"
	"log"

	"github.com/DeepAung/deep-art/pkg/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewTestDB() (*sql.DB, *migrate.Migrate) {
	sourceURL := "file:///home/deepaung/projects/deep-art/migrations"
	databaseURL := "sqlite3:///home/deepaung/projects/deep-art/test.db"
	databaseDir := "/home/deepaung/projects/deep-art/test.db"

	testDB := db.InitDB(databaseDir)
	migrateDB, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}

	return testDB, migrateDB
}

func ResetDB(m *migrate.Migrate) {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("initAndResetDb: m.Down: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("initAndResetDb: m.Up: %v", err)
	}
}
