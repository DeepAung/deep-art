package server

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/handlers"
	"github.com/DeepAung/deep-art/api/middlewares"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/httperror"
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

func (r *Router) UsersRouter() {
	repo := repositories.NewUsersRepo(r.s.db, r.s.cfg.App.Timeout)
	svc := services.NewUsersSvc(repo, r.s.cfg)
	handler := handlers.NewUsersHandler(svc, r.s.cfg)

	r.s.app.POST("/api/signin", handler.SignIn)
	r.s.app.POST("/api/signup", handler.SignUp)
	r.s.app.POST("/api/signout", handler.SignOut, r.mid.OnlyAuthorized)
}

// ------------------------------------------------------------------------- //

func (r *Router) TestRouter() {
	test := r.s.app.Group("/test")
	test.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test route")
	})

	r.testTagsRouter(test)
	r.testCodesRouter(test)
	r.testFilesRouter(test)
}

func (r *Router) testTagsRouter(testGroup *echo.Group) {
	tagsRepo := repositories.NewTagsRepo(r.s.db, r.s.cfg.App.Timeout)

	tagsGroup := testGroup.Group("/tags")
	tagsGroup.GET("/", func(c echo.Context) error {
		tags, err := tagsRepo.FindAllTags()
		if err != nil {
			_, status := httperror.Extract(err)
			return c.JSON(status, err)
		}

		return c.JSON(http.StatusOK, tags)
	})
}

func (r *Router) testCodesRouter(testGroup *echo.Group) {
	repo := repositories.NewCodesRepo(r.s.db, r.s.cfg.App.Timeout)
	handler := handlers.NewTestCodesHandler(repo)

	codesGroup := testGroup.Group("/codes")
	codesGroup.GET("/:id", handler.FindOneCodeById)
	codesGroup.PUT("/:id", handler.UpdateCode)
	codesGroup.DELETE("/:id", handler.DeleteCode)
	codesGroup.GET("/", handler.FindAllCodes)
	codesGroup.POST("/", handler.CreateCode)
}

func (r *Router) testFilesRouter(testGroup *echo.Group) {
	handler := handlers.NewTestFilesHandler(r.storer)

	filesGroup := testGroup.Group("/files")
	filesGroup.POST("/upload", handler.UploadFiles)
	filesGroup.POST("/delete", handler.DeleteFiles)
}
