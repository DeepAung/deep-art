package server

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/handlers"
	"github.com/DeepAung/deep-art/api/middlewares"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/labstack/echo/v4"
)

type Router struct {
	s      *Server
	mid    *middlewares.Middleware
	storer storer.Storer
}

func NewRouter(
	s *Server,
	mid *middlewares.Middleware,
	storer storer.Storer,
) *Router {
	return &Router{
		s:      s,
		mid:    mid,
		storer: storer,
	}
}

func (r *Router) PagesRouter() {
	handler := handlers.NewPagesHandler()

	r.s.app.GET("/", handler.Welcome)
	r.s.app.GET("/home", handler.Home, r.mid.OnlyAuthorized)
	r.s.app.GET("/signin", handler.SignIn)
	r.s.app.GET("/signup", handler.SignUp)
}

func (r *Router) TestRouter() {
	tagsRepo := repositories.NewTagsRepo(r.s.db, r.s.cfg.App.Timeout)
	codesRepo := repositories.NewCodesRepo(r.s.db, r.s.cfg.App.Timeout)
	handler := handlers.NewTestHandler(tagsRepo, codesRepo)

	test := r.s.app.Group("/test")
	test.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test route")
	})

	test.GET("/tags", handler.FindAllTags)

	test.GET("/codes/:id", handler.FindOneCodeById)
	test.PUT("/codes/:id", handler.UpdateCode)
	test.DELETE("/codes/:id", handler.DeleteCode)
	test.GET("/codes", handler.FindAllCodes)
	test.POST("/codes", handler.CreateCode)
}

func (r *Router) UsersRouter() {
	repo := repositories.NewUsersRepo(r.s.db, r.s.cfg.App.Timeout)
	svc := services.NewUsersSvc(repo, r.s.cfg)
	handler := handlers.NewUsersHandler(svc, r.s.cfg)

	r.s.app.POST("/api/signin", handler.SignIn)
	r.s.app.POST("/api/signup", handler.SignUp)
	r.s.app.POST("/api/signout", handler.SignOut, r.mid.OnlyAuthorized)
}
