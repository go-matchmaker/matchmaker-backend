//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/auth/paseto"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/dragonfly"
	psql "github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres"
	adapter_http "github.com/go-matchmaker/matchmaker-server/internal/adapter/transport/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/auth"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/user"
	port_service "github.com/go-matchmaker/matchmaker-server/internal/core/service"
	"github.com/google/wire"
	"go.uber.org/zap"
	"sync"
)

func InitApp(
	ctx context.Context,
	wg *sync.WaitGroup,
	rw *sync.RWMutex,
	Cfg *config.Container,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		dragonflyEngineFunc,
		httpServerFunc,
		psql.UserRepositorySet,
		port_service.UserServiceSet,
		dragonfly.SessionCacheSet,
		port_service.SessionServiceSet,
		paseto.PasetoSet,
	))
}

func dbEngineFunc(
	ctx context.Context,
	Cfg *config.Container) (db.EngineMaker, func(), error) {
	psqlDb := psql.NewDB(Cfg)
	err := psqlDb.Start(ctx)
	if err != nil {
		zap.S().Fatal("failed to start db:", err)
	}

	return psqlDb, func() { psqlDb.Close(ctx) }, nil
}

func dragonflyEngineFunc(
	ctx context.Context,
	Cfg *config.Container) (cache.CacheEngine, func(), error) {
	redisEngine := dragonfly.NewDragonflyCache(Cfg)
	err := redisEngine.Start(ctx)
	if err != nil {
		zap.S().Fatal("failed to start dragonfly:", err)
	}

	return redisEngine, func() { redisEngine.Close(ctx) }, nil
}

func httpServerFunc(
	ctx context.Context,
	Cfg *config.Container,
	UserService user.UserServicePort,
	tokenMaker auth.TokenMaker,
	authService auth.SessionServicePort,
) (http.ServerMaker, func(), error) {
	httpServer := adapter_http.NewHTTPServer(ctx, Cfg, UserService, tokenMaker, authService)
	err := httpServer.Start(ctx)
	if err != nil {
		return nil, nil, err
	}
	return httpServer, func() { httpServer.Close(ctx) }, nil
}
