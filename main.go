package main

import (
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/db"
	"github.com/DeepAung/deep-art/pkg/middlewares"
	"github.com/DeepAung/deep-art/pkg/router"
	"github.com/DeepAung/deep-art/pkg/server"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.NewConfig()
	// cfg.Print()

	db := db.InitDB("db.db")
	app := echo.New()
	mid := middlewares.NewMiddleware()
	router := router.NewRouter(app)
	server := server.NewServer(app, cfg, mid, db, router)

	server.Start()
}
