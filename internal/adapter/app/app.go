package app

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/http"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/repository"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/service"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/token"
	"go.uber.org/zap"
	"sync"
)

type App struct {
	rw          *sync.RWMutex
	Cfg         *config.Container
	HTTP        http.ServerMaker
	Token       token.TokenMaker
	PG          db.EngineMaker
	Dragonfly   cache.EngineMaker
	UserRepo    repository.UserPort
	UserService service.UserPort
}

func New(
	rw *sync.RWMutex,
	cfg *config.Container,
	http http.ServerMaker,
	token token.TokenMaker,
	pg db.EngineMaker,
	dragonfly cache.EngineMaker,
	userRepo repository.UserPort,
	userService service.UserPort,
) *App {
	return &App{
		rw:          rw,
		Cfg:         cfg,
		HTTP:        http,
		Token:       token,
		PG:          pg,
		Dragonfly:   dragonfly,
		UserRepo:    userRepo,
		UserService: userService,
	}
}

func (a *App) Run(ctx context.Context) {
	zap.S().Info("RUNNER!")
}
