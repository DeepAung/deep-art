package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/DeepAung/deep-art/config"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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
	engine := html.New("./views", ".html")

	return &server{
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
			Views:        engine,
		}),
		db:  db,
		cfg: cfg,
	}
}

func (s *server) Start() {
	// Static files
	s.app.Static("/static", "./public")

	// Middlewares
	mid := InitMiddlewares(s.cfg, s.db)
	s.app.Use(mid.Logger())
	s.app.Use(mid.Cors())

	// Modules
	v1 := s.app.Group("/api/v1")
	modules := InitModules(v1, s, mid)

	modules.MonitorModule()
	modules.ViewsModule()
	modules.UsersModule()

	s.app.Use(mid.RouterCheck())

	// Graceful shutdown
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to url
	if err := s.app.Listen(s.cfg.App().Url()); err != nil {
		log.Fatal(err)
	}

	// Clean up tasks go here...
}
