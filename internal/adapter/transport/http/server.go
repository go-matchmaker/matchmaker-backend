package http

import (
	"context"
	"fmt"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	"github.com/bulutcan99/company-matcher/internal/core/port/http"
	"github.com/bulutcan99/company-matcher/internal/core/port/service"
	"github.com/bulutcan99/company-matcher/internal/core/port/token"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	_         http.ServerMaker = (*server)(nil)
	ServerSet                  = wire.NewSet(NewHttpServer)
)

type (
	server struct {
		eg           *errgroup.Group
		cfg          *config.Container
		app          *fiber.App
		cfgFiber     *fiber.Config
		userService  service.UserMaker
		tokenService token.TokenMaker
	}
)

func NewHttpServer(
	eg *errgroup.Group,
	cfg *config.Container,
	userService service.UserMaker,
	tokenService token.TokenMaker,
) http.ServerMaker {
	return &server{
		app:          fiber.New(),
		eg:           eg,
		cfg:          cfg,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *server) Start(ctx context.Context) {
	fiberConnURL := fmt.Sprintf("%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)
	s.eg.Go(func() error {
		if err := s.app.Listen(fiberConnURL); err != nil {
			zap.S().Fatal("server listen error: %w", err)
		}
		return nil
	})

	s.eg.Go(func() error {
		select {
		case <-ctx.Done():
			zap.S().Info("Context is done. Shutting down server...")
			if err := s.app.Shutdown(); err != nil {
				zap.S().Fatal("server shutdown error: %w", err)
			}
			return nil
		}
	})

	if err := s.eg.Wait(); err != nil {
		zap.S().Fatal("errorgroup error: %w", err)
	}
}

func (s *server) Close() {
	zap.S().Info("Shutting down server...")
	if err := s.app.Shutdown(); err != nil {
		zap.S().Fatal("server shutdown error:", err)
	}
	return
}

func (s *server) SetupRouter() {
	route := s.app.Group("/api/v1")
	route.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})

	route.Post("/register", s.RegisterUser)
}
