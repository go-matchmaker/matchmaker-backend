package service

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/cache"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/repository"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/service"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/token"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/google/wire"
)

var (
	_              service.UserPort = (*UserService)(nil)
	UserServiceSet                  = wire.NewSet(NewAuthService)
)

type UserService struct {
	userRepo repository.UserPort
	cache    cache.EngineMaker
	token    token.TokenMaker
}

func NewAuthService(userRepo repository.UserPort, cache cache.EngineMaker, token token.TokenMaker) service.UserPort {
	return &UserService{
		userRepo,
		cache,
		token,
	}
}

func (as *UserService) Register(ctx context.Context, userModel *entity.User) (*uuid.UUID, error) {

	id, err := as.userRepo.Insert(ctx, userModel)
	if err != nil {
		return nil, err
	}

	cachingKey := util.GenerateCacheKey("user", userModel.ID)
	userSerialized, err := json.Marshal(userModel)
	if err != nil {
		return nil, err
	}

	err = as.cache.Set(ctx, cachingKey, userSerialized, 0)
	if err != nil {
		return nil, err
	}
	err = as.cache.DeleteByPrefix(ctx, "users:*") // delete all users cache because of new one
	if err != nil {
		return nil, err

	}

	return id, nil
}

func (as *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := as.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	cachingKey := util.GenerateCacheKey("user", id)
	err = as.cache.Delete(ctx, cachingKey)
	if err != nil {
		return err
	}

	return nil
}

//
// func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
// 	user, err := as.userRepo.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return "", &entity.Error{
// 			Code:    entity.InvalidCredentials,
// 			Message: "Invalid credentials",
// 			Data:    email,
// 		}
// 	}
//
// 	passErr := util.ComparePassword(password, user.Password)
// 	if passErr != nil {
// 		return "", &entity.Error{
// 			Code:    entity.InvalidCredentials,
// 			Message: "Invalid credentials",
// 			Data:    user.Password,
// 		}
// 	}
//
// 	token, err := as.token.CreateToken(user)
// 	if err != nil {
// 		return "", &entity.Error{
// 			Code:    entity.TokenCreation,
// 			Message: "Token creation failed",
// 		}
// 	}
//
// 	return token, nil
// }
