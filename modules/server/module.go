package server

import (
	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresHandler"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresRepository"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresUsecase"
	"github.com/DeepAung/deep-art/modules/monitor/monitorHandler"
	"github.com/DeepAung/deep-art/modules/users/usersHandler"
	"github.com/DeepAung/deep-art/modules/users/usersRepository"
	"github.com/DeepAung/deep-art/modules/users/usersUsecase"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IModuleFactory interface {
	MonitorModule()
	ViewsModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandler.IMiddlewaresHandler
}

func InitModules(
	r fiber.Router,
	s *server,
	mid middlewaresHandler.IMiddlewaresHandler,
) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(cfg config.IConfig, db *sqlx.DB) middlewaresHandler.IMiddlewaresHandler {
	repo := middlewaresRepository.NewMiddlewaresRepository(db)
	usecase := middlewaresUsecase.NewMiddlewaresUsecase(repo)
	return middlewaresHandler.NewMiddlewaresHandler(cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandler.NewMonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) ViewsModule() {
	router := m.r.Group("views")

	router.Get("/index", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})
}

func (m *moduleFactory) UsersModule() {
	repo := usersRepository.NewUsersRepository(m.s.db)
	usecase := usersUsecase.NewUsersUsecase(m.s.cfg, repo)
	handler := usersHandler.NewUsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
	router.Post("/logout", m.mid.JwtAuth(), handler.Logout)
	router.Post("/refresh", m.mid.JwtAuth(), handler.RefreshTokens)

	router.Post("/admin/token", m.mid.JwtAuth(), m.mid.AdminAuth(), handler.GenerateAdminToken)
	// router.Post("/admin/register", m.mid.AdminTokenAuth())

	router.Get("/:provider/login", handler.OAuthLoginOrRegister)
	router.Get("/:provider/register", handler.OAuthLoginOrRegister)
	router.Get(
		"/:provider/connect",
		handler.OAuthConnect,
	) // TODO: make user login before connect to oauths
	// router.Get("/:provider/connect", m.mid.JwtAuth(), handler.OAuthConnect)
	router.Get("/:provider/disconnect", m.mid.JwtAuth(), handler.OAuthDisconnect)
	router.Get("/:provider/callback", handler.OAuthCallback)
}
