package server

import (
	"database/sql"
	"net/http"

	"github.com/DeepAung/deep-art/api/middlewares"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	app *echo.Echo
	db  *sql.DB
	cfg *config.Config
}

func NewServer(
	app *echo.Echo,
	db *sql.DB,
	cfg *config.Config,
) *Server {
	return &Server{
		app: app,
		db:  db,
		cfg: cfg,
	}
}

func (s *Server) Start() {
	storer := storer.NewLocalStorer(s.cfg)
	mid := s.InitMiddleware()

	s.app.Static("/static", "static")
	s.app.Static("/node_modules", "node_modules")

	s.InitRouter(mid, storer)

	s.app.Start(":3000")
}

func (s *Server) InitMiddleware() *middlewares.Middleware {
	usersRepo := repositories.NewUsersRepo(s.db, s.cfg.App.Timeout)
	usersSvc := services.NewUsersSvc(usersRepo, s.cfg)
	mid := middlewares.NewMiddleware(usersSvc, s.cfg)

	s.app.Use(mid.Logger())
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.BodyLimit(s.cfg.App.BodyLimit))
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		Timeout: s.cfg.App.Timeout,
	}))
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.cfg.App.CorsOrigins,
	}))

	return mid
}

func (s *Server) InitRouter(mid *middlewares.Middleware, storer storer.Storer) {
	r := NewRouter(s, mid, storer)

	r.UsersRouter()
	r.ArtsRouter()
	r.TagsRouter()
	r.CodesRouter()
	r.TestRouter()
	r.PagesRouter()

	r.s.app.GET("*", func(c echo.Context) error {
		return utils.Render(c, pages.Error("Page Not Found"), http.StatusNotFound)
	})
}
