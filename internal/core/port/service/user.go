package service

import (
	"context"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/google/uuid"
)

type UserMaker interface {
	Register(ctx context.Context, req *entity.User) (*uuid.UUID, error)
	// Login(ctx context.Context, req *dto.UserRegister) (*domain.User, string, string, error)
	// GetUserByID(ctx context.Context, id string) (*model.User, error)
	// RefreshToken(ctx context.Context, userID string) (string, error)
	// ChangePassword(ctx context.Context, id string, req *dto.ChangePasswordReq) error
}
