package middlewaresHandler

import (
	"strings"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresUsecase"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	routerCheckErr response.TraceId = "middlwares-001"
	jwtAuthErr     response.TraceId = "middlwares-002"
	paramsCheckErr response.TraceId = "middlwares-003"
	authorizeErr   response.TraceId = "middlwares-004"
	apiKeyErr      response.TraceId = "middlwares-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecase.IMiddlewaresUsecase
}

func NewMiddlewaresHandler(
	cfg config.IConfig,
	middlewaresUsecase middlewaresUsecase.IMiddlewaresUsecase,
) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New()
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusNotFound, routerCheckErr, "router not found")
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		claims, err := mytoken.ParseToken(h.cfg.Jwt(), tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, jwtAuthErr, err.Error())
		}

		if !h.middlewaresUsecase.FindAccessToken(claims.Payload.UserId, tokenString) {
			return response.Error(c, fiber.StatusUnauthorized, jwtAuthErr, "invalid token")
		}

		c.Locals("userId", claims.Payload.UserId)
		return c.Next()
	}
}
