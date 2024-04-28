package http

import (
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/http"
	"github.com/gofiber/fiber/v3"
)

type FiberResponseFactory struct {
	ctx fiber.Ctx
}

func NewFiberResponseFactory(ctx fiber.Ctx) http.ResponseFactory {
	return &FiberResponseFactory{ctx: ctx}
}

func (f *FiberResponseFactory) Response(isError bool, msg string, responseType int) error {
	if isError {
		return f.ctx.Status(responseType).JSON(fiber.Map{
			"error": true,
			"msg":   msg,
		})
	}
	return f.ctx.Status(responseType).JSON(fiber.Map{
		"error": false,
		"msg":   msg,
	})
}
