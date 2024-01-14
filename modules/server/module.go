package server

import (
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresHandler"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresRepository"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresUsecase"
	"github.com/DeepAung/deep-art/modules/monitor/monitorHandler"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
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
