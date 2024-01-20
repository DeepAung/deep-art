package server

import (
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresHandler"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresRepository"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresUsecase"
	"github.com/DeepAung/deep-art/modules/monitor/monitorHandler"
	"github.com/DeepAung/deep-art/modules/users/usersHandler"
	"github.com/DeepAung/deep-art/modules/users/usersRepository"
	"github.com/DeepAung/deep-art/modules/users/usersUsecase"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
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

func InitMiddlewares() middlewaresHandler.IMiddlewaresHandler {
	repo := middlewaresRepository.NewMiddlewaresRepository()
	usecase := middlewaresUsecase.NewMiddlewaresUsecase(repo)
	return middlewaresHandler.NewMiddlewaresHandler(usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandler.NewMonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repo := usersRepository.NewUsersRepository(m.s.db)
	usecase := usersUsecase.NewUsersUsecase(m.s.cfg, repo)
	handler := usersHandler.NewUsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
	router.Get("/:provider", handler.AuthenticateOAuth)
	router.Get("/:provider/callback", handler.CallbackOAuth)
}
