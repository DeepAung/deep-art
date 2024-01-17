package main

import (
	"os"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/server"
	"github.com/DeepAung/deep-art/pkg/databases"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
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

	// TODO: should this be in server.go???
	goth.UseProviders(
		google.New(
			cfg.OAuth().GoogleKey(),
			cfg.OAuth().GoogleSecret(),
			"http://127.0.0.1:3000/api/v1/users/google/callback",
		),
		github.New(
			cfg.OAuth().GithubKey(),
			cfg.OAuth().GithubSecret(),
			"http://127.0.0.1:3000/api/v1/users/github/callback",
		),
	)

	server.NewServer(db, cfg).Start()
}
