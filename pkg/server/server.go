package server

import (
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/middlewares"
	"github.com/DeepAung/deep-art/pkg/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	app *echo.Echo
	cfg *config.Config
	mid *middlewares.Middleware
	r   *router.Router
}

func NewServer(
	app *echo.Echo,
	cfg *config.Config,
	mid *middlewares.Middleware,
	r *router.Router,
) *Server {
	return &Server{
		app: app,
		cfg: cfg,
		mid: mid,
		r:   r,
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

	public := s.app.Group("/public")
	public.Use(middleware.Static("/public"))

	s.r.TestRouter()

	s.app.Start(":3000")
}
