package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func (s *server) HTTPMiddleware() error {
	s.app.Use(
		cors.New(),
		compress.New(),
		favicon.New(),
		logger.New(),
		recover.New(),
		etag.New(),
		requestid.New(),
		s.security,
	)

	return nil
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
