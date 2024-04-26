package app

import (
	"context"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	"github.com/bulutcan99/company-matcher/internal/core/port/cache"
	"github.com/bulutcan99/company-matcher/internal/core/port/db"
	"github.com/bulutcan99/company-matcher/internal/core/port/http"
	"github.com/bulutcan99/company-matcher/internal/core/port/repository"
	"github.com/bulutcan99/company-matcher/internal/core/port/service"
	"github.com/bulutcan99/company-matcher/internal/core/port/token"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sync"
)

type App struct {
	rw          *sync.RWMutex
	eg          *errgroup.Group
	Cfg         *config.Container
	HTTP        http.ServerMaker
	Token       token.TokenMaker
	PG          db.EngineMaker
	Redis       cache.EngineMaker
	UserRepo    repository.UserMaker
	UserService service.UserMaker
}

func New(
	rw *sync.RWMutex,
	eg *errgroup.Group,
	Cfg *config.Container,
	HTTP http.ServerMaker,
	Token token.TokenMaker,
	PG db.EngineMaker,
	Redis cache.EngineMaker,
	UserRepo repository.UserMaker,
	UserService service.UserMaker,
) *App {
	return &App{
		rw:          rw,
		eg:          eg,
		Cfg:         Cfg,
		HTTP:        HTTP,
		Token:       Token,
		PG:          PG,
		Redis:       Redis,
		UserRepo:    UserRepo,
		UserService: UserService,
	}
}

func (a *App) Run(ctx context.Context) {
	zap.S().Info("RUNNER!")
}
