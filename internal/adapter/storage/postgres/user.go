package psql

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/converter"
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/user"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_                 user.UserRepositoryPort = (*UserRepository)(nil)
	UserRepositorySet                         = wire.NewSet(NewUserRepository)
)

type UserRepository struct {
	querier user_sql.Querier
	db      *pgxpool.Pool
}

func NewUserRepository(db db.EngineMaker) user.UserRepositoryPort {
	return &UserRepository{
		querier: user_sql.New(),
		db:      db.GetDB(),
	}
}

func (r *UserRepository) Insert(ctx context.Context, userModel *entity.User) (*uuid.UUID, error) {
	userConv := converter.UserModelToArg(userModel)
	userArg := user_sql.InsertParams(*userConv)
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	userArg.ID = id
	userEntity, err := r.querier.Insert(ctx, r.db, &userArg)
	return &userEntity.ID, err
}

func (r *UserRepository) Update(ctx context.Context, userModel *entity.User) (*entity.User, error) {
	userConv := converter.UserModelToUpdateArg(userModel)
	userM, err := r.querier.Update(ctx, r.db, userConv)
	if err != nil {
		return nil, err
	}
	return converter.ArgToUserModel(userM), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	userM, err := r.querier.GetByID(ctx, r.db, id)
	if err != nil {
		return entity.User{}, err
	}
	return *converter.ArgToUserModel(userM), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	userM, err := r.querier.GetByEmail(ctx, r.db, email)
	if err != nil {
		return entity.User{}, err
	}
	return *converter.ArgToUserModel(userM), nil
}

func (r *UserRepository) DeleteOne(ctx context.Context, id uuid.UUID) error {
	return r.querier.DeleteOne(ctx, r.db, id)
}

func (r *UserRepository) DeleteAll(ctx context.Context) error {
	return r.querier.DeleteAll(ctx, r.db)
}
