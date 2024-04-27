package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/service"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/token"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go.uber.org/zap"
	std_http "net/http"
)

var (
	_         http.ServerMaker = (*server)(nil)
	ServerSet                  = wire.NewSet(NewHTTPServer)
)

type (
	server struct {
		ctx          context.Context
		cfg          *config.Container
		app          *fiber.App
		cfgFiber     *fiber.Config
		userService  service.UserMaker
		tokenService token.TokenMaker
	}
)

func NewHTTPServer(
	ctx context.Context,
	cfg *config.Container,
	userService service.UserMaker,
	tokenService token.TokenMaker,
) http.ServerMaker {
	return &server{
		app:          fiber.New(),
		ctx:          ctx,
		cfg:          cfg,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *server) Start(ctx context.Context) error {
	fiberConnURL := fmt.Sprintf("%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)
	go func() {
		zap.S().Info("Starting HTTP server on ", fiberConnURL)
		if err := s.app.Listen(fiberConnURL); err != nil {
			if errors.Is(err, std_http.ErrServerClosed) {
				return
			}
			zap.S().Fatal("server listen error: %w", err)
		}
	}()
	return nil
}

func (s *server) Close(ctx context.Context) error {
	zap.S().Info("HTTP-Server Context is done. Shutting down server...")
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		zap.S().Error("server shutdown error: %w", err)
		return err
	}
	return nil
}

func (s *server) SetupRouter() {
	route := s.app.Group("/api/v1")
	route.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})

	route.Post("/register", s.RegisterUser)
}
