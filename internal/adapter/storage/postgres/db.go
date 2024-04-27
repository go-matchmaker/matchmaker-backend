package psql

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/config"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

var (
	_           db.EngineMaker = (*postgres)(nil)
	PostgresSet                = wire.NewSet(NewDB)
)

type (
	postgres struct {
		cfg          *config.Container
		queryBuilder *squirrel.StatementBuilderType
		pool         *pgxpool.Pool
	}
)

func NewDB(cfg *config.Container) db.EngineMaker {
	psqlDB := &postgres{
		cfg: cfg,
	}
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	psqlDB.queryBuilder = &queryBuilder

	return psqlDB
}

func (ps *postgres) Start(ctx context.Context) error {
	url := ps.getURL()

	go func() {
		err := ps.connect(ctx, url)
		if err != nil {
			zap.S().Fatal("PostgreSQL connection failed", err)
		}
	}()

	zap.S().Info("Connected to PostgreSQL 🎉")
	return nil
}

func (ps *postgres) getURL() string {
	url := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		ps.cfg.PSQL.Conn,
		ps.cfg.PSQL.User,
		ps.cfg.PSQL.Password,
		ps.cfg.PSQL.Host,
		ps.cfg.PSQL.Port,
		ps.cfg.PSQL.Name,
	)
	return url
}

func (ps *postgres) ping(ctx context.Context) error {
	if ps.pool != nil {
		if err := ps.pool.Ping(ctx); err != nil {
			return err
		}
	}
	zap.S().Info("PostgreSQL is ready to serve")
	return nil
}

func (ps *postgres) connect(ctx context.Context, url string) error {
	var lastErr error
	for ps.cfg.Settings.PSQLConnAttempts > 0 {
		zap.S().Info("Connecting to PostgreSQL...")
		ps.pool, lastErr = pgxpool.New(ctx, url)
		if lastErr == nil {
			err := ps.ping(ctx)
			if err == nil {
				zap.S().Info("PostgreSQL Pong! 🐘")
				return nil
			}
		}

		ps.cfg.Settings.PSQLConnAttempts--
		zap.S().Warnf("PostgreSQL connection failed, attempts left: %d", ps.cfg.Settings.PSQLConnAttempts)
		time.Sleep(time.Duration(ps.cfg.Settings.PSQLConnTimeout) * time.Second)
	}
	return lastErr
}

func (ps *postgres) Close(ctx context.Context) error {
	zap.S().Info("Postgres Context is done. Shutting down server...")
	ps.pool.Close()
	return nil
}

func (ps *postgres) GetDB() *pgxpool.Pool {
	return ps.pool
}
