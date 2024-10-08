package server

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/handlers"
	"github.com/DeepAung/deep-art/api/middlewares"
	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/pages"
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
	usersRepo := repositories.NewUsersRepo(r.s.db, r.s.cfg.App.Timeout)
	usersSvc := services.NewUsersSvc(usersRepo, r.storer, r.s.cfg)
	artsRepo := repositories.NewArtsRepo(r.storer, r.s.db, r.s.cfg.App.Timeout)
	artsSvc := services.NewArtsSvc(artsRepo, r.storer, r.s.cfg)
	tagsRepo := repositories.NewTagsRepo(r.s.db, r.s.cfg.App.Timeout)
	tagsSvc := services.NewTagsSvc(tagsRepo)
	handler := handlers.NewPagesHandler(usersSvc, artsSvc, tagsSvc)

	setUserData := middlewares.SetUserData

	r.s.app.GET("/", handler.Welcome)
	r.s.app.GET("/signin", handler.SignIn)
	r.s.app.GET("/signup", handler.SignUp)
	r.s.app.GET("/home", handler.Home, r.mid.OnlyAuthorized(setUserData()))
	r.s.app.GET("/arts/:id", handler.ArtDetail, r.mid.OnlyAuthorized(setUserData()))
	r.s.app.GET("/me", handler.MyProfile, r.mid.OnlyAuthorized(setUserData()))
	r.s.app.GET("/creators/:id", handler.CreatorProfile, r.mid.OnlyAuthorized(setUserData()))

	r.s.app.GET("/creator", handler.CreatorHomePage, r.mid.OnlyAuthorized(setUserData()))
	r.s.app.GET(
		"/creator/arts/create",
		handler.CreatorCreateArt,
		r.mid.OnlyAuthorized(setUserData()),
	)
	r.s.app.GET(
		"/creator/arts/:id",
		handler.CreatorArtDetail,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OwnedArt("id"),
	)

	r.s.app.GET(
		"/admin",
		handler.AdminHomePage,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
}

func (r *Router) UsersRouter() {
	repo := repositories.NewUsersRepo(r.s.db, r.s.cfg.App.Timeout)
	svc := services.NewUsersSvc(repo, r.storer, r.s.cfg)
	handler := handlers.NewUsersHandler(svc, r.mid, r.s.cfg)

	setPayload := middlewares.SetPayload

	// Auth
	r.s.app.POST("/api/auth/signin", handler.SignIn)
	r.s.app.POST("/api/auth/signup", handler.SignUp)
	r.s.app.POST("/api/auth/signout", handler.SignOut, r.mid.OnlyAuthorized(setPayload()))
	r.s.app.POST(
		"/api/auth/update-tokens",
		handler.UpdateTokens,
		r.mid.JwtRefreshToken(setPayload()),
	)
	r.s.app.POST(
		"/api/auth/set-password-and-disconnect",
		handler.SetPasswordAndDisconnect,
		r.mid.OnlyAuthorized(setPayload()),
	)

	// OAuth
	r.s.app.GET("/api/auth/:provider", handler.OAuthHandler)
	r.s.app.GET("/api/auth/:provider/callback", handler.OAuthCallback)

	// Users CRUD
	r.s.app.PUT("/api/users", handler.UpdateUser, r.mid.OnlyAuthorized(setPayload()))
	r.s.app.DELETE("/api/users/me", handler.DeleteMyUser, r.mid.OnlyAuthorized(setPayload()))

	// Etc.
	r.s.app.POST(
		"/api/creators/:id/toggle-follow",
		handler.ToggleFollow,
		r.mid.OnlyAuthorized(setPayload()),
	)
}

func (r *Router) ArtsRouter() {
	repo := repositories.NewArtsRepo(r.storer, r.s.db, r.s.cfg.App.Timeout)
	svc := services.NewArtsSvc(repo, r.storer, r.s.cfg)
	handler := handlers.NewArtsHandler(svc, r.storer, r.s.cfg)

	setPayload := middlewares.SetPayload

	r.s.app.GET("/api/arts", handler.FindManyArts)
	r.s.app.GET(
		"/api/arts-with-art-type",
		handler.FindManyArtsWithArtType,
		r.mid.OnlyAuthorized(setPayload()),
	)

	r.s.app.POST("/api/arts", handler.CreateArt, r.mid.OnlyAuthorized(setPayload()))
	r.s.app.PUT(
		"/api/arts/:id",
		handler.UpdateArt,
		r.mid.OnlyAuthorized(setPayload()),
		r.mid.OwnedArt("id"),
	)
	r.s.app.DELETE(
		"/api/arts/:id",
		handler.DeleteArt,
		r.mid.OnlyAuthorized(setPayload()),
		r.mid.OwnedArt("id"),
	)
	r.s.app.GET(
		"/api/arts/:id/download",
		handler.DownloadArt,
		r.mid.OnlyAuthorized(setPayload()),
		r.mid.CanDownload("id"),
	)
	r.s.app.POST(
		"/api/arts/:id/files",
		handler.UploadFiles,
		r.mid.OnlyAuthorized(setPayload()),
		r.mid.OwnedArt("id"),
	)
	r.s.app.DELETE(
		"/api/arts/:artId/files/:fileId",
		handler.DeleteFile,
		r.mid.OnlyAuthorized(setPayload()),
		r.mid.OwnedArt("artId"),
	)
	r.s.app.PUT(
		"/api/arts/:id/cover",
		handler.ReplaceCover,
		r.mid.OnlyAuthorized(setPayload()),
		r.mid.OwnedArt("id"),
	)

	r.s.app.POST(
		"/api/arts/:id/toggle-star",
		handler.ToggleStar,
		r.mid.OnlyAuthorized(setPayload()),
	)
	r.s.app.POST("/api/arts/:id/buy", handler.BuyArt, r.mid.OnlyAuthorized(setPayload()))
}

func (r *Router) TagsRouter() {
	repo := repositories.NewTagsRepo(r.s.db, r.s.cfg.App.Timeout)
	svc := services.NewTagsSvc(repo)
	handler := handlers.NewTagsHandler(svc)

	r.s.app.GET("/api/tags/filter", handler.TagsFilter)
	r.s.app.GET("/api/tags/options", handler.TagsOptions)

	setUserData := middlewares.SetUserData

	r.s.app.POST(
		"/api/tags",
		handler.CreateTag,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
	r.s.app.GET(
		"/api/tags",
		handler.GetTags,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
	r.s.app.PUT(
		"/api/tags/:id",
		handler.UpdateTag,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
	r.s.app.DELETE(
		"/api/tags/:id",
		handler.DeleteTag,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
}

func (r *Router) CodesRouter() {
	repo := repositories.NewCodesRepo(r.s.db, r.s.cfg.App.Timeout)
	svc := services.NewCodesSvc(repo, r.s.cfg)
	handler := handlers.NewCodesHandler(svc, r.s.cfg)

	setPayload := middlewares.SetPayload
	setUserData := middlewares.SetUserData

	r.s.app.POST("/api/codes/use", handler.UseCode, r.mid.OnlyAuthorized(setPayload()))

	r.s.app.POST(
		"/api/codes",
		handler.CreateCode,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
	r.s.app.GET(
		"/api/codes",
		handler.GetCodes,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
	r.s.app.PUT(
		"/api/codes/:id",
		handler.UpdateCode,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
	r.s.app.DELETE(
		"/api/codes/:id",
		handler.DeleteCode,
		r.mid.OnlyAuthorized(setUserData()),
		r.mid.OnlyAdmin,
	)
}

// ------------------------------------------------------------------------- //

func (r *Router) TestRouter() {
	test := r.s.app.Group("/test")
	test.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "test route")
	})

	r.testTagsRouter(test)
	r.testCodesRouter(test)
	r.testFilesRouter(test)
	r.testArtsRouter(test)
	r.testUsersRouter(test)
	r.testPagesRouter(test)
}

func (r *Router) testTagsRouter(testGroup *echo.Group) {
	tagsRepo := repositories.NewTagsRepo(r.s.db, r.s.cfg.App.Timeout)

	tagsGroup := testGroup.Group("/tags")
	tagsGroup.GET("/", func(c echo.Context) error {
		tags, err := tagsRepo.FindAllTags()
		if err != nil {
			_, status := httperror.Extract(err)
			return c.JSON(status, err.Error())
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
	codesGroup.POST("/use", handler.UseCode)
}

func (r *Router) testFilesRouter(testGroup *echo.Group) {
	handler := handlers.NewTestFilesHandler(r.storer)

	filesGroup := testGroup.Group("/files")
	filesGroup.POST("/upload", handler.UploadFiles)
	filesGroup.POST("/delete", handler.DeleteFiles)
}

func (r *Router) testArtsRouter(testGroup *echo.Group) {
	repo := repositories.NewArtsRepo(r.storer, r.s.db, r.s.cfg.App.Timeout)
	handler := handlers.NewTestArtsHandler(repo)

	artsGroup := testGroup.Group("/arts")
	artsGroup.GET("", handler.FindManyArts)
	artsGroup.GET("/starred", handler.FindManyStarredArts)
	artsGroup.GET("/bought", handler.FindManyBoughtArts)
	artsGroup.GET("/created", handler.FindManyCreatedArts)
	artsGroup.GET("/:id", handler.FindOneArt)
}

func (r *Router) testUsersRouter(testGroup *echo.Group) {
	repo := repositories.NewUsersRepo(r.s.db, r.s.cfg.App.Timeout)
	handler := handlers.NewTestUsersHandler(repo)

	usersGroup := testGroup.Group("/users")
	usersGroup.GET("/creators/:id", handler.GetCreator)
}

func (r *Router) testPagesRouter(testGroup *echo.Group) {
	pagesGroup := testGroup.Group("/pages")
	pagesGroup.GET("/set-password-and-disconnect", func(c echo.Context) error {
		return utils.Render(c, pages.SetPasswordAndDisconnect("github"), http.StatusOK)
	})
}
