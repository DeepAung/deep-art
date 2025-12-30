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
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
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
	// Init OAuth
	gothic.Store = sessions.NewCookieStore([]byte(s.cfg.OAuth.SessionSecret))
	goth.UseProviders(
		google.New(
			s.cfg.OAuth.GoogleKey,
			s.cfg.OAuth.GoogleSecret,
			s.cfg.App.Address+"/api/auth/google/callback",
		),
		github.New(
			s.cfg.OAuth.GithubKey,
			s.cfg.OAuth.GithubSecret,
			s.cfg.App.Address+"/api/auth/github/callback",
		),
	)

	myStorer := storer.NewGCPStorer(s.cfg)
	mid := s.InitMiddleware(myStorer)

	s.app.Static("/static", "static")

	s.InitRouter(mid, myStorer)

	s.app.Start(":3000")
}

func (s *Server) InitMiddleware(storer storer.Storer) *middlewares.Middleware {
	usersRepo := repositories.NewUsersRepo(s.db, s.cfg.App.Timeout)
	usersSvc := services.NewUsersSvc(usersRepo, storer, s.cfg)
	artsRepo := repositories.NewArtsRepo(storer, s.db, s.cfg.App.Timeout)
	artsSvc := services.NewArtsSvc(artsRepo, storer, s.cfg)
	mid := middlewares.NewMiddleware(usersSvc, artsSvc, s.cfg)

	s.app.Use(mid.Logger())
	s.app.Use(middleware.Recover())
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
