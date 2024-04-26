package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	"github.com/bulutcan99/company-matcher/internal/core/port/http"
	"github.com/bulutcan99/company-matcher/internal/core/port/service"
	"github.com/bulutcan99/company-matcher/internal/core/port/token"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	std_http "net/http"
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

func (s *server) Start(ctx context.Context) error {
	fiberConnURL := fmt.Sprintf("%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)
	s.eg.Go(func() error {
		if err := s.app.Listen(fiberConnURL); err != nil {
			if errors.Is(err, std_http.ErrServerClosed) {
				return nil
			}
			zap.S().Error("server listen error: %w", err)
			return err
		}
		return nil
	})
	return nil
}

func (s *server) Close(ctx context.Context) error {
	s.eg.Go(func() error {
		select {
		case <-ctx.Done():
			zap.S().Info("Context is done. Shutting down server...")
			if err := s.app.Shutdown(); err != nil {
				zap.S().Error("server shutdown error: %w", err)
				return err
			}
			return nil
		}
	})
	return nil
}

func (s *server) SetupRouter() {
	route := s.app.Group("/api/v1")
	route.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})

	route.Post("/register", s.RegisterUser)
}
