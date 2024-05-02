//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/auth/paseto"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	adapter_service "github.com/go-matchmaker/matchmaker-server/internal/adapter/service"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/dragonfly"
	psql "github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres"
	adapter_http "github.com/go-matchmaker/matchmaker-server/internal/adapter/transport/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/service"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/token"
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
		psql.UserRepositorySet,
		adapter_service.UserServiceSet,
		paseto.PasetoSet,
		httpServerFunc,
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
	err = psqlDb.Migration()
	if err != nil {
		zap.S().Fatal("failed to migrate db:", err)
	}

	return psqlDb, func() { psqlDb.Close(ctx) }, nil
}

func dragonflyEngineFunc(
	ctx context.Context,
	Cfg *config.Container) (cache.EngineMaker, func(), error) {
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
	UserService service.UserPort,
	tokenMaker token.TokenMaker,
) (http.ServerMaker, func(), error) {
	httpServer := adapter_http.NewHTTPServer(ctx, Cfg, UserService, tokenMaker)
	err := httpServer.Start(ctx)
	if err != nil {
		zap.S().Fatal("failed to start http server:", err)
	}
	httpServer.Config()
	//err = httpServer.HTTPMiddleware()
	//if err != nil {
	//	zap.S().Fatal("middleware error:", err)
	//}

	httpServer.SetupRouter()
	return httpServer, func() { httpServer.Close(ctx) }, nil
}
