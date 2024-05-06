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
	"time"
)

var (
	_              service.UserPort = (*UserService)(nil)
	UserServiceSet                  = wire.NewSet(NewUserService)
)

type UserService struct {
	tokenDuration   time.Duration
	refreshDuration time.Duration
	userRepo        repository.UserPort
	cache           cache.EngineMaker
	token           token.TokenMaker
}

func NewUserService(tokenTTL time.Duration, refreshTTL time.Duration, userRepo repository.UserPort, cache cache.EngineMaker, token token.TokenMaker) service.UserPort {
	return &UserService{
		tokenTTL,
		refreshTTL,
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

func (as *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = util.ComparePassword(password, user.PasswordHash)
	if err != nil {
		return "", err
	}

	accessToken, accessPayload, err := as.token.CreateToken(user.Email, user.Role, as.tokenDuration)
	if err != nil {
		return "", err
	}

	refreshToken, refreshPayload, err := as.token.CreateToken(user.Email, user.Role, as.refreshDuration)
	if err != nil {
		return "", err
	}

}
