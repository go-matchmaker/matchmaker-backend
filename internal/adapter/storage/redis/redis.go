package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	"github.com/bulutcan99/company-matcher/internal/core/port/cache"
	"github.com/google/wire"
	go_redis "github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	_        cache.EngineMaker = (*redis)(nil)
	RedisSet                   = wire.NewSet(NewRedisCache)
)

type redis struct {
	eg     *errgroup.Group
	cfg    *config.Container
	client *go_redis.Client
}

func NewRedisCache(eg *errgroup.Group, cfg *config.Container) cache.EngineMaker {
	r := &redis{
		eg:  eg,
		cfg: cfg,
	}

	return r
}

func (r *redis) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", r.cfg.Redis.Host, r.cfg.Redis.Port)
	dbNumber := r.cfg.Redis.DbNumber
	password := r.cfg.Redis.Password

	var client *go_redis.Client
	var pingErr error
	r.eg.Go(func() error {
		client = go_redis.NewClient(&go_redis.Options{
			Addr:     address,
			DB:       dbNumber,
			Password: password,
		})
		return nil
	})

	r.eg.Go(func() error {
		if client != nil {
			_, pingErr = r.client.Ping(ctx).Result()
			if pingErr != nil {
				return pingErr
			}
			return nil
		}
		return errors.New("redis client is nil")
	})

	if err := r.eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (r *redis) Close() error {
	return r.client.Close()
}

func (r *redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Result()
	bytes := []byte(res)
	return bytes, err
}

func (r *redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.client.Del(ctx, key).Err()
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

func (r *redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}
