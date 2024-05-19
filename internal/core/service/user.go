package service

import (
	"context"
	"errors"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/aggregate"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/auth"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/user"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/google/uuid"
	"github.com/google/wire"
)

var (
	_              user.UserServicePort = (*UserService)(nil)
	UserServiceSet                      = wire.NewSet(NewUserService)
)

type UserService struct {
	userRepo user.UserRepositoryPort
	cache    cache.CacheEngine
	token    auth.TokenMaker
}

func NewUserService(userRepo user.UserRepositoryPort, cache cache.CacheEngine, token auth.TokenMaker) user.UserServicePort {
	return &UserService{
		userRepo,
		cache,
		token,
	}
}

func (us *UserService) Register(ctx context.Context, userModel *entity.User) (*uuid.UUID, error) {
	id, err := us.userRepo.Insert(ctx, userModel)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (us *UserService) Login(ctx context.Context, email, password, ip string) (*aggregate.Session, error) {
	userModel, err := us.userRepo.GetByEmail(ctx, email)
	sessionModel := new(aggregate.Session)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = util.ComparePassword(password, userModel.PasswordHash)
	if err != nil {
		return nil, errors.New("password not match")
	}

	// We will take it from db later.
	isBlocked := false
	accessToken, publicKey, accessPayload, err := us.token.CreateToken(userModel.ID, userModel.Email, userModel.Role, isBlocked)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshPublicKey, refreshPayload, err := us.token.CreateRefreshToken(accessPayload)
	if err != nil {
		return nil, err
	}

	sessionModel = aggregate.NewSession(&userModel, refreshPayload, accessToken, publicKey, refreshToken, refreshPublicKey, ip)

	return sessionModel, nil
}

func (us *UserService) UpdatePassword(ctx context.Context, id uuid.UUID, password string) (entity.User, error) {
	userModel := new(entity.User)
	userModel.PasswordHash = password
	userUpdate, err := us.userRepo.Update(ctx, userModel)
	if err != nil {
		return entity.User{}, err
	}

	return *userUpdate, nil
}
