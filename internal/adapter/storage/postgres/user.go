package psql

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
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	userArg.ID = id
	userEntity, err := r.querier.InsertUser(ctx, r.db, &userArg)
	return &userEntity.ID, err
}

func (r *UserRepository) Update(ctx context.Context, userModel *entity.User) (*entity.User, error) {
	userConv := converter.UserModelToUpdateArg(userModel)
	userM, err := r.querier.UpdateUser(ctx, r.db, userConv)
	if err != nil {
		return nil, err
	}
	return converter.ArgToUserModel(userM), nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	userM, err := r.querier.GetUser(ctx, r.db, id)
	if err != nil {
		return entity.User{}, err
	}
	return *converter.ArgToUserModel(userM), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	userM, err := r.querier.GetUserByEmail(ctx, r.db, email)
	if err != nil {
		return entity.User{}, err
	}
	return *converter.ArgToUserModel(userM), nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.querier.DeleteUser(ctx, r.db, id)
}

func (r *UserRepository) DeleteAllUsers(ctx context.Context) error {
	return r.querier.DeleteAllUsers(ctx, r.db)
}
