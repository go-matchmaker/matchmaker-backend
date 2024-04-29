package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EngineMaker interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	GetDB() *pgxpool.Pool
	Migration() error
	Drop() error
}
