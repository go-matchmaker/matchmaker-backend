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
	UserRepo    repository.UserMaker
	UserService service.UserMaker
}

func New(
	rw *sync.RWMutex,
	Cfg *config.Container,
	HTTP http.ServerMaker,
	Token token.TokenMaker,
	PG db.EngineMaker,
	Dragonfly cache.EngineMaker,
	UserRepo repository.UserMaker,
	UserService service.UserMaker,
) *App {
	return &App{
		rw:          rw,
		Cfg:         Cfg,
		HTTP:        HTTP,
		Token:       Token,
		PG:          PG,
		Dragonfly:   Dragonfly,
		UserRepo:    UserRepo,
		UserService: UserService,
	}
}

func (a *App) Run(ctx context.Context) {
	zap.S().Info("RUNNER!")
}
