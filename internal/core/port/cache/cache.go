package cache

import (
	"context"
	"time"
)

type EngineMaker interface {
	Start(ctx context.Context) error
	Close() error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
}
