package psql

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bulutcan99/company-matcher/internal/adapter/config"
	"github.com/bulutcan99/company-matcher/internal/core/port/db"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	_           db.EngineMaker = (*postgres)(nil)
	PostgresSet                = wire.NewSet(NewDB)
)

type (
	postgres struct {
		eg           *errgroup.Group
		cfg          *config.Container
		queryBuilder *squirrel.StatementBuilderType
		pool         *pgxpool.Pool
	}
)

func NewDB(eg *errgroup.Group, cfg *config.Container) db.EngineMaker {
	psqlDB := &postgres{
		eg:  eg,
		cfg: cfg,
	}
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	psqlDB.queryBuilder = &queryBuilder

	return psqlDB
}

func (ps *postgres) Start(ctx context.Context) error {
	url := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		ps.cfg.PSQL.Conn,
		ps.cfg.PSQL.User,
		ps.cfg.PSQL.Password,
		ps.cfg.PSQL.Host,
		ps.cfg.PSQL.Port,
		ps.cfg.PSQL.Name,
	)

	ps.eg.Go(func() error {
		return ps.connect(ctx, url)
	})

	ps.eg.Go(func() error {
		return ps.ping(ctx)
	})

	if err := ps.eg.Wait(); err != nil {
		zap.S().Fatal("PostgreSQL connection failed: ", err)
		return err
	}

	zap.S().Info("Connected to PostgreSQL 🎉")
	return nil
}

func (ps *postgres) ping(ctx context.Context) error {
	if ps.pool != nil {
		if err := ps.pool.Ping(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ps *postgres) connect(ctx context.Context, url string) error {
	var lastErr error
	for ps.cfg.Settings.PSQLConnAttempts > 0 {
		ps.pool, lastErr = pgxpool.New(ctx, url)
		if lastErr == nil {
			return nil
		}

		ps.cfg.Settings.PSQLConnAttempts--
		zap.S().Warnf("PostgreSQL connection failed, attempts left: %d", ps.cfg.Settings.PSQLConnAttempts)
		time.Sleep(time.Duration(ps.cfg.Settings.PSQLConnTimeout) * time.Second)
	}
	return lastErr
}

func (ps *postgres) Close() error {
	if ps.pool != nil {
		ps.pool.Close()
	}

	zap.S().Info("Disconnected from PostgreSQL")
	return nil
}

func (ps *postgres) GetDB() *pgxpool.Pool {
	return ps.pool
}

func (ps *postgres) GetURL() string {
	return ps.pool.Config().ConnString()
}
