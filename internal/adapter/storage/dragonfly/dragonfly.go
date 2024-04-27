package dragonfly

import (
	"context"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var (
	_            cache.EngineMaker = (*dragonfly)(nil)
	dragonflySet                   = wire.NewSet(NewDragonflyCache)
)

type dragonfly struct {
	cfg    *config.Container
	client *redis.Client
}

func NewDragonflyCache(cfg *config.Container) cache.EngineMaker {
	d := &dragonfly{
		cfg: cfg,
	}

	return d
}

func (d *dragonfly) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", d.cfg.Dragonfly.Host, d.cfg.Dragonfly.Port)
	dbNumber := d.cfg.Dragonfly.DBNumber

	var pingErr error
	go func() {
		d.client = redis.NewClient(&redis.Options{
			Addr: address,
			DB:   dbNumber,
		})
		zap.S().Info("Connecting to Dragonfly...")
		_, pingErr = d.client.Ping(ctx).Result()
		zap.S().Info("DragonFly Pong! 🐉")
		if pingErr != nil {
			zap.S().Fatal("Dragonfly ping failed", pingErr)
		}
	}()

	return nil
}

func (d *dragonfly) Close(ctx context.Context) error {
	zap.S().Info("Dragonfly Context is done. Shutting down server...")
	if err := d.client.Close(); err != nil {
		zap.S().Error("server shutdown error: %w", err)
		return err
	}
	return nil
}

func (d *dragonfly) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return d.client.Set(ctx, key, value, ttl).Err()
}

func (d *dragonfly) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := d.client.Get(ctx, key).Result()
	bytes := []byte(res)
	return bytes, err
}

func (d *dragonfly) Delete(ctx context.Context, key string) error {
	return d.client.Del(ctx, key).Err()
}

func (d *dragonfly) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = d.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := d.client.Del(ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (d *dragonfly) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := d.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}
