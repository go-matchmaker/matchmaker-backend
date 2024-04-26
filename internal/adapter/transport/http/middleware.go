package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"go.uber.org/zap"
)

func (s *server) HttpMiddleware() error {
	err := s.app.Use(
		cors.New(*s.getCorsConfig()),
		logger.New(),
		s.security,
	)

	if err != nil {
		zap.S().Fatal("middleware setup error: %w", err)
	}

	return nil
}

func (s *server) getCorsConfig() *cors.Config {
	return &cors.Config{
		AllowCredentials: true,
	}
}

func (s *server) security(c fiber.Ctx) error {
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Download-Options", "noopen")
	c.Set("Strict-Transport-Security", "max-age=5184000")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
	c.Set("Content-Security-Policy", "default-src https:")
	return c.Next()
}

func (s *server) AuthMiddleware() error {
	return nil
}
