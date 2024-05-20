package server

import (
	"context"
	"database/sql"

	"github.com/DeepAung/deep-art/api/middlewares"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/views/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	app *echo.Echo
	mid *middlewares.Middleware
	cfg *config.Config
	db  *sql.DB
}

func NewServer(
	app *echo.Echo,
	cfg *config.Config,
	db *sql.DB,
) *Server {
	usersRepo := repositories.NewUsersRepo(db, cfg.App.Timeout)
	usersSvc := services.NewUsersSvc(usersRepo, cfg)
	mid := middlewares.NewMiddleware(usersSvc, cfg)

	return &Server{
		app: app,
		mid: mid,
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Start() {
	s.app.Use(s.mid.Logger())

	s.app.Use(middleware.Recover())
	s.app.Use(middleware.BodyLimit(s.cfg.App.BodyLimit))
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		Timeout: s.cfg.App.Timeout,
	}))
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.cfg.App.CorsOrigins,
	}))

	s.app.Static("/static", "static")

	s.UsersRouter()
	s.TestRouter()
	s.PagesRouter()

	s.app.GET("*", func(c echo.Context) error {
		return pages.Error("Page Not Found").Render(context.Background(), c.Response())
	})

	s.app.Start(":3000")
}
