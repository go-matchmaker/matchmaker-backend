package app

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/auth"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/user"
	"go.uber.org/zap"
	"sync"
)

type App struct {
	rw          *sync.RWMutex
	Cfg         *config.Container
	HTTP        http.ServerMaker
	Token       auth.TokenMaker
	PG          db.EngineMaker
	Dragonfly   cache.CacheEngine
	UserRepo    user.UserRepositoryPort
	UserService user.UserServicePort
	AuthRepo    auth.SessionRepositoryPort
	AuthService auth.SessionServicePort
}

func New(
	rw *sync.RWMutex,
	Cfg *config.Container,
	HTTP http.ServerMaker,
	Token auth.TokenMaker,
	PG db.EngineMaker,
	Dragonfly cache.CacheEngine,
	UserRepo user.UserRepositoryPort,
	UserService user.UserServicePort,
	AuthRepo auth.SessionRepositoryPort,
	AuthService auth.SessionServicePort) *App {
	return &App{
		rw:          rw,
		Cfg:         Cfg,
		HTTP:        HTTP,
		Token:       Token,
		PG:          PG,
		Dragonfly:   Dragonfly,
		UserRepo:    UserRepo,
		UserService: UserService,
		AuthRepo:    AuthRepo,
		AuthService: AuthService,
	}
}

func (a *App) Run(ctx context.Context) {
	zap.S().Info("RUNNER!")
}
