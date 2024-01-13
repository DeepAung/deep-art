package main

import (
	"os"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/server"
	"github.com/DeepAung/deep-art/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env.prod"
	}
	return os.Args[1]
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.ConnectDb(cfg.Db())
	defer db.Close()

	server.NewServer(db, cfg).Start()
}
