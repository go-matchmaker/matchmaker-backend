package http

import (
	"context"
)

type ServerMaker interface {
	Start(ctx context.Context)
	Close()
	Config()
	HttpMiddleware() error
	AuthMiddleware() error
	SetupRouter()
}
