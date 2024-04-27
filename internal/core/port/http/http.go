package http

import (
	"context"
)

type ServerMaker interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Config()
	HTTPMiddleware() error
	AuthMiddleware() error
	SetupRouter()
}
