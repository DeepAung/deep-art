package server

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/handlers"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/labstack/echo/v4"
)

func (s *Server) PagesRouter() {
	handler := handlers.NewPagesHandler()

	s.app.GET("/", handler.Welcome)
	s.app.GET("/home", handler.Home, s.mid.OnlyAuthorized)
	s.app.GET("/signin", handler.SignIn)
	s.app.GET("/signup", handler.SignUp)
}

func (s *Server) TestRouter() {
	s.app.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test route")
	})
}

func (s *Server) UsersRouter() {
	repo := repositories.NewUsersRepo(s.db, s.cfg.App.Timeout)
	svc := services.NewUsersSvc(repo, s.cfg)
	handler := handlers.NewUsersHandler(svc, s.cfg)

	s.app.POST("/api/signin", handler.SignIn)
	s.app.POST("/api/signup", handler.SignUp)
	s.app.POST("/api/signout", handler.SignOut, s.mid.OnlyAuthorized)
}
