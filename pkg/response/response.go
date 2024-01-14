package response

import (
	"github.com/DeepAung/deep-art/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type TraceId string

type errorRes struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func Success(c *fiber.Ctx, code int, data any) error {
	logger.NewLogger(c, code, data).Print().Save()
	return c.Status(code).JSON(data)
}

func Error(c *fiber.Ctx, code int, traceId TraceId, msg string) error {
	errRes := &errorRes{
		TraceId: string(traceId),
		Msg:     msg,
	}
	logger.NewLogger(c, code, errRes).Print().Save()
	return c.Status(code).JSON(errRes)
}
