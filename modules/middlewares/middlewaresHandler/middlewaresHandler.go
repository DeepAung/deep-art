package middlewaresHandler

import (
	"strconv"
	"strings"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/modules/middlewares/middlewaresUsecase"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	routerCheckErr response.TraceId = "middlwares-001"
	jwtAuthErr     response.TraceId = "middlwares-002"
	onlyAdminErr   response.TraceId = "middlwares-003"
	adminAuthErr   response.TraceId = "middlwares-004"
	apiKeyAuthErr  response.TraceId = "middlwares-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	OnlyAdmin() fiber.Handler
	AdminAuth() fiber.Handler
	ApiKeyAuth() fiber.Handler
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
		if tokenString == "" {
			return response.Error(c, fiber.StatusUnauthorized, jwtAuthErr, "token required")
		}

		claims, err := mytoken.ParseToken(h.cfg.Jwt(), &mytoken.Access, tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, jwtAuthErr, err.Error())
		}

		if !h.middlewaresUsecase.FindAccessToken(claims.Payload.UserId, tokenString) {
			return response.Error(c, fiber.StatusUnauthorized, jwtAuthErr, "invalid token")
		}

		// TODO: should i use cookie???
		userId := claims.Payload.UserId
		c.Locals("userId", userId)
		utils.SetCookie(c, "userId", strconv.Itoa(userId), h.cfg.Jwt().AccessExpires())

		isAdmin := claims.Payload.IsAdmin
		c.Locals("isAdmin", isAdmin)
		utils.SetCookie(c, "isAdmin", strconv.FormatBool(isAdmin), h.cfg.Jwt().AccessExpires())

		return c.Next()
	}
}

func (h *middlewaresHandler) OnlyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin := c.Locals("isAdmin").(bool)
		if !isAdmin {
			return response.Error(
				c,
				fiber.StatusUnauthorized,
				onlyAdminErr,
				"no permission to access",
			)
		}

		return c.Next()
	}
}

func (h *middlewaresHandler) AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			return response.Error(c, fiber.StatusUnauthorized, adminAuthErr, "token required")
		}

		err := mytoken.VerifyToken(h.cfg.Jwt(), &mytoken.Admin, tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, adminAuthErr, err.Error())
		}

		return c.Next()

	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("X-Api-Key")
		err := mytoken.VerifyToken(h.cfg.Jwt(), &mytoken.ApiKey, tokenString)
		if err != nil {
			return response.Error(
				c,
				fiber.StatusUnauthorized,
				apiKeyAuthErr,
				"invalid or no apikey",
			)
		}

		return c.Next()
	}
}
