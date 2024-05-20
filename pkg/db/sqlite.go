package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
