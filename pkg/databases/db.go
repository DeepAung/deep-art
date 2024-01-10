package databases

import (
	"log"

	"github.com/DeepAung/deep-art/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectDb(cfg config.IDbConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatal("connect to db failed: ", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())

	return db
}
