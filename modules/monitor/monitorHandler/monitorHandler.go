package monitorHandler

import (
	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/monitor"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func NewMonitorHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return response.Success(c, fiber.StatusOK, res)
}
