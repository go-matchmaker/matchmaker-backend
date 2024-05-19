package http

import (
	"time"
)

func (s *server) SetupRouter() {
	s.authRouter()
}

func (s *server) authRouter() {
	route := s.app.Group("/api/auth")
	route.Post("/register", s.RegisterUser, s.RateLimiter(2, time.Minute), s.registerValidate)
	route.Post("/login", s.Login, s.RateLimiter(2, time.Minute), s.loginValidate)
	route.Patch("/update-password", s.UpdatePassword, s.RateLimiter(2, time.Minute), s.updatePasswordValidate, s.IsAuthorized)
}
