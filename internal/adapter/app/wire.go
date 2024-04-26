//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"github.com/bulutcan99/company-matcher/internal/adapter/auth/paseto"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	adapter_service "github.com/bulutcan99/company-matcher/internal/adapter/service"
	psql "github.com/bulutcan99/company-matcher/internal/adapter/storage/postgres"
	"github.com/bulutcan99/company-matcher/internal/adapter/storage/postgres/repository"
	"github.com/bulutcan99/company-matcher/internal/adapter/storage/redis"
	adapter_http "github.com/bulutcan99/company-matcher/internal/adapter/transport/http"
	"github.com/bulutcan99/company-matcher/internal/core/port/cache"
	"github.com/bulutcan99/company-matcher/internal/core/port/db"
	"github.com/bulutcan99/company-matcher/internal/core/port/http"
	"github.com/bulutcan99/company-matcher/internal/core/port/service"
	"github.com/bulutcan99/company-matcher/internal/core/port/token"
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sync"
)

func InitApp(
	ctx context.Context,
	wg *sync.WaitGroup,
	rw *sync.RWMutex,
	eg *errgroup.Group,
	Cfg *config.Container,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		redisEngineFunc,
		repository.UserRepositorySet,
		adapter_service.UserServiceSet,
		paseto.PasetoSet,
		httpServerFunc,
	))
}

func dbEngineFunc(
	ctx context.Context,
	eg *errgroup.Group,
	Cfg *config.Container) (db.EngineMaker, func(), error) {
	psqlDb := psql.NewDB(eg, Cfg)
	err := psqlDb.Start(ctx)
	if err != nil {
		zap.S().Fatal("failed to start db:", err)
	}

	psqlDb.Migration()

	return psqlDb, func() { psqlDb.Close() }, nil
}

func redisEngineFunc(
	ctx context.Context,
	eg *errgroup.Group,
	Cfg *config.Container) (cache.EngineMaker, func(), error) {
	redisEngine := redis.NewRedisCache(eg, Cfg)
	err := redisEngine.Start(ctx)
	if err != nil {
		zap.S().Fatal("failed to start redis:", err)
	}

	return redisEngine, func() { redisEngine.Close() }, nil
}

func httpServerFunc(
	ctx context.Context,
	eg *errgroup.Group,
	Cfg *config.Container,
	UserService service.UserMaker,
	tokenMaker token.TokenMaker,
) (http.ServerMaker, func(), error) {
	httpServer := adapter_http.NewHttpServer(eg, Cfg, UserService, tokenMaker)
	httpServer.Start(ctx)
	httpServer.Config()
	err := httpServer.HttpMiddleware()
	if err != nil {
		zap.S().Fatal("middleware error:", err)
	}

	httpServer.SetupRouter()
	return httpServer, func() { httpServer.Close() }, nil
}
