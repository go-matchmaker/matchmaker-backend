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

func (f *FiberResponseFactory) Response(err error, msg string, responseType int) error {
	if err != nil {
		return f.ctx.Status(responseType).JSON(fiber.Map{
			"error": true,
			"msg":   msg + err.Error(),
		})
	}
	return f.ctx.Status(responseType).JSON(fiber.Map{
		"error": false,
		"msg":   msg,
	})
}
