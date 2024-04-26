package repository

import (
	"context"
	"github.com/bulutcan99/company-matcher/internal/adapter/converter"
	"github.com/bulutcan99/company-matcher/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/bulutcan99/company-matcher/internal/core/domain/entity"
	"github.com/bulutcan99/company-matcher/internal/core/port/db"
	"github.com/bulutcan99/company-matcher/internal/core/port/repository"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_                 repository.UserMaker = (*UserRepository)(nil)
	UserRepositorySet                      = wire.NewSet(NewUserRepository)
)

type UserRepository struct {
	querier user.Querier
	db      *pgxpool.Pool
}

func NewUserRepository(db db.EngineMaker) repository.UserMaker {
	return &UserRepository{
		querier: user.New(),
		db:      db.GetDB(),
	}
}

func (r *UserRepository) Insert(ctx context.Context, userModel *entity.User) error {
	userConv := converter.UserModelToArg(userModel)
	userArg := user.InsertUserParams(*userConv)
	_, err := r.querier.InsertUser(ctx, r.db, &userArg)
	return err
}
