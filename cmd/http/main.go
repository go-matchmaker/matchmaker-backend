package main

import (
	"context"
	"github.com/bulutcan99/company-matcher/internal/adapter/app"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	"github.com/bulutcan99/company-matcher/internal/adapter/logger"
	"github.com/bulutcan99/company-matcher/internal/core/util"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"sync"
)

func main() {
	_, err := maxprocs.Set()
	if err != nil {
		panic("failed set max procs")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), util.InterruptSignals...)
	defer cancel()
	wg := new(sync.WaitGroup)
	rw := new(sync.RWMutex)
	eg, ctx := errgroup.WithContext(ctx)
	cfg, err := config.NewConfig()
	if err != nil {
		panic("failed get config: " + err.Error())
	}

	Logger := logger.InitLogger(cfg.Log.Level)
	defer Logger.Sync()

	cleanup := prepareApp(ctx, wg, rw, eg, cfg)
	zap.S().Info("⚡ Service name:", cfg.Name)
	<-ctx.Done()
	zap.S().Info("Context done signal received, shutting down")
	wg.Wait()
	err = eg.Wait()
	if err != nil {
		panic("error from wait group: " + err.Error())
	}
	zap.S().Info("Waiting for all goroutines to finish")
	cleanup()
	zap.S().Info("Clean-up done")
	zap.S().Info("Shutting down successfully")

}

func prepareApp(ctx context.Context, wg *sync.WaitGroup, rw *sync.RWMutex, eg *errgroup.Group, cfg *config.Container) func() {
	var errMsg error
	a, cleanUp, errMsg := app.InitApp(ctx, wg, rw, eg, cfg)
	if errMsg != nil {
		zap.S().Error("failed init app", errMsg)
		<-ctx.Done()
	}
	a.Run(ctx)
	return cleanUp
}
