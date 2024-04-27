package repository

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/converter"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/repository"
	"github.com/google/uuid"
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

func (r *UserRepository) Insert(ctx context.Context, userModel *entity.User) (*uuid.UUID, error) {
	userConv := converter.UserModelToArg(userModel)
	userArg := user.InsertUserParams(*userConv)
	userEntity, err := r.querier.InsertUser(ctx, r.db, &userArg)
	return &userEntity.ID, err
}
