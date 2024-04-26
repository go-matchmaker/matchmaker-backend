package http

import (
	go_json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"time"
)

func (s *server) Config() {
	s.cfgFiber = &fiber.Config{
		ReadTimeout:   time.Minute * time.Duration(s.cfg.Settings.ServerReadTimeout),
		StrictRouting: false,
		CaseSensitive: false,
		BodyLimit:     4 * 1024 * 1024,
		JSONEncoder:   go_json.Marshal,
		JSONDecoder:   go_json.Unmarshal,
		AppName:       "Go-Matchmaker",
		Immutable:     true,
	}
}
