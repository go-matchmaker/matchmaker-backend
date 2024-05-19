package http

import (
	"context"
)

type ServerMaker interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	HTTPMiddleware() error
	SetupRouter()
}
