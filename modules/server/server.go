package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/DeepAung/deep-art/config"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	db  *sqlx.DB
	cfg config.IConfig
}

func NewServer(db *sqlx.DB, cfg config.IConfig) IServer {
	return &server{
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
		db:  db,
		cfg: cfg,
	}
}

func (s *server) Start() {
	v1 := s.app.Group("/api/v1")
	modules := InitModuleFactory(v1, s)

	modules.MonitorModule()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = s.app.Shutdown()
	}()

	if err := s.app.Listen(s.cfg.App().Url()); err != nil {
		log.Fatal(err)
	}

	// clean up tasks go here...
}
