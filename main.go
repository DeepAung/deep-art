package main

import (
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/db"
	"github.com/DeepAung/deep-art/pkg/server"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.NewConfig()

	db := db.InitDB("db.db")
	app := echo.New()

	server.NewServer(app, db, cfg).Start()
}
